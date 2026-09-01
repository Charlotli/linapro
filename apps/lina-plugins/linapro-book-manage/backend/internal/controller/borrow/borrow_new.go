// This file wires the linapro-book-manage borrow controller and shared
// response mappers.

package borrow

import (
	"time"

	"lina-core/pkg/apitime"
	borrowapi "lina-plugin-linapro-book-manage/backend/api/borrow"
	v1 "lina-plugin-linapro-book-manage/backend/api/borrow/v1"
	borrowsvc "lina-plugin-linapro-book-manage/backend/internal/service/borrow"
)

// dateLayout is the date-only wire format of date-only fields.
const dateLayout = "2006-01-02"

// ControllerV1 is the borrow controller.
type ControllerV1 struct {
	borrowSvc borrowsvc.Service // borrow service
}

// NewV1 creates and returns a new borrow controller instance.
func NewV1(borrowSvc borrowsvc.Service) borrowapi.IBorrowV1 {
	return &ControllerV1{borrowSvc: borrowSvc}
}

// toAPIBorrowItem converts the service-layer list item into the API DTO
// projection returned through HTTP responses.
func toAPIBorrowItem(item *borrowsvc.BorrowItem) v1.BorrowItem {
	if item == nil || item.BorrowEntity == nil {
		return v1.BorrowItem{}
	}
	return v1.BorrowItem{
		Id:         item.Id,
		BookId:     item.BookId,
		BookTitle:  item.BookTitle,
		Borrower:   item.Borrower,
		BorrowDate: formatDate(item.BorrowDate),
		DueDate:    formatDate(item.DueDate),
		ReturnDate: formatDate(item.ReturnDate),
		Status:     item.Status,
		Remark:     item.Remark,
		CreatedAt:  apitime.Milli(item.CreatedAt),
		UpdatedAt:  apitime.Milli(item.UpdatedAt),
	}
}

// formatDate renders the stored date in the date-only wire format.
func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(dateLayout)
}
