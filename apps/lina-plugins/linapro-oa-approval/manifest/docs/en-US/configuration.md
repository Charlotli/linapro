# Configuration

`linapro-oa-approval` currently introduces no plugin-scoped runtime configuration. All business data is stored in plugin-owned tables, and enum labels are governed by the host dictionary module.

To extend plugin configuration, maintain the generation tool config in the plugin-root `hack/config.yaml`, and provide configuration sources through the host main config section `plugin.linapro-oa-approval` or `plugins/<plugin-id>/config.yaml`.
