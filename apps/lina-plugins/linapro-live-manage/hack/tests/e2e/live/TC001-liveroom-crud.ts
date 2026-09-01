import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

test.describe("TC001 直播间管理 CRUD", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
  });

  const roomCode = `ROOM-TEST-${Date.now()}`;
  const roomName = `测试直播间_${Date.now()}`;
  const roomNameRenamed = `${roomName}_修改`;

  test("TC001a: 直播间管理页展示翻译后的表格标题", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();

    // Assert the translated table title renders instead of a raw i18n key.
    await expect(
      adminPage.getByText(/直播间列表|Live room list/i),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      adminPage.getByText(/plugin\.linapro-live-manage\.roomTableTitle/),
    ).toHaveCount(0);
  });

  test("TC001b: 创建新直播间", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    await expect(
      adminPage.getByText(/新增成功|创建成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC001c: 直播间列表中可见新创建的记录", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();

    const hasRoom = await roomPage.hasRoom(roomName);
    expect(hasRoom).toBeTruthy();
  });

  test("TC001d: 编辑直播间", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.editRoom(roomName, roomNameRenamed);

    await expect(
      adminPage.getByText(/更新成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC001e: 删除直播间", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoom(roomNameRenamed);

    await expect(
      adminPage.getByText(/删除成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });
});
