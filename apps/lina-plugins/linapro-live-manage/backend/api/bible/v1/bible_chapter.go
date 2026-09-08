// This file declares the bible chapter DTOs used by the
// linapro-live-manage source plugin. It intentionally carries no permission
// tag: the endpoint binds inside the anonymous viewer route group and reads
// global static content.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Bible Chapter API

// ChapterReq defines the anonymous request for one bible chapter.
type ChapterReq struct {
	g.Meta   `path:"/bible/chapter" method:"get" tags:"LiveViewer" summary:"Get one bible chapter" dc:"Return the verse list of one bible chapter ordered by verse number. Unknown books and chapters beyond the book chapter count answer with an empty list so probing cannot distinguish data gaps."`
	VolumeSn int `json:"volumeSn" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Book sequence number from 1 to 66" eg:"1"`
	Chapter  int `json:"chapter" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Chapter number within the book, 1-based" eg:"1"`
}

// ChapterRes Bible chapter response
type ChapterRes struct {
	List    []*VerseItem `json:"list" dc:"Verses of the chapter ordered by verse number" eg:"[]"`
	Book    string       `json:"book" dc:"Book full name for the reader header; empty for unknown books" eg:"Genesis"`
	Chapter int          `json:"chapter" dc:"Chapter number echoed for the reader header" eg:"1"`
}

// VerseItem defines one verse of the chapter payload.
type VerseItem struct {
	VerseSn int    `json:"verseSn" dc:"Verse number within the chapter, 1-based" eg:"1"`
	Lection string `json:"lection" dc:"Verse text; empty when the source marks the verse without text" eg:"In the beginning God created the heaven and the earth."`
}
