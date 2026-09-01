// equipment_impl.go implements tenant-scoped equipment CRUD with unique code
// checks, date-only parsing, referenced-equipment delete protection, and
// batched creator name assembly. Tenant filters are applied at the database
// query stage.

package equipment

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-equipment-manage/backend/internal/dao"
	"lina-plugin-linapro-equipment-manage/backend/internal/model/do"
)

// dateFormat is the date-only wire format of purchase_date.
const dateFormat = "2006-01-02"

// List queries equipment with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	columns := dao.Equipment.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "")
	if name := strings.TrimSpace(in.Name); name != "" {
		m = m.WhereOrLike(columns.EquipmentName, "%"+name+"%").
			WhereOrLike(columns.EquipmentCode, "%"+name+"%")
	}
	if in.Type > 0 {
		m = m.Where(columns.EquipmentType, in.Type)
	}
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*EquipmentEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	userNames, err := s.resolveCreatorNameMap(ctx, list)
	if err != nil {
		return nil, err
	}

	items := make([]*ListItem, 0, len(list))
	for _, item := range list {
		items = append(items, &ListItem{
			EquipmentEntity: item,
			CreatedByName:   userNames[item.CreatedBy],
		})
	}
	return &ListOutput{List: items, Total: total}, nil
}

// GetById retrieves one equipment record by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*ListItem, error) {
	columns := dao.Equipment.Columns()

	var record *EquipmentEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		Where(columns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeEquipmentNotFound)
	}

	item := &ListItem{EquipmentEntity: record}
	if record.CreatedBy > 0 {
		names, err := s.resolveCreatorNameMap(ctx, []*EquipmentEntity{record})
		if err != nil {
			return nil, err
		}
		item.CreatedByName = names[record.CreatedBy]
	}
	return item, nil
}

// Create creates a new equipment record.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !EquipmentType(in.EquipmentType) {
		return 0, bizerr.NewCode(CodeEquipmentTypeInvalid)
	}
	if !EquipmentStatus(in.Status) {
		return 0, bizerr.NewCode(CodeEquipmentStatusInvalid)
	}
	purchaseDate, err := parseDate(in.PurchaseDate)
	if err != nil {
		return 0, err
	}
	if err = s.ensureCodeFree(ctx, in.EquipmentCode, 0); err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// GoFrame auto-fills created_at and updated_at.
	return dao.Equipment.Ctx(ctx).Data(do.Equipment{
		TenantId:      tenantID,
		EquipmentCode: strings.TrimSpace(in.EquipmentCode),
		EquipmentName: in.EquipmentName,
		EquipmentType: in.EquipmentType,
		BrandModel:    in.BrandModel,
		PurchaseDate:  purchaseDate,
		PurchasePrice: in.PurchasePrice,
		Location:      in.Location,
		Owner:         in.Owner,
		Status:        in.Status,
		Remark:        in.Remark,
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
	}).InsertAndGetId()
}

// Update updates equipment information.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	columns := dao.Equipment.Columns()

	existing := &EquipmentEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		Where(columns.Id, in.Id).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 {
		return bizerr.NewCode(CodeEquipmentNotFound)
	}

	if in.EquipmentType != nil && !EquipmentType(*in.EquipmentType) {
		return bizerr.NewCode(CodeEquipmentTypeInvalid)
	}
	if in.Status != nil && !EquipmentStatus(*in.Status) {
		return bizerr.NewCode(CodeEquipmentStatusInvalid)
	}
	var purchaseDate *time.Time
	if in.PurchaseDate != nil {
		parsed, parseErr := parseDate(*in.PurchaseDate)
		if parseErr != nil {
			return parseErr
		}
		purchaseDate = parsed
	}
	if in.EquipmentCode != nil {
		if err = s.ensureCodeFree(ctx, *in.EquipmentCode, in.Id); err != nil {
			return err
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Equipment{UpdatedBy: updatedBy}
	if in.EquipmentCode != nil {
		data.EquipmentCode = strings.TrimSpace(*in.EquipmentCode)
	}
	if in.EquipmentName != nil {
		data.EquipmentName = *in.EquipmentName
	}
	if in.EquipmentType != nil {
		data.EquipmentType = *in.EquipmentType
	}
	if in.BrandModel != nil {
		data.BrandModel = *in.BrandModel
	}
	if purchaseDate != nil {
		data.PurchaseDate = *purchaseDate
	}
	if in.PurchasePrice != nil {
		data.PurchasePrice = *in.PurchasePrice
	}
	if in.Location != nil {
		data.Location = *in.Location
	}
	if in.Owner != nil {
		data.Owner = *in.Owner
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	if in.Remark != nil {
		data.Remark = *in.Remark
	}

	_, err = dao.Equipment.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(columns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// Delete soft-deletes equipment by IDs with referenced-equipment protection.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeEquipmentDeleteRequired)
	}

	if err := s.ensureNoneReferenced(ctx, idList); err != nil {
		return err
	}

	columns := dao.Equipment.Columns()
	_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		WhereIn(columns.Id, idList).
		Delete()
	return err
}

// Options returns bounded equipment candidates for select controls.
func (s *serviceImpl) Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error) {
	columns := dao.Equipment.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "")
	if keyword := strings.TrimSpace(in.Keyword); keyword != "" {
		m = m.WhereOrLike(columns.EquipmentName, "%"+keyword+"%").
			WhereOrLike(columns.EquipmentCode, "%"+keyword+"%")
	}
	m = m.WhereNotIn(columns.Status, []int{EquipmentStatusScrapped})

	list := make([]*EquipmentEntity, 0, OptionsLimit)
	err := m.Fields(columns.Id, columns.EquipmentCode, columns.EquipmentName, columns.Status).
		Limit(OptionsLimit).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]*OptionItem, 0, len(list))
	for _, record := range list {
		items = append(items, &OptionItem{
			Id:            record.Id,
			EquipmentCode: record.EquipmentCode,
			EquipmentName: record.EquipmentName,
			Status:        record.Status,
		})
	}
	return items, nil
}

