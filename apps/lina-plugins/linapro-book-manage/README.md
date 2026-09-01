# Book Management Plugin

`linapro-book-manage` is the official source plugin for book registry and borrow management.

## Capabilities

| Capability | Description |
| --- | --- |
| Book registry | Tenant-unique ISBN, dictionary-driven categories, copy counts with available tracking. |
| Borrow management | Lending decrements available copies, returning increments them, linked in one transaction. |
| Reference protection | Books still referenced by borrow records cannot be deleted. |

## Routes

| Menu | Path | Permission |
| --- | --- | --- |
| Books | `/book/register` | `book:list` |
| Borrow records | `/book/borrow` | `borrow:list` |

## Data Tables

`plugin_linapro_book_manage_book`, `plugin_linapro_book_manage_borrow`.

## Dictionary Types

`plugin_book_category`, `plugin_book_borrow_status`.
