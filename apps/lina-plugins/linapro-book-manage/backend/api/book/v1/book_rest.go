// This file declares the book CRUD request/response DTOs used by the
// linapro-book-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Book List API

// ListReq defines the request for listing books.
type ListReq struct {
	g.Meta   `path:"/book" method:"get" tags:"Books" summary:"Get book list" dc:"Query books by page, with filtering by title, author, category, and status." permission:"book:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	Title    string `json:"title" dc:"Filter by title or author (fuzzy match)" eg:"Science"`
	Category int    `json:"category" dc:"Filter by category: 1=technology, 2=literature, 3=management, 4=children, 5=other; 0 means no filter" eg:"1"`
	Status   *int   `json:"status" dc:"Filter by book status: 1=on shelf, 2=off shelf; omitted means no filter" eg:"1"`
}

// ListRes Book list response
type ListRes struct {
	List  []*BookItem `json:"list" dc:"Book list" eg:"[]"`
	Total int         `json:"total" dc:"Total number of items" eg:"20"`
}

// Book Create API

// CreateReq defines the request for creating one book.
type CreateReq struct {
	g.Meta        `path:"/book" method:"post" tags:"Books" summary:"Create book" dc:"Create one book owned by the current tenant. ISBN must be unique within the tenant when filled. Available quantity defaults to the total quantity." permission:"book:add"`
	Title         string `json:"title" v:"required|length:1,256#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Book title" eg:"The history of science"`
	Author        string `json:"author" dc:"Author" eg:"Wu Guosheng"`
	Isbn          string `json:"isbn" dc:"ISBN, unique per tenant when filled" eg:"9787301281"`
	Category      int    `json:"category" v:"required|in:1,2,3,4,5#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Category: 1=technology, 2=literature, 3=management, 4=children, 5=other" eg:"1"`
	Publisher     string `json:"publisher" dc:"Publisher" eg:"Peking University Press"`
	PublishDate   string `json:"publishDate" v:"date#gf.gvalid.rule.date" dc:"Publish date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2018-01-01"`
	TotalQuantity int    `json:"totalQuantity" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Total copy count" eg:"3"`
	Location      string `json:"location" dc:"Shelf location" eg:"Shelf A-01"`
	Status        *int   `json:"status" d:"1" dc:"Book status: 1=on shelf, 2=off shelf" eg:"1"`
	CoverUrl      string `json:"coverUrl" dc:"Cover image URL address" eg:"https://example.com/cover.jpg"`
	Remark        string `json:"remark" dc:"Remark" eg:""`
}

// CreateRes Book create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Book ID" eg:"1"`
}

// Book Update API

// UpdateReq defines the request for updating one book.
type UpdateReq struct {
	g.Meta        `path:"/book/{id}" method:"put" tags:"Books" summary:"Update book" dc:"Update the specified book. All fields are optional; omitted fields keep their current values. Total quantity adjustments cannot go below the number of copies currently lent out." permission:"book:edit"`
	Id            int64   `json:"id" v:"required" dc:"Book ID" eg:"1"`
	Title         *string `json:"title" v:"length:1,256#gf.gvalid.rule.length" dc:"Book title" eg:"The history of science"`
	Author        *string `json:"author" dc:"Author" eg:"Wu Guosheng"`
	Isbn          *string `json:"isbn" dc:"ISBN, unique per tenant when filled" eg:"9787301281"`
	Category      *int    `json:"category" dc:"Category: 1=technology, 2=literature, 3=management, 4=children, 5=other" eg:"1"`
	Publisher     *string `json:"publisher" dc:"Publisher" eg:"Peking University Press"`
	PublishDate   *string `json:"publishDate" v:"date#gf.gvalid.rule.date" dc:"Publish date in YYYY-MM-DD format with date-only semantics" eg:"2018-01-01"`
	TotalQuantity *int    `json:"totalQuantity" v:"min:1#gf.gvalid.rule.min" dc:"Total copy count; cannot be below the lent-out count" eg:"3"`
	Location      *string `json:"location" dc:"Shelf location" eg:"Shelf A-01"`
	Status        *int    `json:"status" dc:"Book status: 1=on shelf, 2=off shelf" eg:"1"`
	CoverUrl      *string `json:"coverUrl" dc:"Cover image URL address" eg:"https://example.com/cover.jpg"`
	Remark        *string `json:"remark" dc:"Remark" eg:""`
}

// UpdateRes Book update response
type UpdateRes struct{}

// Book Get API

// GetReq defines the request for retrieving book details.
type GetReq struct {
	g.Meta `path:"/book/{id}" method:"get" tags:"Books" summary:"Get book details" dc:"Get book details by ID. Books outside the current tenant are reported as not found." permission:"book:query"`
	Id     int64 `json:"id" v:"required" dc:"Book ID" eg:"1"`
}

// GetRes Book detail response
type GetRes struct {
	BookItem
}

// Book Delete API

// DeleteReq defines the request for deleting books.
type DeleteReq struct {
	g.Meta `path:"/book" method:"delete" tags:"Books" summary:"Delete books" dc:"Soft-delete one or more books by query array ids[]. Deleting books still referenced by borrow records is rejected." permission:"book:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Book ID list as a query array" eg:"[1,2,3]"`
}

// DeleteRes Book delete response
type DeleteRes struct{}
