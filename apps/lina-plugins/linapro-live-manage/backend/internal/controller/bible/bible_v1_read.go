// bible_v1_read.go implements the controller methods serving the anonymous
// bible endpoints consumed by the plugin H5 bible reader.

package bible

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/bible/v1"
	biblesvc "lina-plugin-linapro-live-manage/backend/internal/service/bible"
)

// Books returns the whole bible book catalog.
func (c *ControllerV1) Books(ctx context.Context, req *v1.BooksReq) (res *v1.BooksRes, err error) {
	books, err := c.bibleSvc.Books(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.BookItem, 0, len(books))
	for _, book := range books {
		list = append(list, &v1.BookItem{
			Sn:           book.Sn,
			ShortName:    book.ShortName,
			FullName:     book.FullName,
			ChapterCount: book.ChapterCount,
			Testament:    book.Testament,
		})
	}
	return &v1.BooksRes{List: list}, nil
}

// Chapter returns the verse list of one chapter with reader-header metadata.
func (c *ControllerV1) Chapter(ctx context.Context, req *v1.ChapterReq) (res *v1.ChapterRes, err error) {
	result, err := c.bibleSvc.Chapter(ctx, biblesvc.ChapterInput{
		VolumeSn: req.VolumeSn,
		Chapter:  req.Chapter,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &biblesvc.ChapterResult{Verses: []*biblesvc.Verse{}}
	}
	list := make([]*v1.VerseItem, 0, len(result.Verses))
	for _, verse := range result.Verses {
		list = append(list, &v1.VerseItem{
			VerseSn: verse.VerseSn,
			Lection: verse.Lection,
		})
	}
	return &v1.ChapterRes{
		List:    list,
		Book:    result.Book,
		Chapter: result.Chapter,
	}, nil
}
