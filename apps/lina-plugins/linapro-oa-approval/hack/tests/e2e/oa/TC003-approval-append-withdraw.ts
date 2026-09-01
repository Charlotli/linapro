import { expect, test } from "@host-tests/fixtures/auth";
import { prepareSourcePluginsBaseline } from "@host-tests/fixtures/plugin";
import { OaApprovalPage } from "../../pages/OaApprovalPage";
import { OaFlowPage } from "../../pages/OaFlowPage";

test.describe("TC003 OA 审批加签与撤回", () => {
  test.beforeAll(async () => {
    await prepareSourcePluginsBaseline(["linapro-oa-approval"]);
  });

  const flowName = `加签流程_${Date.now()}`;
  const title = `测试核销_${Date.now()}`;

  test("TC003a: 提交后加签审批人", async ({ adminPage }) => {
    const flowPage = new OaFlowPage(adminPage);
    await flowPage.goto();
    await flowPage.createFlow(flowName);

    const approvalPage = new OaApprovalPage(adminPage);
    await approvalPage.goto();
    await approvalPage.submitRequest(flowName, title, "300.00");
    await approvalPage.openDetail(title);
    await approvalPage.appendApproverInDrawer("admin");

    await expect(
      adminPage.getByText(/审批人添加成功|success/i),
    ).toBeVisible({ timeout: 5000 });
  });

  test("TC003b: 申请人撤回审批单", async ({ adminPage }) => {
    const approvalPage = new OaApprovalPage(adminPage);
    await approvalPage.goto();
    await approvalPage.withdraw(title);

    await expect(
      adminPage.getByText(/审批撤回成功|Withdrawn|success/i),
    ).toBeVisible({ timeout: 5000 });
  });
});
