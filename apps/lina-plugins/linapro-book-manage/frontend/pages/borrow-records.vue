<script lang="ts">
export const pluginPageMeta = {
  routePath: '/book/borrow',
  title: 'Borrow Records',
};
</script>

<script setup lang="ts">
import type { BookOption, BorrowRecord } from './book-client';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Popconfirm, Space } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { DictTag } from '#/components/dict';
import { useUserStore } from '@vben/stores';
import { useDictStore } from '#/store/dict';

import {
  bookOptions,
  borrowList,
  borrowReturn,
} from './book-client';
import {
  buildBorrowColumns,
  buildBorrowQuerySchema,
} from './borrow-data';
import BorrowModal from './borrow-modal.vue';

const dictStore = useDictStore();
const userStore = useUserStore();
const statusDicts = ref<any[]>([]);
const bookOptionList = ref<BookOption[]>([]);

const currentUserId = computed(() => {
  const info: any = userStore.userInfo || {};
  return Number(info.userId ?? info.id ?? 0);
});

onMounted(async () => {
  statusDicts.value = await dictStore.getDictOptionsAsync(
    'plugin_book_borrow_status',
  );
  const optionRes = await bookOptions();
  bookOptionList.value = optionRes?.list ?? [];
  gridApi.formApi.updateSchema([
    {
      fieldName: 'bookId',
      componentProps: {
        options: bookOptionList.value.map((book) => ({
          label: book.title,
          value: book.id,
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

const [BorrowModalRef, borrowModalApi] = useVbenModal({
  connectedComponent: BorrowModal,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: buildBorrowQuerySchema(),
    commonConfig: {
      labelWidth: 80,
      componentProps: {
        allowClear: true,
      },
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: buildBorrowColumns(),
    height: 'auto',
    keepSource: true,
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues: Record<string, any> = {},
        ) => {
          return await borrowList({
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
    id: 'book-borrow-index',
  },
});

function handleAdd() {
  borrowModalApi.setData({});
  borrowModalApi.open();
}

function canReturn(row: BorrowRecord): boolean {
  return row.status === 1 && Number(row.createdBy) === currentUserId.value;
}

async function handleReturn(row: BorrowRecord) {
  await borrowReturn(row.id);
  message.success($t('plugin.linapro-book-manage.messages.returnSuccess'));
  await gridApi.query();
}

async function handleDelete(row: BorrowRecord) {
  await borrowDelete(String(row.id));
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}

function onReload() {
  gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true">
    <Grid :table-title="$t('plugin.linapro-book-manage.borrowTableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button type="primary" @click="handleAdd">
            {{ $t('plugin.linapro-book-manage.actions.borrow') }}
          </a-button>
        </Space>
      </template>

      <template #status="{ row }">
        <DictTag :dicts="statusDicts" :value="String(row.status)" />
      </template>

      <template #action="{ row }">
        <Space>
          <Popconfirm
            v-if="canReturn(row)"
            placement="left"
            :title="$t('plugin.linapro-book-manage.messages.returnConfirm')"
            @confirm="handleReturn(row)"
          >
            <ghost-button @click.stop="">
              {{ $t('plugin.linapro-book-manage.actions.return') }}
            </ghost-button>
          </Popconfirm>
          <Popconfirm
            v-if="row.status === 2"
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

    <BorrowModalRef @reload="onReload" />
  </Page>
</template>
