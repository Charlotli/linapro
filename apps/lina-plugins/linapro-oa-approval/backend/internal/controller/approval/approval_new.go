// This file wires the linapro-oa-approval request controller and shared
// response mappers.

package approval

import (
	"lina-core/pkg/apitime"
	approvalapi "lina-plugin-linapro-oa-approval/backend/api/approval"
	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// ControllerV1 is the approval-request controller.
type ControllerV1 struct {
	approvalSvc approvalsvc.Service // approval-request service
}

// NewV1 creates and returns a new approval-request controller instance.
func NewV1(approvalSvc approvalsvc.Service) approvalapi.IApprovalV1 {
	return &ControllerV1{approvalSvc: approvalSvc}
}

// toAPIRequestItem converts the service-layer list projection into the API
// DTO projection returned through HTTP responses.
func toAPIRequestItem(item *approvalsvc.RequestItem) v1.RequestItem {
	if item == nil || item.RequestEntity == nil {
		return v1.RequestItem{}
	}
	return v1.RequestItem{
		Id:                  item.Id,
		FlowId:              item.FlowId,
		FlowType:            item.FlowType,
		Title:               item.Title,
		Amount:              item.Amount,
		Content:             item.Content,
		Attachments:         item.Attachments,
		Status:              item.Status,
		CurrentNodeOrder:    item.CurrentNodeOrder,
		CurrentApproverId:   item.CurrentApproverId,
		CurrentApproverName: item.CurrentApproverName,
		ApplicantId:         item.ApplicantId,
		ApplicantName:       item.ApplicantName,
		CreatedAt:           apitime.Milli(item.CreatedAt),
		UpdatedAt:           apitime.Milli(item.UpdatedAt),
	}
}

// toAPINodeItems converts service-layer snapshot nodes into API DTO items.
func toAPINodeItems(nodes []*approvalsvc.NodeItem) []*v1.RequestNodeItem {
	items := make([]*v1.RequestNodeItem, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, &v1.RequestNodeItem{
			Order:        node.Order,
			ApproverId:   node.ApproverId,
			ApproverName: node.ApproverName,
		})
	}
	return items
}

// toAPIRecordItems converts service-layer timeline records into API DTO items.
func toAPIRecordItems(records []*approvalsvc.RecordItem) []*v1.RequestRecordItem {
	items := make([]*v1.RequestRecordItem, 0, len(records))
	for _, record := range records {
		items = append(items, &v1.RequestRecordItem{
			Id:        record.Id,
			NodeOrder: record.NodeOrder,
			Action:    record.Action,
			ActorId:   record.ActorId,
			ActorName: record.ActorName,
			Comment:   record.Comment,
			CreatedAt: apitime.Milli(record.CreatedAt),
		})
	}
	return items
}
