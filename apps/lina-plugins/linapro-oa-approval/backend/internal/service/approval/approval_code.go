// This file defines linapro-oa-approval approval business error codes.

package approval

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeRequestNotFound reports that the requested approval record does not
	// exist or that the caller is not a participant of the request. Both cases
	// share one error so existence never leaks to non-participants.
	CodeRequestNotFound = bizerr.MustDefine(
		"OA_APPROVAL_REQUEST_NOT_FOUND",
		"Approval request does not exist",
		gcode.CodeNotFound,
	)
	// CodeRequestNotPending reports that the request is not in the pending
	// state required by the requested action.
	CodeRequestNotPending = bizerr.MustDefine(
		"OA_APPROVAL_REQUEST_NOT_PENDING",
		"Approval request is not pending",
		gcode.CodeInvalidOperation,
	)
	// CodeApproveNoAuthority reports that the caller is not the current node
	// approver required by the approve, reject, or append action.
	CodeApproveNoAuthority = bizerr.MustDefine(
		"OA_APPROVAL_APPROVE_NO_AUTHORITY",
		"Only the current approver can perform this action",
		gcode.CodeNotAuthorized,
	)
	// CodeWithdrawNoAuthority reports that the caller is not the applicant
	// required by the withdraw action.
	CodeWithdrawNoAuthority = bizerr.MustDefine(
		"OA_APPROVAL_WITHDRAW_NO_AUTHORITY",
		"Only the applicant can withdraw the request",
		gcode.CodeNotAuthorized,
	)
	// CodeFlowUnavailable reports that the referenced approval flow does not
	// exist, is disabled, or has no approval nodes.
	CodeFlowUnavailable = bizerr.MustDefine(
		"OA_APPROVAL_FLOW_UNAVAILABLE",
		"Approval flow does not exist, is disabled, or has no approver nodes",
		gcode.CodeInvalidParameter,
	)
	// CodeApproverInvalid reports that one referenced approver user does not
	// exist.
	CodeApproverInvalid = bizerr.MustDefine(
		"OA_APPROVAL_APPROVER_INVALID",
		"Approver user does not exist",
		gcode.CodeInvalidParameter,
	)
	// CodeApproverDuplicate reports that the appended approver already exists
	// in the frozen node snapshot.
	CodeApproverDuplicate = bizerr.MustDefine(
		"OA_APPROVAL_APPROVER_DUPLICATE",
		"Approver already exists in the approval chain",
		gcode.CodeInvalidParameter,
	)
	// CodeResubmitNoAuthority reports that the caller is not the applicant
	// required by the resubmit action.
	CodeResubmitNoAuthority = bizerr.MustDefine(
		"OA_APPROVAL_RESUBMIT_NO_AUTHORITY",
		"Only the applicant can resubmit the request",
		gcode.CodeNotAuthorized,
	)
	// CodeAttachmentsInvalid reports that the attachment list is not a valid
	// JSON array of URL strings or exceeds the allowed size.
	CodeAttachmentsInvalid = bizerr.MustDefine(
		"OA_APPROVAL_ATTACHMENTS_INVALID",
		"Attachments must be a JSON array of at most 10 URL addresses",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldsInvalid reports that the configured form field list
	// exceeds the allowed size.
	CodeFormFieldsInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELDS_INVALID",
		"Form fields must contain at most 30 entries",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldKeyInvalid reports a missing, duplicated, or illegal form
	// field storage key.
	CodeFormFieldKeyInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELD_KEY_INVALID",
		"Form field keys must be unique identifier-safe strings",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldLabelRequired reports a missing form field display name.
	CodeFormFieldLabelRequired = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELD_LABEL_REQUIRED",
		"Form field display name is required",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldOptionsRequired reports a select field without options.
	CodeFormFieldOptionsRequired = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELD_OPTIONS_REQUIRED",
		"Select fields must define at least one option",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldColumnsRequired reports an invalid detail sub-table
	// column configuration.
	CodeFormFieldColumnsRequired = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELD_COLUMNS_REQUIRED",
		"Detail fields must define 1 to 30 valid columns",
		gcode.CodeInvalidParameter,
	)
	// CodeFormFieldTypeInvalid reports an unknown form field type.
	CodeFormFieldTypeInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FORM_FIELD_TYPE_INVALID",
		"Form field type is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeFormValueInvalid reports a submitted form value that violates its
	// field type or option set.
	CodeFormValueInvalid = bizerr.MustDefine(
		"OA_APPROVAL_FORM_VALUE_INVALID",
		"Form field value is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeFormValueRequired reports a required form field left empty.
	CodeFormValueRequired = bizerr.MustDefine(
		"OA_APPROVAL_FORM_VALUE_REQUIRED",
		"Please fill in all required form fields",
		gcode.CodeInvalidParameter,
	)
)
