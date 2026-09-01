# 书籍管理插件

`linapro-book-manage`是书籍登记与借阅管理的官方源码插件。

## 能力清单

| 能力 | 说明 |
| --- | --- |
| 书籍台账 | ISBN 同租户唯一，分类/状态字典治理，库存可借数联动跟踪。 |
| 借阅管理 | 借出扣减可借数量、归还回增，同一事务内联动防超借。 |
| 引用保护 | 仍被借阅记录引用的书籍禁止删除。 |

## 路由与权限

| 菜单 | 路径 | 权限 |
| --- | --- | --- |
| 书籍登记 | `/book/register` | `book:list` |
| 借阅记录 | `/book/borrow` | `borrow:list` |

## 数据表

`plugin_linapro_book_manage_book`、`plugin_linapro_book_manage_borrow`。

## 字典类型

`plugin_book_category`、`plugin_book_borrow_status`。
