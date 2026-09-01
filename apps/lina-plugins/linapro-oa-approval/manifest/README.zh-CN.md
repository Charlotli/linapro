# Manifest 资源

该目录存放 `linapro-oa-approval` 的安装与卸载 SQL 资源。

## 内容

- `sql/001-linapro-oa-approval-schema.sql`: 创建流程、节点、审批单与记录表并初始化索引与字典
- `sql/uninstall/001-linapro-oa-approval-schema.sql`: 卸载清理数据时删除字典并移除数据表
- `sql/mock-data/001-linapro-oa-approval-mock-data.sql`: 本地演示用的可选流程与审批单数据
