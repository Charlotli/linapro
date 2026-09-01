-- Mock data: approval flows and demo requests for OA approval demos.
-- 模拟数据：OA 审批演示使用的审批流程与审批单。

-- One fund-request flow with two approver nodes resolved from seeded admin users.
-- 一条两级审批人的请款流程，审批人通过稳定用户名解析种子账号。
INSERT INTO plugin_linapro_oa_approval_flow ("tenant_id", "flow_type", "flow_name", "description", "status", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    1,
    '标准请款审批',
    '请款默认两级审批流程',
    1,
    admin."id",
    admin."id",
    '2026-04-20 09:00:00',
    '2026-04-20 09:00:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow flow
    WHERE flow."tenant_id" = 0 AND flow."flow_name" = '标准请款审批'
  );

-- Flow nodes resolve approvers from stable usernames so repeated loads stay idempotent.
-- 流程节点通过稳定用户名解析审批人，重复加载保持幂等。
INSERT INTO plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order", "approver_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    1,
    leader."id",
    leader."id",
    leader."id",
    '2026-04-20 09:05:00',
    '2026-04-20 09:05:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user leader ON leader."username" = 'admin'
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '标准请款审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow_node node
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."node_order" = 1
  );

INSERT INTO plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order", "approver_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    2,
    finance."id",
    finance."id",
    finance."id",
    '2026-04-20 09:06:00',
    '2026-04-20 09:06:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user finance ON finance."username" = 'admin'
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '标准请款审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow_node node
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."node_order" = 2
  );

-- One pending demo request frozen to the flow nodes at submit time.
-- 一条演示审批单，提交时已冻结流程节点快照。
INSERT INTO plugin_linapro_oa_approval_request ("tenant_id", "flow_id", "flow_type", "title", "amount", "content", "attachments", "status", "current_node_order", "current_approver_id", "node_snapshot", "applicant_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    flow."flow_type",
    '青年营活动经费请款',
    3500.00,
    '<p>青年营活动场地与物料经费，明细见附件地址。</p>',
    '["https://example.com/docs/youth-camp-budget.xlsx"]',
    1,
    1,
    first_node."approver_id",
    jsonNodes.snapshot,
    admin."id",
    admin."id",
    admin."id",
    '2026-04-21 10:30:00',
    '2026-04-21 10:30:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user admin ON admin."username" = 'admin'
JOIN plugin_linapro_oa_approval_flow_node first_node
  ON first_node."tenant_id" = flow."tenant_id" AND first_node."flow_id" = flow."id" AND first_node."node_order" = 1
CROSS JOIN LATERAL (
    SELECT '[' || COALESCE(string_agg(
        '{"order":' || node."node_order" ||
        ',"approverId":' || node."approver_id" ||
        ',"approverName":"' || COALESCE(approver."username", '') || '"}', ',' ORDER BY node."node_order"), '') || ']'
        AS snapshot
    FROM plugin_linapro_oa_approval_flow_node node
    LEFT JOIN sys_user approver ON approver."id" = node."approver_id"
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."deleted_at" IS NULL
) jsonNodes
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '标准请款审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_request request
    WHERE request."tenant_id" = 0 AND request."title" = '青年营活动经费请款'
  );

-- Submit timeline entry for the demo request.
-- 演示审批单的发起时间线。
INSERT INTO plugin_linapro_oa_approval_record ("tenant_id", "request_id", "node_order", "action", "actor_id", "comment", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    request."tenant_id",
    request."id",
    0,
    1,
    request."applicant_id",
    '提交申请',
    request."applicant_id",
    request."applicant_id",
    request."created_at",
    request."updated_at"
FROM plugin_linapro_oa_approval_request request
WHERE request."tenant_id" = 0 AND request."title" = '青年营活动经费请款'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_record record
    WHERE record."tenant_id" = request."tenant_id" AND record."request_id" = request."id" AND record."action" = 1
  );

