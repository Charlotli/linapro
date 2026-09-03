<script lang="ts">
export const pluginPageMeta = {
  routePath: '/live/room',
  title: 'Live Rooms',
};
</script>

<script setup lang="ts">
import type { LiveRoom } from './live-client';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Modal, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import {
  liveroomDelete,
  liveroomList,
} from './live-client';
import LiveRoomCalendarModal from './live-room-calendar-modal.vue';
import { buildRoomColumns, buildRoomQuerySchema } from './live-room-data';
import LiveRoomModal from './live-room-modal.vue';

const dictStore = useDictStore();
const roomTypeDicts = ref<any[]>([]);
const roomStatusDicts = ref<any[]>([]);

onMounted(async () => {
  [roomTypeDicts.value, roomStatusDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_live_room_type'),
    dictStore.getDictOptionsAsync('plugin_live_room_status'),
  ]);
  gridApi.formApi.updateSchema([
    {
      fieldName: 'roomType',
      componentProps: {
        options: roomTypeDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
    {
      fieldName: 'status',
      componentProps: {
        options: roomStatusDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [LiveRoomModalRef, roomModalApi] = useVbenModal({
  connectedComponent: LiveRoomModal,
});

const [CalendarModalRef, calendarModalApi] = useVbenModal({
  connectedComponent: LiveRoomCalendarModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildRoomQuerySchema(),
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
    columns: buildRoomColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await liveroomList({
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
    id: 'live-room-index',
  },
  gridEvents: {
    checkboxChange: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as LiveRoom[];
    },
    checkboxAll: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as LiveRoom[];
    },
  },
});

const checkedRows = ref<LiveRoom[]>([]);
const hasChecked = computed(() => checkedRows.value.length > 0);

function handleAdd() {
  roomModalApi.setData({});
  roomModalApi.open();
}

function handleEdit(row: LiveRoom) {
  roomModalApi.setData({ id: row.id });
  roomModalApi.open();
}

function handleCalendar(row: LiveRoom) {
  calendarModalApi.setData({ roomCode: row.roomCode, roomName: row.roomName });
  calendarModalApi.open();
}

async function handleDelete(row: LiveRoom) {
  await liveroomDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function handleMultiDelete() {
  const rows = gridApi.grid.getCheckboxRecords() as LiveRoom[];
  const ids = rows.map((row) => row.id);
  Modal.confirm({
    title: $t('pages.common.confirmTitle'),
    okType: 'danger',
    content: $t('plugin.linapro-live-manage.messages.deleteSelectedRoomConfirm', {
      count: ids.length,
    }),
    onOk: async () => {
      await liveroomDelete(ids);
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
    <Grid :table-title="$t('plugin.linapro-live-manage.roomTableTitle')">
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

      <template #roomType="{ row }">
        <DictTag :dicts="roomTypeDicts" :value="String(row.roomType)" />
      </template>

      <template #status="{ row }">
        <DictTag :dicts="roomStatusDicts" :value="String(row.status)" />
      </template>

      <template #action="{ row }">
        <Space>
          <ghost-button @click.stop="handleCalendar(row)">
            {{ $t('plugin.linapro-live-manage.actions.calendar') }}
          </ghost-button>
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

    <LiveRoomModalRef @reload="onReload" />
    <CalendarModalRef />
  </Page>
</template>
