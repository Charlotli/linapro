// This file declares the borrow request/response DTOs used by the
// linapro-book-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Borrow List API

// ListReq defines the request for listing borrow records.
type ListReq struct {
	g.Meta   `path:"/borrow" method:"get" tags:"BorrowRecords" summary:"Get borrow record list" dc:"Query borrow records by page, with filtering by book, borrower, and status." permission:"borrow:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	BookId   int64  `json:"bookId" dc:"Filter by book ID; 0 means no filter" eg:"1"`
	Borrower string `json:"borrower" dc:"Filter by borrower name (fuzzy match)" eg:"Mary"`
	Status   *int   `json:"status" dc:"Filter by borrow status: 1=borrowed, 2=returned; omitted means no filter" eg:"1"`
}

// ListRes Borrow record list response
type ListRes struct {
	List  []*BorrowItem `json:"list" dc:"Borrow record list" eg:"[]"`
	Total int           `json:"total" dc:"Total number of items" eg:"5"`
}

// Borrow Create API

// CreateReq defines the request for borrowing one book.
type CreateReq struct {
	g.Meta     `path:"/borrow" method:"post" tags:"BorrowRecords" summary:"Borrow book" dc:"Borrow one book within the current tenant. The book must be on shelf with available copies; the available count is decremented in the same transaction." permission:"borrow:add"`
	BookId     int64  `json:"bookId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Borrowed book ID within the same tenant" eg:"1"`
	Borrower   string `json:"borrower" v:"required|length:1,64#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Borrower name" eg:"Mary"`
	BorrowDate string `json:"borrowDate" v:"date#gf.gvalid.rule.date" dc:"Borrow date in YYYY-MM-DD format with date-only semantics; defaults to today when omitted" eg:"2026-04-22"`
	DueDate    string `json:"dueDate" v:"date#gf.gvalid.rule.date" dc:"Expected return date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2026-05-22"`
	Remark     string `json:"remark" dc:"Remark" eg:""`
}

// CreateRes Borrow create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Borrow record ID" eg:"1"`
}

// Borrow Return API

// ReturnReq defines the request for returning one borrowed book.
type ReturnReq struct {
	g.Meta     `path:"/borrow/{id}/return" method:"put" tags:"BorrowRecords" summary:"Return book" dc:"Return one borrowed book. The record status becomes returned, the return date is recorded, and the available count is incremented in the same transaction. Only the creator of the borrow record can return it." permission:"borrow:add"`
	Id         int64  `json:"id" v:"required" dc:"Borrow record ID" eg:"1"`
	ReturnDate string `json:"returnDate" v:"date#gf.gvalid.rule.date" dc:"Actual return date in YYYY-MM-DD format with date-only semantics; defaults to today when omitted" eg:"2026-05-20"`
}

// ReturnRes Borrow return response
type ReturnRes struct{}

// Borrow Update API

// UpdateReq defines the request for updating one borrow record.
type UpdateReq struct {
	g.Meta   `path:"/borrow/{id}" method:"put" tags:"BorrowRecords" summary:"Update borrow record" dc:"Update the specified borrow record. Only borrower, due date, and remark are editable; borrow date and status are managed by borrow and return actions." permission:"borrow:edit"`
	Id       int64   `json:"id" v:"required" dc:"Borrow record ID" eg:"1"`
	Borrower *string `json:"borrower" dc:"Borrower name" eg:"Mary"`
	DueDate  *string `json:"dueDate" v:"date#gf.gvalid.rule.date" dc:"Expected return date in YYYY-MM-DD format with date-only semantics" eg:"2026-05-22"`
	Remark   *string `json:"remark" dc:"Remark" eg:""`
}

// UpdateRes Borrow update response
type UpdateRes struct{}

// Borrow Delete API

// DeleteReq defines the request for deleting borrow records.
type DeleteReq struct {
	g.Meta `path:"/borrow" method:"delete" tags:"BorrowRecords" summary:"Delete borrow records" dc:"Soft-delete one or more borrow records by query array ids[]. Deleting a record currently in borrowed status is rejected; return the book first." permission:"borrow:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Borrow record ID list as a query array" eg:"[1,2,3]"`
}

// DeleteRes Borrow delete response
type DeleteRes struct{}

// Borrow Get API

// GetReq defines the request for retrieving borrow record details.
type GetReq struct {
	g.Meta `path:"/borrow/{id}" method:"get" tags:"BorrowRecords" summary:"Get borrow record details" dc:"Get borrow record details by ID. Records outside the current tenant are reported as not found." permission:"borrow:query"`
	Id     int64 `json:"id" v:"required" dc:"Borrow record ID" eg:"1"`
}

// GetRes Borrow record detail response
type GetRes struct {
	BorrowItem
}
