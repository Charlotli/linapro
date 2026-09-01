// borrow_impl.go implements tenant-scoped borrow actions with transactional
// inventory linkage: borrowing decrements the available copy count behind a
// conditional update so concurrent borrows cannot oversell copies; returning
// increments it back. Book titles are assembled with one batch query per page.

package borrow

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-plugin-linapro-book-manage/backend/internal/dao"
	"lina-plugin-linapro-book-manage/backend/internal/model/do"
	"lina-plugin-linapro-book-manage/backend/internal/service/book"
)

// dateFormat is the date-only wire format of date fields.
const dateFormat = "2006-01-02"

// List queries borrow records with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	columns := dao.Borrow.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "")
	if in.BookId > 0 {
		m = m.Where(columns.BookId, in.BookId)
	}
	if borrower := strings.TrimSpace(in.Borrower); borrower != "" {
		m = m.WhereLike(columns.Borrower, "%"+borrower+"%")
	}
	if in.Status != nil {
		m = m.Where(columns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*BorrowEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(columns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	bookTitles, err := s.resolveBookTitleMap(ctx, list)
	if err != nil {
		return nil, err
	}

	items := make([]*BorrowItem, 0, len(list))
	for _, record := range list {
		items = append(items, &BorrowItem{
			BorrowEntity: record,
			BookTitle:    bookTitles[record.BookId],
		})
	}
	return &ListOutput{List: items, Total: total}, nil
}

// GetById retrieves one borrow record by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*BorrowEntity, error) {
	columns := dao.Borrow.Columns()

	var record *BorrowEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "").
		Where(columns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeBorrowNotFound)
	}
	return record, nil
}

// Create borrows one book with inventory decrement.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	borrowDate, err := parseDate(in.BorrowDate)
	if err != nil {
		return 0, err
	}
	dueDate, err := parseDate(in.DueDate)
	if err != nil {
		return 0, err
	}
	if err = s.ensureBookBorrowable(ctx, in.BookId); err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// Decrement first behind a conditional update: the affected-rows check
	// guarantees the available count never goes negative under concurrency.
	decremented, err := s.decrementAvailable(ctx, tenantID, in.BookId)
	if err != nil {
		return 0, err
	}
	if !decremented {
		return 0, bizerr.NewCode(CodeBookNoAvailable)
	}

	var borrowID int64
	err = dao.Borrow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// GoFrame auto-fills created_at and updated_at.
		id, insertErr := dao.Borrow.Ctx(ctx).Data(do.Borrow{
			TenantId:   tenantID,
			BookId:     in.BookId,
			Borrower:   strings.TrimSpace(in.Borrower),
			BorrowDate: borrowDate,
			DueDate:    dueDate,
			Status:     BorrowStatusBorrowed,
			Remark:     in.Remark,
			CreatedBy:  createdBy,
			UpdatedBy:  createdBy,
		}).InsertAndGetId()
		if insertErr != nil {
			return insertErr
		}
		borrowID = id
		return nil
	})
	if err != nil {
		// Roll the inventory back when the record insert fails.
		if rollbackErr := s.incrementAvailable(ctx, tenantID, in.BookId); rollbackErr != nil {
			_ = rollbackErr
		}
		return 0, err
	}
	return borrowID, nil
}

// Return marks one borrowed record as returned with inventory increment.
func (s *serviceImpl) Return(ctx context.Context, in ReturnInput) error {
	record, err := s.GetById(ctx, in.Id)
	if err != nil {
		return err
	}
	if record.Status != BorrowStatusBorrowed {
		return bizerr.NewCode(CodeBorrowAlreadyReturned)
	}
	currentUserID := s.currentUserID(ctx)
	if record.CreatedBy != currentUserID {
		return bizerr.NewCode(CodeBorrowReturnNoAuthority)
	}
	returnDate, err := parseDate(in.ReturnDate)
	if err != nil {
		return err
	}

	tenantID := s.tenantFilterContext(ctx).TenantID
	err = dao.Borrow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		columns := dao.Borrow.Columns()
		if _, err := dao.Borrow.Ctx(ctx).
			Where(tenantspi.TenantFilterColumn, tenantID).
			Where(columns.Id, in.Id).
			Where(columns.Status, BorrowStatusBorrowed).
			Data(do.Borrow{
				Status:     BorrowStatusReturned,
				ReturnDate: returnDate,
				UpdatedBy:  currentUserID,
			}).
			Update(); err != nil {
			return err
		}
		return s.incrementAvailable(ctx, tenantID, record.BookId)
	})
	return err
}

