<script lang="ts">
export const pluginPageMeta = {
  routePath: '/book/register',
  title: 'Books',
};
</script>

<script setup lang="ts">
import type { Book } from './book-client';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';

import { bookDelete, bookList } from './book-client';
import { buildBookColumns, buildBookQuerySchema } from './book-data';
import BookModal from './book-modal.vue';

const dictStore = useDictStore();
const categoryDicts = ref<any[]>([]);

onMounted(async () => {
  categoryDicts.value = await dictStore.getDictOptionsAsync(
    'plugin_book_category',
  );
  gridApi.formApi.updateSchema([
    {
      fieldName: 'category',
      componentProps: {
        options: categoryDicts.value.map((item: any) => ({
          label: item.label,
          value: Number(item.value),
        })),
      },
    },
  ]);
});

const [BookModalRef, bookModalApi] = useVbenModal({
  connectedComponent: BookModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildBookQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: buildBookColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await bookList({
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
    id: 'book-register-index',
  },
});

function handleAdd() {
  bookModalApi.setData({});
  bookModalApi.open();
}

function handleEdit(row: Book) {
  bookModalApi.setData({ id: row.id });
  bookModalApi.open();
}

async function handleDelete(row: Book) {
  await bookDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function onReload() {
  gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-book-manage.tableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button type="primary" @click="handleAdd">
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>

      <template #category="{ row }">
        <DictTag :dicts="categoryDicts" :value="String(row.category)" />
      </template>

      <template #available="{ row }">
        <span
          :class="row.availableQuantity > 0 ? 'text-green-600' : 'text-red-500'"
          class="font-medium tabular-nums"
        >
          {{ row.availableQuantity }} / {{ row.totalQuantity }}
        </span>
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

    <BookModalRef @reload="onReload" />
  </Page>
</template>
