<script lang="ts">
export const pluginPageMeta = {
  routePath: '/oa/approval-flow',
  title: 'Approval Flows',
};
</script>

<script setup lang="ts">
import type { ApprovalFlow } from './oa-client';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Popconfirm, Space, Switch } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import { flowDelete, flowList, flowUpdate } from './oa-client';
import { buildFlowColumns, buildFlowQuerySchema } from './oa-flow-data';
import OaFlowModal from './oa-flow-modal.vue';

const dictStore = useDictStore();
const flowTypeDicts = ref<any[]>([]);
const flowStatusDicts = ref<any[]>([]);

onMounted(async () => {
  [flowTypeDicts.value, flowStatusDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_oa_approval_flow_type'),
    dictStore.getDictOptionsAsync('plugin_oa_approval_status'),
  ]);
  gridApi.formApi.updateSchema([
    {
      fieldName: 'flowType',
      componentProps: {
        options: flowTypeDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [OaFlowModalRef, flowModalApi] = useVbenModal({
  connectedComponent: OaFlowModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildFlowQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: buildFlowColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await flowList({
            pageNum: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          });
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    id: 'oa-flow-index',
  },
});

function handleAdd() {
  flowModalApi.setData({});
  flowModalApi.open();
}

function handleEdit(row: ApprovalFlow) {
  flowModalApi.setData({ id: row.id });
  flowModalApi.open();
}

async function handleStatusChange(row: ApprovalFlow, checked: boolean) {
  await flowUpdate(row.id, { status: checked ? 1 : 0 });
  message.success($t('pages.common.updateSuccess'));
  await gridApi.query();
}

async function handleDelete(row: ApprovalFlow) {
  await flowDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function onReload() {
  gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-oa-approval.flowTableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button type="primary" @click="handleAdd">
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>

      <template #flowType="{ row }">
        <DictTag :dicts="flowTypeDicts" :value="String(row.flowType)" />
      </template>

      <template #status="{ row }">
        <Switch
          :checked="row.status === 1"
          :checked-children="$t('plugin.linapro-oa-approval.status.enabled')"
          :un-checked-children="$t('plugin.linapro-oa-approval.status.disabled')"
          @change="(checked: any) => handleStatusChange(row, Boolean(checked))"
        />
      </template>

      <template #action="{ row }">
        <Space>
          <ghost-button @click.stop="handleEdit(row)">
            {{ $t('pages.common.edit') }}
          </ghost-button>
          <Popconfirm
            placement="left"
            :title="$t('pages.common.deleteConfirm')"
            @confirm="handleDelete(row)"
          >
            <ghost-button danger @click.stop="">
              {{ $t('pages.common.delete') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </Grid>

    <OaFlowModalRef @reload="onReload" />
  </Page>
</template>
