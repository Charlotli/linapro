// approval_impl.go implements the tenant-scoped approval-request state
// machine: submit with frozen node snapshots, approve/reject/comment/append/
// withdraw actions with participant boundaries, batched name assembly, and
// post-commit notification fan-out. Tenant filters are applied at the database
// query stage and every action runs with its records in one transaction.

package approval

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-oa-approval/backend/internal/dao"
	"lina-plugin-linapro-oa-approval/backend/internal/model/do"
)

// List queries approval requests with pagination, scope, and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	requestColumns := dao.Request.Columns()
	currentUserID := s.currentUserID(ctx)

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "")
	switch in.Scope {
	case scopePening:
		m = m.Where(requestColumns.Status, requestStatusPending).
			Where(requestColumns.CurrentApproverId, currentUserID)
	default:
		m = m.Where(requestColumns.ApplicantId, currentUserID)
	}
	if title := strings.TrimSpace(in.Title); title != "" {
		m = m.WhereLike(requestColumns.Title, "%"+title+"%")
	}
	if in.FlowType > 0 {
		m = m.Where(requestColumns.FlowType, in.FlowType)
	}
	if in.Status != nil {
		m = m.Where(requestColumns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	requests := make([]*RequestEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(requestColumns.Id).
		Scan(&requests)
	if err != nil {
		return nil, err
	}

	userNames, err := s.resolveUserNames(ctx, requestUserIDs(requests))
	if err != nil {
		return nil, err
	}

	items := make([]*RequestItem, 0, len(requests))
	for _, request := range requests {
		items = append(items, &RequestItem{
			RequestEntity:       request,
			ApplicantName:       userNames[request.ApplicantId],
			CurrentApproverName: userNames[request.CurrentApproverId],
		})
	}

	return &ListOutput{
		List:  items,
		Total: total,
	}, nil
}

// GetById returns one request with its snapshot nodes and timeline.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*RequestDetail, error) {
	request, err := s.loadRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	nodes, err := parseSnapshotNodes(request.NodeSnapshot)
	if err != nil {
		return nil, err
	}
	if !isRequestParticipant(request, nodes, s.currentUserID(ctx)) {
		return nil, bizerr.NewCode(CodeRequestNotFound)
	}

	nodeItems := make([]*NodeItem, 0, len(nodes))
	nodeApproverIDs := make([]usercap.UserID, 0, len(nodes))
	seen := make(map[int64]struct{}, len(nodes))
	for _, node := range nodes {
		nodeItems = append(nodeItems, &NodeItem{
			Order:        node.Order,
			ApproverId:   node.ApproverId,
			ApproverName: node.ApproverName,
		})
		if _, ok := seen[node.ApproverId]; ok {
			continue
		}
		seen[node.ApproverId] = struct{}{}
		nodeApproverIDs = append(nodeApproverIDs, usercap.UserID(strconv.FormatInt(node.ApproverId, 10)))
	}

	recordColumns := dao.Record.Columns()
	records := make([]*RecordEntity, 0, MaxRecords)
	err = tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Record.Ctx(ctx), "").
		Where(recordColumns.RequestId, id).
		OrderAsc(recordColumns.CreatedAt).
		OrderAsc(recordColumns.Id).
		Limit(MaxRecords).
		Scan(&records)
	if err != nil {
		return nil, err
	}

	timelineUserIDs := nodeApproverIDs
	for _, record := range records {
		if record == nil || record.ActorId <= 0 {
			continue
		}
		timelineUserIDs = append(timelineUserIDs, usercap.UserID(strconv.FormatInt(record.ActorId, 10)))
	}
	userNames, err := s.resolveUserNames(ctx, timelineUserIDs)
	if err != nil {
		return nil, err
	}
	// Snapshot names are frozen at submit time; refresh them so renamed users
	// still display consistently with the timeline actors.
	for _, item := range nodeItems {
		if name := userNames[item.ApproverId]; name != "" {
			item.ApproverName = name
		}
	}

	recordItems := make([]*RecordItem, 0, len(records))
	for _, record := range records {
		recordItems = append(recordItems, &RecordItem{
			RecordEntity: record,
			ActorName:    userNames[record.ActorId],
		})
	}

	formFields, err := parseFormSnapshot(request.FormSnapshot)
	if err != nil {
		return nil, err
	}
	formValues := make(map[string]any)
	if trimmedForm := strings.TrimSpace(request.FormData); trimmedForm != "" {
		if err = json.Unmarshal([]byte(trimmedForm), &formValues); err != nil {
			return nil, errors.New("parse approval form values failed: " + err.Error())
		}
	}

	return &RequestDetail{
		RequestEntity:       request,
		ApplicantName:       userNames[request.ApplicantId],
		CurrentApproverName: userNames[request.CurrentApproverId],
		Nodes:               nodeItems,
		Records:             recordItems,
		FormFields:          formFields,
		FormValues:          formValues,
	}, nil
}

