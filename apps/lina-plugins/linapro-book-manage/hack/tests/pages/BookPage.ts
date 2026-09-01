import type { Page } from "@host-tests/support/playwright";

import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-book-manage book page. */
export class BookPage {
  constructor(private page: Page) {}

  private get modal() {
    return this.page.locator('[role="dialog"]');
  }

  async goto() {
    await this.page.goto("/book/register");
    await waitForTableReady(this.page);
  }

  async createBook(title: string) {
    await this.page.getByRole("button", { name: /新\s*增/ }).first().click();
    await waitForDialogReady(this.modal);
    await this.modal.getByPlaceholder("请输入书名").first().fill(title);
    await this.modal.getByRole("button", { name: /确\s*认/ }).click();
    await waitForRouteReady(this.page);
    await this.modal.waitFor({ state: "hidden", timeout: 10000 }).catch(() => {});
  }

  async hasBook(title: string): Promise<boolean> {
    const input = this.page
      .getByLabel(/书名|Title|plugin\.linapro-book-manage\.fields\.title/i)
      .first();
    await input.clear();
    await input.fill(title);
    await this.page.getByRole("button", { name: /搜\s*索|Search/i }).first().click();
    await waitForRouteReady(this.page);
    return this.page
      .locator(".vxe-body--row:visible", { hasText: title })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  async deleteBook(title: string) {
    const input = this.page
      .getByLabel(/书名|Title|plugin\.linapro-book-manage\.fields\.title/i)
      .first();
    await input.clear();
    await input.fill(title);
    await this.page.getByRole("button", { name: /搜\s*索|Search/i }).first().click();
    await waitForRouteReady(this.page);
    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
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
