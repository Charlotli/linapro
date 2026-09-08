import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import type { APIRequestContext } from "@host-tests/support/playwright";
import { execPgSQL } from "@host-tests/support/postgres";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LivePlayH5Page } from "../../pages/LivePlayH5Page";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

// Deterministic stamp keeps lookups stable across worker restarts; the
// beforeAll hard-delete makes reruns idempotent (see TC004 for the pattern).
const stamp = "430000000000";
const roomCode = `ROOM-SHARE-${stamp}`;
const roomName = `分享统计测试直播间_${stamp}`;
const liveTitle = `分享统计测试直播_${stamp}`;
const liveUrl = `https://example.com/hls/tc006-${stamp}.m3u8`;
const heartbeatPath = "/x/linapro-live-manage/api/v1/view/heartbeat";

interface HeartbeatPayload {
  code: number;
  data?: { onlineCount: number; totalViews: number };
}

async function sendHeartbeat(
  request: APIRequestContext,
  body: Record<string, unknown>,
) {
  const response = await request.post(heartbeatPath, { data: body });
  expect(response.status()).toBe(200);
  return (await response.json()) as HeartbeatPayload;
}

test.describe("TC006 观看统计与分享", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
    execPgSQL(
      `DELETE FROM plugin_linapro_live_manage_live WHERE title = '${liveTitle}';\n` +
        `DELETE FROM plugin_linapro_live_manage_room WHERE room_code = '${roomCode}';`,
    );
  });

  test("TC006a: 心跳接口按会话键去重且无效房间返回零值不落库", async ({ adminPage }) => {
    // Prepare an ongoing public live so heartbeats have a playable target.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName, liveUrl);
    await livePage.startLive(liveTitle);

    const first = await sendHeartbeat(adminPage.request, {
      roomCode,
      tenantId: 0,
      sessionKey: `tc006-session-a-${stamp}`,
    });
    expect(first.code).toBe(0);
    expect(first.data?.totalViews).toBe(1);

    // Same browser session key again: deduplicated, no total growth.
    const repeat = await sendHeartbeat(adminPage.request, {
      roomCode,
      tenantId: 0,
      sessionKey: `tc006-session-a-${stamp}`,
    });
    expect(repeat.data?.totalViews).toBe(1);

    // A different session key counts as one more viewer.
    const second = await sendHeartbeat(adminPage.request, {
      roomCode,
      tenantId: 0,
      sessionKey: `tc006-session-b-${stamp}`,
    });
    expect(second.data?.totalViews).toBe(2);
    expect(second.data?.onlineCount).toBeGreaterThanOrEqual(2);

    // Unknown rooms stay indistinguishable from empty content: zero counters
    // returned, and nothing recorded.
    const unknown = await sendHeartbeat(adminPage.request, {
      roomCode: `NO-SUCH-ROOM-${stamp}`,
      tenantId: 0,
      sessionKey: `tc006-session-x-${stamp}`,
    });
    expect(unknown.code).toBe(0);
    expect(unknown.data?.onlineCount).toBe(0);
    expect(unknown.data?.totalViews).toBe(0);
  });

  test("TC006b: H5 直播中展示在线人数徽标与分享操作条", async ({ adminPage }) => {
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/直播中/);
    // The page heartbeats on load; the badge reflects the sessions recorded
    // in TC006a and the session created by this page load.
    const badgeText = await h5.expectViewersBadge();
    expect(badgeText).toMatch(/\d+ 人在看/);

    // Share bar is rendered for presentable content; desktop environments
    // without the Web Share API degrade to the copy-link label.
    expect(await h5.isShareBarVisible()).toBe(true);
    const hasNativeShare = await adminPage.evaluate(
      () => typeof navigator.share === "function",
    );
    if (!hasNativeShare) {
      expect(await h5.shareButtonLabel()).toMatch(/复制链接/);
      await h5.clickShareButton();
      await h5.expectToast(/链接已复制/);
    }
  });

  test("TC006c: 管理端直播内容列表展示统计列", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();

    expect(await livePage.hasColumnHeader("在线人数")).toBe(true);
    expect(await livePage.hasColumnHeader("累计观看")).toBe(true);
  });

  test("TC006d: 直播间列表二维码弹窗展示观播链接", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.openQrcodeModal(roomName);

    expect(await roomPage.qrcodeModalHasImage()).toBe(true);
    const link = await roomPage.qrcodeModalWatchLink();
    expect(link).toContain("/x/linapro-live-manage/h5");
    expect(link).toContain(`room=${encodeURIComponent(roomCode)}`);
    expect(await roomPage.qrcodeModalHasCopyButton()).toBe(true);
  });

  test("TC006e: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.stopLive(liveTitle);
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
