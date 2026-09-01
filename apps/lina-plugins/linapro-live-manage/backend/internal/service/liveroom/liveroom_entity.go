// This file exposes the plugin-local generated room entity through the
// service package so controller code can keep a stable type name.

package liveroom

import entitymodel "lina-plugin-linapro-live-manage/backend/internal/model/entity"

// RoomEntity mirrors the generated plugin_linapro_live_manage_room entity owned by this plugin.
type RoomEntity = entitymodel.Room
