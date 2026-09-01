import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { EquipmentPage } from "../../pages/EquipmentPage";

test.describe("TC001 设备登记 CRUD", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-equipment-manage"]);
  });

  const code = `EQ-E2E-${Date.now()}`;
  const name = `测试设备_${Date.now()}`;

  test("TC001a: 设备页展示翻译后的表格标题", async ({ adminPage }) => {
    const page = new EquipmentPage(adminPage);
    await page.goto();
    await expect(
      adminPage.getByText(/设备台账|Equipment registry/i),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      adminPage.getByText(/plugin\.linapro-equipment-manage\.tableTitle/),
    ).toHaveCount(0);
  });

  test("TC001b: 登记设备", async ({ adminPage }) => {
    const page = new EquipmentPage(adminPage);
    await page.goto();
    await page.createEquipment(code, name);
    await expect(
      adminPage.getByText(/新增成功|创建成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC001c: 列表中可见新设备", async ({ adminPage }) => {
    const page = new EquipmentPage(adminPage);
    await page.goto();
    expect(await page.hasEquipment(name)).toBeTruthy();
  });

  test("TC001d: 删除设备", async ({ adminPage }) => {
    const page = new EquipmentPage(adminPage);
    await page.goto();
    await page.deleteEquipment(name);
    await expect(
      adminPage.getByText(/删除成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });
});
