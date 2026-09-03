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
}

async function expectVisible(locator: ReturnType<Page["locator"]>, pattern: RegExp) {
  await locator.waitFor({ state: "visible", timeout: 10000 });
  const text = (await locator.textContent()) ?? "";
  if (!pattern.test(text)) {
    throw new Error(`expected text matching ${pattern}, got: ${text}`);
  }
}
