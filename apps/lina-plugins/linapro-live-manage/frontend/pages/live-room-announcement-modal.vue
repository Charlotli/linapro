<script setup lang="ts">
import type { Announcement } from './live-client';

import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { Input, InputNumber, message, Popconfirm, Switch } from 'ant-design-vue';

import { $t } from '#/locales';

import {
  announcementAdd,
  announcementDelete,
  announcementList,
  announcementUpdate,
} from './live-client';

const roomName = ref('');
const roomId = ref(0);
const items = ref<Announcement[]>([]);
const loading = ref(false);

// 编辑态：editingId 为 null 表示新增，否则为编辑既有公告。
const editing = reactive({
  id: null as number | null,
  title: '',
  content: '',
  sort: 0,
});

const editingVisible = ref(false);

const editingTitle = computed(() =>
  editing.id === null
    ? $t('plugin.linapro-live-manage.announcement.addTitle')
    : $t('plugin.linapro-live-manage.announcement.editTitle'),
);

async function refresh() {
  if (!roomId.value) return;
  loading.value = true;
  try {
    const res = await announcementList({ roomId: roomId.value, pageSize: 100 });
    items.value = res.items;
  } finally {
    loading.value = false;
  }
}

function startAdd() {
  editing.id = null;
  editing.title = '';
  editing.content = '';
  editing.sort = 0;
  editingVisible.value = true;
}

function startEdit(row: Announcement) {
  editing.id = row.id;
  editing.title = row.title;
  editing.content = row.content;
  editing.sort = row.sort;
  editingVisible.value = true;
}

function cancelEdit() {
  editingVisible.value = false;
}

async function submitEdit() {
  if (!editing.title.trim()) {
    message.warning($t('plugin.linapro-live-manage.announcement.titleRequired'));
    return;
  }
  if (editing.id === null) {
    await announcementAdd({
      roomId: roomId.value,
      title: editing.title,
      content: editing.content,
      sort: editing.sort,
    });
    message.success($t('pages.common.createSuccess'));
  } else {
    await announcementUpdate(editing.id, {
      title: editing.title,
      content: editing.content,
      sort: editing.sort,
    });
    message.success($t('pages.common.updateSuccess'));
  }
  editingVisible.value = false;
  await refresh();
}

async function toggleEnabled(row: Announcement, checked: boolean | string | number) {
  await announcementUpdate(row.id, { enabled: Boolean(checked) });
  row.enabled = Boolean(checked);
}

async function handleDelete(row: Announcement) {
  await announcementDelete(row.id);
  message.success($t('pages.common.deleteSuccess'));
  await refresh();
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[680px]',
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const data = modalApi.getData<{ roomId?: number; roomName?: string }>();
    roomId.value = data?.roomId ?? 0;
    roomName.value = data?.roomName ?? '';
    editingVisible.value = false;
    await refresh();
  },
});
</script>

<template>
  <Modal :title="$t('plugin.linapro-live-manage.announcement.title')">
    <div class="flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <span class="text-muted-foreground text-sm">{{ roomName }}</span>
        <a-button
          v-if="!editingVisible"
          type="primary"
          size="small"
          @click="startAdd"
        >
          {{ $t('plugin.linapro-live-manage.announcement.add') }}
        </a-button>
      </div>

      <div v-if="editingVisible" class="rounded-lg border p-3">
        <p class="mb-2 text-sm font-medium">{{ editingTitle }}</p>
        <Input
          v-model:value="editing.title"
          :placeholder="$t('plugin.linapro-live-manage.announcement.titlePlaceholder')"
          :maxlength="256"
          class="mb-2"
        />
        <Input.TextArea
          v-model:value="editing.content"
          :placeholder="$t('plugin.linapro-live-manage.announcement.contentPlaceholder')"
          :rows="3"
          :maxlength="4000"
          class="mb-2"
        />
        <div class="flex items-center justify-between">
          <InputNumber
            v-model:value="editing.sort"
            :min="0"
            :precision="0"
            size="small"
          />
          <div class="flex gap-2">
            <a-button size="small" @click="cancelEdit">
              {{ $t('pages.common.cancel') }}
            </a-button>
            <a-button size="small" type="primary" @click="submitEdit">
              {{ $t('pages.common.save') }}
            </a-button>
          </div>
        </div>
      </div>

      <a-spin :spinning="loading">
        <div v-if="items.length === 0" class="text-muted-foreground py-6 text-center text-sm">
          {{ $t('plugin.linapro-live-manage.announcement.empty') }}
        </div>
        <div
          v-for="row in items"
          :key="row.id"
          class="announcement-row mb-2 rounded-lg border p-3"
        >
          <div class="flex items-center justify-between gap-2">
            <span class="text-sm font-medium">{{ row.title }}</span>
            <Switch
              size="small"
              :checked="row.enabled"
              @change="(checked: any) => toggleEnabled(row, checked)"
            />
          </div>
          <p class="text-muted-foreground mt-1 whitespace-pre-wrap text-sm">
            {{ row.content || $t('plugin.linapro-live-manage.announcement.noContent') }}
          </p>
          <div class="mt-2 flex items-center justify-between">
            <span class="text-muted-foreground text-xs">
              {{ $t('plugin.linapro-live-manage.announcement.sort') }}: {{ row.sort }}
            </span>
            <div class="flex gap-2">
              <a-button size="small" @click="startEdit(row)">
                {{ $t('pages.common.edit') }}
              </a-button>
              <Popconfirm
                :title="$t('pages.common.deleteConfirm')"
                @confirm="handleDelete(row)"
              >
                <a-button size="small" danger>
                  {{ $t('pages.common.delete') }}
                </a-button>
              </Popconfirm>
            </div>
          </div>
        </div>
      </a-spin>
    </div>
  </Modal>
</template>
