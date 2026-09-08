import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

/** 查询表单schema */
export function buildLiveQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'title',
      label: $t('plugin.linapro-live-manage.fields.title'),
    },
    {
      component: 'Select',
      fieldName: 'roomId',
      label: $t('plugin.linapro-live-manage.fields.roomId'),
      componentProps: {
        allowClear: true,
        fieldNames: { label: 'roomName', value: 'id' },
        options: [] as { label: string; value: number }[],
        placeholder: $t('plugin.linapro-live-manage.placeholders.roomId'),
        showSearch: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'state',
      label: $t('plugin.linapro-live-manage.fields.state'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'Select',
      fieldName: 'isPublic',
      label: $t('plugin.linapro-live-manage.fields.isPublic'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildLiveColumns(): VxeGridProps['columns'] {
  return [
    { type: 'checkbox', width: 60 },
    {
      field: 'title',
      title: $t('plugin.linapro-live-manage.fields.title'),
      minWidth: 200,
    },
    {
      field: 'roomName',
      title: $t('plugin.linapro-live-manage.fields.roomId'),
      minWidth: 140,
    },
    {
      field: 'liveDate',
      title: $t('plugin.linapro-live-manage.fields.liveDate'),
      minWidth: 120,
    },
    {
      field: 'startTime',
      title: $t('plugin.linapro-live-manage.fields.startTime'),
      formatter: ({ cellValue }) => formatTimestamp(cellValue),
      minWidth: 180,
    },
    {
      field: 'state',
      title: $t('plugin.linapro-live-manage.fields.state'),
      minWidth: 100,
      slots: { default: 'state' },
    },
    {
      field: 'isPublic',
      title: $t('plugin.linapro-live-manage.fields.isPublic'),
      minWidth: 100,
      slots: { default: 'isPublic' },
    },
    {
      field: 'replayEnabled',
      title: $t('plugin.linapro-live-manage.fields.replayEnabled'),
      minWidth: 100,
      slots: { default: 'replayEnabled' },
    },
    {
      field: 'onlineCount',
      title: $t('plugin.linapro-live-manage.fields.onlineCount'),
      minWidth: 100,
    },
    {
      field: 'totalViews',
      title: $t('plugin.linapro-live-manage.fields.totalViews'),
      minWidth: 100,
    },
    {
      field: 'preacher',
      title: $t('plugin.linapro-live-manage.fields.preacher'),
      minWidth: 120,
    },
    {
      field: 'createdByName',
      title: $t('plugin.linapro-live-manage.fields.createdBy'),
      minWidth: 120,
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
