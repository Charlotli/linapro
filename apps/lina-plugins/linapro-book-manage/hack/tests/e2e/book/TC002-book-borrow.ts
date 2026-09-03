import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { BookPage } from "../../pages/BookPage";

test.describe("TC002 书籍登记与借阅", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-book-manage"]);
  });

  const title = `测试书籍_${Date.now()}`;

  test("TC001a: 书籍页展示翻译后的表格标题", async ({ adminPage }) => {
    const bookPage = new BookPage(adminPage);
    await bookPage.goto();
    await expect(
      adminPage.getByText(/书籍台账|Book registry/i),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      adminPage.getByText(/plugin\.linapro-book-manage\.tableTitle/),
    ).toHaveCount(0);
  });

  test("TC001b: 登记书籍", async ({ adminPage }) => {
    const bookPage = new BookPage(adminPage);
    await bookPage.goto();
    await bookPage.createBook(title);
    await expect(
      adminPage.getByText(/新增成功|创建成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC001c: 列表中可见新书籍", async ({ adminPage }) => {
    const bookPage = new BookPage(adminPage);
    await bookPage.goto();
    expect(await bookPage.hasBook(title)).toBeTruthy();
  });

  test("TC001d: 删除书籍", async ({ adminPage }) => {
    const bookPage = new BookPage(adminPage);
    await bookPage.goto();
    await bookPage.deleteBook(title);
    await expect(
      adminPage.getByText(/删除成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });
});
