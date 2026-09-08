import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { execPgSQL } from "@host-tests/support/postgres";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LivePlayH5Page } from "../../pages/LivePlayH5Page";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

// Deterministic stamp keeps lookups stable across worker restarts; the
// beforeAll hard-delete makes reruns idempotent (see TC004 for the pattern).
const stamp = "440000000000";
const roomCode = `ROOM-REPLAY-${stamp}`;
const roomName = `回放库测试直播间_${stamp}`;
const liveTitle = `回放库测试直播_${stamp}`;
const liveUrl = `https://example.com/hls/tc007-${stamp}.m3u8`;
const replaysPath =
  `/x/linapro-live-manage/api/v1/replays` +
  `?roomCode=${encodeURIComponent(roomCode)}&tenantId=0&page=1&pageSize=10`;

interface ReplaysPayload {
  code: number;
  data?: { list: Array<Record<string, unknown>>; total: number };
}

test.describe("TC007 往期回放库", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
    // A failed previous run may have left the shared records behind; hard
    // delete first so the deterministic setup in TC007a can rerun cleanly.
    execPgSQL(
      `DELETE FROM plugin_linapro_live_manage_live WHERE title = '${liveTitle}';\n` +
        `DELETE FROM plugin_linapro_live_manage_room WHERE room_code = '${roomCode}';`,
    );
  });

  test("TC007a: 回放列表接口返回已结束直播且播放回退可用", async ({ adminPage }) => {
    // One finished public live with a recorded play URL: the replay-library
    // entry under test. The replay switch defaults to enabled. Page-scoped
    // fixtures only, so the setup lives in the first test body.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName, liveUrl);
    await livePage.startLive(liveTitle);
    await livePage.stopLive(liveTitle);

    const response = await adminPage.request.get(replaysPath);
    expect(response.status()).toBe(200);
    const payload = (await response.json()) as ReplaysPayload;
    expect(payload.code).toBe(0);
    expect(payload.data?.total).toBe(1);

    const entry = payload.data?.list?.[0];
    expect(String(entry?.["title"] ?? "")).toContain(liveTitle.slice(0, 8));
    expect(String(entry?.["liveUrl"] ?? "")).toContain(".m3u8");

    // The public play endpoint falls back to the finished replay.
    const play = await adminPage.request.get(
      `/x/linapro-live-manage/api/v1/play?roomCode=${encodeURIComponent(roomCode)}&tenantId=0`,
    );
    const playPayload = (await play.json()) as {
      code: number;
      data?: Record<string, unknown>;
    };
    expect(playPayload.code).toBe(0);
    expect(String(playPayload.data?.["liveUrl"] ?? "")).toContain(".m3u8");
  });

  test("TC007b: H5 回放态展示往期回放区块并支持切换", async ({ adminPage }) => {
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);

    await h5.expectStateChip(/已结束|回放/);
    expect(await h5.isReplaysSectionVisible()).toBe(true);
    expect(await h5.replayCardCount()).toBe(1);
    expect(await h5.replayCardTitle(0)).toContain(liveTitle.slice(0, 8));

    await h5.clickReplayCard(0);
    expect(await h5.isReplayBadgeVisible()).toBe(true);
  });

  test("TC007c: 关闭回放开关后列表与播放回退均不可见", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    expect(await livePage.hasColumnHeader("公开回放")).toBe(true);
    expect(await livePage.replaySwitchChecked(liveTitle)).toBe(true);

    await livePage.toggleReplaySwitch(liveTitle);

    const response = await adminPage.request.get(replaysPath);
    const payload = (await response.json()) as ReplaysPayload;
    expect(payload.data?.total).toBe(0);
    expect(payload.data?.list?.length ?? 0).toBe(0);

    // With no ongoing or preview live, the play fallback disappears too.
    const play = await adminPage.request.get(
      `/x/linapro-live-manage/api/v1/play?roomCode=${encodeURIComponent(roomCode)}&tenantId=0`,
    );
    const playPayload = (await play.json()) as { errorCode?: string };
    expect(playPayload.errorCode).toBe("LIVE_MANAGE_PLAY_NOT_FOUND");

    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);
    expect(await h5.isReplaysSectionVisible()).toBe(false);
  });

  test("TC007d: 重新开启回放开关后恢复可见", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.toggleReplaySwitch(liveTitle);
    expect(await livePage.replaySwitchChecked(liveTitle)).toBe(true);

    const response = await adminPage.request.get(replaysPath);
    const payload = (await response.json()) as ReplaysPayload;
    expect(payload.data?.total).toBe(1);
  });

  test("TC007e: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
