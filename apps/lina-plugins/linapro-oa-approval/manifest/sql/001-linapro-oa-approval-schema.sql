-- 001: linapro-oa-approval schema
-- 001：linapro-oa-approval 数据结构

-- Purpose: Stores tenant-scoped approval flows, flow nodes, approval requests, and approval record timelines.
-- 用途：存储租户级审批流程、流程节点、审批单与审批记录时间线。
CREATE TABLE IF NOT EXISTS plugin_linapro_oa_approval_flow (
    "id"          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"   INT          NOT NULL DEFAULT 0,
    "flow_type"   SMALLINT     NOT NULL DEFAULT 1,
    "flow_name"   VARCHAR(128) NOT NULL,
    "description" VARCHAR(512) NOT NULL DEFAULT '',
    "status"      SMALLINT     NOT NULL DEFAULT 1,
    "created_by"  BIGINT       NOT NULL DEFAULT 0,
    "updated_by"  BIGINT       NOT NULL DEFAULT 0,
    "created_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"  TIMESTAMPTZ  NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_oa_approval_flow IS 'OA approval flow configuration';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."id" IS 'Approval flow ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."flow_type" IS 'Flow type: 1=fund request, 2=write-off, 3=expense reimbursement';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."flow_name" IS 'Flow name, unique per tenant';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."description" IS 'Flow description';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."status" IS 'Flow status: 1=enabled, 0=disabled';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_oa_approval_flow_tenant_name
    ON plugin_linapro_oa_approval_flow ("tenant_id", "flow_name");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_flow_tenant_type ON plugin_linapro_oa_approval_flow ("tenant_id", "flow_type");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_flow_tenant_status ON plugin_linapro_oa_approval_flow ("tenant_id", "status");

CREATE TABLE IF NOT EXISTS plugin_linapro_oa_approval_flow_node (
    "id"          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"   INT          NOT NULL DEFAULT 0,
    "flow_id"     BIGINT       NOT NULL,
    "node_order"  INT          NOT NULL DEFAULT 1,
    "approver_id" BIGINT       NOT NULL,
    "created_by"  BIGINT       NOT NULL DEFAULT 0,
    "updated_by"  BIGINT       NOT NULL DEFAULT 0,
    "created_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"  TIMESTAMPTZ  NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_oa_approval_flow_node IS 'OA approval flow ordered approver nodes';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."id" IS 'Flow node ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."flow_id" IS 'Owning approval flow ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."node_order" IS 'Approval order starting from 1, unique per flow excluding soft-deleted rows';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."approver_id" IS 'Approver user ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow_node."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_oa_approval_flow_node_order
    ON plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order")
    WHERE "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_flow_node_flow ON plugin_linapro_oa_approval_flow_node ("tenant_id", "flow_id", "node_order");

CREATE TABLE IF NOT EXISTS plugin_linapro_oa_approval_request (
    "id"                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"           INT           NOT NULL DEFAULT 0,
    "flow_id"             BIGINT        NOT NULL,
    "flow_type"           SMALLINT      NOT NULL DEFAULT 1,
    "title"               VARCHAR(256)  NOT NULL,
    "amount"              NUMERIC(14,2) NOT NULL DEFAULT 0,
    "content"             TEXT          NOT NULL DEFAULT '',
    "attachments"         TEXT          NOT NULL DEFAULT '',
    "status"              SMALLINT      NOT NULL DEFAULT 1,
    "current_node_order"  INT           NOT NULL DEFAULT 1,
    "current_approver_id" BIGINT        NOT NULL DEFAULT 0,
    "node_snapshot"       TEXT          NOT NULL DEFAULT '',
    "applicant_id"        BIGINT        NOT NULL DEFAULT 0,
    "created_by"          BIGINT        NOT NULL DEFAULT 0,
    "updated_by"          BIGINT        NOT NULL DEFAULT 0,
    "created_at"          TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"          TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"          TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_oa_approval_request IS 'OA approval request document';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."id" IS 'Approval request ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."flow_id" IS 'Originating approval flow ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."flow_type" IS 'Flow type: 1=fund request, 2=write-off, 3=expense reimbursement';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."title" IS 'Request title';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."amount" IS 'Requested amount with two decimal places';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."content" IS 'Applicant statement content';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."attachments" IS 'Attachment URL list stored as a JSON string array of URL addresses, max 10 items';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."status" IS 'Request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."current_node_order" IS 'Pending node order while the request is pending';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."current_approver_id" IS 'Pending approver user ID while the request is pending';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."node_snapshot" IS 'Frozen node list as a JSON array with order, approverId and approverName fields at submit time';
-- Iteration 002: dynamic form columns for per-flow configurable request fields.
-- Comments for these columns run after the conditional ALTER so existing
-- databases gain the columns before the comments reference them.
-- 迭代 002：动态表单列，支持按流程自定义审批单字段；注释在条件加列之后执行以保证存量库顺序正确。
ALTER TABLE plugin_linapro_oa_approval_flow ADD COLUMN IF NOT EXISTS "form_fields" TEXT NOT NULL DEFAULT '';
ALTER TABLE plugin_linapro_oa_approval_request ADD COLUMN IF NOT EXISTS "form_data" TEXT NOT NULL DEFAULT '';
ALTER TABLE plugin_linapro_oa_approval_request ADD COLUMN IF NOT EXISTS "form_snapshot" TEXT NOT NULL DEFAULT '';
COMMENT ON COLUMN plugin_linapro_oa_approval_flow."form_fields" IS 'Configurable form field definitions as a JSON array';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."form_data" IS 'Submitted dynamic form values as a JSON object keyed by field key';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."form_snapshot" IS 'Frozen form field definitions as a JSON array at submit time';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."applicant_id" IS 'Applicant user ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_oa_approval_request."deleted_at" IS 'Deletion time';

CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_request_tenant_applicant ON plugin_linapro_oa_approval_request ("tenant_id", "applicant_id");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_request_tenant_status ON plugin_linapro_oa_approval_request ("tenant_id", "status");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_request_tenant_approver ON plugin_linapro_oa_approval_request ("tenant_id", "current_approver_id", "status");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_request_tenant_type ON plugin_linapro_oa_approval_request ("tenant_id", "flow_type");

CREATE TABLE IF NOT EXISTS plugin_linapro_oa_approval_record (
    "id"          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"   INT          NOT NULL DEFAULT 0,
    "request_id"  BIGINT       NOT NULL,
    "node_order"  INT          NOT NULL DEFAULT 0,
    "action"      SMALLINT     NOT NULL DEFAULT 1,
    "actor_id"    BIGINT       NOT NULL DEFAULT 0,
    "comment"     TEXT         NOT NULL DEFAULT '',
    "created_by"  BIGINT       NOT NULL DEFAULT 0,
    "updated_by"  BIGINT       NOT NULL DEFAULT 0,
    "created_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"  TIMESTAMPTZ  NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_oa_approval_record IS 'OA approval record timeline';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."id" IS 'Approval record ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."request_id" IS 'Owning approval request ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."node_order" IS 'Related node order; 0 means request-level actions';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."action" IS 'Action: 1=submit, 2=approve, 3=reject, 4=comment, 5=append approver, 6=withdraw';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."actor_id" IS 'Actor user ID';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."comment" IS 'Reply or comment content attached to the action';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_oa_approval_record."deleted_at" IS 'Deletion time';

CREATE INDEX IF NOT EXISTS idx_plugin_linapro_oa_approval_record_tenant_request ON plugin_linapro_oa_approval_record ("tenant_id", "request_id", "created_at");

-- Dictionary seeds for OA approval enums. Labels are governed by the host dict module.
-- OA 审批枚举的字典 seed。标签由宿主字典模块治理。
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('审批类型', 'plugin_oa_approval_flow_type', 1, 1, 'OA审批-审批类型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('审批单状态', 'plugin_oa_approval_status', 1, 1, 'OA审批-审批单状态', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('审批动作', 'plugin_oa_approval_action', 1, 1, 'OA审批-审批动作', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_flow_type', '请款', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_flow_type', '核销', '2', 2, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_flow_type', '报销', '3', 3, 'info', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_status', '审批中', '1', 1, 'processing', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_status', '已通过', '2', 2, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_status', '已驳回', '3', 3, 'error', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_status', '已撤回', '4', 4, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '发起', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '通过', '2', 2, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '驳回', '3', 3, 'error', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '回复', '4', 4, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '加签', '5', 5, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_oa_approval_action', '撤回', '6', 6, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
