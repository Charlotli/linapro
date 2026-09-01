<script lang="ts">
export const pluginPageMeta = {
  routePath: '/equipment/register',
  title: 'Equipment',
};
</script>

<script setup lang="ts">
import type { Equipment } from './equipment-client';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import { equipmentDelete, equipmentList } from './equipment-client';
import { buildEquipmentColumns, buildEquipmentQuerySchema } from './equipment-data';
import EquipmentModal from './equipment-modal.vue';

const dictStore = useDictStore();
const equipmentTypeDicts = ref<any[]>([]);
const equipmentStatusDicts = ref<any[]>([]);

onMounted(async () => {
  [equipmentTypeDicts.value, equipmentStatusDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_equipment_type'),
    dictStore.getDictOptionsAsync('plugin_equipment_status'),
  ]);
  gridApi.formApi.updateSchema([
    {
      fieldName: 'type',
      componentProps: {
        options: equipmentTypeDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
    {
      fieldName: 'status',
      componentProps: {
        options: equipmentStatusDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [EquipmentModalRef, equipmentModalApi] = useVbenModal({
  connectedComponent: EquipmentModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildEquipmentQuerySchema(),
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
    columns: buildEquipmentColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await equipmentList({
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
    id: 'equipment-register-index',
  },
  gridEvents: {
    checkboxChange: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as Equipment[];
    },
    checkboxAll: () => {
      checkedRows.value = (gridApi.grid?.getCheckboxRecords() || []) as Equipment[];
    },
  },
});

const checkedRows = ref<Equipment[]>([]);

function handleAdd() {
  equipmentModalApi.setData({});
  equipmentModalApi.open();
}

function handleEdit(row: Equipment) {
  equipmentModalApi.setData({ id: row.id });
  equipmentModalApi.open();
}

async function handleDelete(row: Equipment) {
  await equipmentDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function handleMultiDelete() {
  const rows = gridApi.grid.getCheckboxRecords() as Equipment[];
  const ids = rows.map((row) => row.id);
  Modal.confirm({
    title: $t('pages.common.confirmTitle'),
    okType: 'danger',
    content: $t('plugin.linapro-equipment-manage.messages.deleteSelectedConfirm', {
      count: ids.length,
    }),
    onOk: async () => {
      await equipmentDelete(ids);
      checkedRows.value = [];
      await gridApi.query();
    },
  });
}

function onReload() {
  gridApi.query();
}

import { Modal } from 'ant-design-vue';
import { computed } from 'vue';

const hasChecked = computed(() => checkedRows.value.length > 0);
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-equipment-manage.tableTitle')">
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

      <template #equipmentType="{ row }">
        <DictTag :dicts="equipmentTypeDicts" :value="String(row.equipmentType)" />
      </template>

      <template #status="{ row }">
        <DictTag :dicts="equipmentStatusDicts" :value="String(row.status)" />
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

    <EquipmentModalRef @reload="onReload" />
  </Page>
</template>
