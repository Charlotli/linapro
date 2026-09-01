import type { Page } from "@host-tests/support/playwright";

import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-oa-approval approval center page. */
export class OaApprovalPage {
  constructor(private page: Page) {}

  private get modal() {
    return this.page.locator('[role="dialog"]');
  }

  private get drawer() {
    return this.page.locator(".vben-drawer, [role='dialog']").last();
  }

  async goto() {
    await this.page.goto("/oa/approval");
    await waitForTableReady(this.page);
  }

  /** Submit one approval request bound to the given flow. */
  async submitRequest(flowName: string, title: string, amount: string) {
    await this.page
      .getByRole("button", { name: /提交审批|Submit/i })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    const flowSelect = this.modal.getByPlaceholder("请选择审批流程").first();
    await flowSelect.click();
    await this.page
      .locator(".ant-select-dropdown:visible .ant-select-item-option", {
        hasText: flowName,
      })
      .first()
      .click();

    await this.modal.getByPlaceholder("请输入审批标题").first().fill(title);
    const amountInput = this.modal.locator(".ant-input-number-input").first();
    if (await amountInput.isVisible().catch(() => false)) {
      await amountInput.fill(amount);
    }

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /** Open the detail drawer of one request searched by title. */
  async openDetail(title: string) {
    const input = this.page
      .getByLabel(/审批标题|Title|plugin\.linapro-oa-approval\.fields\.title/i)
      .first();
    await input.clear();
    await input.fill(title);
    await this.page
      .getByRole("button", { name: /搜\s*索|Search/i })
      .first()
      .click();
    await waitForRouteReady(this.page);

    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /详\s*情|Detail/i })
      .first()
      .click();

    await this.drawer.waitFor({ state: "visible", timeout: 10000 });
  }

  /** Approve inside the open detail drawer with an optional reply. */
  async approveInDrawer(reply: string) {
    const replyBox = this.drawer.getByPlaceholder("请输入回复内容");
    if (await replyBox.isVisible().catch(() => false)) {
      await replyBox.fill(reply);
    }
    await this.drawer
      .getByRole("button", { name: /通\s*过|Approve/i })
      .first()
      .click();
    await waitForRouteReady(this.page);
  }

  /** Append one approver inside the open detail drawer. */
  async appendApproverInDrawer(keyword: string) {
    const select = this.drawer.getByPlaceholder("请选择审批人").first();
    await select.click();
    await this.page.keyboard.type(keyword, { delay: 20 });
    await this.page
      .locator(".ant-select-dropdown:visible .ant-select-item-option")
      .first()
      .click();
    await this.drawer
      .getByRole("button", { name: /添加审批人|Add approver/i })
      .first()
      .click();
    await waitForRouteReady(this.page);
  }

  /** Assert one timeline action label is visible in the drawer. */
  async expectTimelineAction(label: RegExp) {
    await this.drawer.getByText(label).first().waitFor({
      state: "visible",
      timeout: 10000,
    });
  }

  /** Withdraw one pending request from the list by title. */
  async withdraw(title: string) {
    const input = this.page
      .getByLabel(/审批标题|Title|plugin\.linapro-oa-approval\.fields\.title/i)
      .first();
    await input.clear();
    await input.fill(title);
    await this.page
      .getByRole("button", { name: /搜\s*索|Search/i })
      .first()
      .click();
    await waitForRouteReady(this.page);

    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /撤\s*回|Withdraw/i })
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
