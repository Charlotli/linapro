// book_impl.go implements tenant-scoped book CRUD with ISBN uniqueness,
// date-only parsing, referenced-book delete protection, and lent-count aware
// total quantity adjustments. Tenant filters are applied at the database
// query stage.

package book

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-book-manage/backend/internal/dao"
	"lina-plugin-linapro-book-manage/backend/internal/model/do"
)

// dateFormat is the date-only wire format of date fields.
const dateFormat = "2006-01-02"

// List queries books with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	columns := dao.Book.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "")
	if keyword := strings.TrimSpace(in.Title); keyword != "" {
		m = m.WhereOrLike(columns.Title, "%"+keyword+"%").
			WhereOrLike(columns.Author, "%"+keyword+"%")
	}
	if in.Category > 0 {
		m = m.Where(columns.Category, in.Category)
	}
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*BookEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return &ListOutput{List: list, Total: total}, nil
}

// GetById retrieves one book by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*BookEntity, error) {
	columns := dao.Book.Columns()

	var record *BookEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(columns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeBookNotFound)
	}
	return record, nil
}

// Create creates a new book. Available quantity defaults to the total
// quantity.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !Category(in.Category) {
		return 0, bizerr.NewCode(CodeBookCategoryInvalid)
	}
	status := BookStatusOnShelf
	publishDate, err := parseDate(in.PublishDate)
	if err != nil {
		return 0, err
	}
	if in.TotalQuantity <= 0 {
		return 0, bizerr.NewCode(CodeBookQuantityInvalid)
	}
	if err = s.ensureIsbnFree(ctx, in.Isbn, 0); err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// GoFrame auto-fills created_at and updated_at.
	return dao.Book.Ctx(ctx).Data(do.Book{
		TenantId:          tenantID,
		Title:             in.Title,
		Author:            in.Author,
		Isbn:              strings.TrimSpace(in.Isbn),
		Category:          in.Category,
		Publisher:         in.Publisher,
		PublishDate:       publishDate,
		TotalQuantity:     in.TotalQuantity,
		AvailableQuantity: in.TotalQuantity,
		Location:          in.Location,
		Status:            status,
		CoverUrl:          in.CoverUrl,
		Remark:            in.Remark,
		CreatedBy:         createdBy,
		UpdatedBy:         createdBy,
	}).InsertAndGetId()
}

// Update updates book information.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	columns := dao.Book.Columns()

	existing := &BookEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(columns.Id, in.Id).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 {
		return bizerr.NewCode(CodeBookNotFound)
	}

	if in.Category != nil && !Category(*in.Category) {
		return bizerr.NewCode(CodeBookCategoryInvalid)
	}
	if in.Status != nil && !BookStatus(*in.Status) {
		return bizerr.NewCode(CodeBookStatusInvalid)
	}
	var publishDate *time.Time
	if in.PublishDate != nil {
		parsed, parseErr := parseDate(*in.PublishDate)
		if parseErr != nil {
			return parseErr
		}
		publishDate = parsed
	}
	if in.Isbn != nil {
		if err = s.ensureIsbnFree(ctx, *in.Isbn, in.Id); err != nil {
			return err
		}
	}
	// The total count cannot go below the copies currently lent out;
	// otherwise the available count would go negative.
	lentOut := existing.TotalQuantity - existing.AvailableQuantity
	if in.TotalQuantity != nil && *in.TotalQuantity < lentOut {
		return bizerr.NewCode(CodeBookQuantityInvalid)
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Book{UpdatedBy: updatedBy}
	if in.Title != nil {
		data.Title = *in.Title
	}
	if in.Author != nil {
		data.Author = *in.Author
	}
	if in.Isbn != nil {
		data.Isbn = strings.TrimSpace(*in.Isbn)
	}
	if in.Category != nil {
		data.Category = *in.Category
	}
	if in.Publisher != nil {
		data.Publisher = *in.Publisher
	}
	if publishDate != nil {
		data.PublishDate = *publishDate
	}
	if in.TotalQuantity != nil {
		data.TotalQuantity = *in.TotalQuantity
		// Keep the lent-out copy count unchanged while shifting the total.
		data.AvailableQuantity = *in.TotalQuantity - lentOut
	}
	if in.Location != nil {
		data.Location = *in.Location
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	if in.CoverUrl != nil {
		data.CoverUrl = *in.CoverUrl
	}
	if in.Remark != nil {
		data.Remark = *in.Remark
	}

	_, err = dao.Book.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(columns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// Delete soft-deletes books by IDs with referenced-book protection.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeBookDeleteRequired)
	}

	if err := s.ensureNoneReferenced(ctx, idList); err != nil {
		return err
	}

	columns := dao.Book.Columns()
	_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		WhereIn(columns.Id, idList).
		Delete()
	return err
}

// Options returns bounded on-shelf book candidates for select controls.
func (s *serviceImpl) Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error) {
	columns := dao.Book.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "")
	if keyword := strings.TrimSpace(in.Keyword); keyword != "" {
		m = m.WhereOrLike(columns.Title, "%"+keyword+"%").
			WhereOrLike(columns.Author, "%"+keyword+"%")
	}
	m = m.Where(columns.Status, BookStatusOnShelf)

	list := make([]*BookEntity, 0, OptionsLimit)
	err := m.Fields(columns.Id, columns.Title, columns.Author, columns.AvailableQuantity).
		Limit(OptionsLimit).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]*OptionItem, 0, len(list))
	for _, record := range list {
		items = append(items, &OptionItem{
			Id:                record.Id,
			Title:             record.Title,
			Author:            record.Author,
			AvailableQuantity: record.AvailableQuantity,
		})
	}
	return items, nil
}

