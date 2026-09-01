// This file exposes the plugin-local generated borrow entity through the
// service package so controller code can keep a stable type name.

package borrow

import entitymodel "lina-plugin-linapro-book-manage/backend/internal/model/entity"

// BorrowEntity mirrors the generated plugin_linapro_book_manage_borrow entity owned by this plugin.
type BorrowEntity = entitymodel.Borrow
