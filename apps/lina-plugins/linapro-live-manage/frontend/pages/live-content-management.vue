<script lang="ts">
export const pluginPageMeta = {
  routePath: '/live/content',
  title: 'Live Content',
};
</script>

<script setup lang="ts">
import type { LiveContent } from './live-client';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import {
  liveDelete,
  liveList,
  liveStart,
  liveStop,
  liveroomOptions,
} from './live-client';
import { buildLiveColumns, buildLiveQuerySchema } from './live-content-data';
import LiveContentModal from './live-content-modal.vue';

const dictStore = useDictStore();
const liveStateDicts = ref<any[]>([]);
const livePublicDicts = ref<any[]>([]);

onMounted(async () => {
  [liveStateDicts.value, livePublicDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_live_state'),
    dictStore.getDictOptionsAsync('plugin_live_public'),
  ]);
  const roomOptions = await loadRoomOptions();
  gridApi.formApi.updateSchema([
    {
      fieldName: 'roomId',
      componentProps: {
        fieldNames: { label: 'roomName', value: 'id' },
        options: roomOptions,
      },
    },
    {
      fieldName: 'state',
      componentProps: {
        options: liveStateDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
    {
      fieldName: 'isPublic',
      componentProps: {
        options: livePublicDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

async function loadRoomOptions(keyword?: string) {
  const res = await liveroomOptions(keyword);
  return res?.list ?? [];
}

const [LiveContentModalRef, liveModalApi] = useVbenModal({
  connectedComponent: LiveContentModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildLiveQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    checkboxConfig: {
      highlight: true,
      reserve: true,
    },
    columns: buildLiveColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await liveList({
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
    id: 'live-content-index',
  },
  gridEvents: {
    checkboxChange: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as LiveContent[];
    },
    checkboxAll: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as LiveContent[];
    },
  },
});

const checkedRows = ref<LiveContent[]>([]);
const hasChecked = computed(() => checkedRows.value.length > 0);

function handleAdd() {
  liveModalApi.setData({});
  liveModalApi.open();
}

function handleEdit(row: LiveContent) {
  liveModalApi.setData({ id: row.id });
  liveModalApi.open();
}

async function handleDelete(row: LiveContent) {
  await liveDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

async function handleStart(row: LiveContent) {
  await liveStart(row.id);
  message.success($t('plugin.linapro-live-manage.messages.startSuccess'));
  await gridApi.query();
}

async function handleStop(row: LiveContent) {
  await liveStop(row.id);
  message.success($t('plugin.linapro-live-manage.messages.stopSuccess'));
  await gridApi.query();
}

function handleMultiDelete() {
  const rows = gridApi.grid.getCheckboxRecords() as LiveContent[];
  const ids = rows.map((row) => row.id);
  Modal.confirm({
    title: $t('pages.common.confirmTitle'),
    okType: 'danger',
    content: $t('plugin.linapro-live-manage.messages.deleteSelectedLiveConfirm', {
      count: ids.length,
    }),
    onOk: async () => {
      await liveDelete(ids);
      checkedRows.value = [];
      await gridApi.query();
    },
  });
}

function onReload() {
  gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-live-manage.liveTableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button
            :disabled="!hasChecked"
            danger
            type="primary"
            @click="handleMultiDelete"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
          <a-button type="primary" @click="handleAdd">
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>

      <template #state="{ row }">
        <DictTag :dicts="liveStateDicts" :value="String(row.state)" />
      </template>

      <template #isPublic="{ row }">
        <DictTag :dicts="livePublicDicts" :value="String(row.isPublic)" />
      </template>

      <template #action="{ row }">
        <Space>
          <Popconfirm
            v-if="row.state === 0"
            placement="left"
            :title="$t('plugin.linapro-live-manage.messages.startConfirm')"
            @confirm="handleStart(row)"
          >
            <ghost-button @click.stop="">
              {{ $t('plugin.linapro-live-manage.actions.start') }}
            </ghost-button>
          </Popconfirm>
          <Popconfirm
            v-else-if="row.state === 1"
            placement="left"
            :title="$t('plugin.linapro-live-manage.messages.stopConfirm')"
            @confirm="handleStop(row)"
          >
            <ghost-button danger @click.stop="">
              {{ $t('plugin.linapro-live-manage.actions.stop') }}
            </ghost-button>
          </Popconfirm>
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

    <LiveContentModalRef @reload="onReload" />
  </Page>
</template>
