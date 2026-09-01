// This file defines linapro-book-manage borrow business error codes.

package borrow

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeBorrowNotFound reports that the requested borrow record does not
	// exist in the current tenant.
	CodeBorrowNotFound = bizerr.MustDefine(
		"BOOK_BORROW_NOT_FOUND",
		"Borrow record does not exist",
		gcode.CodeNotFound,
	)
	// CodeBookNoAvailable reports that the book is off shelf or has no
	// available copies left.
	CodeBookNoAvailable = bizerr.MustDefine(
		"BOOK_NO_AVAILABLE",
		"Book has no available copies",
		gcode.CodeInvalidOperation,
	)
	// CodeBorrowAlreadyReturned reports that the borrow record is already in
	// returned status.
	CodeBorrowAlreadyReturned = bizerr.MustDefine(
		"BOOK_BORROW_ALREADY_RETURNED",
		"Borrow record is already returned",
		gcode.CodeInvalidOperation,
	)
	// CodeBorrowReturnNoAuthority reports that the caller is not the record
	// creator required by the return action.
	CodeBorrowReturnNoAuthority = bizerr.MustDefine(
		"BOOK_BORROW_RETURN_NO_AUTHORITY",
		"Only the record creator can return the book",
		gcode.CodeNotAuthorized,
	)
	// CodeBorrowDateInvalid reports a malformed date value.
	CodeBorrowDateInvalid = bizerr.MustDefine(
		"BOOK_BORROW_DATE_INVALID",
		"Date must be in YYYY-MM-DD format",
		gcode.CodeInvalidParameter,
	)
	// CodeBorrowDeleteRequired reports that a delete operation received no
	// record IDs.
	CodeBorrowDeleteRequired = bizerr.MustDefine(
		"BOOK_BORROW_DELETE_REQUIRED",
		"Select at least one borrow record to delete",
		gcode.CodeInvalidParameter,
	)
	// CodeBorrowDeleteBorrowedRejected reports that a delete operation
	// targeted records still in borrowed status.
	CodeBorrowDeleteBorrowedRejected = bizerr.MustDefine(
		"BOOK_BORROW_DELETE_BORROWED_REJECTED",
		"Return the book before deleting the borrow record",
		gcode.CodeInvalidOperation,
	)
)
