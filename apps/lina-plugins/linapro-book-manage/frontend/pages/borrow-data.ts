import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';

/** 查询表单schema */
export function buildBorrowQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Select',
      fieldName: 'bookId',
      label: $t('plugin.linapro-book-manage.fields.title'),
      componentProps: {
        allowClear: true,
        options: [] as { label: string; value: number }[],
        showSearch: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'borrower',
      label: $t('plugin.linapro-book-manage.fields.borrower'),
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('plugin.linapro-book-manage.fields.borrowStatus'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildBorrowColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'bookTitle',
      title: $t('plugin.linapro-book-manage.fields.title'),
      minWidth: 200,
    },
    {
      field: 'borrower',
      title: $t('plugin.linapro-book-manage.fields.borrower'),
      minWidth: 120,
    },
    {
      field: 'borrowDate',
      title: $t('plugin.linapro-book-manage.fields.borrowDate'),
      minWidth: 120,
    },
    {
      field: 'dueDate',
      title: $t('plugin.linapro-book-manage.fields.dueDate'),
      minWidth: 120,
    },
    {
      field: 'returnDate',
      title: $t('plugin.linapro-book-manage.fields.returnDate'),
      minWidth: 120,
    },
    {
      field: 'status',
      title: $t('plugin.linapro-book-manage.fields.borrowStatus'),
      minWidth: 100,
      slots: { default: 'status' },
    },
    {
      field: 'remark',
      title: $t('plugin.linapro-book-manage.fields.remark'),
      minWidth: 150,
    },
    {
      field: 'action',
      slots: { default: 'action' },
      title: $t('pages.common.actions'),
      fixed: 'right',
      resizable: false,
      width: 'auto',
    },
  ];
}
