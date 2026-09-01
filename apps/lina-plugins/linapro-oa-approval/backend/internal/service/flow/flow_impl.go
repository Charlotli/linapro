// flow_impl.go implements tenant-scoped approval-flow CRUD, ordered node
// replacement, bounded approver candidates, and creator/node name resolution
// for the linapro-oa-approval plugin. Tenant filters are applied at the
// database query stage and all read-side assembly stays batched per page.

package flow

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-plugin-linapro-oa-approval/backend/internal/service/approval"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/capmodel"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-oa-approval/backend/internal/dao"
	"lina-plugin-linapro-oa-approval/backend/internal/model/do"
)

// MaxUserOptions bounds the approver candidate list.
const MaxUserOptions = 200

// userSearchLimit bounds one creator or candidate keyword search.
const userSearchLimit = MaxUserOptions

// List queries approval flows with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	flowColumns := dao.Flow.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "")

	if name := strings.TrimSpace(in.FlowName); name != "" {
		m = m.WhereLike(flowColumns.FlowName, "%"+name+"%")
	}
	if in.FlowType > 0 {
		m = m.Where(flowColumns.FlowType, in.FlowType)
	}
	if in.Status != nil {
		m = m.Where(flowColumns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	flows := make([]*FlowEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(flowColumns.Id).
		Scan(&flows)
	if err != nil {
		return nil, err
	}

	nodeCounts, err := s.resolveNodeCountMap(ctx, flows)
	if err != nil {
		return nil, err
	}
	userNames, err := s.resolveUserNames(ctx, creatorUserIDs(flows))
	if err != nil {
		return nil, err
	}

	items := make([]*FlowItem, 0, len(flows))
	for _, flow := range flows {
		items = append(items, &FlowItem{
			FlowEntity:    flow,
			NodeCount:     nodeCounts[flow.Id],
			FieldCount:    countConfiguredFields(flow.FormFields),
			CreatedByName: userNames[flow.CreatedBy],
		})
	}

	return &ListOutput{
		List:  items,
		Total: total,
	}, nil
}

// GetById retrieves one approval flow with its ordered nodes.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*FlowDetail, error) {
	flowColumns := dao.Flow.Columns()

	var flow *FlowEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "").
		Where(flowColumns.Id, id).
		Scan(&flow)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, bizerr.NewCode(CodeFlowNotFound)
	}

	nodes, err := s.loadFlowNodes(ctx, id)
	if err != nil {
		return nil, err
	}
	fields, err := s.loadFlowFields(ctx, flow)
	if err != nil {
		return nil, err
	}
	return &FlowDetail{FlowEntity: flow, Nodes: nodes, Fields: fields}, nil
}

// loadFlowFields parses the stored flow form-field definitions; an empty
// storage value means the flow uses the built-in standard form.
func loadFlowFieldsFromString(raw string) ([]approval.FieldConfig, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return []approval.FieldConfig{}, nil
	}
	fields := make([]approval.FieldConfig, 0)
	if err := json.Unmarshal([]byte(trimmed), &fields); err != nil {
		return nil, gerror.Wrap(err, "parse flow form fields failed")
	}
	return fields, nil
}

// countConfiguredFields counts entries in the stored field-definition JSON
// without any extra query; zero means the built-in standard form.
func countConfiguredFields(raw string) int {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0
	}
	fields := make([]approval.FieldConfig, 0)
	if err := json.Unmarshal([]byte(trimmed), &fields); err != nil {
		return 0
	}
	return len(fields)
}

// loadFlowFields resolves the form-field definitions of one loaded flow.
func (s *serviceImpl) loadFlowFields(ctx context.Context, flow *FlowEntity) ([]approval.FieldConfig, error) {
	if flow == nil {
		return []approval.FieldConfig{}, nil
	}
	return loadFlowFieldsFromString(flow.FormFields)
}

