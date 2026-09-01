# Live Management Plugin

`linapro-live-manage` is the official source plugin for live-streaming back-office management. It provides live-room maintenance and live-content maintenance with tenant isolation, dictionary-driven enums, and plugin menu governance.

## Capabilities

| Capability | Description |
| --- | --- |
| Live room management | Create, query, update, and delete live rooms with tenant-unique room codes, types, and statuses. |
| Live content management | Maintain live events, including titles, push/play/cover/page URLs, song lists, sermon info, scripture, and states. |
| Room availability governance | Live content must reference a live room that exists and is not disabled; referenced rooms cannot be deleted. |
| Bounded room options | Live content forms consume a bounded room candidate API with minimal projection. |

## Routes

| Menu | Path | Permission |
| --- | --- | --- |
| Live room management | `/live/room` | `live:room:list` |
| Live content | `/live/content` | `live:live:list` |

Button permissions: `live:room:query`, `live:room:add`, `live:room:edit`, `live:room:remove`, `live:live:query`, `live:live:add`, `live:live:edit`, `live:live:remove`.

## Data Tables

| Table | Purpose |
| --- | --- |
| `plugin_linapro_live_manage_room` | Tenant-scoped live rooms with audit and soft-delete fields. |
| `plugin_linapro_live_manage_live` | Tenant-scoped live content with audit and soft-delete fields. |

## Dictionary Types

| Type | Values |
| --- | --- |
| `plugin_live_room_type` | `1` gathering, `2` event, `3` other |
| `plugin_live_room_status` | `0` idle, `1` live, `2` disabled |
| `plugin_live_state` | `0` not started, `1` ongoing, `2` finished |
| `plugin_live_public` | `1` public, `0` private |
