<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { Input, message } from 'ant-design-vue';
import { useClipboard } from '@vueuse/core';

import { $t } from '#/locales';
import { useTenantStore } from '#/store/tenant';

import { buildSubscribeUrl } from './live-client';

const tenantStore = useTenantStore();
const roomCode = ref('');
const roomName = ref('');

// 订阅链接与公开接口同源同参：租户启用时必须携带当前租户 ID，
// 未启用时省略参数（公开接口会按平台租户处理）。
const tenantId = computed(() => {
  const current = tenantStore.currentTenant;
  return tenantStore.enabled && current && current.id > 0 ? current.id : undefined;
});

const subscribeUrl = computed(() => buildSubscribeUrl(roomCode.value, tenantId.value));

const title = computed(() => $t('plugin.linapro-live-manage.actions.calendar'));

const { copy } = useClipboard({ legacy: true });

async function handleCopy() {
  await copy(subscribeUrl.value);
  message.success($t('plugin.linapro-live-manage.messages.calendarCopySuccess'));
}

function selectAll(event: FocusEvent) {
  (event.target as HTMLTextAreaElement).select();
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[560px]',
  footer: false,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const data = modalApi.getData<{ roomCode?: string; roomName?: string }>();
    roomCode.value = data?.roomCode ?? '';
    roomName.value = data?.roomName ?? '';
  },
});
</script>

<template>
  <Modal :title="title">
    <div class="flex flex-col gap-3">
      <p class="text-muted-foreground text-sm">
        {{ $t('plugin.linapro-live-manage.messages.calendarHint') }}
      </p>
      <Input.TextArea :value="subscribeUrl" readonly :rows="3" @focus="selectAll" />
      <div class="flex items-center justify-between">
        <span class="text-muted-foreground text-xs">
          {{ roomName || roomCode }}
        </span>
        <a-button type="primary" @click="handleCopy">
          {{ $t('plugin.linapro-live-manage.messages.calendarCopy') }}
        </a-button>
      </div>
    </div>
  </Modal>
</template>
