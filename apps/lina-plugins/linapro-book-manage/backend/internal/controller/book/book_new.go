// This file wires the linapro-book-manage book controller and shared response
// mappers.

package book

import (
	"time"

	"lina-core/pkg/apitime"
	bookapi "lina-plugin-linapro-book-manage/backend/api/book"
	v1 "lina-plugin-linapro-book-manage/backend/api/book/v1"
	booksvc "lina-plugin-linapro-book-manage/backend/internal/service/book"
)

// dateLayout is the date-only wire format of date-only fields.
const dateLayout = "2006-01-02"

// ControllerV1 is the book controller.
type ControllerV1 struct {
	bookSvc booksvc.Service // book service
}

// NewV1 creates and returns a new book controller instance.
func NewV1(bookSvc booksvc.Service) bookapi.IBookV1 {
	return &ControllerV1{bookSvc: bookSvc}
}

// toAPIBookItem converts the service-layer book entity into the API DTO
// projection returned through HTTP responses.
func toAPIBookItem(e *booksvc.BookEntity) v1.BookItem {
	if e == nil {
		return v1.BookItem{}
	}
	return v1.BookItem{
		Id:                e.Id,
		Title:             e.Title,
		Author:            e.Author,
		Isbn:              e.Isbn,
		Category:          e.Category,
		Publisher:         e.Publisher,
		PublishDate:       formatDate(e.PublishDate),
		TotalQuantity:     e.TotalQuantity,
		AvailableQuantity: e.AvailableQuantity,
		Location:          e.Location,
		Status:            e.Status,
		CoverUrl:          e.CoverUrl,
		Remark:            e.Remark,
		CreatedAt:         apitime.Milli(e.CreatedAt),
		UpdatedAt:         apitime.Milli(e.UpdatedAt),
	}
}

// formatDate renders the stored date in the date-only wire format.
func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(dateLayout)
}
