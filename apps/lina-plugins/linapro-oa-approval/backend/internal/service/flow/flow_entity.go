// This file exposes the plugin-local generated flow entities through the
// service package so controller code can keep stable type names.

package flow

import entitymodel "lina-plugin-linapro-oa-approval/backend/internal/model/entity"

// FlowEntity mirrors the generated plugin_linapro_oa_approval_flow entity owned by this plugin.
type FlowEntity = entitymodel.Flow

// FlowNodeEntity mirrors the generated plugin_linapro_oa_approval_flow_node entity owned by this plugin.
type FlowNodeEntity = entitymodel.FlowNode
