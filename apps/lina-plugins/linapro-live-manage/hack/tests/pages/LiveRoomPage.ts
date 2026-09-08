import type { Page } from "@host-tests/support/playwright";

import { workspacePath } from "@host-tests/fixtures/config";
import {
  waitForConfirmOverlay,
  waitForDialogReady,
  waitForRouteReady,
  waitForTableReady,
} from "@host-tests/support/ui";

/** Page object for the linapro-live-manage live-room management page. */
export class LiveRoomPage {
  constructor(private page: Page) {}

  private resolveLocalizedLabel(label: string) {
    const labelMap: Record<string, RegExp> = {
      直播间名称:
        /直播间名称|Room Name|plugin\.linapro-live-manage\.fields\.roomName/i,
      直播间编码:
        /直播间编码|Room Code|plugin\.linapro-live-manage\.fields\.roomCode/i,
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
    // The hosted admin SPA uses hash routing: the workspace path stays
    // constant and the SPA route lives in the URL fragment.
    await this.page.goto(`${workspacePath("/live/room")}#/live/room`);
    await waitForTableReady(this.page);
  }

  /** Create a new live room with a gathering type. */
  async createRoom(roomCode: string, roomName: string) {
    await this.page
      .getByRole("button", { name: /新\s*增/ })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    await this.modal
      .getByPlaceholder("请输入直播间编码")
      .first()
      .fill(roomCode);
    await this.modal
      .getByPlaceholder("请输入直播间名称")
      .first()
      .fill(roomName);

    // Select room type from the first ant-select in the modal.
    const typeSelect = this.modal.locator(".ant-select").first();
    await typeSelect.click();
    await this.page
      .locator(".ant-select-dropdown:visible .ant-select-item-option", {
        hasText: /正常聚会|Gathering/i,
      })
      .first()
      .click();

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /**
   * Locate all visible table rows sharing one vxe rowid. vxe-table renders
   * the action column in a separate fixed-right table whose rows repeat the
   * source row's `rowid`, so action lookups must go through that id.
   */
  private rowById(rowId: string) {
    return this.page.locator(`tr.vxe-body--row[rowid="${rowId}"]`);
  }

  private async findRowIdByText(text: string) {
    const row = this.page.locator(".vxe-body--row:visible", { hasText: text });
    await row.first().waitFor({ state: "visible", timeout: 10000 });
    const rowId = await row.first().getAttribute("rowid");
    if (!rowId) {
      throw new Error(`row for ${text} has no vxe rowid`);
    }
    return rowId;
  }

  /** Edit a live room: search by name, update the name. */
  async editRoom(searchName: string, newName: string) {
    await this.fillSearchField("直播间名称", searchName);
    await this.clickSearch();

    const rowId = await this.findRowIdByText(searchName);
    await this.rowById(rowId)
      .locator("button:visible")
      .filter({ hasText: /编\s*辑/ })
      .first()
      .click();

    await waitForDialogReady(this.modal);

    const nameInput = this.modal
      .getByPlaceholder("请输入直播间名称")
      .first();
    await nameInput.clear();
    await nameInput.fill(newName);

    await this.modal.getByRole("button", { name: /确\s*认/ }).click();

    await waitForRouteReady(this.page);
    await this.modal
      .waitFor({ state: "hidden", timeout: 10000 })
      .catch(() => {});
  }

  /** Delete a live room: search by name, click delete, confirm. */
  async deleteRoom(roomName: string) {
    await this.fillSearchField("直播间名称", roomName);
    await this.clickSearch();

    const rowId = await this.findRowIdByText(roomName);
    await this.rowById(rowId)
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

  async deleteRoomIfExists(roomName: string) {
    if (await this.hasRoom(roomName)) {
      await this.deleteRoom(roomName);
    }
  }

  /** Check whether a live room with the given name is visible. */
  async hasRoom(roomName: string): Promise<boolean> {
    await this.fillSearchField("直播间名称", roomName);
    await this.clickSearch();
    return this.page
      .locator(".vxe-body--row:visible", { hasText: roomName })
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

  /** Open the calendar-subscription modal of the room found by name. */
  async openCalendarModal(roomName: string) {
    await this.fillSearchField("直播间名称", roomName);
    await this.clickSearch();

    const rowId = await this.findRowIdByText(roomName);
    await this.rowById(rowId)
      .locator("button:visible")
      .filter({ hasText: /日历订阅|Calendar/i })
      .first()
      .click();

    await waitForDialogReady(this.modal);
  }

  /** Full subscription link text rendered inside the calendar modal. */
  async calendarModalSubscribeLink(): Promise<string> {
    const textarea = this.modal.locator("textarea").first();
    await textarea.waitFor({ state: "visible", timeout: 10000 });
    return textarea.inputValue();
  }

  /** Whether the calendar modal renders the localized copy action. */
  async calendarModalHasCopyButton(): Promise<boolean> {
    return this.modal
      .getByRole("button", { name: /复制链接|Copy link/i })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Localized subscription how-to hint inside the calendar modal. */
  async calendarModalHintText(): Promise<string> {
    // The vben modal renders a hidden aria-description paragraph first, so
    // the lookup must target the visible hint paragraph only.
    const hint = this.modal.locator("p:visible").first();
    await hint.waitFor({ state: "visible", timeout: 10000 });
    return (await hint.textContent()) ?? "";
  }

  /** Open the QR-code watch-link modal of the room found by name. */
  async openQrcodeModal(roomName: string) {
    await this.fillSearchField("直播间名称", roomName);
    await this.clickSearch();

    const rowId = await this.findRowIdByText(roomName);
    await this.rowById(rowId)
      .locator("button:visible")
      .filter({ hasText: /二维码|QR Code/i })
      .first()
      .click();

    await waitForDialogReady(this.modal);
  }

  /** Whether the QR modal renders a generated QR-code image. */
  async qrcodeModalHasImage(): Promise<boolean> {
    return this.modal
      .locator("img[src^='data:image']")
      .first()
      .isVisible({ timeout: 10000 })
      .catch(() => false);
  }

  /** Full H5 watch link rendered inside the QR modal. */
  async qrcodeModalWatchLink(): Promise<string> {
    const textarea = this.modal.locator("textarea").first();
    await textarea.waitFor({ state: "visible", timeout: 10000 });
    return textarea.inputValue();
  }

  /** Whether the QR modal renders the localized copy action. */
  async qrcodeModalHasCopyButton(): Promise<boolean> {
    return this.modal
      .getByRole("button", { name: /复制链接|Copy link/i })
      .first()
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Open the announcement management modal of the room found by name. */
  async openAnnouncementModal(roomName: string) {
    await this.fillSearchField("直播间名称", roomName);
    await this.clickSearch();

    const rowId = await this.findRowIdByText(roomName);
    await this.rowById(rowId)
      .locator("button:visible")
      // Ghost buttons render two-character CJK labels with an inserted
      // space ("公 告"), so the filter must tolerate optional whitespace.
      .filter({ hasText: /公\s*告|Announcement/i })
      .first()
      .click();

    await waitForDialogReady(this.modal);
  }

  /** Create one announcement inside the open modal. */
  async createAnnouncement(title: string, content: string) {
    await this.modal
      .getByRole("button", { name: /新增公告|New Announcement/i })
      .first()
      .click();
    await this.modal
      .getByPlaceholder("请输入公告标题")
      .fill(title);
    await this.modal
      .getByPlaceholder("请输入公告内容，支持换行")
      .fill(content);
    await this.modal
      .getByRole("button", { name: /保\s*存|Save/i })
      .first()
      .click();
    await this.modal
      .locator(".announcement-row", { hasText: title })
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Announcement card count in the open modal. */
  async announcementRowCount(): Promise<number> {
    return this.modal.locator(".announcement-row").count();
  }

  /** Whether one announcement card (found by title) shows the empty switch. */
  async announcementSwitchChecked(title: string): Promise<boolean> {
    const row = this.modal
      .locator(".announcement-row", { hasText: title })
      .first();
    await row.waitFor({ state: "visible", timeout: 10000 });
    const button = row.locator(".ant-switch");
    await button.waitFor({ state: "visible", timeout: 5000 });
    return (await button.getAttribute("aria-checked")) === "true";
  }

  /** Toggle the enable switch of one announcement card. */
  async toggleAnnouncementSwitch(title: string) {
    const row = this.modal
      .locator(".announcement-row", { hasText: title })
      .first();
    await row.locator(".ant-switch").click();
    await this.page.waitForTimeout(500);
  }

  /** Delete one announcement card by title and confirm the popconfirm. */
  async deleteAnnouncement(title: string) {
    const row = this.modal
      .locator(".announcement-row", { hasText: title })
      .first();
    await row
      .locator("button")
      .filter({ hasText: /删\s*除|Delete/i })
      .first()
      .click();
    const popconfirm = await waitForConfirmOverlay(this.page);
    const confirmBtn = popconfirm.getByRole("button", {
      name: /确\s*定|OK|是/i,
    });
    if (await confirmBtn.isVisible({ timeout: 2000 }).catch(() => false)) {
      await confirmBtn.click();
    }
    await row.waitFor({ state: "hidden", timeout: 10000 }).catch(() => {});
  }
}
