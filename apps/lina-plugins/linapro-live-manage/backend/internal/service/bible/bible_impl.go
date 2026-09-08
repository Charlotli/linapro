// bible_impl.go implements the bible reading service: one bounded catalog
// query and one chapter range query. Out-of-range volume or chapter numbers
// answer with an empty list, keeping the anonymous probe surface flat.

package bible

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
)

// Books returns the whole book catalog ordered by sequence number.
func (s *serviceImpl) Books(ctx context.Context) ([]*Book, error) {
	cols := dao.BibleBook.Columns()
	var rows []*entity.BibleBook
	err := dao.BibleBook.Ctx(ctx).
		Fields(cols.Sn, cols.ShortName, cols.FullName, cols.ChapterCount, cols.Testament).
		OrderAsc(cols.Sn).
		Limit(BooksLimit).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "list bible books")
	}
	books := make([]*Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, &Book{
			Sn:           row.Sn,
			ShortName:    row.ShortName,
			FullName:     row.FullName,
			ChapterCount: row.ChapterCount,
			Testament:    row.Testament,
		})
	}
	return books, nil
}

// Chapter returns the verse list of one chapter with reader-header metadata.
func (s *serviceImpl) Chapter(ctx context.Context, in ChapterInput) (*ChapterResult, error) {
	if in.VolumeSn <= 0 || in.Chapter <= 0 {
		return &ChapterResult{Verses: []*Verse{}}, nil
	}
	cols := dao.BibleBook.Columns()
	var book *entity.BibleBook
	err := dao.BibleBook.Ctx(ctx).
		Fields(cols.Sn, cols.FullName, cols.ChapterCount).
		Where(cols.Sn, in.VolumeSn).
		Scan(&book)
	if err != nil {
		return nil, gerror.Wrap(err, "load bible book for chapter")
	}
	// Unknown volume or an out-of-range chapter both answer with an empty
	// payload: the response cannot distinguish catalog gaps from blanks.
	if book == nil || in.Chapter > book.ChapterCount {
		return &ChapterResult{Verses: []*Verse{}}, nil
	}

	verseCols := dao.BibleVerse.Columns()
	var verseRows []*entity.BibleVerse
	err = dao.BibleVerse.Ctx(ctx).
		Fields(verseCols.VerseSn, verseCols.Lection).
		Where(verseCols.VolumeSn, in.VolumeSn).
		Where(verseCols.ChapterSn, in.Chapter).
		OrderAsc(verseCols.VerseSn).
		Scan(&verseRows)
	if err != nil {
		return nil, gerror.Wrap(err, "list bible verses of chapter")
	}
	verses := make([]*Verse, 0, len(verseRows))
	for _, row := range verseRows {
		verses = append(verses, &Verse{
			VerseSn: row.VerseSn,
			Lection: row.Lection,
		})
	}
	return &ChapterResult{
		Book:    book.FullName,
		Chapter: in.Chapter,
		Verses:  verses,
	}, nil
}
