// book_v1_rest.go implements the book controller methods (list, get, create,
// update, delete).

package book

import (
	"context"

	v1 "lina-plugin-linapro-book-manage/backend/api/book/v1"
	booksvc "lina-plugin-linapro-book-manage/backend/internal/service/book"
)

// List queries books
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.bookSvc.List(ctx, booksvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Title:    req.Title,
		Category: req.Category,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.BookItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIBookItem(item)
		items = append(items, &projection)
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}

// Get returns book details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.bookSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	projection := toAPIBookItem(item)
	return &v1.GetRes{BookItem: projection}, nil
}

// Create creates a book
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	id, err := c.bookSvc.Create(ctx, booksvc.CreateInput{
		Title:         req.Title,
		Author:        req.Author,
		Isbn:          req.Isbn,
		Category:      req.Category,
		Publisher:     req.Publisher,
		PublishDate:   req.PublishDate,
		TotalQuantity: req.TotalQuantity,
		Location:      req.Location,
		Status:        status,
		CoverUrl:      req.CoverUrl,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

// Update updates a book
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.bookSvc.Update(ctx, booksvc.UpdateInput{
		Id:            req.Id,
		Title:         req.Title,
		Author:        req.Author,
		Isbn:          req.Isbn,
		Category:      req.Category,
		Publisher:     req.Publisher,
		PublishDate:   req.PublishDate,
		TotalQuantity: req.TotalQuantity,
		Location:      req.Location,
		Status:        req.Status,
		CoverUrl:      req.CoverUrl,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Delete deletes books by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.bookSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
