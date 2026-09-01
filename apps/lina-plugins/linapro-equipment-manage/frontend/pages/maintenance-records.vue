<script lang="ts">
export const pluginPageMeta = {
  routePath: '/equipment/maintenance',
  title: 'Maintenance Records',
};
</script>

<script setup lang="ts">
import type { EquipmentOption, MaintenanceRecord } from './equipment-client';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import { formatYuan } from './equipment-format';
import {
  equipmentOptions,
  maintenanceDelete,
  maintenanceList,
} from './equipment-client';
import {
  buildMaintenanceColumns,
  buildMaintenanceQuerySchema,
} from './maintenance-data';
import MaintenanceModal from './maintenance-modal.vue';

const dictStore = useDictStore();
const maintTypeDicts = ref<any[]>([]);
const equipmentOptionList = ref<EquipmentOption[]>([]);

onMounted(async () => {
  maintTypeDicts.value = await dictStore.getDictOptionsAsync(
    'plugin_equipment_maint_type',
  );
  const optionRes = await equipmentOptions();
  equipmentOptionList.value = optionRes?.list ?? [];
  gridApi.formApi.updateSchema([
    {
      fieldName: 'equipmentId',
      componentProps: {
        options: equipmentOptionList.value.map((equipment) => ({
          label: `${equipment.equipmentName}（${equipment.equipmentCode}）`,
          value: equipment.id,
        })),
      },
    },
    {
      fieldName: 'maintType',
      componentProps: {
        options: maintTypeDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [MaintenanceModalRef, maintenanceModalApi] = useVbenModal({
  connectedComponent: MaintenanceModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildMaintenanceQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: buildMaintenanceColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          {
            page,
          }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          const { maintDateRange, ...rest } = formValues;
          const [dateStart, dateEnd] = maintDateRange ?? ['', ''];
          return await maintenanceList({
            pageNum: page.currentPage,
            pageSize: page.pageSize,
            dateStart: dateStart || undefined,
            dateEnd: dateEnd || undefined,
            ...rest,
          });
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    id: 'equipment-maintenance-index',
  },
});

function handleAdd() {
  maintenanceModalApi.setData({});
  maintenanceModalApi.open();
}

function handleEdit(row: MaintenanceRecord) {
  maintenanceModalApi.setData({ id: row.id });
  maintenanceModalApi.open();
}

async function handleDelete(row: MaintenanceRecord) {
  await maintenanceDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function onReload() {
  gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-equipment-manage.maintTableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button type="primary" @click="handleAdd">
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>

      <template #maintType="{ row }">
        <DictTag :dicts="maintTypeDicts" :value="String(row.maintType)" />
      </template>

      <template #cost="{ row }">
        {{ formatYuan(row.cost) }}
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

    <MaintenanceModalRef @reload="onReload" />
  </Page>
</template>
