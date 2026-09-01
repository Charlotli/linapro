// approval_v1_get.go implements the controller method that serves the
// approval-request detail endpoint with snapshot nodes and the timeline.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Get returns approval request details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	detail, err := c.approvalSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{
		RequestItem: toAPIRequestItem(&approvalsvc.RequestItem{
			RequestEntity:       detail.RequestEntity,
			ApplicantName:       detail.ApplicantName,
			CurrentApproverName: detail.CurrentApproverName,
		}),
		Nodes:      toAPINodeItems(detail.Nodes),
		Records:    toAPIRecordItems(detail.Records),
		FormFields: toAPIFormFieldItems(detail.FormFields),
		Form:       detail.FormValues,
	}, nil
}

// toAPIFormFieldItems converts service form field definitions into API DTO
// items.
func toAPIFormFieldItems(fields []approvalsvc.FieldConfig) []*v1.FormFieldItem {
	items := make([]*v1.FormFieldItem, 0, len(fields))
	for _, field := range fields {
		columns := make([]*v1.FormFieldColumnItem, 0, len(field.Columns))
		for _, column := range field.Columns {
			columns = append(columns, &v1.FormFieldColumnItem{
				Key:   column.Key,
				Label: column.Label,
				Type:  column.Type,
			})
		}
		items = append(items, &v1.FormFieldItem{
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
