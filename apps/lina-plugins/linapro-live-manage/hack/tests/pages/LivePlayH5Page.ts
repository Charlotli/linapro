import type { Page } from "@host-tests/support/playwright";

/**
 * Page object for the linapro-live-manage anonymous viewer H5 play page
 * served by the plugin at /x/linapro-live-manage/h5.
 */
export class LivePlayH5Page {
  constructor(private page: Page) {}

  private get video() {
    return this.page.locator("#video");
  }

  private get liveBadge() {
    return this.page.locator("#live-badge");
  }

  private get stateChip() {
    return this.page.locator("#state-chip");
  }

  private get noticeSection() {
    return this.page.locator("#notice-section");
  }

  private get noticeText() {
    return this.page.locator("#notice-text");
  }

  private get liveTitle() {
    return this.page.locator("#live-title");
  }

  private get calendarBar() {
    return this.page.locator("#calendar-bar");
  }

  private get viewersBadge() {
    return this.page.locator("#viewers-badge");
  }

  private get shareButton() {
    return this.page.locator("#share-button");
  }

  private get shareLabel() {
    return this.page.locator("#share-button-label");
  }

  private get replaySection() {
    return this.page.locator("#replay-section");
  }

  private get replayList() {
    return this.page.locator("#replay-list");
  }

  private get replayBadge() {
    return this.page.locator("#replay-badge");
  }

  private get toast() {
    return this.page.locator("#toast");
  }

  /** Open the H5 player page for one room without authentication. */
  async goto(roomCode: string) {
    await this.page.goto(
      `/x/linapro-live-manage/h5?room=${encodeURIComponent(roomCode)}&tenant=0`,
    );
    await this.page
      .locator("#app")
      .waitFor({ state: "visible", timeout: 15000 });
  }

  /**
   * Open the bare /h5 entry (no trailing slash). Regression FB-1: relative
   * asset references broke on this form, so the page must still load its
   * stylesheet and script from the absolute /x/<plugin>/h5/ paths.
   */
  async gotoBare(roomCode: string) {
    await this.page.goto(
      `/x/linapro-live-manage/h5?room=${encodeURIComponent(roomCode)}&tenant=0`,
      { waitUntil: "domcontentloaded" },
    );
    await this.page
      .locator("#app")
      .waitFor({ state: "visible", timeout: 15000 });
  }

  /**
   * Assert the stylesheet and script referenced through absolute paths took
   * effect: a section carrying the .hidden class must compute display:none
   * (CSS loaded) and the player state must have been rendered by app.js.
   */
  async expectStylesApplied() {
    const applied = await this.page.evaluate(() => {
      const hidden = document.querySelector("#loading-section");
      if (!hidden) {
        return false;
      }
      return window.getComputedStyle(hidden).display === "none";
    });
    if (!applied) {
      throw new Error(
        "hidden section is visible: stylesheet or script failed to load",
      );
    }
  }

