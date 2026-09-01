import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { OaFlowPage } from "../../pages/OaFlowPage";

test.describe("TC001 OA 流程配置", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-oa-approval"]);
  });

  const flowName = `测试流程_${Date.now()}`;

  test("TC001a: 流程配置页展示翻译后的表格标题", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();

    await expect(
      adminPage.getByText(/审批流程列表|Approval flow list/i),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      adminPage.getByText(/plugin\.linapro-oa-approval\.flowTableTitle/),
    ).toHaveCount(0);
  });

  test("TC001b: 创建带审批人的流程", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();
    await flowPage.createFlow(flowName);

    await expect(
      adminPage.getByText(/新增成功|创建成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC001c: 列表中可见新创建的流程", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();

    expect(await flowPage.hasFlow(flowName)).toBeTruthy();
  });

  test("TC001d: 删除流程", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();
    await flowPage.deleteFlow(flowName);

    await expect(
      adminPage.getByText(/删除成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });
});
