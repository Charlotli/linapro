<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { Input, message } from 'ant-design-vue';
import { useClipboard } from '@vueuse/core';
import { useQRCode } from '@vueuse/integrations/useQRCode';

import { $t } from '#/locales';
import { useTenantStore } from '#/store/tenant';

import { buildWatchUrl } from './live-client';

const tenantStore = useTenantStore();
const roomCode = ref('');
const roomName = ref('');

// 观播链接与公开播放接口同源同参：租户启用时必须携带当前租户 ID，
// 未启用时省略参数（H5 页面会按平台租户处理）。
const tenantId = computed(() => {
  const current = tenantStore.currentTenant;
  return tenantStore.enabled && current && current.id > 0 ? current.id : undefined;
});

const watchUrl = computed(() => buildWatchUrl(roomCode.value, tenantId.value));

const qrcode = useQRCode(watchUrl, {
  errorCorrectionLevel: 'M',
  margin: 2,
  width: 220,
});

const title = computed(() => $t('plugin.linapro-live-manage.actions.qrcode'));

const { copy } = useClipboard({ legacy: true });

async function handleCopy() {
  await copy(watchUrl.value);
  message.success($t('plugin.linapro-live-manage.messages.watchCopySuccess'));
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
    <div class="flex flex-col items-center gap-3">
      <p class="text-muted-foreground text-sm">
        {{ $t('plugin.linapro-live-manage.messages.qrcodeHint') }}
      </p>
      <img
        v-if="qrcode"
        :src="qrcode"
        :alt="$t('plugin.linapro-live-manage.actions.qrcode')"
        class="h-[220px] w-[220px]"
      />
      <Input.TextArea :value="watchUrl" readonly :rows="2" @focus="selectAll" />
      <div class="flex w-full items-center justify-between">
        <span class="text-muted-foreground text-xs">
          {{ roomName || roomCode }}
        </span>
        <a-button type="primary" @click="handleCopy">
          {{ $t('plugin.linapro-live-manage.messages.watchCopy') }}
        </a-button>
      </div>
    </div>
  </Modal>
</template>
