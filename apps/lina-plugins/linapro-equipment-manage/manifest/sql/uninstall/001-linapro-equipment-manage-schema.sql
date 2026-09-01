-- 001: linapro-equipment-manage schema uninstall
-- 001：linapro-equipment-manage 数据结构卸载

DELETE FROM sys_dict_data WHERE "dict_type" IN ('plugin_equipment_type', 'plugin_equipment_status', 'plugin_equipment_maint_type');
DELETE FROM sys_dict_type WHERE "type" IN ('plugin_equipment_type', 'plugin_equipment_status', 'plugin_equipment_maint_type');
DROP TABLE IF EXISTS plugin_linapro_equipment_manage_maintenance;
DROP TABLE IF EXISTS plugin_linapro_equipment_manage_equipment;
