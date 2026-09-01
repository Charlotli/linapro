// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Book is the golang structure of table plugin_linapro_book_manage_book for DAO operations like Where/Data.
type Book struct {
	g.Meta            `orm:"table:plugin_linapro_book_manage_book, do:true"`
	Id                any        // Book ID
	TenantId          any        // Owning tenant ID, 0 means PLATFORM
	Title             any        // Book title
	Author            any        // Author
	Isbn              any        // ISBN, unique per tenant when filled
	Category          any        // Category: 1=technology, 2=literature, 3=management, 4=children, 5=other
	Publisher         any        // Publisher
	PublishDate       any        // Publish date with date-only semantics
	TotalQuantity     any        // Total copy count
	AvailableQuantity any        // Currently available copy count
	Location          any        // Shelf location
	Status            any        // Book status: 1=on shelf, 2=off shelf
	CoverUrl          any        // Cover image URL address
	Remark            any        // Remark
	CreatedBy         any        // Creator
	UpdatedBy         any        // Updater
	CreatedAt         *time.Time // Creation time
	UpdatedAt         *time.Time // Update time
	DeletedAt         *time.Time // Deletion time
}
