// This file defines shared approval-request response DTOs for the linapro-oa-approval API.
package v1

// RequestItem exposes approval-request fields visible to callers.
type RequestItem struct {
	Id                  int64   `json:"id" dc:"Approval request ID" eg:"1"`
	FlowId              int64   `json:"flowId" dc:"Originating approval flow ID" eg:"1"`
	FlowType            int     `json:"flowType" dc:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement" eg:"1"`
	Title               string  `json:"title" dc:"Request title" eg:"Youth camp fund request"`
	Amount              float64 `json:"amount" dc:"Requested amount with two decimal places" eg:"3500.00"`
	Content             string  `json:"content" dc:"Applicant statement content" eg:"Venue and material budget"`
	Attachments         string  `json:"attachments" dc:"Attachment URL list as a JSON string array" eg:"[\"https://example.com/a.xlsx\"]"`
	Status              int     `json:"status" dc:"Request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn" eg:"1"`
	CurrentNodeOrder    int     `json:"currentNodeOrder" dc:"Pending node order while the request is pending" eg:"1"`
	CurrentApproverId   int64   `json:"currentApproverId" dc:"Pending approver user ID while the request is pending" eg:"2"`
	CurrentApproverName string  `json:"currentApproverName" dc:"Pending approver display name resolved in batch" eg:"leader"`
	ApplicantId         int64   `json:"applicantId" dc:"Applicant user ID" eg:"1"`
	ApplicantName       string  `json:"applicantName" dc:"Applicant display name resolved in batch" eg:"admin"`
	CreatedAt           *int64  `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt           *int64  `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}

// RequestNodeItem exposes one frozen approval node from the request snapshot.
type RequestNodeItem struct {
	Order        int64  `json:"order" dc:"Approval order starting from 1" eg:"1"`
	ApproverId   int64  `json:"approverId" dc:"Approver user ID" eg:"2"`
	ApproverName string `json:"approverName" dc:"Approver display name" eg:"leader"`
}

// RequestRecordItem exposes one approval timeline record.
type RequestRecordItem struct {
	Id        int64  `json:"id" dc:"Approval record ID" eg:"1"`
	NodeOrder int    `json:"nodeOrder" dc:"Related node order; 0 means request-level actions" eg:"1"`
	Action    int    `json:"action" dc:"Action: 1=submit, 2=approve, 3=reject, 4=comment, 5=append approver, 6=withdraw" eg:"2"`
	ActorId   int64  `json:"actorId" dc:"Actor user ID" eg:"2"`
	ActorName string `json:"actorName" dc:"Actor display name resolved in batch" eg:"leader"`
	Comment   string `json:"comment" dc:"Reply or comment content attached to the action" eg:"Approved, please proceed"`
	CreatedAt *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756100000"`
}
