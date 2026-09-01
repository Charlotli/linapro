import type { VbenFormSchema } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

import { formatYuan } from './oa-format';

/** 查询表单schema */
export function buildRequestQuerySchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'title',
      label: $t('plugin.linapro-oa-approval.fields.title'),
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
export function buildRequestColumns(): VxeGridProps['columns'] {
  return [
    {
      field: 'title',
      title: $t('plugin.linapro-oa-approval.fields.title'),
      minWidth: 200,
    },
    {
      field: 'flowType',
      title: $t('plugin.linapro-oa-approval.fields.flowType'),
      minWidth: 100,
      slots: { default: 'flowType' },
    },
    {
      field: 'amount',
      title: $t('plugin.linapro-oa-approval.fields.amount'),
      minWidth: 140,
      align: 'right',
      headerAlign: 'center',
      formatter: ({ cellValue }) => formatYuan(cellValue),
      className: 'oa-amount-cell',
    },
    {
      field: 'applicantName',
      title: $t('plugin.linapro-oa-approval.fields.applicant'),
      minWidth: 120,
    },
    {
      field: 'currentApproverName',
      title: $t('plugin.linapro-oa-approval.fields.currentApprover'),
      minWidth: 140,
    },
    {
      field: 'status',
      title: $t('plugin.linapro-oa-approval.fields.status'),
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
