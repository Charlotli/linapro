-- Mock data: book registry and borrow records for demos.
-- 模拟数据：书籍管理演示使用的书籍与借阅记录。
INSERT INTO plugin_linapro_book_manage_book ("tenant_id", "title", "author", "isbn", "category", "publisher", "publish_date", "total_quantity", "available_quantity", "location", "status", "cover_url", "remark", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    '科学的历程',
    '吴国盛',
    '9787301281',
    1,
    '北京大学出版社',
    '2018-01-01',
    3,
    3,
    'A区01架',
    1,
    '',
    '',
    admin."id",
    admin."id",
    '2026-04-20 09:00:00',
    '2026-04-20 09:00:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_book_manage_book b
    WHERE b."tenant_id" = 0 AND b."isbn" = '9787301281'
  );

INSERT INTO plugin_linapro_book_manage_borrow ("tenant_id", "book_id", "borrower", "borrow_date", "due_date", "return_date", "status", "remark", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    b."tenant_id",
    b."id",
    '王姊妹',
    '2026-04-22',
    '2026-05-22',
    NULL,
    1,
    '',
    admin."id",
    admin."id",
    '2026-04-22 10:00:00',
    '2026-04-22 10:00:00'
FROM plugin_linapro_book_manage_book b
JOIN sys_user admin ON admin."username" = 'admin'
WHERE b."tenant_id" = 0 AND b."isbn" = '9787301281'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_book_manage_borrow br
    WHERE br."tenant_id" = b."tenant_id" AND br."book_id" = b."id" AND br."status" = 1
  );

UPDATE plugin_linapro_book_manage_book b
SET "available_quantity" = "available_quantity" - 1
FROM sys_user admin
WHERE admin."username" = 'admin' AND b."tenant_id" = 0 AND b."isbn" = '9787301281'
  AND EXISTS (
    SELECT 1 FROM plugin_linapro_book_manage_borrow br
    WHERE br."tenant_id" = b."tenant_id" AND br."book_id" = b."id" AND br."status" = 1
  )
  AND b."available_quantity" = 3;
