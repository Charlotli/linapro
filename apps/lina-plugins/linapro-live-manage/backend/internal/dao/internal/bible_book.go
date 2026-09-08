// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BibleBookDao is the data access object for the table plugin_linapro_live_manage_bible_book.
type BibleBookDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  BibleBookColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// BibleBookColumns defines and stores column names for the table plugin_linapro_live_manage_bible_book.
type BibleBookColumns struct {
	Sn           string // Stable book sequence number 1..66 from the source catalogue
	KindSn       string // Source book group sequence (law/history/poetry/prophets/gospels/epistles and so on)
	ChapterCount string // Total chapter count of the book
	Testament    string // Testament: 1=old testament, 2=new testament (source 0/1 normalized)
	Pinyin       string // Book name pinyin abbreviation from the source catalogue
	ShortName    string // Book short name shown in selectors
	FullName     string // Book full name shown in the reader header
}

// bibleBookColumns holds the columns for the table plugin_linapro_live_manage_bible_book.
var bibleBookColumns = BibleBookColumns{
	Sn:           "sn",
	KindSn:       "kind_sn",
	ChapterCount: "chapter_count",
	Testament:    "testament",
	Pinyin:       "pinyin",
	ShortName:    "short_name",
	FullName:     "full_name",
}

// NewBibleBookDao creates and returns a new DAO object for table data access.
func NewBibleBookDao(handlers ...gdb.ModelHandler) *BibleBookDao {
	return &BibleBookDao{
		group:    "default",
		table:    "plugin_linapro_live_manage_bible_book",
		columns:  bibleBookColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BibleBookDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BibleBookDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BibleBookDao) Columns() BibleBookColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BibleBookDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BibleBookDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BibleBookDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
