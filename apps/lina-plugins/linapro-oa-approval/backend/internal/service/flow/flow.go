// Package flow implements tenant-scoped approval-flow configuration services
// for the linapro-oa-approval source plugin. It owns the
// plugin_linapro_oa_approval_flow and plugin_linapro_oa_approval_flow_node
// table access and consumes host capability seams for business context,
// tenant filtering, and user projections.
package flow

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Flow type values (matching the plugin_oa_approval_flow_type dictionary).
const (
	FlowTypeFundRequest = 1 // Fund request
	FlowTypeWriteOff    = 2 // Write-off
	FlowTypeReimburse   = 3 // Expense reimbursement
)

// Flow status values (matching the flow status column semantics).
const (
	FlowStatusDisabled = 0 // Disabled
	FlowStatusEnabled  = 1 // Enabled
)

// MaxFlowNodes bounds the ordered approver node list of one flow.
const MaxFlowNodes = 20

// FlowType validates one flow type value and reports whether it is known to
// the plugin.
func FlowType(value int) bool {
	switch value {
	case FlowTypeFundRequest, FlowTypeWriteOff, FlowTypeReimburse:
		return true
	default:
		return false
	}
}

// FlowStatus validates one flow status value.
func FlowStatus(value int) bool {
	switch value {
	case FlowStatusDisabled, FlowStatusEnabled:
		return true
	default:
		return false
	}
}

// NodeInput defines one ordered approver node submitted by callers.
type NodeInput struct {
	ApproverId int64 // Approver user ID
}

// Service defines tenant-scoped approval-flow configuration for the plugin.
type Service interface {
	// List returns approval flows visible to ctx's tenant using the supplied
	// page and name/type/status filters. Each item carries the flow node count
	// resolved with one batched node query for the current page. Status
	// filtering is skipped when Status is nil.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible flow by primary key with its ordered
	// approver nodes, or CodeFlowNotFound when the row is outside scope or
	// absent.
	GetById(ctx context.Context, id int64) (*FlowDetail, error)
	// Create creates one approval flow with the ordered node list. The name
	// must be unique within the tenant; violations return CodeFlowNameExists.
	// Every approver must resolve through the host user capability;
	// unresolvable approvers return CodeApproverInvalid.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the flow fields and, when Nodes is provided, replaces the
	// whole ordered node list in one transaction. Replaced nodes only affect
	// approval requests submitted afterwards.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes approval flows by IDs together with their nodes.
	// Submitted requests keep their frozen snapshots. Empty IDs return a
	// business error.
	Delete(ctx context.Context, ids []int64) error
	// UserOptions returns bounded approver candidates resolved through the
	// host user capability with the minimal id and name projection. The result
	// never exceeds MaxUserOptions items.
	UserOptions(ctx context.Context, in UserOptionsInput) ([]*UserOptionItem, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service // Business context bridge
	tenantSvc tenantcap.Service // Tenant capability bridge
	userSvc   usercap.Service   // User domain projection capability
}

// New creates and returns a new Service instance.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
	userSvc usercap.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc: bizCtxSvc,
		tenantSvc: tenantSvc,
		userSvc:   userSvc,
	}
}

// ListInput defines input for List function.
type ListInput struct {
	PageNum  int    // Page number, starting from 1
	PageSize int    // Page size
	FlowName string // Flow name, supports fuzzy search
	FlowType int    // Flow type: 1=fund request 2=write-off 3=reimburse; 0 means no filter
	Status   *int   // Flow status: 1=enabled 0=disabled; nil means no filter
}

// FlowDetail defines one flow with its ordered approver nodes.
type FlowDetail struct {
	*FlowEntity                        // Flow entity
	Nodes       []*NodeItem            // Ordered approver nodes
	Fields      []approval.FieldConfig // Configurable form field definitions
}

// NodeItem defines one ordered approver node projection.
type NodeItem struct {
	Order        int64  // Approval order starting from 1
	ApproverId   int64  // Approver user ID
	ApproverName string // Approver display name
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*FlowItem // Flow list items
	Total int         // Total count
}

// FlowItem defines one list projection item.
type FlowItem struct {
	*FlowEntity
	NodeCount     int // Number of approval nodes
	FieldCount    int // Number of configured form fields; 0 means the built-in standard form
	CreatedByName string
}

// CreateInput defines input for Create function.
type CreateInput struct {
	FlowType    int                    // Flow type: 1=fund request 2=write-off 3=reimburse
	FlowName    string                 // Flow name, unique per tenant
	Description string                 // Flow description
	Status      int                    // Flow status: 1=enabled 0=disabled
	Nodes       []NodeInput            // Ordered approver nodes
	Fields      []approval.FieldConfig // Optional configurable form fields
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id          int64                   // Approval flow ID
	FlowType    *int                    // Flow type; nil keeps current value
	FlowName    *string                 // Flow name; nil keeps current value
	Description *string                 // Flow description; nil keeps current value
	Status      *int                    // Flow status; nil keeps current value
	Nodes       []NodeInput             // When non-nil the whole node list is replaced
	Fields      *[]approval.FieldConfig // When non-nil the whole form field list is replaced
}

// UserOptionsInput defines input for UserOptions function.
type UserOptionsInput struct {
	Keyword string // Keyword, matches username or nickname fuzzily
}

// UserOptionItem defines one minimal approver candidate projection.
type UserOptionItem struct {
	Id   int64  // User ID
	Name string // User display name
}