// marshalFlowFormFields serializes validated field definitions for storage.
func marshalFlowFormFields(fields []approval.FieldConfig) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	content, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Create creates one approval flow with its ordered nodes.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !FlowType(in.FlowType) {
		return 0, bizerr.NewCode(CodeFlowTypeInvalid)
	}
	if !FlowStatus(in.Status) {
		return 0, bizerr.NewCode(CodeFlowStatusInvalid)
	}
	if err := s.validateNodeInputs(in.Nodes); err != nil {
		return 0, err
	}
	if err := s.ensureFlowNameFree(ctx, in.FlowName, 0); err != nil {
		return 0, err
	}
	if _, err := s.resolveApproverNames(ctx, in.Nodes); err != nil {
		return 0, err
	}
	fields, err := approval.ValidateFormFields(in.Fields)
	if err != nil {
		return 0, err
	}
	formFieldsJSON, err := marshalFlowFormFields(fields)
	if err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	var flowID int64
	err = dao.Flow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// GoFrame auto-fills created_at and updated_at.
		id, insertErr := dao.Flow.Ctx(ctx).Data(do.Flow{
			TenantId:    tenantID,
			FlowType:    in.FlowType,
			FlowName:    strings.TrimSpace(in.FlowName),
			Description: in.Description,
			Status:      in.Status,
			FormFields:  formFieldsJSON,
			CreatedBy:   createdBy,
			UpdatedBy:   createdBy,
		}).InsertAndGetId()
		if insertErr != nil {
			return insertErr
		}
		flowID = id
		return s.replaceFlowNodes(ctx, tenantID, createdBy, id, in.Nodes)
	})
	if err != nil {
		return 0, err
	}
	return flowID, nil
}

// Update updates the flow fields and optionally replaces the node list.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	flowColumns := dao.Flow.Columns()

	var existing *FlowEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "").
		Where(flowColumns.Id, in.Id).
		Scan(&existing)
	if err != nil {
		return err
	}
	if existing == nil {
		return bizerr.NewCode(CodeFlowNotFound)
	}

	if in.FlowType != nil {
		if !FlowType(*in.FlowType) {
			return bizerr.NewCode(CodeFlowTypeInvalid)
		}
	}
	if in.Status != nil && !FlowStatus(*in.Status) {
		return bizerr.NewCode(CodeFlowStatusInvalid)
	}
	if in.FlowName != nil {
		if err := s.ensureFlowNameFree(ctx, *in.FlowName, in.Id); err != nil {
			return err
		}
	}
	if len(in.Nodes) > MaxFlowNodes {
		return bizerr.NewCode(CodeFlowNodesInvalid)
	}
	if in.Nodes != nil {
		if err := s.validateNodeInputs(in.Nodes); err != nil {
			return err
		}
		if _, err = s.resolveApproverNames(ctx, in.Nodes); err != nil {
			return err
		}
	}
	var formFieldsJSON string
	if in.Fields != nil {
		fields, err := approval.ValidateFormFields(*in.Fields)
		if err != nil {
			return err
		}
		if formFieldsJSON, err = marshalFlowFormFields(fields); err != nil {
			return err
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Flow{UpdatedBy: updatedBy}
	if in.FlowType != nil {
		data.FlowType = *in.FlowType
	}
	if in.FlowName != nil {
		data.FlowName = strings.TrimSpace(*in.FlowName)
	}
	if in.Description != nil {
		data.Description = *in.Description
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	if in.Fields != nil {
		data.FormFields = formFieldsJSON
	}

	return dao.Flow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := dao.Flow.Ctx(ctx).
			OmitNilData().
			Where(tenantspi.TenantFilterColumn, tenantID).
			Where(flowColumns.Id, in.Id).
			Data(data).
			Update(); err != nil {
			return err
		}
		if in.Nodes != nil {
			return s.replaceFlowNodes(ctx, tenantID, updatedBy, in.Id, in.Nodes)
		}
		return nil
	})
}