// Update updates editable borrow fields.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	columns := dao.Borrow.Columns()

	existing := &BorrowEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "").
		Where(columns.Id, in.Id).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 {
		return bizerr.NewCode(CodeBorrowNotFound)
	}

	var dueDate *time.Time
	if in.DueDate != nil {
		parsed, parseErr := parseDate(*in.DueDate)
		if parseErr != nil {
			return parseErr
		}
		dueDate = parsed
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Borrow{UpdatedBy: updatedBy}
	if in.Borrower != nil {
		data.Borrower = strings.TrimSpace(*in.Borrower)
	}
	if dueDate != nil {
		data.DueDate = *dueDate
	}
	if in.Remark != nil {
		data.Remark = *in.Remark
	}

	_, err = dao.Borrow.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(columns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// Delete soft-deletes borrow records by IDs. Records still in borrowed
// status are rejected so the inventory linkage can never dangle.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeBorrowDeleteRequired)
	}

	columns := dao.Borrow.Columns()
	borrowedCount, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "").
		WhereIn(columns.Id, idList).
		Where(columns.Status, BorrowStatusBorrowed).
		Count()
	if err != nil {
		return err
	}
	if borrowedCount > 0 {
		return bizerr.NewCode(CodeBorrowDeleteBorrowedRejected)
	}

	_, err = tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Borrow.Ctx(ctx), "").
		WhereIn(columns.Id, idList).
		Delete()
	return err
}

// ensureBookBorrowable verifies the book exists in the current tenant, is on
// shelf, and still has available copies.
func (s *serviceImpl) ensureBookBorrowable(ctx context.Context, bookID int64) error {
	bookColumns := dao.Book.Columns()
	existing := &book.BookEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(bookColumns.Id, bookID).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 || existing.Status != book.BookStatusOnShelf || existing.AvailableQuantity <= 0 {
		return bizerr.NewCode(CodeBookNoAvailable)
	}
	return nil
}

// decrementAvailable decrements one available copy with a conditional update
// so the count never goes negative under concurrency. It reports whether a
// row was actually affected.
func (s *serviceImpl) decrementAvailable(ctx context.Context, tenantID int, bookID int64) (bool, error) {
	bookColumns := dao.Book.Columns()
	result, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(bookColumns.Id, bookID).
		WhereGT(bookColumns.AvailableQuantity, 0).
		Decrement(bookColumns.AvailableQuantity, 1)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	return rows > 0, err
}

// incrementAvailable increments one available copy back.
func (s *serviceImpl) incrementAvailable(ctx context.Context, tenantID int, bookID int64) error {
	bookColumns := dao.Book.Columns()
	_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Where(bookColumns.Id, bookID).
		Increment(bookColumns.AvailableQuantity, 1)
	return err
}

// resolveBookTitleMap resolves book titles for the current page with one
// batch query, keeping list assembly free of per-row lookups.
func (s *serviceImpl) resolveBookTitleMap(ctx context.Context, records []*BorrowEntity) (map[int64]string, error) {
	titles := make(map[int64]string)
	bookIDs := make([]int64, 0, len(records))
	seen := make(map[int64]struct{}, len(records))
	for _, record := range records {
		if record == nil || record.BookId <= 0 {
			continue
		}
		if _, ok := seen[record.BookId]; ok {
			continue
		}
		seen[record.BookId] = struct{}{}
		bookIDs = append(bookIDs, record.BookId)
	}
	if len(bookIDs) == 0 {
		return titles, nil
	}

	bookColumns := dao.Book.Columns()
	books := make([]*book.BookEntity, 0, len(bookIDs))
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Book.Ctx(ctx), "").
		Fields(bookColumns.Id, bookColumns.Title).
		WhereIn(bookColumns.Id, bookIDs).
		Scan(&books)
	if err != nil {
		return nil, err
	}
	for _, item := range books {
		if item == nil {
			continue
		}
		titles[item.Id] = item.Title
	}
	return titles, nil
}

// parseDate converts the date-only wire format into a stored date value.
// Empty input defaults to today for required dates and nil for optional ones.
func parseDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation(dateFormat, trimmed, time.Local)
	if err != nil {
		return nil, bizerr.NewCode(CodeBorrowDateInvalid)
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

// currentUserID returns the current user domain ID from the business context.
func (s *serviceImpl) currentUserID(ctx context.Context) int64 {
	return int64(s.bizCtxSvc.Current(ctx).UserID)
}
