// This file declares the get-approval-request request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Get API

// GetReq defines the request for retrieving approval-request details.
type GetReq struct {
	g.Meta `path:"/approval/request/{id}" method:"get" tags:"ApprovalRequests" summary:"Get approval request details" dc:"Get one approval request with the frozen node snapshot, the frozen form field definitions, submitted form values, and the full approval timeline. Only the applicant and snapshot approvers can view the record; other callers receive a not-found error without existence leaking." permission:"oa:request:query"`
	Id     int64 `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
}

// FormFieldColumnItem exposes one column of a detail field definition.
type FormFieldColumnItem struct {
	Key   string `json:"key" dc:"Storage key of the column" eg:"amount"`
	Label string `json:"label" dc:"Display name of the column" eg:"Amount"`
	Type  string `json:"type" dc:"Column type: text/number/date" eg:"number"`
}

// FormFieldItem exposes one frozen form field definition.
type FormFieldItem struct {
	Key      string                 `json:"key" dc:"Storage key of the field" eg:"expenseType"`
	Label    string                 `json:"label" dc:"Display name of the field" eg:"Expense type"`
	Type     string                 `json:"type" dc:"Field type: text/textarea/number/date/select/attachment/detail" eg:"select"`
	Required bool                   `json:"required" dc:"Whether the field is required" eg:"true"`
	AsAmount bool                   `json:"asAmount" dc:"Whether a number field projects to the request amount" eg:"false"`
	Options  []string               `json:"options" dc:"Options for select fields" eg:"[]"`
	Columns  []*FormFieldColumnItem `json:"columns" dc:"Column definitions for detail fields" eg:"[]"`
}

// GetRes Approval request detail response
type GetRes struct {
	RequestItem
	Nodes      []*RequestNodeItem   `json:"nodes" dc:"Frozen approval node snapshot" eg:"[]"`
	Records    []*RequestRecordItem `json:"records" dc:"Approval timeline ordered by time" eg:"[]"`
	FormFields []*FormFieldItem     `json:"formFields" dc:"Frozen form field definitions at submit time" eg:"[]"`
	Form       map[string]any       `json:"form" dc:"Submitted dynamic form values keyed by field key" eg:"{}"`
}
