// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// BibleBook is the golang structure for table bible_book.
type BibleBook struct {
	Sn           int    `json:"sn"           orm:"sn"            description:"Stable book sequence number 1..66 from the source catalog"`
	KindSn       int    `json:"kindSn"       orm:"kind_sn"       description:"Source book group sequence (law/history/poetry/prophets/gospels/epistles and so on)"`
	ChapterCount int    `json:"chapterCount" orm:"chapter_count" description:"Total chapter count of the book"`
	Testament    int    `json:"testament"    orm:"testament"     description:"Testament: 1=old testament, 2=new testament (source 0/1 normalized)"`
	Pinyin       string `json:"pinyin"       orm:"pinyin"        description:"Book name pinyin abbreviation from the source catalog"`
	ShortName    string `json:"shortName"    orm:"short_name"    description:"Book short name shown in selectors"`
	FullName     string `json:"fullName"     orm:"full_name"     description:"Book full name shown in the reader header"`
}
