// book_v1_options.go implements the controller method that serves the
// bounded book candidate endpoint.

package book

import (
	"context"

	v1 "lina-plugin-linapro-book-manage/backend/api/book/v1"
	booksvc "lina-plugin-linapro-book-manage/backend/internal/service/book"
)

// Options returns bounded book candidates
func (c *ControllerV1) Options(ctx context.Context, req *v1.OptionsReq) (res *v1.OptionsRes, err error) {
	items, err := c.bookSvc.Options(ctx, booksvc.OptionsInput{
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*v1.OptionItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.OptionItem{
			Id:                item.Id,
			Title:             item.Title,
			Author:            item.Author,
			AvailableQuantity: item.AvailableQuantity,
		})
	}
	return &v1.OptionsRes{List: list}, nil
}
