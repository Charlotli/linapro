// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LiveDao is the data access object for the table plugin_linapro_live_manage_live.
type LiveDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LiveColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LiveColumns defines and stores column names for the table plugin_linapro_live_manage_live.
type LiveColumns struct {
	Id               string // Live content ID
	TenantId         string // Owning tenant ID, 0 means PLATFORM
	RoomId           string // Owning live room ID within the same tenant
	Title            string // Live title shown on list and cards
	Content          string // Live description or rich content body
	LiveDate         string // Calendar date of the live event for list sorting and filtering
	PageUrl          string // External page link of the live event
	PushUrl          string // Stream push URL for the publisher
	LiveUrl          string // Play URL for viewers, such as HLS or RTMP
	CoverUrl         string // Cover image URL for cards and previews
	SongName         string // Featured song or hymn name
	SongList         string // Song list as a JSON array text, e.g. [{"name":"Song A","singer":"Alice","order":1}]
	LeadSinger       string // Lead singer or worship team
	Accompaniment    string // Accompaniment info such as band or format
	Host             string // Host name or role
	SermonTitle      string // Sermon title when the live has a standalone sermon item
	Preacher         string // Preacher name
	PreacherIdentity string // Preacher identity or title, such as pastor
	ScriptureRef     string // Scripture reference, e.g. John 3:16
	ScriptureContent string // Scripture content for display
	Outline          string // Sermon outline or agenda for preview
	DeviceInfo       string // Device info free text, e.g. camera and microphone
	Reception        string // Reception arrangement or contact info
	State            string // Live state: 0=not started, 1=ongoing, 2=finished
	IsPublic         string // Visibility: 1=public, 0=private
	StartTime        string // Actual start time of the live for schedule and reminders
	CreatedBy        string // Creator
	UpdatedBy        string // Updater
	CreatedAt        string // Creation time
	UpdatedAt        string // Update time
	DeletedAt        string // Deletion time
}

// liveColumns holds the columns for the table plugin_linapro_live_manage_live.
var liveColumns = LiveColumns{
	Id:               "id",
	TenantId:         "tenant_id",
	RoomId:           "room_id",
	Title:            "title",
	Content:          "content",
	LiveDate:         "live_date",
	PageUrl:          "page_url",
	PushUrl:          "push_url",
	LiveUrl:          "live_url",
	CoverUrl:         "cover_url",
	SongName:         "song_name",
	SongList:         "song_list",
	LeadSinger:       "lead_singer",
	Accompaniment:    "accompaniment",
	Host:             "host",
	SermonTitle:      "sermon_title",
	Preacher:         "preacher",
	PreacherIdentity: "preacher_identity",
	ScriptureRef:     "scripture_ref",
	ScriptureContent: "scripture_content",
	Outline:          "outline",
	DeviceInfo:       "device_info",
	Reception:        "reception",
	State:            "state",
	IsPublic:         "is_public",
	StartTime:        "start_time",
	CreatedBy:        "created_by",
	UpdatedBy:        "updated_by",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewLiveDao creates and returns a new DAO object for table data access.
func NewLiveDao(handlers ...gdb.ModelHandler) *LiveDao {
	return &LiveDao{
		group:    "default",
		table:    "plugin_linapro_live_manage_live",
		columns:  liveColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LiveDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LiveDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LiveDao) Columns() LiveColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LiveDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LiveDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LiveDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
