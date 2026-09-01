# Manifest 资源

该目录存放 `linapro-live-manage` 的安装与卸载 SQL 资源。

## 内容

- `sql/001-linapro-live-manage-schema.sql`：创建直播间与直播内容表、索引并初始化字典
- `sql/uninstall/001-linapro-live-manage-schema.sql`：在卸载且选择清理数据时删除相关字典并移除两张表
- `sql/mock-data/001-linapro-live-manage-mock-data.sql`：本地演示用的可选直播间与直播内容数据
