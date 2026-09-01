// This file declares the update-approval-flow request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Flow Update API

// UpdateNodeItem defines one ordered approver node in an update request.
type UpdateNodeItem struct {
	ApproverId int64 `json:"approverId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Approver user ID" eg:"2"`
}

// UpdateColumnItem defines one detail sub-table column in an update request.
type UpdateColumnItem struct {
	Key   string `json:"key" v:"required#gf.gvalid.rule.required" dc:"Storage key" eg:"amount"`
	Label string `json:"label" v:"required#gf.gvalid.rule.required" dc:"Display name" eg:"Amount"`
	Type  string `json:"type" v:"required|in:text,number,date#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Column type: text/number/date" eg:"number"`
}

// UpdateFieldItem defines one configurable form field in an update request.
type UpdateFieldItem struct {
	Key      string              `json:"key" v:"required#gf.gvalid.rule.required" dc:"Storage key, identifier-safe and unique within the flow" eg:"expenseType"`
	Label    string              `json:"label" v:"required#gf.gvalid.rule.required" dc:"Display name" eg:"Expense type"`
	Type     string              `json:"type" v:"required|in:text,textarea,number,date,select,attachment,detail#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Field type: text/textarea/number/date/select/attachment/detail" eg:"select"`
	Required bool                `json:"required" dc:"Whether the field must be filled" eg:"true"`
	AsAmount bool                `json:"asAmount" dc:"Whether a number field projects to the request amount" eg:"false"`
	Options  []string            `json:"options" dc:"Options for select fields" eg:"[\"Travel\",\"Material\"]"`
	Columns  []*UpdateColumnItem `json:"columns" dc:"Column definitions for detail fields" eg:"[]"`
}

// UpdateReq defines the request for updating an approval flow.
type UpdateReq struct {
	g.Meta      `path:"/approval/flow/{id}" method:"put" tags:"ApprovalFlows" summary:"Update approval flow" dc:"Update the specified approval flow. Fields are optional; omitted fields keep their current values. When Nodes is provided the whole approver node list is replaced in order. When Fields is provided the whole form field list is replaced; changes only affect approval requests submitted afterwards." permission:"oa:flow:edit"`
	Id          int64              `json:"id" v:"required" dc:"Approval flow ID" eg:"1"`
	FlowType    *int               `json:"flowType" dc:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement" eg:"1"`
	FlowName    *string            `json:"flowName" v:"length:1,128#gf.gvalid.rule.length" dc:"Flow name, unique per tenant" eg:"Standard fund approval (update)"`
	Description *string            `json:"description" dc:"Flow description" eg:"Updated description"`
	Status      *int               `json:"status" dc:"Flow status: 1=enabled, 0=disabled" eg:"1"`
	Nodes       []*UpdateNodeItem  `json:"nodes" v:"max-length:20#gf.gvalid.rule.max-length" dc:"Optional ordered approver nodes; when provided the whole node list is replaced" eg:"[{}]"`
	Fields      []*UpdateFieldItem `json:"fields" v:"max-length:30#gf.gvalid.rule.max-length" dc:"Optional form field definitions; when provided the whole field list is replaced; empty array restores the built-in standard form" eg:"[]"`
}

// UpdateRes Approval flow update response
type UpdateRes struct{}