// Create submits one approval request with frozen node and form snapshots.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	flow, nodes, formFields, err := s.loadEnabledFlowDefinition(ctx, in.FlowId)
	if err != nil {
		return 0, err
	}
	formValues, err := validateFormValues(formFields, in.Form)
	if err != nil {
		return 0, err
	}
	nodeSnapshot, err := marshalSnapshotNodes(nodes)
	if err != nil {
		return 0, err
	}
	formSnapshot, err := marshalFormSnapshot(formFields)
	if err != nil {
		return 0, err
	}
	requestAmount := extractAmountValue(formFields, in.Form)

	var (
		bizCtx      = s.bizCtxSvc.Current(ctx)
		applicantID = int64(bizCtx.UserID)
		tenantID    = s.tenantFilterContext(ctx).TenantID
	)

	firstApproverID := nodes[0].ApproverId

	var requestID int64
	err = dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// GoFrame auto-fills created_at and updated_at.
		id, insertErr := dao.Request.Ctx(ctx).Data(do.Request{
			TenantId:          tenantID,
			FlowId:            in.FlowId,
			FlowType:          flow.FlowType,
			Title:             strings.TrimSpace(in.Title),
			Amount:            requestAmount,
			Content:           formTextValue(in.Form, "content"),
			Status:            requestStatusPending,
			CurrentNodeOrder:  1,
			CurrentApproverId: firstApproverID,
			NodeSnapshot:      nodeSnapshot,
			FormData:          formValues,
			FormSnapshot:      formSnapshot,
			ApplicantId:       applicantID,
			CreatedBy:         applicantID,
			UpdatedBy:         applicantID,
		}).InsertAndGetId()
		if insertErr != nil {
			return insertErr
		}
		requestID = id
		return s.insertRecord(ctx, tenantID, id, 0, recordActionSubmit, applicantID, "")
	})
	if err != nil {
		return 0, err
	}

	s.notifier.send(ctx, notifyKeySubmit, firstApproverID, requestID, strings.TrimSpace(in.Title))
	return requestID, nil
}

// Approve approves one pending request as the current node approver.
func (s *serviceImpl) Approve(ctx context.Context, in ActionInput) error {
	nextApproverID, requestID, requestTitle, err := s.transitionPendingRequest(ctx, in, recordActionApprove)
	if err != nil {
		return err
	}
	if nextApproverID > 0 {
		s.notifier.send(ctx, notifyKeyNext, nextApproverID, requestID, requestTitle)
		// Keep the applicant informed about the intermediate progress.
		request, loadErr := s.loadRequest(ctx, in.Id)
		if loadErr == nil {
			s.notifier.send(ctx, notifyKeyAdvance, request.ApplicantId, request.Id, request.Title)
		}
		return nil
	}
	// The last node finished: notify the applicant about the final approval.
	request, loadErr := s.loadRequest(ctx, in.Id)
	if loadErr != nil {
		return nil
	}
	s.notifier.send(ctx, notifyKeyFinal, request.ApplicantId, request.Id, request.Title)
	return nil
}

// Reject rejects one pending request as the current node approver.
func (s *serviceImpl) Reject(ctx context.Context, in ActionInput) error {
	_, requestID, requestTitle, err := s.transitionPendingRequest(ctx, in, recordActionReject)
	if err != nil {
		return err
	}
	request, loadErr := s.loadRequest(ctx, in.Id)
	if loadErr != nil {
		return nil
	}
	s.notifier.send(ctx, notifyKeyReject, request.ApplicantId, requestID, requestTitle)
	return nil
}

