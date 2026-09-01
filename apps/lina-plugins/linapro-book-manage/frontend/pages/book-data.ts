import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

/** 查询表单schema */
export function buildBookQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'title',
      label: $t('plugin.linapro-book-manage.fields.title'),
    },
    {
      component: 'Select',
      fieldName: 'category',
      label: $t('plugin.linapro-book-manage.fields.category'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('plugin.linapro-book-manage.fields.status'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildBookColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'title',
      title: $t('plugin.linapro-book-manage.fields.title'),
      minWidth: 200,
    },
    {
      field: 'author',
      title: $t('plugin.linapro-book-manage.fields.author'),
      minWidth: 130,
    },
    {
      field: 'isbn',
      title: $t('plugin.linapro-book-manage.fields.isbn'),
      minWidth: 140,
    },
    {
      field: 'category',
      title: $t('plugin.linapro-book-manage.fields.category'),
      minWidth: 100,
      slots: { default: 'category' },
    },
    {
      field: 'publisher',
      title: $t('plugin.linapro-book-manage.fields.publisher'),
      minWidth: 160,
    },
    {
      field: 'availableQuantity',
      title: $t('plugin.linapro-book-manage.fields.available'),
      minWidth: 100,
      slots: { default: 'available' },
    },
    {
      field: 'totalQuantity',
      title: $t('plugin.linapro-book-manage.fields.totalQuantity'),
      minWidth: 100,
    },
    {
      field: 'location',
      title: $t('plugin.linapro-book-manage.fields.location'),
      minWidth: 110,
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
