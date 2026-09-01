// This file exposes the plugin-local generated live entity through the
// service package so controller code can keep a stable type name.

package live

import entitymodel "lina-plugin-linapro-live-manage/backend/internal/model/entity"

// LiveEntity mirrors the generated plugin_linapro_live_manage_live entity owned by this plugin.
type LiveEntity = entitymodel.Live
