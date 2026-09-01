// This file defines shared approval-flow response DTOs for the linapro-oa-approval API.
package v1

// FlowItem exposes approval-flow fields visible to callers.
type FlowItem struct {
	Id            int64  `json:"id" dc:"Approval flow ID" eg:"1"`
	FlowType      int    `json:"flowType" dc:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement" eg:"1"`
	FlowName      string `json:"flowName" dc:"Flow name, unique per tenant" eg:"Standard fund approval"`
	Description   string `json:"description" dc:"Flow description" eg:"Two-level fund approval"`
	Status        int    `json:"status" dc:"Flow status: 1=enabled, 0=disabled" eg:"1"`
	NodeCount     int    `json:"nodeCount" dc:"Number of approval nodes in the flow" eg:"2"`
	FieldCount    int    `json:"fieldCount" dc:"Number of configured form fields; 0 means the built-in standard form" eg:"3"`
	CreatedBy     int64  `json:"createdBy" dc:"Creator user ID" eg:"1"`
	CreatedByName string `json:"createdByName" dc:"Creator username" eg:"admin"`
	CreatedAt     *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt     *int64 `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}

// FlowNodeItem exposes one ordered approver node of a flow.
type FlowNodeItem struct {
	Order        int64  `json:"order" dc:"Approval order starting from 1" eg:"1"`
	ApproverId   int64  `json:"approverId" dc:"Approver user ID" eg:"2"`
	ApproverName string `json:"approverName" dc:"Approver display name" eg:"leader"`
}
