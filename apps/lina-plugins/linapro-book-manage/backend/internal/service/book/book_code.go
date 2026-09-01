// This file defines linapro-book-manage book business error codes.

package book

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeBookNotFound reports that the requested book does not exist in the
	// current tenant.
	CodeBookNotFound = bizerr.MustDefine(
		"BOOK_NOT_FOUND",
		"Book does not exist",
		gcode.CodeNotFound,
	)
	// CodeBookIsbnExists reports that the ISBN is already used in the current
	// tenant.
	CodeBookIsbnExists = bizerr.MustDefine(
		"BOOK_ISBN_EXISTS",
		"Book ISBN already exists",
		gcode.CodeInvalidParameter,
	)
	// CodeBookCategoryInvalid reports an unknown book category value.
	CodeBookCategoryInvalid = bizerr.MustDefine(
		"BOOK_CATEGORY_INVALID",
		"Book category is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeBookStatusInvalid reports an unknown book status value.
	CodeBookStatusInvalid = bizerr.MustDefine(
		"BOOK_STATUS_INVALID",
		"Book status is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeBookQuantityInvalid reports an invalid copy count, such as a total
	// below the lent-out count or a non-positive total.
	CodeBookQuantityInvalid = bizerr.MustDefine(
		"BOOK_QUANTITY_INVALID",
		"Copy count is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeBookDateInvalid reports a malformed date value.
	CodeBookDateInvalid = bizerr.MustDefine(
		"BOOK_DATE_INVALID",
		"Date must be in YYYY-MM-DD format",
		gcode.CodeInvalidParameter,
	)
	// CodeBookReferenced reports that a book is still referenced by borrow
	// records and cannot be deleted.
	CodeBookReferenced = bizerr.MustDefine(
		"BOOK_REFERENCED",
		"Book is still referenced by borrow records",
		gcode.CodeInvalidOperation,
	)
	// CodeBookDeleteRequired reports that a delete operation received no IDs.
	CodeBookDeleteRequired = bizerr.MustDefine(
		"BOOK_DELETE_REQUIRED",
		"Select at least one book to delete",
		gcode.CodeInvalidParameter,
	)
)
