## 1. 骨架与数据契约

- [x] 1.1 插件骨架（`plugin.yaml`两菜单八按钮、`plugin_embed.go`、`go.mod`、`Makefile`、`hack/config.yaml`、`backend/plugin.go`）。
- [x] 1.2 幂等`SQL`（两表、部分唯一索引、字典 seed）、卸载与`mock-data`。
- [x] 1.3 `backend/api/book`与`backend/api/borrow`的`RESTful DTO`。
- [x] 1.4 建表并生成`dao/ctrl`。

## 2. 后端实现

- [x] 2.1 书籍服务：租户过滤`CRUD`、`ISBN`唯一校验、删除引用保护。
- [x] 2.2 借阅服务：借出/归还事务库存联动、书籍存在校验、列表书名批量装配、`bizerr`错误码。
- [x] 2.3 控制器与路由绑定，契约测试与单测通过。

## 3. 前端实现

- [x] 3.1 书籍登记页（列表/筛选/新增编辑弹窗/删除确认）。
- [x] 3.2 借阅记录页（借出登记书籍候选联动、归还动作）。

## 4. 资源与文档

- [x] 4.1 双语`i18n`资源与`apidoc`翻译（`en-US`空占位）。
- [x] 4.2 `README`双语镜像与`manifest/docs`。

## 5. 测试与验证

- [x] 5.1 `E2E`用例与页面对象落位。
- [ ] 5.2 编译门禁与`plugins.check`/`i18n.check`治理检查。
- [ ] 5.3 `openspec validate`；工具未安装时记录阻断。

## 6. 影响记录与审查

- [x] 6.1 影响判断：`i18n`、缓存一致性、数据权限、跨平台、测试策略、`DI`来源。
- [ ] 6.2 `lina-review`审查通过后标记完成。
