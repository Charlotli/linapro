import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { execPgSQL } from "@host-tests/support/postgres";
import { LiveContentPage } from "../../pages/LiveContentPage";
import { LivePlayH5Page } from "../../pages/LivePlayH5Page";
import { LiveRoomPage } from "../../pages/LiveRoomPage";

// Deterministic stamp keeps lookups stable across worker restarts; the
// beforeAll hard-delete makes reruns idempotent (see TC004 for the pattern).
// Each test receives a fresh adminPage, so every goto is a real navigation
// (hash-only repeats inside one test do not rerender the SPA shell).
const stamp = "450000000000";
const roomCode = `ROOM-ANNOUNCE-${stamp}`;
const roomName = `公告测试直播间_${stamp}`;
const liveTitle = `公告测试直播_${stamp}`;
const liveUrl = `https://example.com/hls/tc008-${stamp}.m3u8`;
const announcementTitle = `主日聚会通知_${stamp}`;
const announcementContent = "本周主日聚会九点开始\n请提前十分钟入场";

const announcementsPath =
  `/x/linapro-live-manage/api/v1/announcements` +
  `?roomCode=${encodeURIComponent(roomCode)}&tenantId=0`;
const bibleBooksPath = `/x/linapro-live-manage/api/v1/bible/books`;
const bibleChapterPath = (volumeSn: number, chapter: number) =>
  `/x/linapro-live-manage/api/v1/bible/chapter` +
  `?volumeSn=${volumeSn}&chapter=${chapter}`;

interface ApiPayload {
  code: number;
  errorCode?: string;
  data?: { list?: Array<Record<string, unknown>>; total?: number } & Record<
    string,
    unknown
  >;
}

