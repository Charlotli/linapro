-- 001: linapro-book-manage schema uninstall
-- 001：linapro-book-manage 数据结构卸载

DELETE FROM sys_dict_data WHERE "dict_type" IN ('plugin_book_category', 'plugin_book_borrow_status');
DELETE FROM sys_dict_type WHERE "type" IN ('plugin_book_category', 'plugin_book_borrow_status');
DROP TABLE IF EXISTS plugin_linapro_book_manage_borrow;
DROP TABLE IF EXISTS plugin_linapro_book_manage_book;
