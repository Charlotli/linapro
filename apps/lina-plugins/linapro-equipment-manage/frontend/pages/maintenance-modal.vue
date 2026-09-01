<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  DatePicker,
  Form,
  FormItem,
  Input,
  InputNumber,
  message,
  Select,
  Textarea,
} from 'ant-design-vue';

import { $t } from '#/locales';
import {
  equipmentOptions,
  maintenanceAdd,
  maintenanceInfo,
  maintenanceUpdate,
} from './equipment-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface FormData {
  id?: number;
  equipmentId?: number;
  maintType: number;
  maintDate: string;
  maintainer: string;
  cost: number;
  content: string;
  result: string;
  remark: string;
}

const defaultValues: FormData = {
  id: undefined,
  equipmentId: undefined,
  maintType: 1,
  maintDate: '',
  maintainer: '',
  cost: 0,
  content: '',
  result: '',
  remark: '',
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-equipment-manage.drawer.maintEditTitle')
    : $t('plugin.linapro-equipment-manage.drawer.maintCreateTitle'),
);

const formRules = reactive({
  equipmentId: [
    { message: $t('plugin.linapro-equipment-manage.validation.equipment'), required: true },
  ],
  maintType: [
    { message: $t('plugin.linapro-equipment-manage.validation.maintType'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const dictStore = useDictStore();
const typeOptions = ref<{ label: string; value: number }[]>([]);
interface EquipmentListOption {
  id: number;
  equipmentName: string;
  equipmentCode: string;
}

interface EquipmentListOption {
  id: number;
  equipmentCode: string;
  equipmentName: string;
}

const equipmentOptionList = ref<EquipmentListOption[]>([]);

async function loadEquipmentOptions() {
  const res = await equipmentOptions();
  equipmentOptionList.value = res?.list ?? [];
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[680px]',
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const dicts = await dictStore.getDictOptionsAsync(
      'plugin_equipment_maint_type',
    );
    typeOptions.value = dicts.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    await loadEquipmentOptions();
    const data = modalApi.getData();
    if (data?.id) {
      modalApi.setState({ confirmLoading: true });
      try {
        const record = await maintenanceInfo(data.id);
        formData.value = {
          id: record.id,
          equipmentId: record.equipmentId,
          maintType: record.maintType,
          maintDate: record.maintDate || '',
          maintainer: record.maintainer || '',
          cost: record.cost,
          content: record.content || '',
          result: record.result || '',
          remark: record.remark || '',
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
      await maintenanceUpdate(id, values);
      message.success($t('pages.common.updateSuccess'));
    } else {
      await maintenanceAdd(values);
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
          :label="$t('plugin.linapro-equipment-manage.fields.equipment')"
          v-bind="validateInfos.equipmentId"
        >
          <Select
            v-model:value="formData.equipmentId"
            :disabled="!!formData.id"
            :filterOption="false"
            :options="
              equipmentOptionList.map((equipment) => ({
                label: `${equipment.equipmentName}（${equipment.equipmentCode}）`,
                value: equipment.id,
              }))
            "
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.equipment')"
            show-search
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.maintType')"
          v-bind="validateInfos.maintType"
        >
          <Select
            v-model:value="formData.maintType"
            :options="typeOptions"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.maintType')"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.maintDate')"
        >
          <DatePicker
            v-model:value="formData.maintDate"
            class="w-full"
            :placeholder="
              $t('plugin.linapro-equipment-manage.placeholders.maintDate')
            "
            value-format="YYYY-MM-DD"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.maintainer')"
        >
          <Input
            v-model:value="formData.maintainer"
            :maxlength="64"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.cost')">
          <InputNumber
            v-model:value="formData.cost"
            class="w-full"
            :min="0"
            :precision="2"
            :formatter="
              (value: any) => {
                const text = String(value ?? '');
                return text === ''
                  ? ''
                  : `¥ ${text}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
              }
            "
            :parser="(value: any) => Number(String(value).replace(/[^\d.]/g, ''))"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-equipment-manage.fields.content')">
        <Textarea v-model:value="formData.content" :rows="3" />
      </FormItem>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.result')">
          <Input
            v-model:value="formData.result"
            :maxlength="256"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.result')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.remark')">
          <Input v-model:value="formData.remark" :maxlength="512" />
        </FormItem>
      </div>
    </Form>
  </Modal>
</template>
