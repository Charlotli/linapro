import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

test.describe("TC003 直播开启与关闭", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
  });

  const roomCode = `ROOM-STATE-${Date.now()}`;
  const roomName = `状态测试直播间_${Date.now()}`;
  const liveTitle = `状态测试直播_${Date.now()}`;

  test("TC003a: 开启直播后状态变为进行中", async ({ adminPage }) => {
    // Prepare one room and one not-started live bound to it.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName);

    await livePage.startLive(liveTitle);

    // Assert the success toast and the translated ongoing state tag.
    await expect(
      adminPage.getByText(/直播已开启|Live started/i),
    ).toBeVisible({ timeout: 5000 });
    await livePage.clickSearch();
    const row = adminPage.locator(".vxe-body--row:visible", {
      hasText: liveTitle,
    });
    await expect(row.first()).toBeVisible({ timeout: 10000 });
    await expect(row.first().getByText(/进行中|Ongoing/i)).toBeVisible({
      timeout: 5000,
    });
  });

  test("TC003b: 重复开启被状态机拒绝", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();

    // The row now shows the stop action; starting again is only possible via
    // another not-started record, so assert the ongoing row offers no start.
    await livePage.fillSearchField("直播主题", liveTitle);
    await livePage.clickSearch();
    const row = adminPage.locator(".vxe-body--row:visible", {
      hasText: liveTitle,
    });
    await expect(row.first()).toBeVisible({ timeout: 10000 });
    await expect(
      row.first().getByRole("button", { name: /开\s*启直播|Start Live/i }),
    ).toHaveCount(0);
  });

  test("TC003c: 关闭直播后状态变为已结束", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.stopLive(liveTitle);

    await expect(
      adminPage.getByText(/直播已关闭|Live stopped/i),
    ).toBeVisible({ timeout: 5000 });
    await livePage.clickSearch();
    const row = adminPage.locator(".vxe-body--row:visible", {
      hasText: liveTitle,
    });
    await expect(row.first()).toBeVisible({ timeout: 10000 });
    await expect(row.first().getByText(/已结束|Finished/i)).toBeVisible({
      timeout: 5000,
    });
  });

  test("TC003d: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
