import type { Page } from "@host-tests/support/playwright";

import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-equipment-manage equipment page. */
export class EquipmentPage {
  constructor(private page: Page) {}

  private get modal() {
    return this.page.locator('[role="dialog"]');
  }

  async goto() {
    await this.page.goto("/equipment/register");
    await waitForTableReady(this.page);
  }

  async createEquipment(code: string, name: string) {
    await this.page.getByRole("button", { name: /新\s*增/ }).first().click();
    await waitForDialogReady(this.modal);
    await this.modal.getByPlaceholder("请输入设备编号").first().fill(code);
    await this.modal.getByPlaceholder("请输入设备名称").first().fill(name);
    await this.modal.getByRole("button", { name: /确\s*认/ }).click();
    await waitForRouteReady(this.page);
    await this.modal.waitFor({ state: "hidden", timeout: 10000 }).catch(() => {});
  }

  async hasEquipment(name: string): Promise<boolean> {
    const input = this.page
      .getByLabel(/设备名称|Name|plugin\.linapro-equipment-manage\.fields\.name/i)
      .first();
    await input.clear();
    await input.fill(name);
    await this.page.getByRole("button", { name: /搜\s*索|Search/i }).first().click();
    await waitForRouteReady(this.page);
    return this.page
      .locator(".vxe-body--row:visible", { hasText: name })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  async deleteEquipment(name: string) {
    const input = this.page
      .getByLabel(/设备名称|Name|plugin\.linapro-equipment-manage\.fields\.name/i)
      .first();
    await input.clear();
    await input.fill(name);
    await this.page.getByRole("button", { name: /搜\s*索|Search/i }).first().click();
    await waitForRouteReady(this.page);
    const row = this.page.locator(".vxe-body--row:visible", { hasText: name });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row.locator("button:visible").filter({ hasText: /删\s*除/ }).first().click();
    const popconfirm = await waitForConfirmOverlay(this.page);
    const confirmBtn = popconfirm.getByRole("button", { name: /确\s*定|OK|是/i });
    if (await confirmBtn.isVisible({ timeout: 2000 }).catch(() => false)) {
      await confirmBtn.click();
    }
    await waitForRouteReady(this.page);
  }
}
