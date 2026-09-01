<script lang="ts">
export const pluginPageMeta = {
  routePath: '/oa/approval',
  title: 'Approval Center',
};
</script>

<script setup lang="ts">
import type { ApprovalRequest } from './oa-client';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';

import {
  Badge,
  message,
  Popconfirm,
  RadioButton,
  RadioGroup,
  Space,
} from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import { pendingCount, requestList, requestWithdraw } from './oa-client';
import { buildRequestColumns, buildRequestQuerySchema } from './oa-approval-data';
import OaDetailDrawer from './oa-detail-drawer.vue';
import OaSubmitModal from './oa-submit-modal.vue';

const dictStore = useDictStore();
const flowTypeDicts = ref<any[]>([]);
const statusDicts = ref<any[]>([]);
const listScope = ref<'mine' | 'pending'>('mine');

const pendingBadgeCount = ref(0);

async function refreshPendingCount() {
  try {
    const res = await pendingCount();
    pendingBadgeCount.value = res?.count ?? 0;
  } catch {
    pendingBadgeCount.value = 0;
  }
}

onMounted(async () => {
  [flowTypeDicts.value, statusDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_oa_approval_flow_type'),
    dictStore.getDictOptionsAsync('plugin_oa_approval_status'),
  ]);
  await refreshPendingCount();
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
    {
      fieldName: 'status',
      componentProps: {
        options: statusDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [OaSubmitModalRef, submitModalApi] = useVbenModal({
  connectedComponent: OaSubmitModal,
});

const [OaDetailDrawerRef, detailDrawerApi] = useVbenDrawer({
  connectedComponent: OaDetailDrawer,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildRequestQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: buildRequestColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await requestList({
            pageNum: page.currentPage,
            pageSize: page.pageSize,
            scope: listScope.value,
            ...formValues,
          });
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    id: 'oa-approval-index',
  },
});

function handleScopeChange() {
  gridApi.query();
  refreshPendingCount();
}

function handleAdd() {
  submitModalApi.setData({});
  submitModalApi.open();
}

function handleDetail(row: ApprovalRequest) {
  detailDrawerApi.setData({ id: row.id });
  detailDrawerApi.open();
}

async function handleWithdraw(row: ApprovalRequest) {
  await requestWithdraw(row.id);
  message.success($t('plugin.linapro-oa-approval.messages.withdrawSuccess'));
  await gridApi.query();
}

function onReload() {
  gridApi.query();
  refreshPendingCount();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-oa-approval.requestTableTitle')">
      <template #toolbar-tools>
        <Space>
          <RadioGroup v-model:value="listScope" button-style="solid" size="small" @change="handleScopeChange">
            <RadioButton value="mine">
              {{ $t('plugin.linapro-oa-approval.scope.mine') }}
            </RadioButton>
            <RadioButton value="pending">
              <Badge
                :count="pendingBadgeCount"
                :number-style="{ boxShadow: 'none', marginLeft: '-4px' }"
                :overflow-count="99"
              >
                {{ $t('plugin.linapro-oa-approval.scope.pending') }}
              </Badge>
            </RadioButton>
          </RadioGroup>
          <a-button type="primary" @click="handleAdd">
            {{ $t('plugin.linapro-oa-approval.actions.submit') }}
          </a-button>
        </Space>
      </template>

      <template #flowType="{ row }">
        <DictTag :dicts="flowTypeDicts" :value="String(row.flowType)" />
      </template>

      <template #status="{ row }">
        <DictTag :dicts="statusDicts" :value="String(row.status)" />
      </template>

      <template #action="{ row }">
        <Space>
          <ghost-button @click.stop="handleDetail(row)">
            {{ $t('pages.common.detail') }}
          </ghost-button>
          <Popconfirm
            v-if="row.status === 1"
            placement="left"
            :title="$t('plugin.linapro-oa-approval.messages.withdrawConfirm')"
            @confirm="handleWithdraw(row)"
          >
            <ghost-button danger @click.stop="">
              {{ $t('plugin.linapro-oa-approval.actions.withdraw') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </Grid>

    <OaSubmitModalRef @reload="onReload" />
    <OaDetailDrawerRef @reload="onReload" />
  </Page>
</template>
