// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BibleVerse is the golang structure of table plugin_linapro_live_manage_bible_verse for DAO operations like Where/Data.
type BibleVerse struct {
	g.Meta     `orm:"table:plugin_linapro_live_manage_bible_verse, do:true"`
	Id         any // Surrogate identity, not part of the source data
	VolumeSn   any // Book sequence number referencing the book catalog
	ChapterSn  any // Chapter number within the book, 1-based
	VerseSn    any // Verse number within the chapter, 1-based
	Lection    any // Verse text; null or empty when the source marks the verse without text
	SoundBegin any // Audio offset start in milliseconds reserved for future reading playback
	SoundEnd   any // Audio offset end in milliseconds reserved for future reading playback
}
