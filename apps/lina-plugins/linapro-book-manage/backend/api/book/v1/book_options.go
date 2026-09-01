// This file declares the book options request/response DTOs used by the
// linapro-book-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Book Options API

// OptionsReq defines the request for listing bounded book candidates.
type OptionsReq struct {
	g.Meta  `path:"/book/options" method:"get" tags:"Books" summary:"Get book options" dc:"Return the minimal book projection (id, title, author, availableQuantity) for select controls, filtered by keyword across title and author. Only on-shelf books are returned. At most 200 items are returned; excess results are not paged." permission:"book:query"`
	Keyword string `json:"keyword" dc:"Filter by title or author (fuzzy match); empty means all on-shelf books" eg:"Science"`
}

// OptionsRes Book options response
type OptionsRes struct {
	List []*OptionItem `json:"list" dc:"Book option list" eg:"[]"`
}

// OptionItem Book minimal option projection
type OptionItem struct {
	Id                int64  `json:"id" dc:"Book ID" eg:"1"`
	Title             string `json:"title" dc:"Book title" eg:"The history of science"`
	Author            string `json:"author" dc:"Author" eg:"Wu Guosheng"`
	AvailableQuantity int    `json:"availableQuantity" dc:"Currently available copy count" eg:"3"`
}
