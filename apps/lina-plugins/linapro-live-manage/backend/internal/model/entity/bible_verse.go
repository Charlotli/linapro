// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// BibleVerse is the golang structure for table bible_verse.
type BibleVerse struct {
	Id         int64  `json:"id"         orm:"id"          description:"Surrogate identity, not part of the source data"`
	VolumeSn   int    `json:"volumeSn"   orm:"volume_sn"   description:"Book sequence number referencing the book catalog"`
	ChapterSn  int    `json:"chapterSn"  orm:"chapter_sn"  description:"Chapter number within the book, 1-based"`
	VerseSn    int    `json:"verseSn"    orm:"verse_sn"    description:"Verse number within the chapter, 1-based"`
	Lection    string `json:"lection"    orm:"lection"     description:"Verse text; null or empty when the source marks the verse without text"`
	SoundBegin int    `json:"soundBegin" orm:"sound_begin" description:"Audio offset start in milliseconds reserved for future reading playback"`
	SoundEnd   int    `json:"soundEnd"   orm:"sound_end"   description:"Audio offset end in milliseconds reserved for future reading playback"`
}
