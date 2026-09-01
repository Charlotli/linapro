import type { Page } from "@host-tests/support/playwright";

import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-oa-approval flow configuration page. */
export class OaFlowPage {
  constructor(private page: Page) {}

  private get modal() {
    return this.page.locator('[role="dialog"]');
  }

  async goto() {
    await this.page.goto("/oa/approval-flow");
    await waitForTableReady(this.page);
  }

  /** Create one approval flow with the given name. */
  async createFlow(flowName: string) {
    await this.page
      .getByRole("button", { name: /新\s*增/ })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    await this.modal.getByPlaceholder("请输入流程名称").first().fill(flowName);

    // Fill the first node approver search and pick the admin option.
    const approverSelect = this.modal
      .getByPlaceholder("请选择审批人")
      .first();
    await approverSelect.click();
    await this.page.keyboard.type("admin", { delay: 20 });
    await this.page
      .locator(".ant-select-dropdown:visible .ant-select-item-option")
      .first()
      .click();

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /** Check whether a flow with the given name is visible. */
  async hasFlow(flowName: string): Promise<boolean> {
    const input = this.page
      .getByLabel(/流程名称|Flow Name|plugin\.linapro-oa-approval\.fields\.flowName/i)
      .first();
    await input.clear();
    await input.fill(flowName);
    await this.page
      .getByRole("button", { name: /搜\s*索|Search/i })
      .first()
      .click();
    await waitForRouteReady(this.page);
    return this.page
      .locator(".vxe-body--row:visible", { hasText: flowName })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Delete a flow by name with confirm. */
  async deleteFlow(flowName: string) {
    const input = this.page
      .getByLabel(/流程名称|Flow Name|plugin\.linapro-oa-approval\.fields\.flowName/i)
      .first();
    await input.clear();
    await input.fill(flowName);
    await this.page
      .getByRole("button", { name: /搜\s*索|Search/i })
      .first()
      .click();
    await waitForRouteReady(this.page);

    const row = this.page.locator(".vxe-body--row:visible", {
      hasText: flowName,
    });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /删\s*除/ })
      .first()
      .click();

    const popconfirm = await waitForConfirmOverlay(this.page);
    const confirmBtn = popconfirm.getByRole("button", {
      name: /确\s*定|OK|是/i,
    });
    if (await confirmBtn.isVisible({ timeout: 2000 }).catch(() => false)) {
      await confirmBtn.click();
    }
    await waitForRouteReady(this.page);
  }
}
