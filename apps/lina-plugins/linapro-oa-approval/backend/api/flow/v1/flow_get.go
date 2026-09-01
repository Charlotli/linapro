// This file declares the get-approval-flow request/response DTOs used by the
// linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Flow Get API

// GetReq defines the request for retrieving approval-flow details.
type GetReq struct {
	g.Meta `path:"/approval/flow/{id}" method:"get" tags:"ApprovalFlows" summary:"Get approval flow details" dc:"Get approval flow details by ID, including the ordered approver nodes and the configurable form field definitions. Flows outside the current tenant are reported as not found." permission:"oa:flow:query"`
	Id     int64 `json:"id" v:"required" dc:"Approval flow ID" eg:"1"`
}

// FlowFieldColumn exposes one column of a detail field definition.
type FlowFieldColumn struct {
	Key   string `json:"key" dc:"Storage key of the column" eg:"amount"`
	Label string `json:"label" dc:"Display name of the column" eg:"Amount"`
	Type  string `json:"type" dc:"Column type: text/number/date" eg:"number"`
}

// FlowFieldItem exposes one configurable form field definition of a flow.
type FlowFieldItem struct {
	Key      string             `json:"key" dc:"Storage key of the field" eg:"expenseType"`
	Label    string             `json:"label" dc:"Display name of the field" eg:"Expense type"`
	Type     string             `json:"type" dc:"Field type: text/textarea/number/date/select/attachment/detail" eg:"select"`
	Required bool               `json:"required" dc:"Whether the field is required" eg:"true"`
	AsAmount bool               `json:"asAmount" dc:"Whether a number field projects to the request amount" eg:"false"`
	Options  []string           `json:"options" dc:"Options for select fields" eg:"[]"`
	Columns  []*FlowFieldColumn `json:"columns" dc:"Column definitions for detail fields" eg:"[]"`
}

// GetRes Approval flow detail response
type GetRes struct {
	FlowItem
	Nodes  []*FlowNodeItem  `json:"nodes" dc:"Ordered approver nodes of the flow" eg:"[]"`
	Fields []*FlowFieldItem `json:"fields" dc:"Configurable form field definitions; empty means the built-in standard form" eg:"[]"`
}
