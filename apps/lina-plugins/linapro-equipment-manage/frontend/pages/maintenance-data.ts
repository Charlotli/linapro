import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

/** 查询表单schema */
export function buildMaintenanceQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Select',
      fieldName: 'equipmentId',
      label: $t('plugin.linapro-equipment-manage.fields.equipment'),
      componentProps: {
        allowClear: true,
        options: [] as { label: string; value: number }[],
        showSearch: true,
      },
    },
    {
      component: 'Select',
      fieldName: 'maintType',
      label: $t('plugin.linapro-equipment-manage.fields.maintType'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'RangePicker',
      fieldName: 'maintDateRange',
      label: $t('plugin.linapro-equipment-manage.fields.maintDate'),
    },
  ];
}

/** 表格列定义 */
export function buildMaintenanceColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'equipmentName',
      title: $t('plugin.linapro-equipment-manage.fields.equipment'),
      minWidth: 180,
    },
    {
      field: 'maintType',
      title: $t('plugin.linapro-equipment-manage.fields.maintType'),
      minWidth: 100,
      slots: { default: 'maintType' },
    },
    {
      field: 'maintDate',
      title: $t('plugin.linapro-equipment-manage.fields.maintDate'),
      minWidth: 120,
    },
    {
      field: 'maintainer',
      title: $t('plugin.linapro-equipment-manage.fields.maintainer'),
      minWidth: 120,
    },
    {
      field: 'cost',
      title: $t('plugin.linapro-equipment-manage.fields.cost'),
      minWidth: 130,
      align: 'right',
      headerAlign: 'center',
      slots: { default: 'cost' },
    },
    {
      field: 'content',
      title: $t('plugin.linapro-equipment-manage.fields.content'),
      minWidth: 200,
    },
    {
      field: 'result',
      title: $t('plugin.linapro-equipment-manage.fields.result'),
      minWidth: 130,
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
