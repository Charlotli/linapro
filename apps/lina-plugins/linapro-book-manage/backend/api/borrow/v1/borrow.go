// This file defines shared borrow response DTOs for the linapro-book-manage API.
package v1

// BorrowItem exposes borrow record fields visible to callers.
type BorrowItem struct {
	Id         int64  `json:"id" dc:"Borrow record ID" eg:"1"`
	BookId     int64  `json:"bookId" dc:"Borrowed book ID within the same tenant" eg:"1"`
	BookTitle  string `json:"bookTitle" dc:"Book title resolved in batch" eg:"The history of science"`
	Borrower   string `json:"borrower" dc:"Borrower name" eg:"Mary"`
	BorrowDate string `json:"borrowDate" dc:"Borrow date in YYYY-MM-DD format with date-only semantics" eg:"2026-04-22"`
	DueDate    string `json:"dueDate" dc:"Expected return date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2026-05-22"`
	ReturnDate string `json:"returnDate" dc:"Actual return date in YYYY-MM-DD format with date-only semantics; empty means not returned" eg:""`
	Status     int    `json:"status" dc:"Borrow status: 1=borrowed, 2=returned" eg:"1"`
	Remark     string `json:"remark" dc:"Remark" eg:""`
	CreatedAt  *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt  *int64 `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
