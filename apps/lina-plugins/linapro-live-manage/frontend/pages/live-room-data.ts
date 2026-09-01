import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

/** 查询表单schema */
export function buildRoomQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'roomName',
      label: $t('plugin.linapro-live-manage.fields.roomName'),
    },
    {
      component: 'Select',
      fieldName: 'roomType',
      label: $t('plugin.linapro-live-manage.fields.roomType'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('plugin.linapro-live-manage.fields.status'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildRoomColumns(): VxeGridProps['columns'] {
  return [
    { type: 'checkbox', width: 60 },
    {
      field: 'roomCode',
      title: $t('plugin.linapro-live-manage.fields.roomCode'),
      minWidth: 140,
    },
    {
      field: 'roomName',
      title: $t('plugin.linapro-live-manage.fields.roomName'),
      minWidth: 180,
    },
    {
      field: 'roomType',
      title: $t('plugin.linapro-live-manage.fields.roomType'),
      minWidth: 100,
      slots: { default: 'roomType' },
    },
    {
      field: 'status',
      title: $t('plugin.linapro-live-manage.fields.status'),
      minWidth: 100,
      slots: { default: 'status' },
    },
    {
      field: 'description',
      title: $t('plugin.linapro-live-manage.fields.description'),
      minWidth: 180,
    },
    {
      field: 'createdByName',
      title: $t('plugin.linapro-live-manage.fields.createdBy'),
      minWidth: 120,
    },
    {
      field: 'createdAt',
      title: $t('pages.common.createdAt'),
      formatter: ({ cellValue }) => formatTimestamp(cellValue),
      minWidth: 180,
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
