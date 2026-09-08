-- 001: linapro-live-manage schema uninstall
-- 001：linapro-live-manage 数据结构卸载

DELETE FROM sys_dict_data WHERE "dict_type" IN ('plugin_live_room_type', 'plugin_live_room_status', 'plugin_live_state', 'plugin_live_public');
DELETE FROM sys_dict_type WHERE "type" IN ('plugin_live_room_type', 'plugin_live_room_status', 'plugin_live_state', 'plugin_live_public');
DROP TABLE IF EXISTS plugin_linapro_live_manage_bible_verse;
DROP TABLE IF EXISTS plugin_linapro_live_manage_bible_book;
DROP TABLE IF EXISTS plugin_linapro_live_manage_announcement;
DROP TABLE IF EXISTS plugin_linapro_live_manage_view_session;
DROP TABLE IF EXISTS plugin_linapro_live_manage_live;
DROP TABLE IF EXISTS plugin_linapro_live_manage_room;
