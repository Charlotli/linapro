// maintenance_impl.go implements tenant-scoped maintenance record CRUD with
// equipment existence validation, date-only parsing, and batched equipment
// name assembly. Tenant filters are applied at the database query stage.

package maintenance

import (
	"context"
	"strings"
	"time"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-plugin-linapro-equipment-manage/backend/internal/dao"
	"lina-plugin-linapro-equipment-manage/backend/internal/model/do"
	"lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// dateFormat is the date-only wire format of maint_date.
const dateFormat = "2006-01-02"

// List queries maintenance records with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	columns := dao.Maintenance.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Maintenance.Ctx(ctx), "")
	if in.EquipmentId > 0 {
		m = m.Where(columns.EquipmentId, in.EquipmentId)
	}
	if in.MaintType > 0 {
		m = m.Where(columns.MaintType, in.MaintType)
	}
	if start := strings.TrimSpace(in.DateStart); start != "" {
		m = m.WhereGTE(columns.MaintDate, start)
	}
	if end := strings.TrimSpace(in.DateEnd); end != "" {
		m = m.WhereLTE(columns.MaintDate, end+" 23:59:59")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*MaintenanceEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	equipmentNames, err := s.resolveEquipmentNameMap(ctx, list)
	if err != nil {
		return nil, err
	}

	items := make([]*MaintenanceItem, 0, len(list))
	for _, record := range list {
		items = append(items, &MaintenanceItem{
			MaintenanceEntity: record,
			EquipmentName:     equipmentNames[record.EquipmentId],
		})
	}
	return &ListOutput{List: items, Total: total}, nil
}

// GetById retrieves one maintenance record by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*MaintenanceEntity, error) {
	columns := dao.Maintenance.Columns()

	var record *MaintenanceEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Maintenance.Ctx(ctx), "").
		Where(columns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeMaintenanceNotFound)
	}
	return record, nil
}

// Create creates one maintenance record.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !MaintType(in.MaintType) {
		return 0, bizerr.NewCode(CodeMaintTypeInvalid)
	}
	maintDate, err := parseDate(in.MaintDate)
	if err != nil {
		return 0, err
	}
	if err = s.ensureEquipmentExists(ctx, in.EquipmentId); err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// GoFrame auto-fills created_at and updated_at.
	return dao.Maintenance.Ctx(ctx).Data(do.Maintenance{
		TenantId:    tenantID,
		EquipmentId: in.EquipmentId,
		MaintType:   in.MaintType,
		MaintDate:   maintDate,
		Maintainer:  in.Maintainer,
		Cost:        in.Cost,
		Content:     in.Content,
		Result:      in.Result,
		Remark:      in.Remark,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}).InsertAndGetId()
}

// Update updates maintenance record information.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	if in.MaintType != nil && !MaintType(*in.MaintType) {
		return bizerr.NewCode(CodeMaintTypeInvalid)
	}
	columns := dao.Maintenance.Columns()

	existing := &MaintenanceEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Maintenance.Ctx(ctx), "").
		Where(columns.Id, in.Id).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 {
		return bizerr.NewCode(CodeMaintenanceNotFound)
	}

	var maintDate *time.Time
	if in.MaintDate != nil {
		parsed, parseErr := parseDate(*in.MaintDate)
		if parseErr != nil {
			return parseErr
		}
		maintDate = parsed
	}
	if in.EquipmentId != nil {
		if err = s.ensureEquipmentExists(ctx, *in.EquipmentId); err != nil {
			return err
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Maintenance{UpdatedBy: updatedBy}
	if in.EquipmentId != nil {
		data.EquipmentId = *in.EquipmentId
	}
	if in.MaintType != nil {
		data.MaintType = *in.MaintType
	}
	if maintDate != nil {
		data.MaintDate = *maintDate
	}
	if in.Maintainer != nil {
		data.Maintainer = *in.Maintainer
	}
	if in.Cost != nil {
		data.Cost = *in.Cost
	}
	if in.Content != nil {
		data.Content = *in.Content
	}
	if in.Result != nil {
		data.Result = *in.Result
	}
	if in.Remark != nil {
		data.Remark = *in.Remark
	}

	_, err = dao.Maintenance.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(columns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// Delete soft-deletes maintenance records by IDs.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeMaintenanceDeleteRequired)
	}

	columns := dao.Maintenance.Columns()
	_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Maintenance.Ctx(ctx), "").
		WhereIn(columns.Id, idList).
		Delete()
	return err
}

// ensureEquipmentExists verifies the referenced equipment exists in the
// current tenant. Out-of-scope equipment is reported as not found.
func (s *serviceImpl) ensureEquipmentExists(ctx context.Context, equipmentID int64) error {
	equipmentColumns := dao.Equipment.Columns()
	count, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		Where(equipmentColumns.Id, equipmentID).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return bizerr.NewCode(equipment.CodeEquipmentNotFound)
	}
	return nil
}

// resolveEquipmentNameMap resolves equipment display names for the current
// page with one batch query, keeping list assembly free of per-row lookups.
func (s *serviceImpl) resolveEquipmentNameMap(ctx context.Context, records []*MaintenanceEntity) (map[int64]string, error) {
	names := make(map[int64]string)
	equipmentIDs := make([]int64, 0, len(records))
	seen := make(map[int64]struct{}, len(records))
	for _, record := range records {
		if record == nil || record.EquipmentId <= 0 {
			continue
		}
		if _, ok := seen[record.EquipmentId]; ok {
			continue
		}
		seen[record.EquipmentId] = struct{}{}
		equipmentIDs = append(equipmentIDs, record.EquipmentId)
	}
	if len(equipmentIDs) == 0 {
		return names, nil
	}

	equipmentColumns := dao.Equipment.Columns()
	equipments := make([]*equipment.EquipmentEntity, 0, len(equipmentIDs))
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		Fields(equipmentColumns.Id, equipmentColumns.EquipmentName).
		WhereIn(equipmentColumns.Id, equipmentIDs).
		Scan(&equipments)
	if err != nil {
		return nil, err
	}
	for _, item := range equipments {
		if item == nil {
			continue
		}
		names[item.Id] = item.EquipmentName
	}
	return names, nil
}

// parseDate converts the date-only wire format into a stored date value.
// Empty input defaults to today.
func parseDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		today := time.Now().Format(dateFormat)
		trimmed = today
	}
	parsed, err := time.ParseInLocation(dateFormat, trimmed, time.Local)
	if err != nil {
		return nil, bizerr.NewCode(CodeMaintenanceDateInvalid)
	}
	return &parsed, nil
}

// normalizePositiveInt64IDs drops non-positive identifiers from a batch list.
func normalizePositiveInt64IDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			result = append(result, id)
		}
	}
	return result
}

// pluginTableFilter returns the tenant table-filter slice from the injected tenant service.
func (s *serviceImpl) pluginTableFilter() tenantcap.FilterService {
	if s == nil || s.tenantSvc == nil {
		return nil
	}
	return s.tenantSvc.Filter()
}

// tenantFilterContext returns current tenant metadata for write ownership fields.
func (s *serviceImpl) tenantFilterContext(ctx context.Context) tenantcap.TenantFilterContext {
	if filter := s.pluginTableFilter(); filter != nil {
		return filter.Context(ctx)
	}
	return tenantcap.TenantFilterContext{}
}
