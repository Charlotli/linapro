## Why

组织内部需要管理办公与 IT 设备：登记设备台账、跟踪维护保养历史。需要一个官方源码插件提供设备登记与维护记录的完整管理能力，闭环在插件内部。

## What Changes

- 新增源码插件`linapro-equipment-manage`（设备管理），`distribution: managed`，多租户`tenant_aware`。
- 新增两张插件业务表：
  - `plugin_linapro_equipment_manage_equipment`：设备台账（同租户唯一设备编号、名称、类型、品牌型号、购置日期、购置价格、存放位置、负责人、状态、备注）。
  - `plugin_linapro_equipment_manage_maintenance`：维护记录（设备关联、维护类型、维护日期、维护人、费用、内容、结果、备注）。
- 新增设备与维护记录的`RESTful CRUD`接口与有界候选接口；删除仍被维护记录引用的设备被拒绝；创建维护记录校验设备存在。
- 新增插件字典 seed：设备类型、设备状态、维护类型。
- 新增插件菜单（宿主`content`目录）：设备登记、维护记录及查询/新增/修改/删除按钮权限点。

## Capabilities

### New Capabilities

- `equipment-manage`：设备台账与维护记录管理能力。

### Modified Capabilities

- 无。

## Impact

- 影响范围限定在`apps/lina-plugins/linapro-equipment-manage/`与`openspec/changes/add-equipment-manage-plugin/`；不修改`lina-core`与既有插件。
- `i18n`影响判断：插件启用`i18n`，维护双语菜单、错误、文案与`apidoc`翻译资源。
- 缓存一致性影响判断：无缓存状态，无影响。
- 数据权限影响判断：读写均在数据库侧注入租户过滤；组织维度不接入（设备负责人为展示字段而非权限边界），例外见`design.md`。
- 开发工具跨平台影响判断：复用共享`codegen`入口，无新增脚本。
- 测试策略影响判断：`API`契约测试、服务层单测与`E2E`用例；`openspec`工具未安装记录阻断。
