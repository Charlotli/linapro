// This file exposes the plugin-local generated approval entities through the
// service package so controller code can keep stable type names.

package approval

import entitymodel "lina-plugin-linapro-oa-approval/backend/internal/model/entity"

// RequestEntity mirrors the generated plugin_linapro_oa_approval_request entity owned by this plugin.
type RequestEntity = entitymodel.Request

// RecordEntity mirrors the generated plugin_linapro_oa_approval_record entity owned by this plugin.
type RecordEntity = entitymodel.Record

// FlowEntity mirrors the generated plugin_linapro_oa_approval_flow entity owned by this plugin.
type FlowEntity = entitymodel.Flow

// FlowNodeEntity mirrors the generated plugin_linapro_oa_approval_flow_node entity owned by this plugin.
type FlowNodeEntity = entitymodel.FlowNode
