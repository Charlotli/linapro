-- 001: linapro-equipment-manage schema
-- 001：linapro-equipment-manage 数据结构

-- Purpose: Stores tenant-scoped equipment registry and maintenance records.
-- 用途：存储租户级设备台账与维护记录。
CREATE TABLE IF NOT EXISTS plugin_linapro_equipment_manage_equipment (
    "id"             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"      INT           NOT NULL DEFAULT 0,
    "equipment_code" VARCHAR(64)   NOT NULL,
    "equipment_name" VARCHAR(128)  NOT NULL,
    "equipment_type" SMALLINT      NOT NULL DEFAULT 1,
    "brand_model"    VARCHAR(128)  NOT NULL DEFAULT '',
    "purchase_date"  DATE          NULL DEFAULT NULL,
    "purchase_price" NUMERIC(12,2) NOT NULL DEFAULT 0,
    "location"       VARCHAR(128)  NOT NULL DEFAULT '',
    "owner"          VARCHAR(64)   NOT NULL DEFAULT '',
    "status"         SMALLINT      NOT NULL DEFAULT 1,
    "remark"         VARCHAR(512)  NOT NULL DEFAULT '',
    "created_by"     BIGINT        NOT NULL DEFAULT 0,
    "updated_by"     BIGINT        NOT NULL DEFAULT 0,
    "created_at"     TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"     TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"     TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_equipment_manage_equipment IS 'Equipment registry table';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."id" IS 'Equipment ID';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."equipment_code" IS 'Equipment code, unique per tenant';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."equipment_name" IS 'Equipment name';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."equipment_type" IS 'Equipment type: 1=office, 2=it, 3=production, 4=other';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."brand_model" IS 'Brand and model';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."purchase_date" IS 'Purchase date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."purchase_price" IS 'Purchase price with two decimal places';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."location" IS 'Storage location';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."owner" IS 'Responsible person';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."status" IS 'Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."remark" IS 'Remark';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_equipment_manage_equipment."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_equipment_manage_eq_tenant_code
    ON plugin_linapro_equipment_manage_equipment ("tenant_id", "equipment_code");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_equipment_manage_eq_tenant_status
    ON plugin_linapro_equipment_manage_equipment ("tenant_id", "status");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_equipment_manage_eq_tenant_type
    ON plugin_linapro_equipment_manage_equipment ("tenant_id", "equipment_type");

CREATE TABLE IF NOT EXISTS plugin_linapro_equipment_manage_maintenance (
    "id"             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"      INT           NOT NULL DEFAULT 0,
    "equipment_id"   BIGINT        NOT NULL,
    "maint_type"     SMALLINT      NOT NULL DEFAULT 1,
    "maint_date"     DATE          NOT NULL DEFAULT CURRENT_DATE,
    "maintainer"     VARCHAR(64)   NOT NULL DEFAULT '',
    "cost"           NUMERIC(12,2) NOT NULL DEFAULT 0,
    "content"        TEXT          NOT NULL DEFAULT '',
    "result"         VARCHAR(256)  NOT NULL DEFAULT '',
    "remark"         VARCHAR(512)  NOT NULL DEFAULT '',
    "created_by"     BIGINT        NOT NULL DEFAULT 0,
    "updated_by"     BIGINT        NOT NULL DEFAULT 0,
    "created_at"     TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"     TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"     TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_equipment_manage_maintenance IS 'Equipment maintenance record table';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."id" IS 'Maintenance record ID';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."equipment_id" IS 'Maintained equipment ID within the same tenant';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."maint_type" IS 'Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."maint_date" IS 'Maintenance date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."maintainer" IS 'Maintainer name';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."cost" IS 'Maintenance cost with two decimal places';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."content" IS 'Maintenance content description';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."result" IS 'Maintenance result';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."remark" IS 'Remark';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_equipment_manage_maintenance."deleted_at" IS 'Deletion time';

CREATE INDEX IF NOT EXISTS idx_plugin_linapro_equipment_manage_mt_tenant_equipment
    ON plugin_linapro_equipment_manage_maintenance ("tenant_id", "equipment_id");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_equipment_manage_mt_tenant_date
    ON plugin_linapro_equipment_manage_maintenance ("tenant_id", "maint_date");

-- Dictionary seeds for equipment enums. Labels are governed by the host dict module.
-- 设备枚举的字典 seed。标签由宿主字典模块治理。
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('设备类型', 'plugin_equipment_type', 1, 1, '设备管理-设备类型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('设备状态', 'plugin_equipment_status', 1, 1, '设备管理-设备状态', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('维护类型', 'plugin_equipment_maint_type', 1, 1, '设备管理-维护类型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_type', '办公设备', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_type', 'IT设备', '2', 2, 'info', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_type', '生产设备', '3', 3, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_type', '其他', '4', 4, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_status', '在用', '1', 1, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_status', '闲置', '2', 2, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_status', '维修中', '3', 3, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_status', '报废', '4', 4, 'error', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_maint_type', '保养', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_maint_type', '维修', '2', 2, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_maint_type', '巡检', '3', 3, 'info', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_equipment_maint_type', '其他', '4', 4, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