// Delete soft-deletes flows and their nodes by IDs.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeFlowDeleteRequired)
	}

	return dao.Flow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		flowColumns := dao.Flow.Columns()
		if _, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "").
			WhereIn(flowColumns.Id, idList).
			Delete(); err != nil {
			return err
		}
		_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.FlowNode.Ctx(ctx), "").
			WhereIn(dao.FlowNode.Columns().FlowId, idList).
			Delete()
		return err
	})
}

// UserOptions returns bounded approver candidates through the host user domain.
func (s *serviceImpl) UserOptions(ctx context.Context, in UserOptionsInput) ([]*UserOptionItem, error) {
	if s.userSvc == nil {
		return nil, gerror.New("linapro-oa-approval requires host user capability")
	}
	result, err := s.userSvc.List(ctx, usercap.ListInput{
		Keyword: strings.TrimSpace(in.Keyword),
		Page: capmodel.PageRequest{
			PageNum:  1,
			PageSize: userSearchLimit,
			Limit:    userSearchLimit,
		},
	})
	if err != nil || result == nil {
		return nil, err
	}
	items := make([]*UserOptionItem, 0, len(result.Items))
	for _, user := range result.Items {
		if user == nil {
			continue
		}
		id, ok := userDomainID(user.ID)
		if !ok {
			continue
		}
		items = append(items, &UserOptionItem{
			Id:   id,
			Name: userDisplayName(user),
		})
	}
	return items, nil
}

// validateNodeInputs enforces the node count boundary and non-positive
// approver rejection; ordering is assigned by list position.
func (s *serviceImpl) validateNodeInputs(nodes []NodeInput) error {
	if len(nodes) == 0 || len(nodes) > MaxFlowNodes {
		return bizerr.NewCode(CodeFlowNodesInvalid)
	}
	for _, node := range nodes {
		if node.ApproverId <= 0 {
			return bizerr.NewCode(CodeApproverInvalid)
		}
	}
	return nil
}

// ensureFlowNameFree verifies the tenant-unique flow name excluding one flow.
func (s *serviceImpl) ensureFlowNameFree(ctx context.Context, name string, excludeID int64) error {
	flowColumns := dao.Flow.Columns()
	trimmed := strings.TrimSpace(name)
	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "").
		Where(flowColumns.FlowName, trimmed)
	if excludeID > 0 {
		m = m.WhereNot(flowColumns.Id, excludeID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeFlowNameExists)
	}
	return nil
}

// replaceFlowNodes soft-deletes existing nodes and inserts the ordered node
// list inside the caller transaction. Orders are assigned from list position.
func (s *serviceImpl) replaceFlowNodes(
	ctx context.Context,
	tenantID int,
	updatedBy int64,
	flowID int64,
	nodes []NodeInput,
) error {
	if _, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.FlowNode.Ctx(ctx), "").
		Where(dao.FlowNode.Columns().FlowId, flowID).
		Delete(); err != nil {
		return err
	}
	for order, node := range nodes {
		if _, err := dao.FlowNode.Ctx(ctx).Data(do.FlowNode{
			TenantId:   tenantID,
			FlowId:     flowID,
			NodeOrder:  order + 1,
			ApproverId: node.ApproverId,
			CreatedBy:  updatedBy,
			UpdatedBy:  updatedBy,
		}).Insert(); err != nil {
			return err
		}
	}
	return nil
}

// loadFlowNodes returns the ordered nodes of one flow with approver names
// resolved through one user batch call.
func (s *serviceImpl) loadFlowNodes(ctx context.Context, flowID int64) ([]*NodeItem, error) {
	nodeColumns := dao.FlowNode.Columns()
	nodes := make([]*FlowNodeEntity, 0)
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.FlowNode.Ctx(ctx), "").
		Where(nodeColumns.FlowId, flowID).
		OrderAsc(nodeColumns.NodeOrder).
		Scan(&nodes)
	if err != nil {
		return nil, err
	}

	approverIDs := make([]usercap.UserID, 0, len(nodes))
	seen := make(map[int64]struct{}, len(nodes))
	for _, node := range nodes {
		if node == nil || node.ApproverId <= 0 {
			continue
		}
		if _, ok := seen[node.ApproverId]; ok {
			continue
		}
		seen[node.ApproverId] = struct{}{}
		approverIDs = append(approverIDs, usercap.UserID(strconv.FormatInt(node.ApproverId, 10)))
	}
	names, err := s.resolveUserNames(ctx, approverIDs)
	if err != nil {
		return nil, err
	}

	items := make([]*NodeItem, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, &NodeItem{
			Order:        int64(node.NodeOrder),
			ApproverId:   node.ApproverId,
			ApproverName: names[node.ApproverId],
		})
	}
	return items, nil
}

