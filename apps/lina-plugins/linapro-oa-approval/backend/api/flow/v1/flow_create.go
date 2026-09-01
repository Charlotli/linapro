// This file declares the create-approval-flow request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Flow Create API

// CreateNodeItem defines one ordered approver node in a create request.
type CreateNodeItem struct {
	ApproverId int64 `json:"approverId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Approver user ID" eg:"2"`
}

// CreateColumnItem defines one detail sub-table column in a create request.
type CreateColumnItem struct {
	Key   string `json:"key" v:"required#gf.gvalid.rule.required" dc:"Storage key" eg:"amount"`
	Label string `json:"label" v:"required#gf.gvalid.rule.required" dc:"Display name" eg:"Amount"`
	Type  string `json:"type" v:"required|in:text,number,date#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Column type: text/number/date" eg:"number"`
}

// CreateFieldItem defines one configurable form field in a create request.
type CreateFieldItem struct {
	Key      string              `json:"key" v:"required#gf.gvalid.rule.required" dc:"Storage key, identifier-safe and unique within the flow" eg:"expenseType"`
	Label    string              `json:"label" v:"required#gf.gvalid.rule.required" dc:"Display name" eg:"Expense type"`
	Type     string              `json:"type" v:"required|in:text,textarea,number,date,select,attachment,detail#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Field type: text/textarea/number/date/select/attachment/detail" eg:"select"`
	Required bool                `json:"required" dc:"Whether the field must be filled" eg:"true"`
	AsAmount bool                `json:"asAmount" dc:"Whether a number field projects to the request amount" eg:"false"`
	Options  []string            `json:"options" dc:"Options for select fields" eg:"[\"Travel\",\"Material\"]"`
	Columns  []*CreateColumnItem `json:"columns" dc:"Column definitions for detail fields" eg:"[]"`
}

// CreateReq defines the request for creating an approval flow.
type CreateReq struct {
	g.Meta      `path:"/approval/flow" method:"post" tags:"ApprovalFlows" summary:"Create approval flow" dc:"Create an approval flow owned by the current tenant with an ordered approver node list and optional configurable form fields. The flow name must be unique within the tenant and the node list must contain 1 to 20 approvers with consecutive orders." permission:"oa:flow:add"`
	FlowType    int                `json:"flowType" v:"required|in:1,2,3#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement" eg:"1"`
	FlowName    string             `json:"flowName" v:"required|length:1,128#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Flow name, unique per tenant" eg:"Standard fund approval"`
	Description string             `json:"description" dc:"Flow description" eg:"Two-level fund approval"`
	Status      *int               `json:"status" d:"1" dc:"Flow status: 1=enabled, 0=disabled" eg:"1"`
	Nodes       []*CreateNodeItem  `json:"nodes" v:"required|min-length:1|max-length:20#gf.gvalid.rule.required|gf.gvalid.rule.min-length|gf.gvalid.rule.max-length" dc:"Ordered approver nodes; the list order defines the approval order" eg:"[{}]"`
	Fields      []*CreateFieldItem `json:"fields" v:"max-length:30#gf.gvalid.rule.max-length" dc:"Optional configurable form field definitions; empty uses the built-in standard form" eg:"[]"`
}

// CreateRes Approval flow create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Approval flow ID" eg:"1"`
}