// ensureIsbnFree verifies the tenant-unique ISBN excluding one book. Empty
// ISBN values skip the check (anonymous copies share no ISBN).
func (s *serviceImpl) ensureIsbnFree(ctx context.Context, isbn string, excludeID int64) error {
	trimmed := strings.TrimSpace(isbn)
	if trimmed == "" {
		return nil
	}
	columns := dao.Book.Columns()
	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(columns.Isbn, trimmed)
	if excludeID > 0 {
		m = m.WhereNot(columns.Id, excludeID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeBookIsbnExists)
	}
	return nil
}

// ensureNoneReferenced verifies no book in the batch is still referenced by
// borrow records. One bounded query covers the whole batch.
func (s *serviceImpl) ensureNoneReferenced(ctx context.Context, bookIDs []int64) error {
	borrowColumns := dao.Borrow.Columns()
	count, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "").
		WhereIn(borrowColumns.BookId, bookIDs).
		Limit(1).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeBookReferenced)
	}
	return nil
}

// parseDate converts the date-only wire format into a stored date value.
// Empty input yields nil so the column keeps its default.
func parseDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation(dateFormat, trimmed, time.Local)
	if err != nil {
		return nil, bizerr.NewCode(CodeBookDateInvalid)
	}
	return &parsed, nil
}

// normalizePositiveInt64IDs drops non-positive identifiers from a batch list.
func normalizePositiveInt64IDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			result = append(result, id)
		}
	}
	return result
}

// pluginTableFilter returns the tenant table-filter slice from the injected tenant service.
func (s *serviceImpl) pluginTableFilter() tenantcap.FilterService {
	if s == nil || s.tenantSvc == nil {
		return nil
	}
	return s.tenantSvc.Filter()
}

// tenantFilterContext returns current tenant metadata for write ownership fields.
func (s *serviceImpl) tenantFilterContext(ctx context.Context) tenantcap.TenantFilterContext {
	if filter := s.pluginTableFilter(); filter != nil {
		return filter.Context(ctx)
	}
	return tenantcap.TenantFilterContext{}
}

// resolveUserNames resolves display names through one user-domain batch call.
// Kept as a seam for future borrower projections.
func (s *serviceImpl) resolveUserNames(ctx context.Context, ids []usercap.UserID) (map[int64]string, error) {
	names := make(map[int64]string)
	if len(ids) == 0 {
		return names, nil
	}
	if s.userSvc == nil {
		return nil, gerror.New("linapro-book-manage requires host user capability")
	}
	result, err := s.userSvc.BatchGet(ctx, ids)
	if err != nil || result == nil {
		return names, err
	}
	for id, projection := range result.Items {
		storageID, ok := userDomainID(id)
		if !ok {
			continue
		}
		names[storageID] = userDisplayName(projection)
	}
	return names, nil
}

// userDomainID parses the host user-domain ID encoding used by plugin columns.
func userDomainID(id usercap.UserID) (int64, bool) {
	storageID, err := strconv.ParseInt(strings.TrimSpace(string(id)), 10, 64)
	return storageID, err == nil && storageID > 0
}

// userDisplayName chooses the stable user display field from the projection.
func userDisplayName(user *usercap.UserInfo) string {
	if user == nil {
		return ""
	}
	if user.Username != "" {
		return user.Username
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.Label != "" {
		return user.Label
	}
	return string(user.ID)
}
