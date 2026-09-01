// borrow_v1_rest.go implements the borrow controller methods (list, get,
// create, return, update, delete).

package borrow

import (
	"context"

	v1 "lina-plugin-linapro-book-manage/backend/api/borrow/v1"
	borrowsvc "lina-plugin-linapro-book-manage/backend/internal/service/borrow"
)

// List queries borrow records
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.borrowSvc.List(ctx, borrowsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		BookId:   req.BookId,
		Borrower: req.Borrower,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.BorrowItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIBorrowItem(item)
		items = append(items, &projection)
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}

// Get returns borrow record details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.borrowSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	projection := toAPIBorrowItem(&borrowsvc.BorrowItem{BorrowEntity: item})
	return &v1.GetRes{BorrowItem: projection}, nil
}

// Create borrows a book
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.borrowSvc.Create(ctx, borrowsvc.CreateInput{
		BookId:     req.BookId,
		Borrower:   req.Borrower,
		BorrowDate: req.BorrowDate,
		DueDate:    req.DueDate,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

// Return returns a borrowed book
func (c *ControllerV1) Return(ctx context.Context, req *v1.ReturnReq) (res *v1.ReturnRes, err error) {
	err = c.borrowSvc.Return(ctx, borrowsvc.ReturnInput{
		Id:         req.Id,
		ReturnDate: req.ReturnDate,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ReturnRes{}, nil
}

// Update updates a borrow record
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.borrowSvc.Update(ctx, borrowsvc.UpdateInput{
		Id:       req.Id,
		Borrower: req.Borrower,
		DueDate:  req.DueDate,
		Remark:   req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Delete deletes borrow records by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.borrowSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
