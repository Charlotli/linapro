// This file exposes the plugin-local generated entities through the service
// packages so controller code can keep stable type names.

package book

import entitymodel "lina-plugin-linapro-book-manage/backend/internal/model/entity"

// BookEntity mirrors the generated plugin_linapro_book_manage_book entity owned by this plugin.
type BookEntity = entitymodel.Book

// BorrowEntity mirrors the generated plugin_linapro_book_manage_borrow entity owned by this plugin.
type BorrowEntity = entitymodel.Borrow
