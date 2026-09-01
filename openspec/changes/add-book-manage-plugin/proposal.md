## Why

组织内部需要管理共享书籍：登记书籍台账、跟踪借阅流转。需要一个官方源码插件提供书籍登记与借阅管理能力，库存随借出/归还自动联动。

## What Changes

- 新增源码插件`linapro-book-manage`（书籍管理），`distribution: managed`，多租户`tenant_aware`。
- 新增两张插件业务表：
  - `plugin_linapro_book_manage_book`：书籍台账（同租户唯一 ISBN、书名、作者、分类、出版社、出版日期、库存总数、可借数量、存放位置、状态、封面地址、备注）。
  - `plugin_linapro_book_manage_borrow`：借阅记录（书籍关联、借阅人、借出日期、应还日期、归还日期、状态、备注）。
- 借出时校验并扣减可借数量，归还时回增（事务内联动）；删除有借阅记录的书籍被拒绝。
- 新增插件字典 seed：书籍分类、借阅状态。
- 新增插件菜单（宿主`content`目录）：书籍登记、借阅记录及查询/新增/修改/删除按钮权限点。

## Capabilities

### New Capabilities

- `book-manage`：书籍台账与借阅管理能力。

### Modified Capabilities

- 无。

## Impact

- 影响范围限定在`apps/lina-plugins/linapro-book-manage/`与`openspec/changes/add-book-manage-plugin/`；不修改`lina-core`与既有插件。
- `i18n`影响判断：插件启用`i18n`，维护双语菜单、错误、文案与`apidoc`翻译资源。
- 缓存一致性影响判断：无缓存状态，无影响。
- 数据权限影响判断：读写均在数据库侧注入租户过滤；借阅人/存放位置为展示字段，组织维度不接入，例外见`design.md`。
- 开发工具跨平台影响判断：复用共享`codegen`入口，无新增脚本。
- 测试策略影响判断：`API`契约测试、服务层单测（库存联动）与`E2E`用例；`openspec`工具未安装记录阻断。
