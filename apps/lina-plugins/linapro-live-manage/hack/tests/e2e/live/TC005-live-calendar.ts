import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { execPgSQL } from "@host-tests/support/postgres";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LivePlayH5Page } from "../../pages/LivePlayH5Page";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

// Deterministic stamp: Playwright re-imports this file when a worker restarts
// after a failure, which would regenerate Date.now() constants mid-suite and
// strand later tests on data they cannot see. A fixed stamp keeps every
// lookup stable, and the beforeAll hard-delete below makes reruns idempotent.
const stamp = "430000000000";
const roomCode = `ROOM-ICS-${stamp}`;
const roomName = `日历订阅测试直播间_${stamp}`;
const liveTitle = `日历订阅测试直播_${stamp}`;
const liveUrl = `https://example.com/hls/tc005-${stamp}.m3u8`;
const subscribePath = "/x/linapro-live-manage/api/v1/subscribe";

function subscribeUrl(room: string, tenantId = 0) {
  return `${subscribePath}?roomCode=${encodeURIComponent(room)}&tenantId=${tenantId}`;
}

/** Parse the VEVENT blocks out of an ICS document. */
function parseEvents(document: string): string[] {
  return document
    .split(/END:VEVENT\r?\n/)
    .map((chunk) => (chunk.includes("BEGIN:VEVENT") ? chunk : ""))
    .filter((chunk) => chunk !== "");
}

test.describe("TC005 直播日程日历订阅", () => {
  test.beforeAll(async () => {
    // Worker-scoped fixtures only: no page objects here, because Playwright
    // creates a fresh adminPage per test and the worker restarts on failure.
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
    // A failed previous run may have left the shared records behind, and
    // soft-deleted rows still occupy the (tenant_id, room_code) unique
    // index, so hard-delete the deterministic suite records first.
    execPgSQL(
      `DELETE FROM plugin_linapro_live_manage_live WHERE title = '${liveTitle}';\n` +
        `DELETE FROM plugin_linapro_live_manage_room WHERE room_code = '${roomCode}';`,
    );
  });

  test("TC005a: 订阅接口返回含稳定UID的日历流", async ({ adminPage }) => {
    // One room plus one not-started public live bound to it. The create
    // form leaves the date picker empty and the service defaults it to
    // today, which lands inside the subscription window.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName, liveUrl);

    const response = await adminPage.request.get(subscribeUrl(roomCode));
    expect(response.status()).toBe(200);
    expect(response.headers()["content-type"]).toContain("text/calendar");
    const document = await response.text();

    const events = parseEvents(document);
    expect(events.length).toBe(1);
    // The live has no recorded start time, so it renders as an all-day entry
    // instead of inventing concrete times.
    expect(events[0]).toContain("DTSTART;VALUE=DATE:");
    expect(events[0]).toMatch(/UID:linapro-live-\d+@linapro-live-manage/);
    expect(events[0]).toContain(`SUMMARY:${liveTitle}`);
    // The description points viewers back to the watch page.
    expect(decodeURIComponent(events[0])).toContain(`/x/linapro-live-manage/h5?room=${roomCode}`);

    // Re-fetching must keep the same UID so subscribers never get duplicates.
    const second = await adminPage.request.get(subscribeUrl(roomCode));
    const secondEvents = parseEvents(await second.text());
    const uid = (line: string) => line.split("\r\n").find((l) => l.startsWith("UID:"));
    expect(uid(secondEvents[0])).toBe(uid(events[0]));
  });

  test("TC005b: 预告态展示日历入口并反馈复制结果", async ({ adminPage }) => {
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/未开始|预告/);
    expect(await h5.isCalendarBarVisible()).toBe(true);

    // Copy feedback: the toast reports success or a graceful clipboard
    // failure, either outcome proves the action produced user-visible
    // feedback without depending on browser clipboard permissions.
    await h5.clickCopyCalendarLink();
    await h5.expectToast(/订阅链接已复制|复制失败/);
  });

  test("TC005c: 非预告态不展示日历入口", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.startLive(liveTitle);

    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);
    await h5.expectStateChip(/直播中/);
    expect(await h5.isCalendarBarVisible()).toBe(false);

    // Cleanup into the replay state so the room stays consistent for the
    // admin modal assertions below.
    await livePage.goto();
    await livePage.stopLive(liveTitle);
  });

  test("TC005d: 管理端直播间列表提供订阅链接弹窗", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.openCalendarModal(roomName);

    const link = await roomPage.calendarModalSubscribeLink();
    // Same origin, same params: the link must be the public endpoint with
    // this room's code so it works directly in a system calendar app.
    expect(link).toContain(subscribePath);
    expect(link).toContain(`roomCode=${roomCode}`);
    expect(roomPage.calendarModalHasCopyButton()).resolves.toBe(true);
    // The how-to hint must render as translated copy, not a raw i18n key.
    const hint = await roomPage.calendarModalHintText();
    expect(hint).toMatch(/订阅日历|Subscribe Calendar/);
    expect(hint).not.toContain("plugin.linapro-live-manage");
  });

  test("TC005e: 无效输入统一返回空日历不泄露存在性", async ({ adminPage }) => {
    const missing = await (await adminPage.request.get(subscribeUrl(`NO-SUCH-ROOM-${stamp}`))).text();
    const invalidTenant = await (
      await adminPage.request.get(subscribeUrl(roomCode, 999))
    ).text();

    for (const document of [missing, invalidTenant]) {
      expect(document).toContain("BEGIN:VCALENDAR");
      expect(document).toContain("END:VCALENDAR");
      expect(parseEvents(document).length).toBe(0);
      // Room existence stays undisclosed: no display name on empty feeds.
      expect(document).not.toContain("X-WR-CALNAME");
      expect(document).not.toContain(roomName);
    }
  });

  test("TC005f: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
