import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

test.describe("TC002 直播内容 CRUD", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
  });

  const roomCode = `ROOM-LIVE-${Date.now()}`;
  const roomName = `内容测试直播间_${Date.now()}`;
  const liveTitle = `测试直播_${Date.now()}`;
  const liveTitleRenamed = `${liveTitle}_修改`;

  test("TC002a: 直播内容页展示翻译后的表格标题", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();

    // Assert the translated table title renders instead of a raw i18n key.
    await expect(
      adminPage.getByText(/直播内容列表|Live content list/i),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      adminPage.getByText(/plugin\.linapro-live-manage\.liveTableTitle/),
    ).toHaveCount(0);
  });

  test("TC002b: 创建直播内容并关联直播间", async ({ adminPage }) => {
    // Prepare a live room for the content to reference.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName);

    await expect(
      adminPage.getByText(/新增成功|创建成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC002c: 直播内容列表中可见新创建的记录", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();

    const hasLive = await livePage.hasLive(liveTitle);
    expect(hasLive).toBeTruthy();
  });

  test("TC002d: 编辑直播内容", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.editLive(liveTitle, liveTitleRenamed);

    await expect(
      adminPage.getByText(/更新成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC002e: 删除直播内容并清理关联直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLive(liveTitleRenamed);

    await expect(
      adminPage.getByText(/删除成功|success/i),
    ).toBeVisible({ timeout: 5000 });

    // Clean up the live room created in TC002b so the suite leaves no data.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
