import { pluginApiPath, requestClient } from '#/api/request';

const pluginID = 'linapro-book-manage';

function api(pathName: string) {
  return pluginApiPath(pluginID, pathName);
}

export interface Book {
  id: number;
  title: string;
  author: string;
  isbn: string;
  category: number;
  publisher: string;
  publishDate: string;
  totalQuantity: number;
  availableQuantity: number;
  location: string;
  status: number;
  coverUrl: string;
  remark: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface BookListParams {
  pageNum?: number;
  pageSize?: number;
  title?: string;
  category?: number;
  status?: number;
}

export interface BorrowRecord {
  id: number;
  bookId: number;
  bookTitle: string;
  borrower: string;
  borrowDate: string;
  dueDate: string;
  returnDate: string;
  status: number;
  remark: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface BorrowListParams {
  pageNum?: number;
  pageSize?: number;
  bookId?: number;
  borrower?: string;
  status?: number;
}

export async function bookList(params?: BookListParams) {
  const res = await requestClient.get<{ list: Book[]; total: number }>(
    api('book'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export interface BookOption {
  id: number;
  title: string;
  author: string;
  availableQuantity: number;
}

export function bookOptions(keyword?: string) {
  return requestClient.get<{ list: BookOption[] }>(api('book/options'), {
    params: { keyword },
  });
}

export function bookInfo(id: number) {
  return requestClient.get<Book>(api(`book/${id}`));
}

export function bookAdd(data: Partial<Book>) {
  return requestClient.post(api('book'), data);
}

export function bookUpdate(id: number, data: Partial<Book>) {
  return requestClient.put(api(`book/${id}`), data);
}

export function bookDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(api('book'), {
    params: { ids: list },
  });
}

export async function borrowList(params?: BorrowListParams) {
  const res = await requestClient.get<{ list: BorrowRecord[]; total: number }>(
    api('borrow'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export function borrowAdd(data: {
  bookId: number;
  borrower: string;
  borrowDate?: string;
  dueDate?: string;
  remark?: string;
}) {
  return requestClient.post(api('borrow'), data);
}

export function borrowReturn(id: number, returnDate?: string) {
  return requestClient.put(api(`borrow/${id}/return`), {
    returnDate,
  });
}

export function borrowUpdate(
  id: number,
  data: { borrower?: string; dueDate?: string; remark?: string },
) {
  return requestClient.put(api(`borrow/${id}`), data);
}

export function borrowDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(api('borrow'), {
    params: { ids: list },
  });
}
