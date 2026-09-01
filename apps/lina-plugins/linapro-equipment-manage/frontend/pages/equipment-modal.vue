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
  RadioGroup,
  Select,
} from 'ant-design-vue';

import { $t } from '#/locales';
import {
  equipmentAdd,
  equipmentInfo,
  equipmentUpdate,
} from './equipment-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface FormData {
  id?: number;
  equipmentCode: string;
  equipmentName: string;
  equipmentType: number;
  brandModel: string;
  purchaseDate: string;
  purchasePrice: number;
  location: string;
  owner: string;
  status: number;
  remark: string;
}

const defaultValues: FormData = {
  id: undefined,
  equipmentCode: '',
  equipmentName: '',
  equipmentType: 1,
  brandModel: '',
  purchaseDate: '',
  purchasePrice: 0,
  location: '',
  owner: '',
  status: 1,
  remark: '',
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-equipment-manage.drawer.editTitle')
    : $t('plugin.linapro-equipment-manage.drawer.createTitle'),
);

const formRules = reactive({
  equipmentCode: [
    { message: $t('plugin.linapro-equipment-manage.validation.code'), required: true },
  ],
  equipmentName: [
    { message: $t('plugin.linapro-equipment-manage.validation.name'), required: true },
  ],
  equipmentType: [
    { message: $t('plugin.linapro-equipment-manage.validation.type'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const dictStore = useDictStore();
const typeOptions = ref<{ label: string; value: number }[]>([]);
const statusOptions = ref<{ label: string; value: number }[]>([]);

const [Modal, modalApi] = useVbenModal({
  class: 'w-[720px]',
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const [types, statuses] = await Promise.all([
      dictStore.getDictOptionsAsync('plugin_equipment_type'),
      dictStore.getDictOptionsAsync('plugin_equipment_status'),
    ]);
    typeOptions.value = types.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    statusOptions.value = statuses.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    const data = modalApi.getData();
    if (data?.id) {
      modalApi.setState({ confirmLoading: true });
      try {
        const record = await equipmentInfo(data.id);
        formData.value = {
          id: record.id,
          equipmentCode: record.equipmentCode,
          equipmentName: record.equipmentName,
          equipmentType: record.equipmentType,
          brandModel: record.brandModel || '',
          purchaseDate: record.purchaseDate || '',
          purchasePrice: record.purchasePrice,
          location: record.location || '',
          owner: record.owner || '',
          status: record.status,
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
      await equipmentUpdate(id, values);
      message.success($t('pages.common.updateSuccess'));
    } else {
      await equipmentAdd(values);
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
          :label="$t('plugin.linapro-equipment-manage.fields.code')"
          v-bind="validateInfos.equipmentCode"
        >
          <Input
            v-model:value="formData.equipmentCode"
            :maxlength="64"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.code')"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.name')"
          v-bind="validateInfos.equipmentName"
        >
          <Input
            v-model:value="formData.equipmentName"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.name')"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.type')"
          v-bind="validateInfos.equipmentType"
        >
          <Select
            v-model:value="formData.equipmentType"
            :options="typeOptions"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.type')"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.brandModel')"
        >
          <Input
            v-model:value="formData.brandModel"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.brandModel')"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.purchaseDate')"
        >
          <DatePicker
            v-model:value="formData.purchaseDate"
            class="w-full"
            :placeholder="
              $t('plugin.linapro-equipment-manage.placeholders.purchaseDate')
            "
            value-format="YYYY-MM-DD"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-equipment-manage.fields.purchasePrice')"
        >
          <InputNumber
            v-model:value="formData.purchasePrice"
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
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.status')">
          <RadioGroup
            v-model:value="formData.status"
            button-style="solid"
            option-type="button"
            :options="statusOptions"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.location')">
          <Input
            v-model:value="formData.location"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.location')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-equipment-manage.fields.owner')">
          <Input
            v-model:value="formData.owner"
            :maxlength="64"
            :placeholder="$t('plugin.linapro-equipment-manage.placeholders.owner')"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-equipment-manage.fields.remark')">
        <Input.TextArea
          v-model:value="formData.remark"
          :maxlength="512"
          :rows="2"
          show-count
        />
      </FormItem>
    </Form>
  </Modal>
</template>
