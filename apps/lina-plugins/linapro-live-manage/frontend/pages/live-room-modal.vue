<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { Form, FormItem, Input, message, RadioGroup, Select } from 'ant-design-vue';

import { $t } from '#/locales';
import { liveroomAdd, liveroomInfo, liveroomUpdate } from './live-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface FormData {
  id?: number;
  roomCode: string;
  roomName: string;
  roomType: number;
  status: number;
  description: string;
}

const defaultValues: FormData = {
  id: undefined,
  roomCode: '',
  roomName: '',
  roomType: 1,
  status: 0,
  description: '',
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-live-manage.drawer.roomEditTitle')
    : $t('plugin.linapro-live-manage.drawer.roomCreateTitle'),
);

const formRules = reactive({
  roomCode: [
    { message: $t('plugin.linapro-live-manage.validation.roomCode'), required: true },
  ],
  roomName: [
    { message: $t('plugin.linapro-live-manage.validation.roomName'), required: true },
  ],
  roomType: [
    { message: $t('plugin.linapro-live-manage.validation.roomType'), required: true },
  ],
  status: [
    { message: $t('plugin.linapro-live-manage.validation.status'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const dictStore = useDictStore();
const roomTypeOptions = ref<{ label: string; value: number }[]>([]);
const roomStatusOptions = ref<{ label: string; value: number }[]>([]);

const [Modal, modalApi] = useVbenModal({
  class: 'w-[640px]',
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const [roomTypes, roomStatuses] = await Promise.all([
      dictStore.getDictOptionsAsync('plugin_live_room_type'),
      dictStore.getDictOptionsAsync('plugin_live_room_status'),
    ]);
    roomTypeOptions.value = roomTypes.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    roomStatusOptions.value = roomStatuses.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    const data = modalApi.getData();
    if (data?.id) {
      modalApi.setState({ confirmLoading: true });
      try {
        const record = await liveroomInfo(data.id);
        formData.value = {
          id: record.id,
          roomCode: record.roomCode,
          roomName: record.roomName,
          roomType: record.roomType,
          status: record.status,
          description: record.description || '',
        };
      } finally {
        modalApi.setState({ confirmLoading: false });
      }
    } else {
      formData.value = { ...defaultValues };
      resetFields();
    }
  },
});

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();

    const { id, ...values } = formData.value;
    if (isEdit.value && id) {
      await liveroomUpdate(id, values);
      message.success($t('pages.common.updateSuccess'));
    } else {
      await liveroomAdd(values);
      message.success($t('pages.common.createSuccess'));
    }
    emit('reload');
    modalApi.close();
  } catch (error) {
    console.error(error);
  } finally {
    modalApi.lock(false);
  }
}
</script>

<template>
  <Modal :title="title">
    <Form layout="vertical">
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.roomCode')"
          v-bind="validateInfos.roomCode"
        >
          <Input
            v-model:value="formData.roomCode"
            :maxlength="64"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.roomCode')"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.roomName')"
          v-bind="validateInfos.roomName"
        >
          <Input
            v-model:value="formData.roomName"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.roomName')"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.roomType')"
          v-bind="validateInfos.roomType"
        >
          <Select
            v-model:value="formData.roomType"
            :options="roomTypeOptions"
            :placeholder="$t('plugin.linapro-live-manage.placeholders.roomType')"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-live-manage.fields.status')"
          v-bind="validateInfos.status"
        >
          <RadioGroup
            v-model:value="formData.status"
            button-style="solid"
            option-type="button"
            :options="roomStatusOptions"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-live-manage.fields.description')">
        <Input.TextArea
          v-model:value="formData.description"
          :maxlength="512"
          :placeholder="$t('plugin.linapro-live-manage.placeholders.description')"
          :rows="3"
          show-count
        />
      </FormItem>
    </Form>
  </Modal>
</template>
