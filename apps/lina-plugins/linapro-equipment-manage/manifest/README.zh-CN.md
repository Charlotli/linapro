# Manifest 资源

该目录存放`linapro-equipment-manage`的安装与卸载 SQL 资源。

## 内容

- `sql/001-linapro-equipment-manage-schema.sql`：创建设备与维护记录表并初始化索引与字典
- `sql/uninstall/001-linapro-equipment-manage-schema.sql`：卸载清理数据时删除字典并移除数据表
- `sql/mock-data/001-linapro-equipment-manage-mock-data.sql`：本地演示用的可选设备与维护记录
