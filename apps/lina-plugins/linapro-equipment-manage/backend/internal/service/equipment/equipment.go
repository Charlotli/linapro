// Package equipment implements tenant-scoped equipment registry services for
// the linapro-equipment-manage source plugin. It owns the
// plugin_linapro_equipment_manage_equipment table access and consumes host
// capability seams for business context, tenant filtering, and user
// projections.
package equipment

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
)

// Equipment type values (matching the plugin_equipment_type dictionary).
const (
	EquipmentTypeOffice     = 1 // Office equipment
	EquipmentTypeIT         = 2 // IT equipment
	EquipmentTypeProduction = 3 // Production equipment
	EquipmentTypeOther      = 4 // Other equipment
)

// Equipment status values (matching the plugin_equipment_status dictionary).
const (
	EquipmentStatusInUse       = 1 // In use
	EquipmentStatusIdle        = 2 // Idle
	EquipmentStatusUnderRepair = 3 // Under repair
	EquipmentStatusScrapped    = 4 // Scrapped
)

// EquipmentType validates one equipment type value.
func EquipmentType(value int) bool {
	switch value {
	case EquipmentTypeOffice, EquipmentTypeIT, EquipmentTypeProduction, EquipmentTypeOther:
		return true
	default:
		return false
	}
}

// EquipmentStatus validates one equipment status value.
func EquipmentStatus(value int) bool {
	switch value {
	case EquipmentStatusInUse, EquipmentStatusIdle, EquipmentStatusUnderRepair, EquipmentStatusScrapped:
		return true
	default:
		return false
	}
}

// Service defines tenant-scoped equipment registry services for the plugin.
type Service interface {
	// List returns equipment visible to ctx's tenant using the supplied page
	// and name/type/status filters. Status filtering is skipped when Status is
	// nil. It returns DAO or tenant-scope errors.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible equipment record by primary key with
	// creator display metadata, or CodeEquipmentNotFound when the row is
	// outside scope or absent.
	GetById(ctx context.Context, id int64) (*ListItem, error)
	// Create creates one equipment record owned by the current tenant and user
	// from ctx. The equipment code must be unique within the tenant.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the equipment fields. It only updates the current
	// tenant's row and returns CodeEquipmentNotFound for missing scope.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes equipment by IDs. Equipment still referenced by
	// maintenance records is rejected as a whole; empty IDs return a business
	// error.
	Delete(ctx context.Context, ids []int64) error
	// Options returns bounded minimal equipment candidates for select
	// controls, filtered by keyword across name and code. Scrapped equipment
	// is excluded. The result never exceeds OptionsLimit items.
	Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error)
}

// OptionsLimit bounds the equipment candidate list.
const OptionsLimit = 200

// OptionsInput defines input for Options function.
type OptionsInput struct {
	Keyword string // Keyword, matches equipment name or code fuzzily
}

// OptionItem defines one minimal equipment option projection.
type OptionItem struct {
	Id            int64  // Equipment ID
	EquipmentCode string // Equipment code
	EquipmentName string // Equipment name
	Status        int    // Equipment status
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
	Name     string // Name or code, supports fuzzy search
	Type     int    // Equipment type; 0 means no filter
	Status   *int   // Equipment status; nil means no filter
}

// ListItem defines a single list item.
type ListItem struct {
	*EquipmentEntity
	CreatedByName string // Creator username
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*ListItem // List items
	Total int         // Total count
}

// CreateInput defines input for Create function.
type CreateInput struct {
	EquipmentCode string  // Equipment code, unique per tenant
	EquipmentName string  // Equipment name
	EquipmentType int     // Equipment type
	BrandModel    string  // Brand and model
	PurchaseDate  string  // Purchase date, YYYY-MM-DD; empty means unset
	PurchasePrice float64 // Purchase price
	Location      string  // Storage location
	Owner         string  // Responsible person
	Status        int     // Equipment status
	Remark        string  // Remark
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id            int64    // Equipment ID
	EquipmentCode *string  // Equipment code; nil keeps current value
	EquipmentName *string  // Equipment name; nil keeps current value
	EquipmentType *int     // Equipment type; nil keeps current value
	BrandModel    *string  // Brand and model; nil keeps current value
	PurchaseDate  *string  // Purchase date, YYYY-MM-DD; nil keeps current value
	PurchasePrice *float64 // Purchase price; nil keeps current value
	Location      *string  // Storage location; nil keeps current value
	Owner         *string  // Responsible person; nil keeps current value
	Status        *int     // Equipment status; nil keeps current value
	Remark        *string  // Remark; nil keeps current value
}
