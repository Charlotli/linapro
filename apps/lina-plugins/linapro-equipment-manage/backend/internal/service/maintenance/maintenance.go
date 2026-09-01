// Package maintenance implements tenant-scoped equipment maintenance record
// services for the linapro-equipment-manage source plugin. It owns the
// plugin_linapro_equipment_manage_maintenance table access, validates the
// referenced equipment on writes, and assembles equipment display names with
// one batch query per page.
package maintenance

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
)

// Maintenance type values (matching the plugin_equipment_maint_type dictionary).
const (
	MaintTypeMaintenance = 1 // Maintenance
	MaintTypeRepair      = 2 // Repair
	MaintTypeInspection  = 3 // Inspection
	MaintTypeOther       = 4 // Other
)

// MaintType validates one maintenance type value.
func MaintType(value int) bool {
	switch value {
	case MaintTypeMaintenance, MaintTypeRepair, MaintTypeInspection, MaintTypeOther:
		return true
	default:
		return false
	}
}

// Service defines tenant-scoped maintenance record services for the plugin.
type Service interface {
	// List returns maintenance records visible to ctx's tenant using the
	// supplied page and equipment/type/date-range filters. Equipment names are
	// resolved with one batch query for the current page.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible maintenance record by primary key, or
	// CodeMaintenanceNotFound when the row is outside scope or absent.
	GetById(ctx context.Context, id int64) (*MaintenanceEntity, error)
	// Create creates one maintenance record owned by the current tenant and
	// user from ctx. The equipment must exist in the same tenant.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the maintenance fields. Equipment existence is re-checked
	// when EquipmentId is updated.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes maintenance records by IDs. Empty IDs return a
	// business error.
	Delete(ctx context.Context, ids []int64) error
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service // Business context bridge
	tenantSvc tenantcap.Service // Tenant capability bridge
}

// New creates and returns a new Service instance.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc: bizCtxSvc,
		tenantSvc: tenantSvc,
	}
}

// ListInput defines input for List function.
type ListInput struct {
	PageNum     int    // Page number, starting from 1
	PageSize    int    // Page size
	EquipmentId int64  // Equipment ID; 0 means no filter
	MaintType   int    // Maintenance type; 0 means no filter
	DateStart   string // Maintenance date range start, YYYY-MM-DD, inclusive; empty means unbounded
	DateEnd     string // Maintenance date range end, YYYY-MM-DD, inclusive; empty means unbounded
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*MaintenanceItem // List items
	Total int                // Total count
}

// MaintenanceItem defines one list projection item.
type MaintenanceItem struct {
	*MaintenanceEntity
	EquipmentName string // Equipment name resolved in batch
}

// CreateInput defines input for Create function.
type CreateInput struct {
	EquipmentId int64   // Maintained equipment ID
	MaintType   int     // Maintenance type
	MaintDate   string  // Maintenance date, YYYY-MM-DD; empty defaults to today
	Maintainer  string  // Maintainer name
	Cost        float64 // Maintenance cost
	Content     string  // Maintenance content
	Result      string  // Maintenance result
	Remark      string  // Remark
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id          int64    // Maintenance record ID
	EquipmentId *int64   // Maintained equipment ID; nil keeps current value
	MaintType   *int     // Maintenance type; nil keeps current value
	MaintDate   *string  // Maintenance date, YYYY-MM-DD; nil keeps current value
	Maintainer  *string  // Maintainer name; nil keeps current value
	Cost        *float64 // Maintenance cost; nil keeps current value
	Content     *string  // Maintenance content; nil keeps current value
	Result      *string  // Maintenance result; nil keeps current value
	Remark      *string  // Remark; nil keeps current value
}