// transitionPendingRequest validates the current-approver boundary, applies
// the approve or reject transition, and writes the timeline entry inside one
// transaction. It returns the next pending approver for approval advancement
// (zero when the request finished or rejected).
func (s *serviceImpl) transitionPendingRequest(
	ctx context.Context,
	in ActionInput,
	action int,
) (int64, int64, string, error) {
	request, err := s.loadRequest(ctx, in.Id)
	if err != nil {
		return 0, 0, "", err
	}
	if request.Status != requestStatusPending {
		return 0, 0, "", bizerr.NewCode(CodeRequestNotPending)
	}
	currentUserID := s.currentUserID(ctx)
	if request.CurrentApproverId != currentUserID {
		return 0, 0, "", bizerr.NewCode(CodeApproveNoAuthority)
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	updatedBy := currentUserID
	targetStatus := requestStatusApproved
	nextApproverID := int64(0)
	nextNodeOrder := request.CurrentNodeOrder
	nextApprover := request.CurrentApproverId

	if action == recordActionApprove {
		nodes, parseErr := parseSnapshotNodes(request.NodeSnapshot)
		if parseErr != nil {
			return 0, 0, "", parseErr
		}
		if request.CurrentNodeOrder < len(nodes) {
			// More nodes follow: advance the pending pointer.
			nextNodeOrder = request.CurrentNodeOrder + 1
			nextApproverID = nodes[nextNodeOrder-1].ApproverId
			nextApprover = nextApproverID
		} else {
			nextNodeOrder = request.CurrentNodeOrder
		}
	} else {
		targetStatus = requestStatusRejected
	}

	err = dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := s.insertRecord(ctx, tenantID, request.Id, request.CurrentNodeOrder, action, updatedBy, in.Comment); err != nil {
			return err
		}
		requestColumns := dao.Request.Columns()
		_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
			Where(requestColumns.Id, request.Id).
			Data(do.Request{
				Status:            targetStatus,
				CurrentNodeOrder:  nextNodeOrder,
				CurrentApproverId: nextApprover,
				UpdatedBy:         updatedBy,
			}).
			Update()
		return err
	})
	if err != nil {
		return 0, 0, "", err
	}
	return nextApproverID, request.Id, request.Title, nil
}

// Comment appends one reply to the pending request timeline.
func (s *serviceImpl) Comment(ctx context.Context, in ActionInput) error {
	request, err := s.loadRequest(ctx, in.Id)
	if err != nil {
		return err
	}
	if request.Status != requestStatusPending {
		return bizerr.NewCode(CodeRequestNotPending)
	}
	nodes, err := parseSnapshotNodes(request.NodeSnapshot)
	if err != nil {
		return err
	}
	if !isRequestParticipant(request, nodes, s.currentUserID(ctx)) {
		return bizerr.NewCode(CodeRequestNotFound)
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	return dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return s.insertRecord(ctx, tenantID, request.Id, 0, recordActionComment, s.currentUserID(ctx), in.Comment)
	})
}

// Appender inserts one new approver node after the current pending node.
func (s *serviceImpl) Appender(ctx context.Context, in AppenderInput) error {
	request, err := s.loadRequest(ctx, in.Id)
	if err != nil {
		return err
	}
	if request.Status != requestStatusPending {
		return bizerr.NewCode(CodeRequestNotPending)
	}
	currentUserID := s.currentUserID(ctx)
	if request.CurrentApproverId != currentUserID {
		return bizerr.NewCode(CodeApproveNoAuthority)
	}
	nodes, err := parseSnapshotNodes(request.NodeSnapshot)
	if err != nil {
		return err
	}
	if in.ApproverId <= 0 {
		return bizerr.NewCode(CodeApproverInvalid)
	}
	for _, node := range nodes {
		if node.ApproverId == in.ApproverId {
			return bizerr.NewCode(CodeApproverDuplicate)
		}
	}
	newApproverName, err := s.resolveRequiredUserName(ctx, in.ApproverId)
	if err != nil {
		return err
	}

	// Insert the new node at the zero-based position right after the current
	// node and shift later nodes back.
	insertAt := request.CurrentNodeOrder
	nodes = append(nodes, snapshotNode{})
	for index := len(nodes) - 1; index > insertAt; index-- {
		nodes[index] = nodes[index-1]
	}
	nodes[insertAt] = snapshotNode{
		Order:        int64(insertAt + 1),
		ApproverId:   in.ApproverId,
		ApproverName: newApproverName,
	}
	snapshot, err := marshalSnapshotNodes(nodes)
	if err != nil {
		return err
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	err = dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := s.insertRecord(ctx, tenantID, request.Id, request.CurrentNodeOrder, recordActionAppend, currentUserID, ""); err != nil {
			return err
		}
		requestColumns := dao.Request.Columns()
		_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
			Where(requestColumns.Id, request.Id).
			Data(do.Request{
				NodeSnapshot: snapshot,
				UpdatedBy:    currentUserID,
			}).
			Update()
		return err
	})
	if err != nil {
		return err
	}
	s.notifier.send(ctx, notifyKeyAppends, in.ApproverId, request.Id, request.Title)
	return nil
}

