import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

/** 查询表单schema */
export function buildFlowQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'flowName',
      label: $t('plugin.linapro-oa-approval.fields.flowName'),
    },
    {
      component: 'Select',
      fieldName: 'flowType',
      label: $t('plugin.linapro-oa-approval.fields.flowType'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('plugin.linapro-oa-approval.fields.status'),
      componentProps: {
        options: [] as { label: string; value: number }[],
      },
    },
  ];
}

/** 表格列定义 */
export function buildFlowColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'flowName',
      title: $t('plugin.linapro-oa-approval.fields.flowName'),
      minWidth: 180,
    },
    {
      field: 'flowType',
      title: $t('plugin.linapro-oa-approval.fields.flowType'),
      minWidth: 100,
      slots: { default: 'flowType' },
    },
    {
      field: 'nodeCount',
      title: $t('plugin.linapro-oa-approval.fields.nodeCount'),
      minWidth: 100,
    },
    {
      field: 'fieldCount',
      title: $t('plugin.linapro-oa-approval.fields.fieldCount'),
      minWidth: 120,
      formatter: ({ cellValue }) =>
        cellValue > 0
          ? String(cellValue)
          : $t('plugin.linapro-oa-approval.messages.standardForm'),
    },
    {
      field: 'description',
      title: $t('plugin.linapro-oa-approval.fields.description'),
      minWidth: 200,
    },
    {
      field: 'status',
      title: $t('plugin.linapro-oa-approval.fields.status'),
      minWidth: 100,
      slots: { default: 'status' },
    },
    {
      field: 'createdByName',
      title: $t('plugin.linapro-oa-approval.fields.createdBy'),
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
