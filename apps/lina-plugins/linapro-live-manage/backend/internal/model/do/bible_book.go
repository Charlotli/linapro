// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BibleBook is the golang structure of table plugin_linapro_live_manage_bible_book for DAO operations like Where/Data.
type BibleBook struct {
	g.Meta       `orm:"table:plugin_linapro_live_manage_bible_book, do:true"`
	Sn           any // Stable book sequence number 1..66 from the source catalog
	KindSn       any // Source book group sequence (law/history/poetry/prophets/gospels/epistles and so on)
	ChapterCount any // Total chapter count of the book
	Testament    any // Testament: 1=old testament, 2=new testament (source 0/1 normalized)
	Pinyin       any // Book name pinyin abbreviation from the source catalog
	ShortName    any // Book short name shown in selectors
	FullName     any // Book full name shown in the reader header
}
