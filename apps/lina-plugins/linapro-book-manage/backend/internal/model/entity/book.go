// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Book is the golang structure for table book.
type Book struct {
	Id                int64      `json:"id"                orm:"id"                 description:"Book ID"`
	TenantId          int        `json:"tenantId"          orm:"tenant_id"          description:"Owning tenant ID, 0 means PLATFORM"`
	Title             string     `json:"title"             orm:"title"              description:"Book title"`
	Author            string     `json:"author"            orm:"author"             description:"Author"`
	Isbn              string     `json:"isbn"              orm:"isbn"               description:"ISBN, unique per tenant when filled"`
	Category          int        `json:"category"          orm:"category"           description:"Category: 1=technology, 2=literature, 3=management, 4=children, 5=other"`
	Publisher         string     `json:"publisher"         orm:"publisher"          description:"Publisher"`
	PublishDate       time.Time  `json:"publishDate"       orm:"publish_date"       description:"Publish date with date-only semantics"`
	TotalQuantity     int        `json:"totalQuantity"     orm:"total_quantity"     description:"Total copy count"`
	AvailableQuantity int        `json:"availableQuantity" orm:"available_quantity" description:"Currently available copy count"`
	Location          string     `json:"location"          orm:"location"           description:"Shelf location"`
	Status            int        `json:"status"            orm:"status"             description:"Book status: 1=on shelf, 2=off shelf"`
	CoverUrl          string     `json:"coverUrl"          orm:"cover_url"          description:"Cover image URL address"`
	Remark            string     `json:"remark"            orm:"remark"             description:"Remark"`
	CreatedBy         int64      `json:"createdBy"         orm:"created_by"         description:"Creator"`
	UpdatedBy         int64      `json:"updatedBy"         orm:"updated_by"         description:"Updater"`
	CreatedAt         *time.Time `json:"createdAt"         orm:"created_at"         description:"Creation time"`
	UpdatedAt         *time.Time `json:"updatedAt"         orm:"updated_at"         description:"Update time"`
	DeletedAt         *time.Time `json:"deletedAt"         orm:"deleted_at"         description:"Deletion time"`
}