// ensureCodeFree verifies the tenant-unique equipment code excluding one ID.
func (s *serviceImpl) ensureCodeFree(ctx context.Context, code string, excludeID int64) error {
	columns := dao.Equipment.Columns()
	trimmed := strings.TrimSpace(code)
	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Equipment.Ctx(ctx), "").
		Where(columns.EquipmentCode, trimmed)
	if excludeID > 0 {
		m = m.WhereNot(columns.Id, excludeID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeEquipmentCodeExists)
	}
	return nil
}

// ensureNoneReferenced verifies no equipment in the batch is still referenced
// by maintenance records. One bounded query covers the whole batch.
func (s *serviceImpl) ensureNoneReferenced(ctx context.Context, equipmentIDs []int64) error {
	maintenanceColumns := dao.Maintenance.Columns()
	count, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Maintenance.Ctx(ctx), "").
		WhereIn(maintenanceColumns.EquipmentId, equipmentIDs).
		Limit(1).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeEquipmentReferenced)
	}
	return nil
}

// parseDate converts the date-only wire format into a stored date value.
// Empty input yields nil so the column keeps its default.
func parseDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation(dateFormat, trimmed, time.Local)
	if err != nil {
		return nil, bizerr.NewCode(CodeEquipmentDateInvalid)
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

// resolveCreatorNameMap resolves current-page creator display names through one
// user-domain batch call and leaves invisible or missing users blank.
func (s *serviceImpl) resolveCreatorNameMap(ctx context.Context, records []*EquipmentEntity) (map[int64]string, error) {
	names := make(map[int64]string)
	if s.userSvc == nil {
		return nil, gerror.New("linapro-equipment-manage requires host user capability")
	}
	ids := make([]usercap.UserID, 0, len(records))
	seen := make(map[int64]struct{}, len(records))
	for _, record := range records {
		if record == nil || record.CreatedBy <= 0 {
			continue
		}
		if _, ok := seen[record.CreatedBy]; ok {
			continue
		}
		seen[record.CreatedBy] = struct{}{}
		ids = append(ids, usercap.UserID(strconv.FormatInt(record.CreatedBy, 10)))
	}
	if len(ids) == 0 {
		return names, nil
	}
	result, err := s.userSvc.BatchGet(ctx, ids)
	if err != nil || result == nil {
		return names, err
	}
	for id, projection := range result.Items {
		storageID, ok := userDomainID(id)
		if !ok {
			continue
		}
		names[storageID] = userDisplayName(projection)
	}
	return names, nil
}

// userDomainID parses the host user-domain ID encoding used by plugin columns.
func userDomainID(id usercap.UserID) (int64, bool) {
	storageID, err := strconv.ParseInt(strings.TrimSpace(string(id)), 10, 64)
	return storageID, err == nil && storageID > 0
}

// userDisplayName chooses the stable user display field from the projection.
func userDisplayName(user *usercap.UserInfo) string {
	if user == nil {
		return ""
	}
	if user.Username != "" {
		return user.Username
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.Label != "" {
		return user.Label
	}
	return string(user.ID)
}
