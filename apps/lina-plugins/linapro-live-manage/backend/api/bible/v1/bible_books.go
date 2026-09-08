// This file declares the bible book list DTOs used by the
// linapro-live-manage source plugin. It intentionally carries no permission
// tag: the endpoint binds inside the anonymous viewer route group and reads
// global static content.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Bible Books API

// BooksReq defines the anonymous request for the bible book catalog.
type BooksReq struct {
	g.Meta `path:"/bible/books" method:"get" tags:"LiveViewer" summary:"Get bible book catalog" dc:"Return all 66 bible books ordered by sequence number with short name, full name, chapter count, and testament, for the H5 bible reader book picker."`
}

// BooksRes Bible book catalog response
type BooksRes struct {
	List []*BookItem `json:"list" dc:"Book catalog ordered by sequence number" eg:"[]"`
}

// BookItem defines one book catalog entry.
type BookItem struct {
	Sn           int    `json:"sn" dc:"Stable book sequence number from 1 to 66" eg:"1"`
	ShortName    string `json:"shortName" dc:"Book short name shown in selectors" eg:"Gen"`
	FullName     string `json:"fullName" dc:"Book full name shown in the reader header" eg:"Genesis"`
	ChapterCount int    `json:"chapterCount" dc:"Total chapter count of the book" eg:"50"`
	Testament    int    `json:"testament" dc:"Testament: 1=old testament, 2=new testament" eg:"1"`
}
