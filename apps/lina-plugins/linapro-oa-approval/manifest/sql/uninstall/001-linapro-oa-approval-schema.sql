-- 001: linapro-oa-approval schema uninstall
-- 001：linapro-oa-approval 数据结构卸载

DELETE FROM sys_dict_data WHERE "dict_type" IN ('plugin_oa_approval_flow_type', 'plugin_oa_approval_status', 'plugin_oa_approval_action');
DELETE FROM sys_dict_type WHERE "type" IN ('plugin_oa_approval_flow_type', 'plugin_oa_approval_status', 'plugin_oa_approval_action');
DROP TABLE IF EXISTS plugin_linapro_oa_approval_record;
DROP TABLE IF EXISTS plugin_linapro_oa_approval_request;
DROP TABLE IF EXISTS plugin_linapro_oa_approval_flow_node;
DROP TABLE IF EXISTS plugin_linapro_oa_approval_flow;
