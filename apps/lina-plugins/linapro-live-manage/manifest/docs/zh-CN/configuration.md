# 配置说明

`linapro-live-manage`当前不引入插件作用域运行期配置。所有业务数据保存在插件数据表中，枚举标签由宿主字典模块治理。

如需扩展插件配置，请在插件根`hack/config.yaml`维护生成工具配置，并在宿主主配置`plugin.linapro-live-manage`段或`plugins/<plugin-id>/config.yaml`中提供配置来源。
