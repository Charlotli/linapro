// This file defines shared book response DTOs for the linapro-book-manage API.
package v1

// BookItem exposes book fields visible to callers.
type BookItem struct {
	Id                int64  `json:"id" dc:"Book ID" eg:"1"`
	Title             string `json:"title" dc:"Book title" eg:"The history of science"`
	Author            string `json:"author" dc:"Author" eg:"Wu Guosheng"`
	Isbn              string `json:"isbn" dc:"ISBN, unique per tenant when filled" eg:"9787301281"`
	Category          int    `json:"category" dc:"Category: 1=technology, 2=literature, 3=management, 4=children, 5=other" eg:"1"`
	Publisher         string `json:"publisher" dc:"Publisher" eg:"Peking University Press"`
	PublishDate       string `json:"publishDate" dc:"Publish date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2018-01-01"`
	TotalQuantity     int    `json:"totalQuantity" dc:"Total copy count" eg:"3"`
	AvailableQuantity int    `json:"availableQuantity" dc:"Currently available copy count" eg:"2"`
	Location          string `json:"location" dc:"Shelf location" eg:"Shelf A-01"`
	Status            int    `json:"status" dc:"Book status: 1=on shelf, 2=off shelf" eg:"1"`
	CoverUrl          string `json:"coverUrl" dc:"Cover image URL address" eg:"https://example.com/cover.jpg"`
	Remark            string `json:"remark" dc:"Remark" eg:""`
	CreatedAt         *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt         *int64 `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