// resolveNodeCountMap resolves node counts for the current page with one
// grouped node query, keeping list assembly free of per-flow lookups.
func (s *serviceImpl) resolveNodeCountMap(ctx context.Context, flows []*FlowEntity) (map[int64]int, error) {
	counts := make(map[int64]int)
	flowIDs := make([]int64, 0, len(flows))
	for _, flow := range flows {
		if flow == nil {
			continue
		}
		flowIDs = append(flowIDs, flow.Id)
	}
	if len(flowIDs) == 0 {
		return counts, nil
	}

	nodeColumns := dao.FlowNode.Columns()
	type nodeCountRow struct {
		FlowId int64 `json:"flowId"`
		Total  int   `json:"total"`
	}
	rows := make([]*nodeCountRow, 0, len(flowIDs))
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.FlowNode.Ctx(ctx), "").
		Fields(nodeColumns.FlowId, "COUNT(*) AS total").
		WhereIn(nodeColumns.FlowId, flowIDs).
		Group(nodeColumns.FlowId).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		counts[row.FlowId] = row.Total
	}
	return counts, nil
}

// resolveApproverNames resolves node approver display names before node
// insertion; every approver must be resolvable through the host user domain.
func (s *serviceImpl) resolveApproverNames(ctx context.Context, nodes []NodeInput) (map[int64]string, error) {
	if s.userSvc == nil {
		return nil, gerror.New("linapro-oa-approval requires host user capability")
	}
	ids := make([]usercap.UserID, 0, len(nodes))
	seen := make(map[int64]struct{}, len(nodes))
	for _, node := range nodes {
		if _, ok := seen[node.ApproverId]; ok {
			continue
		}
		seen[node.ApproverId] = struct{}{}
		ids = append(ids, usercap.UserID(strconv.FormatInt(node.ApproverId, 10)))
	}
	result, err := s.userSvc.BatchGet(ctx, ids)
	if err != nil {
		return nil, err
	}
	names := make(map[int64]string, len(ids))
	for _, node := range nodes {
		if _, ok := names[node.ApproverId]; ok {
			continue
		}
		projection := result.Items[usercap.UserID(strconv.FormatInt(node.ApproverId, 10))]
		if projection == nil {
			return nil, bizerr.NewCode(CodeApproverInvalid)
		}
		names[node.ApproverId] = userDisplayName(projection)
	}
	return names, nil
}

// resolveUserNames resolves display names for a de-duplicated user ID list
// through one user-domain batch call and leaves missing users blank.
func (s *serviceImpl) resolveUserNames(ctx context.Context, ids []usercap.UserID) (map[int64]string, error) {
	names := make(map[int64]string)
	if len(ids) == 0 {
		return names, nil
	}
	if s.userSvc == nil {
		return nil, gerror.New("linapro-oa-approval requires host user capability")
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

// creatorUserIDs collects de-duplicated creator IDs from the current page.
func creatorUserIDs(flows []*FlowEntity) []usercap.UserID {
	ids := make([]usercap.UserID, 0, len(flows))
	seen := make(map[int64]struct{}, len(flows))
	for _, flow := range flows {
		if flow == nil || flow.CreatedBy <= 0 {
			continue
		}
		if _, ok := seen[flow.CreatedBy]; ok {
			continue
		}
		seen[flow.CreatedBy] = struct{}{}
		ids = append(ids, usercap.UserID(strconv.FormatInt(flow.CreatedBy, 10)))
	}
	return ids
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
