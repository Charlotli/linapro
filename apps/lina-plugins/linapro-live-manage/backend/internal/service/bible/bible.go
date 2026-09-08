// Package bible implements the public bible reading service for the
// linapro-live-manage source plugin. The book catalog and verse text are
// global static content seeded by the plugin SQL, so the queries carry no
// tenant scoping by design; every endpoint stays anonymous like the other
// viewer endpoints and answers out-of-range requests with empty payloads
// instead of errors, as documented in the OpenSpec change
// add-live-announcement-and-bible.
package bible

import "context"

// BooksLimit is the expected catalog size; the constant doubles as a query
// bound should the seed ever grow unexpectedly.
const BooksLimit = 100

// Service defines the anonymous bible reading contract.
type Service interface {
	// Books returns the whole book catalog ordered by the stable book
	// sequence number: 66 entries split across the two testaments. Database
	// access is one bounded query.
	Books(ctx context.Context) ([]*Book, error)
	// Chapter returns the verse list of one chapter ordered by verse number
	// together with the book display name for the reader header. An unknown
	// book sequence or a chapter beyond the book chapter count answers with
	// an empty verse list and empty book name, so probing cannot distinguish
	// data gaps. Database access is two point/range queries on indexed keys.
	Chapter(ctx context.Context, in ChapterInput) (*ChapterResult, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct{}

// New creates and returns a new Service instance. The service holds no
// runtime dependencies: both endpoints are static catalog reads.
func New() Service {
	return &serviceImpl{}
}

// Testament enumeration values stored in the book catalog.
const (
	// TestamentOld marks old-testament books.
	TestamentOld = 1
	// TestamentNew marks new-testament books.
	TestamentNew = 2
)

// Book defines one catalog entry of the public book list.
type Book struct {
	Sn           int    // Stable book sequence number 1..66
	ShortName    string // Book short name shown in selectors
	FullName     string // Book full name shown in the reader header
	ChapterCount int    // Total chapter count of the book
	Testament    int    // Testament: 1=old testament, 2=new testament
}

// Verse defines one verse of the public chapter payload.
type Verse struct {
	VerseSn int    // Verse number within the chapter, 1-based
	Lection string // Verse text; empty when the source marks the verse without text
}

// ChapterResult defines the chapter payload with its reader-header metadata.
type ChapterResult struct {
	Book    string   // Book full name for the reader header; empty for unknown books
	Chapter int      // Chapter number echoed for the reader header
	Verses  []*Verse // Verses of the chapter ordered by verse number
}

// ChapterInput defines input for Chapter function.
type ChapterInput struct {
	VolumeSn int // Book sequence number 1..66
	Chapter  int // Chapter number within the book, 1-based
}
