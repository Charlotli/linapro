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
const stamp = "420000000000";
const roomCode = `ROOM-H5-${stamp}`;
const roomName = `H5播放测试直播间_${stamp}`;
const liveTitle = `H5播放测试直播_${stamp}`;
const liveUrl = `https://example.com/hls/tc004-${stamp}.m3u8`;

test.describe("TC004 观众端 H5 播放页", () => {
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

  test("TC004a: 未开始直播显示预告且播放器不出现", async ({ adminPage }) => {
    // One room plus one not-started public live bound to it.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName, liveUrl);

    // Anonymous viewer context: the H5 page must not require login.
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/未开始|预告/);
    await h5.expectTitle(new RegExp(liveTitle.slice(0, 8)));
    expect(await h5.isPlayerVisible()).toBe(false);
    // 三次迭代：直播间名称眉标与封面开播日期是预告态核心信息。
    expect(await h5.roomLineText()).toContain(roomName);
    expect(await h5.coverDateText()).toContain("开播");
    // FB-1 regression: absolute asset paths must keep styles and script
    // working on the canonical entry form.
    await h5.expectStylesApplied();
  });

  test("TC004b: 开启直播后页面呈现播放器并附加播放流", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.startLive(liveTitle);

    // The player polls every 30s; reload to fetch fresh play info instead.
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/直播中/);
    expect(await h5.isPlayerVisible()).toBe(true);
    expect(await h5.isLiveBadgeVisible()).toBe(true);
    // The HLS URL is attached through hls.js (MSE) or a native src; both
    // surface through the visible player with the live badge present.
    expect(await h5.hasAttachedSource(".m3u8")).toBe(true);

    // FB-1 regression: the bare /h5 form (no trailing slash) must render the
    // styled page with the player too.
    await h5.gotoBare(roomCode);
    await h5.expectStylesApplied();
    await h5.expectStateChip(/直播中/);
    expect(await h5.isPlayerVisible()).toBe(true);
  });

  test("TC004c: 公开播放信息接口不泄露推流地址等管理端字段", async ({ adminPage }) => {
    const response = await adminPage.request.get(
      `/x/linapro-live-manage/api/v1/play?roomCode=${encodeURIComponent(roomCode)}&tenantId=0`,
    );
    expect(response.status()).toBe(200);
    const payload = (await response.json()) as {
      code: number;
      data: Record<string, unknown>;
    };
    expect(payload.code).toBe(0);
    const data = payload.data ?? {};
    for (const forbidden of [
      "pushUrl",
      "pageUrl",
      "deviceInfo",
      "reception",
      "isPublic",
      "createdBy",
      "updatedBy",
      "deletedAt",
    ]) {
      expect(data, `must not expose ${forbidden}`).not.toHaveProperty(forbidden);
    }
    expect(String(data["liveUrl"] ?? "")).toContain(".m3u8");
  });

  test("TC004d: 关闭直播后回放可见", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.stopLive(liveTitle);

    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/已结束|回放/);
    expect(await h5.isPlayerVisible()).toBe(true);
  });

  test("TC004e: 私密与不存在的直播间统一返回不存在", async ({ adminPage }) => {
    // Stop-then-private: a finished live set private must disappear for viewers.
    const missing = await adminPage.request.get(
      `/x/linapro-live-manage/api/v1/play?roomCode=NO-SUCH-ROOM-${Date.now()}&tenantId=0`,
    );
    expect(missing.status()).toBe(200);
    const missingPayload = (await missing.json()) as { errorCode?: string };
    expect(missingPayload.errorCode).toBe("LIVE_MANAGE_PLAY_NOT_FOUND");

    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(`NO-SUCH-ROOM-${Date.now()}`);
    await h5.expectNoLiveNotice();
  });

  test("TC004f: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