test.describe("TC008 直播间公告与圣经阅读器", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-live-manage"]);
    // A failed previous run may have left the shared records behind; hard
    // delete first so the deterministic setup in TC008a can rerun cleanly.
    execPgSQL(
      `DELETE FROM plugin_linapro_live_manage_announcement WHERE title = '${announcementTitle}';\n` +
        `DELETE FROM plugin_linapro_live_manage_live WHERE title = '${liveTitle}';\n` +
        `DELETE FROM plugin_linapro_live_manage_room WHERE room_code = '${roomCode}';`,
    );
  });

  test("TC008a: 管理端公告弹窗新增公告且观众端可见", async ({ adminPage }) => {
    // One dedicated room hosts the announcement fixtures.
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.createRoom(roomCode, roomName);

    await roomPage.openAnnouncementModal(roomName);
    await roomPage.createAnnouncement(announcementTitle, announcementContent);
    expect(await roomPage.announcementRowCount()).toBe(1);
    expect(
      await roomPage.announcementSwitchChecked(announcementTitle),
    ).toBe(true);

    // The enabled announcement is visible on the anonymous viewer endpoint.
    const created = await adminPage.request.get(announcementsPath);
    expect(created.status()).toBe(200);
    const createdPayload = (await created.json()) as ApiPayload;
    expect(createdPayload.code).toBe(0);
    expect(createdPayload.data?.list?.length ?? 0).toBe(1);
    expect(String(createdPayload.data?.list?.[0]?.["title"] ?? "")).toBe(
      announcementTitle,
    );
  });

  test("TC008b: 开启直播后 H5 进入播放态", async ({ adminPage }) => {
    // An ongoing live makes the H5 page presentable: the tool row (announcements
    // / bible) only renders in live or replay states, never on notices.
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.createLive(liveTitle, roomName, liveUrl);
    await livePage.startLive(liveTitle);

    const play = await adminPage.request.get(
      `/x/linapro-live-manage/api/v1/play?roomCode=${encodeURIComponent(roomCode)}&tenantId=0`,
    );
    expect(play.status()).toBe(200);
    const playPayload = (await play.json()) as ApiPayload;
    expect(playPayload.code).toBe(0);
  });

  test("TC008c: H5 公告面板展示启用公告并可关闭", async ({ adminPage }) => {
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);
    expect(await h5.isToolRowVisible()).toBe(true);

    await h5.openAnnouncements();
    expect(await h5.announcementCardCount()).toBe(1);
    const cardText = await h5.announcementCardText(0);
    expect(cardText).toContain(announcementTitle);
    expect(cardText).toContain("请提前十分钟入场");

    await h5.closeSheets();
  });

  test("TC008d: 管理端启停删公告并同步观众端", async ({ adminPage }) => {
    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.openAnnouncementModal(roomName);

    // Disabling hides it from viewers immediately (enabled filter).
    await roomPage.toggleAnnouncementSwitch(announcementTitle);
    const disabled = await adminPage.request.get(announcementsPath);
    const disabledPayload = (await disabled.json()) as ApiPayload;
    expect(disabledPayload.data?.list?.length ?? 0).toBe(0);

    // Deleting removes the record entirely.
    await roomPage.deleteAnnouncement(announcementTitle);
    expect(await roomPage.announcementRowCount()).toBe(0);

    // Recreate for the H5 panel test and leave it enabled.
    await roomPage.createAnnouncement(announcementTitle, announcementContent);
    expect(await roomPage.announcementRowCount()).toBe(1);
  });

  test("TC008e: 圣经接口 66 卷且越界返回空", async ({ adminPage }) => {
    // Anonymous catalogue endpoints: 66 books, Genesis 1 = 31 verses.
    const books = await adminPage.request.get(bibleBooksPath);
    expect(books.status()).toBe(200);
    const booksPayload = (await books.json()) as ApiPayload;
    expect(booksPayload.code).toBe(0);
    expect(booksPayload.data?.list?.length ?? 0).toBe(66);

    const chapter = await adminPage.request.get(bibleChapterPath(1, 1));
    const chapterPayload = (await chapter.json()) as ApiPayload;
    expect(chapterPayload.code).toBe(0);
    expect(chapterPayload.data?.list?.length ?? 0).toBe(31);

    // Out-of-range volume and chapter answer empty payloads, never errors.
    const beyond = await adminPage.request.get(bibleChapterPath(67, 1));
    const beyondPayload = (await beyond.json()) as ApiPayload;
    expect(beyondPayload.code).toBe(0);
    expect(beyondPayload.data?.list?.length ?? 0).toBe(0);
  });

  test("TC008f: H5 圣经阅读器三级导航与跨章翻页可用", async ({ adminPage }) => {
    const h5 = new LivePlayH5Page(adminPage);
    await h5.goto(roomCode);
    expect(await h5.isToolRowVisible()).toBe(true);

    await h5.openBible();
    expect(await h5.bibleBookCount()).toBe(66);

    await h5.openBibleBook("创");
    expect(await h5.bibleChapterCount()).toBe(50);

    await h5.openBibleChapter(1);
    expect(await h5.bibleVerseCount()).toBe(31);
    expect(await h5.bibleHeaderText()).toContain("创世记");
    expect(await h5.bibleFirstVerseText()).toContain("起初");

    // Next from Genesis 1 goes to Genesis 2 (within the same volume).
    await h5.bibleNext();
    expect(await h5.bibleHeaderText()).toContain("第2章");

    // Back to chapters, then step back to the book grid.
    await h5.bibleBack();
    await h5.bibleBack();
    expect(await h5.bibleBookCount()).toBe(66);

    await h5.closeSheets();
  });

  test("TC008g: 清理测试直播与直播间", async ({ adminPage }) => {
    const livePage = new LiveContentPage(adminPage);
    await livePage.goto();
    await livePage.deleteLiveIfExists(liveTitle);

    const roomPage = new LiveRoomPage(adminPage);
    await roomPage.goto();
    await roomPage.deleteRoomIfExists(roomName);
  });
});
