-- 001: linapro-live-manage schema
-- 001：linapro-live-manage 数据结构

-- Purpose: Stores tenant-scoped live rooms and live content, including push/play addresses, song lists, sermon info, and audit fields.
-- 用途：存储租户级直播间与直播内容，包括推流/播放地址、歌单、讲道信息与审计字段。
CREATE TABLE IF NOT EXISTS plugin_linapro_live_manage_room (
    "id"          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"   INT           NOT NULL DEFAULT 0,
    "room_code"   VARCHAR(64)   NOT NULL,
    "room_name"   VARCHAR(128)  NOT NULL,
    "room_type"   SMALLINT      NOT NULL DEFAULT 1,
    "status"      SMALLINT      NOT NULL DEFAULT 0,
    "description" VARCHAR(512)  NOT NULL DEFAULT '',
    "created_by"  BIGINT        NOT NULL DEFAULT 0,
    "updated_by"  BIGINT        NOT NULL DEFAULT 0,
    "created_at"  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"  TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_live_manage_room IS 'Live room table for live streaming management';
COMMENT ON COLUMN plugin_linapro_live_manage_room."id" IS 'Live room ID';
COMMENT ON COLUMN plugin_linapro_live_manage_room."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_live_manage_room."room_code" IS 'Room code exposed to external systems, unique per tenant';
COMMENT ON COLUMN plugin_linapro_live_manage_room."room_name" IS 'Room name';
COMMENT ON COLUMN plugin_linapro_live_manage_room."room_type" IS 'Room type: 1=gathering, 2=event, 3=other';
COMMENT ON COLUMN plugin_linapro_live_manage_room."status" IS 'Room status: 0=idle, 1=live, 2=disabled';
COMMENT ON COLUMN plugin_linapro_live_manage_room."description" IS 'Room description';
COMMENT ON COLUMN plugin_linapro_live_manage_room."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_live_manage_room."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_live_manage_room."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_live_manage_room."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_live_manage_room."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_live_manage_room_tenant_code ON plugin_linapro_live_manage_room ("tenant_id", "room_code");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_room_tenant_status ON plugin_linapro_live_manage_room ("tenant_id", "status");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_room_tenant_type ON plugin_linapro_live_manage_room ("tenant_id", "room_type");

CREATE TABLE IF NOT EXISTS plugin_linapro_live_manage_live (
    "id"                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"         INT           NOT NULL DEFAULT 0,
    "room_id"           BIGINT        NOT NULL,
    "title"             VARCHAR(512)  NOT NULL DEFAULT '',
    "content"           TEXT          NOT NULL DEFAULT '',
    "live_date"         DATE          NOT NULL DEFAULT CURRENT_DATE,
    "page_url"          VARCHAR(2048) NOT NULL DEFAULT '',
    "push_url"          TEXT          NOT NULL DEFAULT '',
    "live_url"          VARCHAR(2048) NOT NULL DEFAULT '',
    "cover_url"         VARCHAR(2048) NOT NULL DEFAULT '',
    "song_name"         VARCHAR(512)  NOT NULL DEFAULT '',
    "song_list"         TEXT          NOT NULL DEFAULT '',
    "lead_singer"       VARCHAR(256)  NOT NULL DEFAULT '',
    "accompaniment"     VARCHAR(256)  NOT NULL DEFAULT '',
    "host"              VARCHAR(256)  NOT NULL DEFAULT '',
    "sermon_title"      VARCHAR(512)  NOT NULL DEFAULT '',
    "preacher"          VARCHAR(256)  NOT NULL DEFAULT '',
    "preacher_identity" VARCHAR(256)  NOT NULL DEFAULT '',
    "scripture_ref"     VARCHAR(256)  NOT NULL DEFAULT '',
    "scripture_content" TEXT          NOT NULL DEFAULT '',
    "outline"           TEXT          NOT NULL DEFAULT '',
    "device_info"       VARCHAR(512)  NOT NULL DEFAULT '',
    "reception"         VARCHAR(512)  NOT NULL DEFAULT '',
    "state"             SMALLINT      NOT NULL DEFAULT 0,
    "is_public"         SMALLINT      NOT NULL DEFAULT 1,
    "start_time"        TIMESTAMPTZ   NULL DEFAULT NULL,
    "created_by"        BIGINT        NOT NULL DEFAULT 0,
    "updated_by"        BIGINT        NOT NULL DEFAULT 0,
    "created_at"        TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"        TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"        TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_live_manage_live IS 'Live content table for live streaming management';
COMMENT ON COLUMN plugin_linapro_live_manage_live."id" IS 'Live content ID';
COMMENT ON COLUMN plugin_linapro_live_manage_live."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_live_manage_live."room_id" IS 'Owning live room ID within the same tenant';
COMMENT ON COLUMN plugin_linapro_live_manage_live."title" IS 'Live title shown on list and cards';
COMMENT ON COLUMN plugin_linapro_live_manage_live."content" IS 'Live description or rich content body';
COMMENT ON COLUMN plugin_linapro_live_manage_live."live_date" IS 'Calendar date of the live event for list sorting and filtering';
COMMENT ON COLUMN plugin_linapro_live_manage_live."page_url" IS 'External page link of the live event';
COMMENT ON COLUMN plugin_linapro_live_manage_live."push_url" IS 'Stream push URL for the publisher';
COMMENT ON COLUMN plugin_linapro_live_manage_live."live_url" IS 'Play URL for viewers, such as HLS or RTMP';
COMMENT ON COLUMN plugin_linapro_live_manage_live."cover_url" IS 'Cover image URL for cards and previews';
COMMENT ON COLUMN plugin_linapro_live_manage_live."song_name" IS 'Featured song or hymn name';
COMMENT ON COLUMN plugin_linapro_live_manage_live."song_list" IS 'Song list as a JSON array text, e.g. [{"name":"Song A","singer":"Alice","order":1}]';
COMMENT ON COLUMN plugin_linapro_live_manage_live."lead_singer" IS 'Lead singer or worship team';
COMMENT ON COLUMN plugin_linapro_live_manage_live."accompaniment" IS 'Accompaniment info such as band or format';
COMMENT ON COLUMN plugin_linapro_live_manage_live."host" IS 'Host name or role';
COMMENT ON COLUMN plugin_linapro_live_manage_live."sermon_title" IS 'Sermon title when the live has a standalone sermon item';
COMMENT ON COLUMN plugin_linapro_live_manage_live."preacher" IS 'Preacher name';
COMMENT ON COLUMN plugin_linapro_live_manage_live."preacher_identity" IS 'Preacher identity or title, such as pastor';
COMMENT ON COLUMN plugin_linapro_live_manage_live."scripture_ref" IS 'Scripture reference, e.g. John 3:16';
COMMENT ON COLUMN plugin_linapro_live_manage_live."scripture_content" IS 'Scripture content for display';
COMMENT ON COLUMN plugin_linapro_live_manage_live."outline" IS 'Sermon outline or agenda for preview';
COMMENT ON COLUMN plugin_linapro_live_manage_live."device_info" IS 'Device info free text, e.g. camera and microphone';
COMMENT ON COLUMN plugin_linapro_live_manage_live."reception" IS 'Reception arrangement or contact info';
COMMENT ON COLUMN plugin_linapro_live_manage_live."state" IS 'Live state: 0=not started, 1=ongoing, 2=finished';
COMMENT ON COLUMN plugin_linapro_live_manage_live."is_public" IS 'Visibility: 1=public, 0=private';
COMMENT ON COLUMN plugin_linapro_live_manage_live."start_time" IS 'Actual start time of the live for schedule and reminders';
COMMENT ON COLUMN plugin_linapro_live_manage_live."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_live_manage_live."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_live_manage_live."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_live_manage_live."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_live_manage_live."deleted_at" IS 'Deletion time';

CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_live_tenant_room ON plugin_linapro_live_manage_live ("tenant_id", "room_id");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_live_tenant_state ON plugin_linapro_live_manage_live ("tenant_id", "state");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_live_tenant_date ON plugin_linapro_live_manage_live ("tenant_id", "live_date");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_live_manage_live_tenant_start ON plugin_linapro_live_manage_live ("tenant_id", "start_time");

-- Enforce at most one ongoing live per room per tenant at the database level so
-- concurrent start actions cannot bypass the service-side occupancy check.
-- Soft-deleted rows are excluded from the constraint.
-- 数据库层强制同租户同直播间最多一场进行中的直播，防止并发开启绕过服务端占用校验；软删除行不参与约束。
CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_live_manage_live_tenant_room_ongoing
    ON plugin_linapro_live_manage_live ("tenant_id", "room_id")
    WHERE "state" = 1 AND "deleted_at" IS NULL;

-- Dictionary seeds for live-manage enums. Labels are governed by the host dict module.
-- 直播管理枚举的字典 seed。标签由宿主字典模块治理。
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('直播间类型', 'plugin_live_room_type', 1, 1, '直播后台管理-直播间类型', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('直播间状态', 'plugin_live_room_status', 1, 1, '直播后台管理-直播间状态', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('直播状态', 'plugin_live_state', 1, 1, '直播后台管理-直播状态', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('直播公开性', 'plugin_live_public', 1, 1, '直播后台管理-直播公开性', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_type', '正常聚会', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_type', '活动', '2', 2, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_type', '其他', '3', 3, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_status', '空闲', '0', 1, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_status', '直播中', '1', 2, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_room_status', '禁用', '2', 3, 'error', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_state', '未开始', '0', 1, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_state', '进行中', '1', 2, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_state', '已结束', '2', 3, 'info', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_public', '公开', '1', 1, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_live_public', '私有', '0', 2, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
