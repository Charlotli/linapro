# Live Management Plugin

`linapro-live-manage` is the official source plugin for live-streaming back-office management. It provides live-room maintenance and live-content maintenance with tenant isolation, dictionary-driven enums, and plugin menu governance.

## Capabilities

| Capability | Description |
| --- | --- |
| Live room management | Create, query, update, and delete live rooms with tenant-unique room codes, types, and statuses. |
| Live content management | Maintain live events, including titles, push/play/cover/page URLs, song lists, sermon info, scripture, and states. |
| Room availability governance | Live content must reference a live room that exists and is not disabled; referenced rooms cannot be deleted. |
| Bounded room options | Live content forms consume a bounded room candidate API with minimal projection. |
| Viewer H5 play page | Anonymous mobile page that plays the public HLS stream of one room, with replay and preview states. |
| Live schedule calendar subscription | Anonymous public ICS subscription API so viewers can sync the public live schedule in their system calendar. |

## Routes

| Menu | Path | Permission |
| --- | --- | --- |
| Live room management | `/live/room` | `live:room:list` |
| Live content | `/live/content` | `live:live:list` |

Button permissions: `live:room:query`, `live:room:add`, `live:room:edit`, `live:room:remove`, `live:live:query`, `live:live:add`, `live:live:edit`, `live:live:remove`.

## Viewer H5 Play Page

The plugin serves an anonymous, version-stable viewer page from its own public routes:

| URL | Purpose |
| --- | --- |
| `GET /x/linapro-live-manage/h5?room={roomCode}&tenant={tenantId}` | Mobile H5 player page; `tenant` is required only when the tenant plugin is enabled. |
| `GET /x/linapro-live-manage/api/v1/play?roomCode={roomCode}&tenantId={tenantId}` | Public play-info API behind the same tenant rules. |
| `GET /x/linapro-live-manage/api/v1/subscribe?roomCode={roomCode}&tenantId={tenantId}` | Public calendar subscription API returning an RFC 5545 ICS stream. |

- Selection priority: ongoing public live first, latest finished public live as replay, latest not-started public live as preview (no play URL).
- Private lives, other tenants, and missing rooms all answer with the same not-found error; the push URL and other administrative fields are never exposed.
- Playback uses the embedded `hls.js` (no external CDN) with native HLS fallback on iOS Safari; the page polls every 30 seconds while the live has not started or nothing is available.

## Live Schedule Calendar Subscription

The subscription API covers public lives from the last 7 days through the next 90 days, ordered by live date:

- Every public live renders one `VEVENT` whose `UID` is keyed by the live record, so refreshing or re-subscribing never duplicates entries.
- Lives without a recorded start time render as all-day entries on the live date; all other entries carry UTC start and end times.
- The preview state of the H5 page offers an "Add to calendar" entry, and the admin live-room list offers a subscription-link dialog.
- Missing rooms, no public lives, and failed tenant validation all answer with a valid empty calendar (HTTP 200), so room existence stays undisclosed.

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
