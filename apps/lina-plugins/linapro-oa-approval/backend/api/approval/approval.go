// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package approval

import (
	"context"

	"lina-plugin-linapro-oa-approval/backend/api/approval/v1"
)

type IApprovalV1 interface {
	Appender(ctx context.Context, req *v1.AppenderReq) (res *v1.AppenderRes, err error)
	Approve(ctx context.Context, req *v1.ApproveReq) (res *v1.ApproveRes, err error)
	Comment(ctx context.Context, req *v1.CommentReq) (res *v1.CommentRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error)
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	PendingCount(ctx context.Context, req *v1.PendingCountReq) (res *v1.PendingCountRes, err error)
	Reject(ctx context.Context, req *v1.RejectReq) (res *v1.RejectRes, err error)
	Resubmit(ctx context.Context, req *v1.ResubmitReq) (res *v1.ResubmitRes, err error)
	Withdraw(ctx context.Context, req *v1.WithdrawReq) (res *v1.WithdrawRes, err error)
}
