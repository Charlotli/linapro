import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

import { formatYuan } from './equipment-format';

/** 查询表单schema */
export function buildEquipmentQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('plugin.linapro-equipment-manage.fields.name'),
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: $t('plugin.linapro-equipment-manage.fields.type'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('plugin.linapro-equipment-manage.fields.status'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildEquipmentColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'equipmentCode',
      title: $t('plugin.linapro-equipment-manage.fields.code'),
      minWidth: 150,
    },
    {
      field: 'equipmentName',
      title: $t('plugin.linapro-equipment-manage.fields.name'),
      minWidth: 180,
    },
    {
      field: 'equipmentType',
      title: $t('plugin.linapro-equipment-manage.fields.type'),
      minWidth: 100,
      slots: { default: 'equipmentType' },
    },
    {
      field: 'brandModel',
      title: $t('plugin.linapro-equipment-manage.fields.brandModel'),
      minWidth: 150,
    },
    {
      field: 'purchasePrice',
      title: $t('plugin.linapro-equipment-manage.fields.purchasePrice'),
      minWidth: 130,
      align: 'right',
      headerAlign: 'center',
      formatter: ({ cellValue }) => formatYuan(cellValue),
      className: 'oa-amount-cell',
    },
    {
      field: 'location',
      title: $t('plugin.linapro-equipment-manage.fields.location'),
      minWidth: 130,
    },
    {
      field: 'owner',
      title: $t('plugin.linapro-equipment-manage.fields.owner'),
      minWidth: 100,
    },
    {
      field: 'status',
      title: $t('plugin.linapro-equipment-manage.fields.status'),
      minWidth: 100,
      slots: { default: 'status' },
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
