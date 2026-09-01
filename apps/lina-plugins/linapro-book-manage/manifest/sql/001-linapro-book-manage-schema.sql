-- 001: linapro-book-manage schema
-- 001：linapro-book-manage 数据结构

-- Purpose: Stores tenant-scoped book registry and borrow records.
-- 用途：存储租户级书籍台账与借阅记录。
CREATE TABLE IF NOT EXISTS plugin_linapro_book_manage_book (
    "id"                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"         INT           NOT NULL DEFAULT 0,
    "title"             VARCHAR(256)  NOT NULL,
    "author"            VARCHAR(128)  NOT NULL DEFAULT '',
    "isbn"              VARCHAR(32)   NOT NULL DEFAULT '',
    "category"          SMALLINT      NOT NULL DEFAULT 1,
    "publisher"         VARCHAR(128)  NOT NULL DEFAULT '',
    "publish_date"      DATE          NULL DEFAULT NULL,
    "total_quantity"    INT           NOT NULL DEFAULT 1,
    "available_quantity" INT          NOT NULL DEFAULT 1,
    "location"          VARCHAR(128)  NOT NULL DEFAULT '',
    "status"            SMALLINT      NOT NULL DEFAULT 1,
    "cover_url"         VARCHAR(2048) NOT NULL DEFAULT '',
    "remark"            VARCHAR(512)  NOT NULL DEFAULT '',
    "created_by"        BIGINT        NOT NULL DEFAULT 0,
    "updated_by"        BIGINT        NOT NULL DEFAULT 0,
    "created_at"        TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"        TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"        TIMESTAMPTZ   NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_book_manage_book IS 'Book registry table';
COMMENT ON COLUMN plugin_linapro_book_manage_book."id" IS 'Book ID';
COMMENT ON COLUMN plugin_linapro_book_manage_book."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_book_manage_book."title" IS 'Book title';
COMMENT ON COLUMN plugin_linapro_book_manage_book."author" IS 'Author';
COMMENT ON COLUMN plugin_linapro_book_manage_book."isbn" IS 'ISBN, unique per tenant when filled';
COMMENT ON COLUMN plugin_linapro_book_manage_book."category" IS 'Category: 1=technology, 2=literature, 3=management, 4=children, 5=other';
COMMENT ON COLUMN plugin_linapro_book_manage_book."publisher" IS 'Publisher';
COMMENT ON COLUMN plugin_linapro_book_manage_book."publish_date" IS 'Publish date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_book_manage_book."total_quantity" IS 'Total copy count';
COMMENT ON COLUMN plugin_linapro_book_manage_book."available_quantity" IS 'Currently available copy count';
COMMENT ON COLUMN plugin_linapro_book_manage_book."location" IS 'Shelf location';
COMMENT ON COLUMN plugin_linapro_book_manage_book."status" IS 'Book status: 1=on shelf, 2=off shelf';
COMMENT ON COLUMN plugin_linapro_book_manage_book."cover_url" IS 'Cover image URL address';
COMMENT ON COLUMN plugin_linapro_book_manage_book."remark" IS 'Remark';
COMMENT ON COLUMN plugin_linapro_book_manage_book."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_book_manage_book."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_book_manage_book."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_book_manage_book."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_book_manage_book."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_plugin_linapro_book_manage_book_tenant_isbn
    ON plugin_linapro_book_manage_book ("tenant_id", "isbn")
    WHERE "isbn" <> '' AND "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_book_manage_book_tenant_category
    ON plugin_linapro_book_manage_book ("tenant_id", "category");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_book_manage_book_tenant_status
    ON plugin_linapro_book_manage_book ("tenant_id", "status");

CREATE TABLE IF NOT EXISTS plugin_linapro_book_manage_borrow (
    "id"             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"      INT          NOT NULL DEFAULT 0,
    "book_id"        BIGINT       NOT NULL,
    "borrower"       VARCHAR(64)  NOT NULL DEFAULT '',
    "borrow_date"    DATE         NOT NULL DEFAULT CURRENT_DATE,
    "due_date"       DATE         NULL DEFAULT NULL,
    "return_date"    DATE         NULL DEFAULT NULL,
    "status"         SMALLINT     NOT NULL DEFAULT 1,
    "remark"         VARCHAR(512) NOT NULL DEFAULT '',
    "created_by"     BIGINT       NOT NULL DEFAULT 0,
    "updated_by"     BIGINT       NOT NULL DEFAULT 0,
    "created_at"     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"     TIMESTAMPTZ  NULL DEFAULT NULL
);

COMMENT ON TABLE plugin_linapro_book_manage_borrow IS 'Book borrow record table';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."id" IS 'Borrow record ID';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."tenant_id" IS 'Owning tenant ID, 0 means PLATFORM';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."book_id" IS 'Borrowed book ID within the same tenant';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."borrower" IS 'Borrower name';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."borrow_date" IS 'Borrow date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."due_date" IS 'Expected return date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."return_date" IS 'Actual return date with date-only semantics';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."status" IS 'Borrow status: 1=borrowed, 2=returned';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."remark" IS 'Remark';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."created_by" IS 'Creator';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."updated_by" IS 'Updater';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."created_at" IS 'Creation time';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."updated_at" IS 'Update time';
COMMENT ON COLUMN plugin_linapro_book_manage_borrow."deleted_at" IS 'Deletion time';

CREATE INDEX IF NOT EXISTS idx_plugin_linapro_book_manage_bw_tenant_book
    ON plugin_linapro_book_manage_borrow ("tenant_id", "book_id");
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_book_manage_bw_tenant_status
    ON plugin_linapro_book_manage_borrow ("tenant_id", "status");

-- Dictionary seeds for book enums. Labels are governed by the host dict module.
-- 书籍枚举的字典 seed。标签由宿主字典模块治理。
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('书籍分类', 'plugin_book_category', 1, 1, '书籍管理-书籍分类', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_type ("name", "type", "status", "is_builtin", "remark", "created_at", "updated_at")
VALUES ('借阅状态', 'plugin_book_borrow_status', 1, 1, '书籍管理-借阅状态', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_category', '技术', '1', 1, 'primary', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_category', '文学', '2', 2, 'info', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_category', '管理', '3', 3, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_category', '童书', '4', 4, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_category', '其他', '5', 5, 'default', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_borrow_status', '借出', '1', 1, 'warning', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
INSERT INTO sys_dict_data ("dict_type", "label", "value", "sort", "tag_style", "status", "is_builtin", "created_at", "updated_at")
VALUES ('plugin_book_borrow_status', '已归还', '2', 2, 'success', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
