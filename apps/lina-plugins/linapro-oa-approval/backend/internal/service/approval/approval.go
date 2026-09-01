// Package approval implements the tenant-scoped approval-request state
// machine for the linapro-oa-approval source plugin. It owns the
// plugin_linapro_oa_approval_request and plugin_linapro_oa_approval_record
// table access, freezes flow nodes into per-request snapshots at submit time,
// enforces participant visibility, and fans progress notifications out through
// the host notify capability.
package approval

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/i18ncap"
	"lina-core/pkg/plugin/capability/notifycap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
)

// Request status values (matching the plugin_oa_approval_status dictionary).
const (
	requestStatusPending   = 1 // Pending
	requestStatusApproved  = 2 // Approved
	requestStatusRejected  = 3 // Rejected
	requestStatusWithdrawn = 4 // Withdrawn
)

// Timeline action values (matching the plugin_oa_approval_action dictionary).
const (
	recordActionSubmit  = 1 // Submit
	recordActionApprove = 2 // Approve
	recordActionReject  = 3 // Reject
	recordActionComment = 4 // Comment
	recordActionAppend  = 5 // Append approver
	recordActionWit     = 6 // Withdraw
)

// Scope values for the request list view.
const (
	scopeMine   = "mine"    // Requests submitted by the current user
	scopePening = "pending" // Requests waiting for the current user
)

// MaxRecords bounds one detail timeline read.
const MaxRecords = 200

// Service defines tenant-scoped approval-request actions for the plugin.
type Service interface {
	// List returns one bounded page of tenant-visible requests. Scope mine
	// filters by the current applicant; scope pending filters by the current
	// approver and the pending status. Other filters are optional and applied
	// on the database side.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one request with its frozen node snapshot and the full
	// timeline. Only the applicant and snapshot approvers can read the record;
	// other callers receive CodeRequestNotFound without existence leaking.
	GetById(ctx context.Context, id int64) (*RequestDetail, error)
	// Create submits one approval request: it validates the flow, freezes the
	// node snapshot, records the submit timeline entry, and notifies the first
	// approver.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Approve approves one pending request as the current node approver and
	// advances to the next node or finishes the request.
	Approve(ctx context.Context, in ActionInput) error
	// Reject rejects one pending request as the current node approver.
	Reject(ctx context.Context, in ActionInput) error
	// Comment appends one reply to the timeline without a status change. Only
	// participants can reply while the request is pending.
	Comment(ctx context.Context, in ActionInput) error
	// Appender inserts one new approver node right after the current node of a
	// pending request. Later node orders shift back and the current node stays.
	Appender(ctx context.Context, in AppenderInput) error
	// Withdraw withdraws one pending request as the applicant.
	Withdraw(ctx context.Context, id int64) error
	// Resubmit re-submits one rejected or withdrawn request as the applicant.
	// The frozen node snapshot is preserved, the request restarts from the
	// first node, and the first approver is notified. Non-applicants return
	// CodeResubmitNoAuthority; other states return CodeRequestNotPending.
	Resubmit(ctx context.Context, id int64) error
	// PendingCount returns the number of pending requests waiting for the
	// current user within the tenant, counted on the database side.
	PendingCount(ctx context.Context) (int64, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service // Business context bridge
	tenantSvc tenantcap.Service // Tenant capability bridge
	userSvc   usercap.Service   // User domain projection capability
	notifySvc notifycap.Service // Host notification capability
	i18nSvc   i18ncap.Service   // Host runtime i18n capability
	notifier  *approvalNotifier // Localized notification sender
}

// New creates and returns a new Service instance.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
	userSvc usercap.Service,
	notifySvc notifycap.Service,
	i18nSvc i18ncap.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc: bizCtxSvc,
		tenantSvc: tenantSvc,
		userSvc:   userSvc,
		notifySvc: notifySvc,
		i18nSvc:   i18nSvc,
		notifier:  &approvalNotifier{notifySvc: notifySvc, i18nSvc: i18nSvc},
	}
}

// ListInput defines input for List function.
type ListInput struct {
	PageNum  int    // Page number, starting from 1
	PageSize int    // Page size
	Scope    string // List scope: mine or pending
	Title    string // Title, supports fuzzy search
	FlowType int    // Flow type: 1=fund request 2=write-off 3=reimburse; 0 means no filter
	Status   *int   // Request status: 1=pending 2=approved 3=rejected 4=withdrawn; nil means no filter
}

// RequestDetail defines one request with snapshot nodes and the timeline.
type RequestDetail struct {
	*RequestEntity                     // Request entity
	ApplicantName       string         // Applicant display name
	CurrentApproverName string         // Pending approver display name
	Nodes               []*NodeItem    // Frozen node snapshot
	Records             []*RecordItem  // Approval timeline
	FormFields          []FieldConfig  // Frozen form field definitions
	FormValues          map[string]any // Submitted dynamic form values
}

// NodeItem defines one frozen snapshot node projection.
type NodeItem struct {
	Order        int64  // Approval order starting from 1
	ApproverId   int64  // Approver user ID
	ApproverName string // Approver display name
}

// RecordItem defines one timeline record projection.
type RecordItem struct {
	*RecordEntity
	ActorName string // Actor display name
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*RequestItem // List items
	Total int            // Total count
}

// RequestItem defines one list projection item.
type RequestItem struct {
	*RequestEntity
	ApplicantName       string // Applicant display name
	CurrentApproverName string // Pending approver display name
}

// CreateInput defines input for Create function.
type CreateInput struct {
	FlowId int64          // Originating approval flow ID
	Title  string         // Request title
	Form   map[string]any // Dynamic form values keyed by configured field key
}

// ActionInput defines input for approve, reject, and comment actions.
type ActionInput struct {
	Id      int64  // Approval request ID
	Comment string // Reply content stored on the timeline
}

// AppenderInput defines input for the append-approver action.
type AppenderInput struct {
	Id         int64 // Approval request ID
	ApproverId int64 // New approver user ID inserted after the current node
}
