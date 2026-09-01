import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { OaApprovalPage } from "../../pages/OaApprovalPage";
import { OaFlowPage } from "../../pages/OaFlowPage";

test.describe("TC002 OA 审批全流程", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-oa-approval"]);
  });

  const flowName = `审批流程_${Date.now()}`;
  const title = `测试请款_${Date.now()}`;

  test("TC002a: 提交审批单", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();
    await flowPage.createFlow(flowName);

    const approvalPage = new OaApprovalPage(adminPage);
    await approvalPage.goto();
    await approvalPage.submitRequest(flowName, title, "1200.00");

    await expect(
      adminPage.getByText(/审批提交成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC002b: 详情展示时间线并通过审批", async ({ adminPage }) => {
    const approvalPage = new OaApprovalPage(adminPage);
    await approvalPage.goto();
    await approvalPage.openDetail(title);
    await approvalPage.approveInDrawer("同意，按预算执行");

    await expect(
      adminPage.getByText(/审批通过成功|Approved/i),
    ).toBeVisible({ timeout: 5000 });
    await approvalPage.expectTimelineAction(/通过|Approve/i);
  });
});
