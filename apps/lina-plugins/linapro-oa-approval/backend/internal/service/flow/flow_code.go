// This file defines linapro-oa-approval flow business error codes.

package flow

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeFlowNotFound reports that the requested approval flow does not exist
	// in the current tenant.
	CodeFlowNotFound = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_NOT_FOUND",
		"Approval flow does not exist",
		gcode.CodeNotFound,
	)
	// CodeFlowNameExists reports that the flow name is already used in the
	// current tenant.
	CodeFlowNameExists = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_NAME_EXISTS",
		"Approval flow name already exists",
		gcode.CodeInvalidParameter,
	)
	// CodeFlowTypeInvalid reports an unknown approval flow type value.
	CodeFlowTypeInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_TYPE_INVALID",
		"Approval flow type is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeFlowStatusInvalid reports an unknown approval flow status value.
	CodeFlowStatusInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_STATUS_INVALID",
		"Approval flow status is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeFlowNodesInvalid reports an invalid ordered node list boundary, such
	// as an empty list or more than the allowed node count.
	CodeFlowNodesInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_NODES_INVALID",
		"Approval flow nodes must contain 1 to 20 approvers",
		gcode.CodeInvalidParameter,
	)
	// CodeApproverInvalid reports that one approver user does not exist.
	CodeApproverInvalid = bizerr.MustDefine(
		"OA_APPROVAL_APPROVER_INVALID",
		"Approver user does not exist",
		gcode.CodeInvalidParameter,
	)
	// CodeFlowDeleteRequired reports that a delete operation received no IDs.
	CodeFlowDeleteRequired = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_DELETE_REQUIRED",
		"Select at least one approval flow to delete",
		gcode.CodeInvalidParameter,
	)
)