-- Demo expense-reimbursement flow with configurable dynamic form fields.
-- 演示用报销流程，含可配置的动态表单字段。
INSERT INTO plugin_linapro_oa_approval_flow ("tenant_id", "flow_type", "flow_name", "description", "status", "form_fields", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    3,
    '日常报销审批',
    '日常报销两级审批，含报销明细与发票地址',
    1,
    '[{"key":"expenseDate","label":"费用发生日期","type":"date","required":true,"asAmount":false,"options":[],"columns":[]},{"key":"expenseType","label":"费用类型","type":"select","required":true,"asAmount":false,"options":["差旅费","办公费","材料费","其他"],"columns":[]},{"key":"details","label":"报销明细","type":"detail","required":true,"asAmount":false,"options":[],"columns":[{"key":"date","label":"费用日期","type":"date"},{"key":"type","label":"费用类型","type":"text"},{"key":"amount","label":"金额","type":"number"},{"key":"note","label":"费用说明","type":"text"}]},{"key":"invoiceUrls","label":"发票","type":"attachment","required":false,"asAmount":false,"options":[],"columns":[]},{"key":"owner","label":"归属人","type":"text","required":false,"asAmount":false,"options":[],"columns":[]},{"key":"remark","label":"备注","type":"textarea","required":false,"asAmount":false,"options":[],"columns":[]}]',
    admin."id",
    admin."id",
    '2026-04-20 09:20:00',
    '2026-04-20 09:20:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow flow
    WHERE flow."tenant_id" = 0 AND flow."flow_name" = '日常报销审批'
  );

-- Approval nodes for the reimbursement flow (two levels).
-- 报销流程的两级审批节点。
INSERT INTO plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order", "approver_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    1,
    leader."id",
    leader."id",
    leader."id",
    '2026-04-20 09:25:00',
    '2026-04-20 09:25:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user leader ON leader."username" = 'admin'
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '日常报销审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow_node node
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."node_order" = 1
  );

INSERT INTO plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order", "approver_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    2,
    finance."id",
    finance."id",
    finance."id",
    '2026-04-20 09:26:00',
    '2026-04-20 09:26:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user finance ON finance."username" = 'admin'
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '日常报销审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_flow_node node
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."node_order" = 2
  );

-- One demo reimbursement request with filled dynamic form values.
-- 一条已填写动态表单值的演示报销单。
INSERT INTO plugin_linapro_oa_approval_request ("tenant_id", "flow_id", "flow_type", "title", "amount", "content", "attachments", "status", "current_node_order", "current_approver_id", "node_snapshot", "form_data", "form_snapshot", "applicant_id", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    flow."tenant_id",
    flow."id",
    flow."flow_type",
    '四月日常费用报销',
    860.50,
    '四月差旅与办公费用报销。',
    '',
    1,
    1,
    first_node."approver_id",
    nodeSnap.snapshot,
    formData.data,
    flow."form_fields",
    admin."id",
    admin."id",
    admin."id",
    '2026-04-22 11:00:00',
    '2026-04-22 11:00:00'
FROM plugin_linapro_oa_approval_flow flow
JOIN sys_user admin ON admin."username" = 'admin'
JOIN plugin_linapro_oa_approval_flow_node first_node
  ON first_node."tenant_id" = flow."tenant_id" AND first_node."flow_id" = flow."id" AND first_node."node_order" = 1
CROSS JOIN LATERAL (
    SELECT '[' || COALESCE(string_agg(
        '{"order":' || node."node_order" ||
        ',"approverId":' || node."approver_id" ||
        ',"approverName":"' || COALESCE(approver."username", '') || '"}', ',' ORDER BY node."node_order"), '') || ']' AS snapshot
    FROM plugin_linapro_oa_approval_flow_node node
    LEFT JOIN sys_user approver ON approver."id" = node."approver_id"
    WHERE node."tenant_id" = flow."tenant_id" AND node."flow_id" = flow."id" AND node."deleted_at" IS NULL
) nodeSnap
CROSS JOIN LATERAL (
    SELECT ('{"expenseDate":"2026-04-25","expenseType":"差旅费","details":[' ||
        '{"date":"2026-04-10","type":"高铁票","amount":430.50,"note":"出差往返车票"},' ||
        '{"date":"2026-04-15","type":"办公用品","amount":430.00,"note":"打印耗材"}],' ||
        '"invoiceUrls":["https://example.com/invoices/april.zip"],"owner":"admin","remark":"当月累计报销"}') AS data
) formData
WHERE flow."tenant_id" = 0 AND flow."flow_name" = '日常报销审批'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_request request
    WHERE request."tenant_id" = 0 AND request."title" = '四月日常费用报销'
  );

-- Submit timeline entry for the demo reimbursement request.
-- 演示报销单的发起时间线。
INSERT INTO plugin_linapro_oa_approval_record ("tenant_id", "request_id", "node_order", "action", "actor_id", "comment", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    request."tenant_id",
    request."id",
    0,
    1,
    request."applicant_id",
    '提交申请',
    request."applicant_id",
    request."applicant_id",
    request."created_at",
    request."updated_at"
FROM plugin_linapro_oa_approval_request request
WHERE request."tenant_id" = 0 AND request."title" = '四月日常费用报销'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_oa_approval_record record
    WHERE record."tenant_id" = request."tenant_id" AND record."request_id" = request."id" AND record."action" = 1
  );