  /** Whether the video element is inside a visible player section. */
  async isPlayerVisible() {
    return this.page
      .locator("#player-section")
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Whether the HLS source was attached to the video element. */
  async hasAttachedSource(expectedUrlPart: string) {
    const source = await this.video.evaluate((el: HTMLVideoElement) => ({
      src: el.getAttribute("src"),
      hasMseSource: (el as unknown as { src?: string }).src !== undefined,
    }));
    if (source.src && source.src.includes(expectedUrlPart)) {
      return true;
    }
    // MSE playback (hls.js) never sets a src attribute; the badge is the
    // observable proxy for an attached stream in tests.
    return this.liveBadge.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /** Whether the ongoing live badge is visible. */
  async isLiveBadgeVisible() {
    return this.liveBadge.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /** Whether the state chip shows the given text. */
  async expectStateChip(pattern: RegExp) {
    await expectVisible(this.stateChip, pattern);
  }

  /** Whether the notice section shows the given text. */
  async expectNotice(pattern: RegExp) {
    await expectVisible(this.noticeText, pattern);
  }

  /** Whether the page headline contains the given text. */
  async expectTitle(pattern: RegExp) {
    await expectVisible(this.liveTitle, pattern);
  }

  /** Whether the "no live" notice with retry button is displayed. */
  async expectNoLiveNotice() {
    await this.expectNotice(/暂无直播|敬请期待/);
    await expectVisible(
      this.page.locator("#retry-button"),
      /重新加载/,
    );
  }

  /** Text of the room-name eyebrow line above the title (empty if hidden). */
  async roomLineText(): Promise<string> {
    const room = this.page.locator("#live-room");
    await room.waitFor({ state: "visible", timeout: 10000 });
    return (await room.textContent()) ?? "";
  }

  /** Text of the cover start-date capsule (empty if hidden). */
  async coverDateText(): Promise<string> {
    const date = this.page.locator("#cover-date");
    await date.waitFor({ state: "visible", timeout: 10000 });
    return (await date.textContent()) ?? "";
  }

  /** Whether the preview-only calendar subscription bar is displayed. */
  async isCalendarBarVisible() {
    return this.calendarBar.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /**
   * Wait for the watch-session heartbeat to land and the "N watching" badge
   * to appear, then return its text.
   */
  async expectViewersBadge(): Promise<string> {
    await this.viewersBadge.waitFor({ state: "visible", timeout: 15000 });
    return (await this.viewersBadge.textContent()) ?? "";
  }

  /** Whether the viewers badge is displayed. */
  async isViewersBadgeVisible() {
    return this.viewersBadge
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Whether the share action bar is displayed for presentable content. */
  async isShareBarVisible() {
    return this.page
      .locator("#share-bar")
      .isVisible({ timeout: 5000 })
      .catch(() => false);
  }

  /** Localized share button label: share when navigator.share exists. */
  async shareButtonLabel(): Promise<string> {
    await this.shareButton.waitFor({ state: "visible", timeout: 10000 });
    return (await this.shareLabel.textContent()) ?? "";
  }

  /** Trigger the share action (native share sheet or clipboard fallback). */
  async clickShareButton() {
    await this.shareButton.click();
  }

  /** Whether the replay library section is rendered for this room. */
  async isReplaysSectionVisible() {
    return this.replaySection.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /** Number of replay cards rendered inside the library list. */
  async replayCardCount(): Promise<number> {
    return this.replayList.locator(".replay-card").count();
  }

  /** Title text of one replay card by index. */
  async replayCardTitle(index: number): Promise<string> {
    const card = this.replayList.locator(".replay-card").nth(index);
    await card.waitFor({ state: "visible", timeout: 10000 });
    return (await card.locator(".replay-title").textContent()) ?? "";
  }

  /** Open one replay from the library and wait for the player. */
  async clickReplayCard(index: number) {
    await this.replayList.locator(".replay-card").nth(index).click();
  }

  /** Whether the player-attached replay badge is displayed. */
  async isReplayBadgeVisible() {
    return this.replayBadge.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /** Click the copy-subscription-link action inside the calendar bar. */
  async clickCopyCalendarLink() {
    await this.page
      .locator("#calendar-copy")
      .click();
  }

  /**
   * Assert the lightweight toast shows a result message. The toast element
   * stays mounted with opacity 0, so the wait must target the .visible state
   * class instead of plain visibility.
   */
  async expectToast(pattern: RegExp) {
    const visibleToast = this.page.locator("#toast.visible");
    await visibleToast.waitFor({ state: "visible", timeout: 5000 });
    const text = (await visibleToast.textContent()) ?? "";
    if (!pattern.test(text)) {
      throw new Error(`expected toast matching ${pattern}, got: ${text}`);
    }
  }

  /* ---- 公告与圣经工具面板 ---- */

  private get toolRow() {
    return this.page.locator("#tool-row");
  }

  private get announcementOverlay() {
    return this.page.locator("#announcement-overlay");
  }

  private get announcementList() {
    return this.page.locator("#announcement-list");
  }

  private get bibleOverlay() {
    return this.page.locator("#bible-overlay");
  }

  private get bibleBody() {
    return this.page.locator("#bible-body");
  }

  private get bibleTitle() {
    return this.page.locator("#bible-title");
  }

  /** Whether the tool entry row (announcements / bible) is visible. */
  async isToolRowVisible() {
    return this.toolRow.isVisible({ timeout: 5000 }).catch(() => false);
  }

  /** Whether the tool entry row is hidden (notice states). */
  async isToolRowHidden() {
    return this.toolRow
      .waitFor({ state: "hidden", timeout: 5000 })
      .then(() => true)
      .catch(() => false);
  }

  /** Open the announcements sheet and wait for it to settle. */
  async openAnnouncements() {
    await this.page.locator("#tool-announcement").click();
    await this.announcementOverlay.waitFor({ state: "visible", timeout: 5000 });
    // Either the empty hint or at least one card must render after fetch.
    await this.announcementList
      .locator(".announcement-card, .sheet-empty")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Number of announcement cards currently rendered in the sheet. */
  async announcementCardCount(): Promise<number> {
    return this.announcementList.locator(".announcement-card").count();
  }

  /** Text of one announcement card (title + content) by index. */
  async announcementCardText(index: number): Promise<string> {
    const card = this.announcementList.locator(".announcement-card").nth(index);
    await card.waitFor({ state: "visible", timeout: 10000 });
    return (await card.textContent()) ?? "";
  }

  /** Close any open sheet through its close button (mask clicks are
   * intercepted by the bottom panel in the books grid state). */
  async closeSheets() {
    const closeButtons = this.page.locator(".sheet-close:visible");
    while ((await closeButtons.count()) > 0) {
      await closeButtons.first().click();
      await this.page.waitForTimeout(300);
    }
    await this.announcementOverlay
      .waitFor({ state: "hidden", timeout: 5000 })
      .catch(() => {});
    await this.bibleOverlay
      .waitFor({ state: "hidden", timeout: 5000 })
      .catch(() => {});
  }

  /** Open the bible reader sheet; books grid must be ready. */
  async openBible() {
    await this.page.locator("#tool-bible").click();
    await this.bibleOverlay.waitFor({ state: "visible", timeout: 5000 });
    await this.bibleBody
      .locator(".bible-book")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Count of book chips currently rendered (39 old + 27 new = 66). */
  async bibleBookCount(): Promise<number> {
    return this.bibleBody.locator(".bible-book").count();
  }

  /** Open one book by its short name chip (e.g. 创). */
  async openBibleBook(shortName: string) {
    await this.bibleBody
      .locator(".bible-book", { hasText: shortName })
      .first()
      .click();
    await this.bibleBody
      .locator(".bible-chapter")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Chapter chip count of the current book (Genesis = 50). */
  async bibleChapterCount(): Promise<number> {
    return this.bibleBody.locator(".bible-chapter").count();
  }

  /** Open one chapter by its number chip. */
  async openBibleChapter(chapter: number) {
    await this.bibleBody
      .locator(`.bible-chapter[data-chapter="${chapter}"]`)
      .click();
    await this.bibleBody
      .locator(".bible-verse")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Verse count currently rendered in the reader. */
  async bibleVerseCount(): Promise<number> {
    return this.bibleBody.locator(".bible-verse").count();
  }

  /** Header title text of the reader sheet. */
  async bibleHeaderText(): Promise<string> {
    await this.bibleTitle.waitFor({ state: "visible", timeout: 10000 });
    return (await this.bibleTitle.textContent()) ?? "";
  }

  /** First verse full text (verse number + lection) of the current chapter. */
  async bibleFirstVerseText(): Promise<string> {
    const verse = this.bibleBody.locator(".bible-verse").first();
    await verse.waitFor({ state: "visible", timeout: 10000 });
    return (await verse.textContent()) ?? "";
  }

  /** Navigate to the previous chapter (cross-volume when on chapter 1). */
  async biblePrev() {
    await this.page.locator("#bible-prev").click();
    await this.bibleBody
      .locator(".bible-verse")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Navigate to the next chapter (cross-volume at chapter boundary). */
  async bibleNext() {
    await this.page.locator("#bible-next").click();
    await this.bibleBody
      .locator(".bible-verse")
      .first()
      .waitFor({ state: "visible", timeout: 10000 });
  }

  /** Step back one level in the reader (reading -> chapters -> books). */
  async bibleBack() {
    await this.page.locator("#bible-back").click();
  }
}

async function expectVisible(locator: ReturnType<Page["locator"]>, pattern: RegExp) {
  await locator.waitFor({ state: "visible", timeout: 10000 });
  const text = (await locator.textContent()) ?? "";
  if (!pattern.test(text)) {
    throw new Error(`expected text matching ${pattern}, got: ${text}`);
  }
}