// Withdraw withdraws one pending request as the applicant.
func (s *serviceImpl) Withdraw(ctx context.Context, id int64) error {
	request, err := s.loadRequest(ctx, id)
	if err != nil {
		return err
	}
	if request.Status != requestStatusPending {
		return bizerr.NewCode(CodeRequestNotPending)
	}
	currentUserID := s.currentUserID(ctx)
	if request.ApplicantId != currentUserID {
		return bizerr.NewCode(CodeWithdrawNoAuthority)
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	err = dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := s.insertRecord(ctx, tenantID, request.Id, 0, recordActionWit, currentUserID, ""); err != nil {
			return err
		}
		requestColumns := dao.Request.Columns()
		_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
			Where(requestColumns.Id, request.Id).
			Data(do.Request{
				Status:    requestStatusWithdrawn,
				UpdatedBy: currentUserID,
			}).
			Update()
		return err
	})
	if err != nil {
		return err
	}
	s.notifier.send(ctx, notifyKeyWit, request.CurrentApproverId, request.Id, request.Title)
	return nil
}

// Resubmit restarts one rejected or withdrawn request from the first node.
func (s *serviceImpl) Resubmit(ctx context.Context, id int64) error {
	request, err := s.loadRequest(ctx, id)
	if err != nil {
		return err
	}
	if request.Status != requestStatusRejected && request.Status != requestStatusWithdrawn {
		return bizerr.NewCode(CodeRequestNotPending)
	}
	currentUserID := s.currentUserID(ctx)
	if request.ApplicantId != currentUserID {
		return bizerr.NewCode(CodeResubmitNoAuthority)
	}
	nodes, err := parseSnapshotNodes(request.NodeSnapshot)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return bizerr.NewCode(CodeFlowUnavailable)
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	firstApproverID := nodes[0].ApproverId
	err = dao.Request.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := s.insertRecord(ctx, tenantID, request.Id, 0, recordActionSubmit, currentUserID, ""); err != nil {
			return err
		}
		requestColumns := dao.Request.Columns()
		_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
			Where(requestColumns.Id, request.Id).
			Data(do.Request{
				Status:            requestStatusPending,
				CurrentNodeOrder:  1,
				CurrentApproverId: firstApproverID,
				UpdatedBy:         currentUserID,
			}).
			Update()
		return err
	})
	if err != nil {
		return err
	}
	s.notifier.send(ctx, notifyKeySubmit, firstApproverID, request.Id, request.Title)
	return nil
}

// PendingCount counts pending requests waiting for the current user.
func (s *serviceImpl) PendingCount(ctx context.Context) (int64, error) {
	requestColumns := dao.Request.Columns()
	count, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
		Where(requestColumns.Status, requestStatusPending).
		Where(requestColumns.CurrentApproverId, s.currentUserID(ctx)).
		Count()
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

// insertRecord writes one timeline entry inside the caller transaction.
func (s *serviceImpl) insertRecord(
	ctx context.Context,
	tenantID int,
	requestID int64,
	nodeOrder int,
	action int,
	actorID int64,
	comment string,
) error {
	_, err := dao.Record.Ctx(ctx).Data(do.Record{
		TenantId:  tenantID,
		RequestId: requestID,
		NodeOrder: nodeOrder,
		Action:    action,
		ActorId:   actorID,
		Comment:   comment,
		CreatedBy: actorID,
		UpdatedBy: actorID,
	}).Insert()
	return err
}

// loadRequest loads one tenant-scoped request; out-of-scope or absent rows
// return the shared not-found error so existence never leaks.
func (s *serviceImpl) loadRequest(ctx context.Context, id int64) (*RequestEntity, error) {
	requestColumns := dao.Request.Columns()
	var request *RequestEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Request.Ctx(ctx), "").
		Where(requestColumns.Id, id).
		Scan(&request)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, bizerr.NewCode(CodeRequestNotFound)
	}
	return request, nil
}

// loadEnabledFlowDefinition loads one enabled tenant flow with its ordered
// nodes and frozen-form field definitions for submit-time snapshots.
func (s *serviceImpl) loadEnabledFlowDefinition(ctx context.Context, flowID int64) (*FlowEntity, []snapshotNode, []FieldConfig, error) {
	flow, nodes, err := s.loadEnabledFlowNodes(ctx, flowID)
	if err != nil {
		return nil, nil, nil, err
	}
	fields, err := parseFormSnapshot(flowFormFieldsJSON(flow))
	if err != nil {
		return nil, nil, nil, err
	}
	return flow, nodes, fields, nil
}

