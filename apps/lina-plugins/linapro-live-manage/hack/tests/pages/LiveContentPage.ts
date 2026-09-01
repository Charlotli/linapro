import type { Page } from "@host-tests/support/playwright";

import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-live-manage live-content management page. */
export class LiveContentPage {
  constructor(private page: Page) {}

  private resolveLocalizedLabel(label: string) {
    const labelMap: Record<string, RegExp> = {
      直播主题:
        /直播主题|Title|plugin\.linapro-live-manage\.fields\.title/i,
    };
    const localizedLabel = labelMap[label];
    if (localizedLabel) {
      return this.page.getByLabel(localizedLabel).first();
    }
    return this.page.getByLabel(label, { exact: true }).first();
  }

  private get modal() {
    return this.page.locator('[role="dialog"]');
  }

  async goto() {
    await this.page.goto("/live/content");
    await waitForTableReady(this.page);
  }

  /** Create live content bound to the first available live room. */
  async createLive(title: string, roomName: string) {
    await this.page
      .getByRole("button", { name: /新\s*增/ })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    // Select the live room from the room select control.
    const roomSelect = this.modal
      .getByPlaceholder("请选择直播间")
      .first();
    await roomSelect.click();
    await this.page
      .locator(".ant-select-dropdown:visible .ant-select-item-option", {
        hasText: roomName,
      })
      .first()
      .click();

    await this.modal.getByPlaceholder("请输入直播主题").first().fill(title);

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /** Edit live content: search by title, update the title. */
  async editLive(searchTitle: string, newTitle: string) {
    await this.fillSearchField("直播主题", searchTitle);
    await this.clickSearch();

    const row = this.page.locator(".vxe-body--row:visible", {
      hasText: searchTitle,
    });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /编\s*辑/ })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    const titleInput = this.modal.getByPlaceholder("请输入直播主题").first();
    await titleInput.clear();
    await titleInput.fill(newTitle);

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /** Delete live content: search by title, click delete, confirm. */
  async deleteLive(title: string) {
    await this.fillSearchField("直播主题", title);
    await this.clickSearch();

    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
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

  /** Start live content: search by title, click start, confirm. */
  async startLive(title: string) {
    await this.fillSearchField("直播主题", title);
    await this.clickSearch();

    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /开\s*启直播|Start Live/i })
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

  /** Stop live content: search by title, click stop, confirm. */
  async stopLive(title: string) {
    await this.fillSearchField("直播主题", title);
    await this.clickSearch();

    const row = this.page.locator(".vxe-body--row:visible", { hasText: title });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    await row
      .locator("button:visible")
      .filter({ hasText: /关\s*闭直播|Stop Live/i })
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

  async deleteLiveIfExists(title: string) {
    if (await this.hasLive(title)) {
      await this.deleteLive(title);
    }
  }

  /** Check whether live content with the given title is visible. */
  async hasLive(title: string): Promise<boolean> {
    await this.fillSearchField("直播主题", title);
    await this.clickSearch();
    return this.page
      .locator(".vxe-body--row:visible", { hasText: title })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Fill search form field by label */
  async fillSearchField(label: string, value: string) {
    const input = this.resolveLocalizedLabel(label);
    await input.clear();
    await input.fill(value);
  }

  /** Click search button */
  async clickSearch() {
    await this.page
      .getByRole("button", { name: /搜\s*索|Search/i })
      .first()
      .click();
    await waitForRouteReady(this.page);
  }
}
