// This file wires the linapro-oa-approval flow controller and shared response
// mappers.

package flow

import (
	flowapi "lina-plugin-linapro-oa-approval/backend/api/flow"
	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	"lina-plugin-linapro-oa-approval/backend/internal/service/approval"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// ControllerV1 is the approval-flow controller.
type ControllerV1 struct {
	flowSvc flowsvc.Service // approval-flow service
}

// NewV1 creates and returns a new approval-flow controller instance.
func NewV1(flowSvc flowsvc.Service) flowapi.IFlowV1 {
	return &ControllerV1{flowSvc: flowSvc}
}

// toAPIFlowItem converts the service-layer list projection into the API DTO
// projection returned through HTTP responses.
func toAPIFlowItem(item *flowsvc.FlowItem) v1.FlowItem {
	if item == nil || item.FlowEntity == nil {
		return v1.FlowItem{}
	}
	return v1.FlowItem{
		Id:            item.Id,
		FlowType:      item.FlowType,
		FlowName:      item.FlowName,
		Description:   item.Description,
		Status:        item.Status,
		NodeCount:     item.NodeCount,
		FieldCount:    item.FieldCount,
		CreatedBy:     item.CreatedBy,
		CreatedByName: item.CreatedByName,
	}
}

// toAPINodeItems converts service-layer node projections into API DTO items.
func toAPINodeItems(nodes []*flowsvc.NodeItem) []*v1.FlowNodeItem {
	items := make([]*v1.FlowNodeItem, 0, len(nodes))
	for _, node := range nodes {
		items = append(items, &v1.FlowNodeItem{
			Order:        node.Order,
			ApproverId:   node.ApproverId,
			ApproverName: node.ApproverName,
		})
	}
	return items
}

// toServiceNodeInputs converts API node items into service-layer inputs while
// assigning orders from the submitted list position.
func toServiceNodeInputs(nodes []*v1.CreateNodeItem) []flowsvc.NodeInput {
	inputs := make([]flowsvc.NodeInput, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		inputs = append(inputs, flowsvc.NodeInput{ApproverId: node.ApproverId})
	}
	return inputs
}

// toServiceUpdateNodeInputs converts update API node items into service inputs.
func toServiceUpdateNodeInputs(nodes []*v1.UpdateNodeItem) []flowsvc.NodeInput {
	inputs := make([]flowsvc.NodeInput, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		inputs = append(inputs, flowsvc.NodeInput{ApproverId: node.ApproverId})
	}
	return inputs
}

// toCreateColumnInputs converts create API detail columns into service configs.
func toCreateColumnInputs(columns []*v1.CreateColumnItem) []approval.ColumnConfig {
	inputs := make([]approval.ColumnConfig, 0, len(columns))
	for _, column := range columns {
		if column == nil {
			continue
		}
		inputs = append(inputs, approval.ColumnConfig{
			Key:   column.Key,
			Label: column.Label,
			Type:  column.Type,
		})
	}
	return inputs
}

// toUpdateColumnInputs converts update API detail columns into service configs.
func toUpdateColumnInputs(columns []*v1.UpdateColumnItem) []approval.ColumnConfig {
	inputs := make([]approval.ColumnConfig, 0, len(columns))
	for _, column := range columns {
		if column == nil {
			continue
		}
		inputs = append(inputs, approval.ColumnConfig{
			Key:   column.Key,
			Label: column.Label,
			Type:  column.Type,
		})
	}
	return inputs
}

// toServiceCreateFieldInputs converts create API field items into service
// form field definitions.
func toServiceCreateFieldInputs(fields []*v1.CreateFieldItem) []approval.FieldConfig {
	inputs := make([]approval.FieldConfig, 0, len(fields))
	for _, field := range fields {
		if field == nil {
			continue
		}
		inputs = append(inputs, approval.FieldConfig{
			Key:      field.Key,
			Label:    field.Label,
			Type:     field.Type,
			Required: field.Required,
			AsAmount: field.AsAmount,
			Options:  field.Options,
			Columns:  toCreateColumnInputs(field.Columns),
		})
	}
	return inputs
}

// toServiceUpdateFieldInputs converts update API field items into service
// form field definitions.
func toServiceUpdateFieldInputs(fields []*v1.UpdateFieldItem) []approval.FieldConfig {
	inputs := make([]approval.FieldConfig, 0, len(fields))
	for _, field := range fields {
		if field == nil {
			continue
		}
		inputs = append(inputs, approval.FieldConfig{
			Key:      field.Key,
			Label:    field.Label,
			Type:     field.Type,
			Required: field.Required,
			AsAmount: field.AsAmount,
			Options:  field.Options,
			Columns:  toUpdateColumnInputs(field.Columns),
		})
	}
	return inputs
}

// toAPIFieldItems converts service form field definitions into API DTO items.
func toAPIFieldItems(fields []approval.FieldConfig) []*v1.FlowFieldItem {
	items := make([]*v1.FlowFieldItem, 0, len(fields))
	for _, field := range fields {
		columns := make([]*v1.FlowFieldColumn, 0, len(field.Columns))
		for _, column := range field.Columns {
			columns = append(columns, &v1.FlowFieldColumn{
				Key:   column.Key,
				Label: column.Label,
				Type:  column.Type,
			})
		}
		items = append(items, &v1.FlowFieldItem{
			Key:      field.Key,
			Label:    field.Label,
			Type:     field.Type,
			Required: field.Required,
			AsAmount: field.AsAmount,
			Options:  field.Options,
			Columns:  columns,
		})
	}
	return items
}