// flowFormFieldsJSON reads the flow form-field storage without mutating the
// entity. Empty storage maps to an empty snapshot handled by the parser.
func flowFormFieldsJSON(flow *FlowEntity) string {
	if flow == nil {
		return ""
	}
	return flow.FormFields
}

// formTextValue reads one text value from the submitted form map for legacy
// projection columns.
func formTextValue(values map[string]any, key string) string {
	if values == nil {
		return ""
	}
	text, _ := values[key].(string)
	return strings.TrimSpace(text)
}

// loadEnabledFlowNodes loads one enabled tenant flow with its ordered nodes
// and resolved approver names for snapshot freezing.
func (s *serviceImpl) loadEnabledFlowNodes(ctx context.Context, flowID int64) (*FlowEntity, []snapshotNode, error) {
	flowColumns := dao.Flow.Columns()
	var flow *FlowEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Flow.Ctx(ctx), "").
		Where(flowColumns.Id, flowID).
		Scan(&flow)
	if err != nil {
		return nil, nil, err
	}
	if flow == nil || flow.Status != 1 {
		return nil, nil, bizerr.NewCode(CodeFlowUnavailable)
	}

	nodeColumns := dao.FlowNode.Columns()
	nodes := make([]*FlowNodeEntity, 0)
	err = tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.FlowNode.Ctx(ctx), "").
		Where(nodeColumns.FlowId, flowID).
		OrderAsc(nodeColumns.NodeOrder).
		Scan(&nodes)
	if err != nil {
		return nil, nil, err
	}
	if len(nodes) == 0 {
		return nil, nil, bizerr.NewCode(CodeFlowUnavailable)
	}

	snapshot := make([]snapshotNode, 0, len(nodes))
	approverIDs := make([]usercap.UserID, 0, len(nodes))
	seen := make(map[int64]struct{}, len(nodes))
	for _, node := range nodes {
		snapshot = append(snapshot, snapshotNode{
			Order:        int64(node.NodeOrder),
			ApproverId:   node.ApproverId,
			ApproverName: "",
		})
		if _, ok := seen[node.ApproverId]; ok {
			continue
		}
		seen[node.ApproverId] = struct{}{}
		approverIDs = append(approverIDs, usercap.UserID(strconv.FormatInt(node.ApproverId, 10)))
	}
	userNames, err := s.resolveUserNames(ctx, approverIDs)
	if err != nil {
		return nil, nil, err
	}
	for index := range snapshot {
		if name := userNames[snapshot[index].ApproverId]; name != "" {
			snapshot[index].ApproverName = name
		}
	}
	return flow, snapshot, nil
}

// resolveRequiredUserName resolves one user display name and reports missing
// users as invalid approvers.
func (s *serviceImpl) resolveRequiredUserName(ctx context.Context, userID int64) (string, error) {
	if s.userSvc == nil {
		return "", gerror.New("linapro-oa-approval requires host user capability")
	}
	result, err := s.userSvc.BatchGet(ctx, []usercap.UserID{usercap.UserID(strconv.FormatInt(userID, 10))})
	if err != nil {
		return "", err
	}
	projection := result.Items[usercap.UserID(strconv.FormatInt(userID, 10))]
	if projection == nil {
		return "", bizerr.NewCode(CodeApproverInvalid)
	}
	if projection.Username != "" {
		return projection.Username, nil
	}
	if projection.Nickname != "" {
		return projection.Nickname, nil
	}
	return string(projection.ID), nil
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

// requestUserIDs collects de-duplicated applicant and current approver IDs
// from the current page.
func requestUserIDs(requests []*RequestEntity) []usercap.UserID {
	ids := make([]usercap.UserID, 0, len(requests)*2)
	seen := make(map[int64]struct{}, len(requests)*2)
	for _, request := range requests {
		if request == nil {
			continue
		}
		for _, id := range []int64{request.ApplicantId, request.CurrentApproverId} {
			if id <= 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, usercap.UserID(strconv.FormatInt(id, 10)))
		}
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

// currentUserID returns the current user domain ID from the business context.
func (s *serviceImpl) currentUserID(ctx context.Context) int64 {
	return int64(s.bizCtxSvc.Current(ctx).UserID)
}
