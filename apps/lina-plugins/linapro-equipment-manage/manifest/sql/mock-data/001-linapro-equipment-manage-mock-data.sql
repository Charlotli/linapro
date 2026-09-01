-- Mock data: equipment registry and maintenance records for demos.
-- 模拟数据：设备管理演示使用的设备与维护记录。
INSERT INTO plugin_linapro_equipment_manage_equipment ("tenant_id", "equipment_code", "equipment_name", "equipment_type", "brand_model", "purchase_date", "purchase_price", "location", "owner", "status", "remark", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    0,
    'EQ-PROJECTOR-01',
    '会议室投影仪',
    2,
    'Epson CB-X06',
    '2025-06-15',
    3299.00,
    '主堂会议室',
    '张弟兄',
    1,
    '',
    admin."id",
    admin."id",
    '2026-04-20 09:00:00',
    '2026-04-20 09:00:00'
FROM sys_user admin
WHERE admin."username" = 'admin'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_equipment_manage_equipment eq
    WHERE eq."tenant_id" = 0 AND eq."equipment_code" = 'EQ-PROJECTOR-01'
  );

INSERT INTO plugin_linapro_equipment_manage_maintenance ("tenant_id", "equipment_id", "maint_type", "maint_date", "maintainer", "cost", "content", "result", "remark", "created_by", "updated_by", "created_at", "updated_at")
SELECT
    eq."tenant_id",
    eq."id",
    3,
    '2026-04-18',
    '李弟兄',
    0,
    '季度巡检：灯泡亮度与接口检查',
    '正常',
    '',
    admin."id",
    admin."id",
    '2026-04-18 10:00:00',
    '2026-04-18 10:00:00'
FROM plugin_linapro_equipment_manage_equipment eq
JOIN sys_user admin ON admin."username" = 'admin'
WHERE eq."tenant_id" = 0 AND eq."equipment_code" = 'EQ-PROJECTOR-01'
  AND NOT EXISTS (
    SELECT 1 FROM plugin_linapro_equipment_manage_maintenance mt
    WHERE mt."tenant_id" = eq."tenant_id" AND mt."equipment_id" = eq."id" AND mt."maint_date" = '2026-04-18'
  );
