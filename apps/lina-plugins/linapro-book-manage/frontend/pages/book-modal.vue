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
import { bookAdd, bookInfo, bookUpdate } from './book-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface FormData {
  id?: number;
  title: string;
  author: string;
  isbn: string;
  category: number;
  publisher: string;
  publishDate: string;
  totalQuantity: number;
  location: string;
  status: number;
  coverUrl: string;
  remark: string;
}

const defaultValues: FormData = {
  id: undefined,
  title: '',
  author: '',
  isbn: '',
  category: 1,
  publisher: '',
  publishDate: '',
  totalQuantity: 1,
  location: '',
  status: 1,
  coverUrl: '',
  remark: '',
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-book-manage.drawer.editTitle')
    : $t('plugin.linapro-book-manage.drawer.createTitle'),
);

const formRules = reactive({
  title: [
    { message: $t('plugin.linapro-book-manage.validation.title'), required: true },
  ],
  category: [
    { message: $t('plugin.linapro-book-manage.validation.category'), required: true },
  ],
  totalQuantity: [
    {
      message: $t('plugin.linapro-book-manage.validation.totalQuantity'),
      required: true,
    },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const dictStore = useDictStore();
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const statusOptions = ref<{ label: string; value: number }[]>([]);

const [Modal, modalApi] = useVbenModal({
  class: 'w-[720px]',
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const [categories, statuses] = await Promise.all([
      dictStore.getDictOptionsAsync('plugin_book_category'),
      dictStore.getDictOptionsAsync('plugin_book_status'),
    ]);
    categoryOptions.value = categories.map((item: any) => ({
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
        const record = await bookInfo(data.id);
        formData.value = {
          id: record.id,
          title: record.title,
          author: record.author || '',
          isbn: record.isbn || '',
          category: record.category,
          publisher: record.publisher || '',
          publishDate: record.publishDate || '',
          totalQuantity: record.totalQuantity,
          location: record.location || '',
          status: record.status,
          coverUrl: record.coverUrl || '',
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
      await bookUpdate(id, values);
      message.success($t('pages.common.updateSuccess'));
    } else {
      await bookAdd(values);
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
      <FormItem
        :label="$t('plugin.linapro-book-manage.fields.title')"
        v-bind="validateInfos.title"
      >
        <Input
          v-model:value="formData.title"
          :maxlength="256"
          :placeholder="$t('plugin.linapro-book-manage.placeholders.title')"
        />
      </FormItem>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem :label="$t('plugin.linapro-book-manage.fields.author')">
          <Input
            v-model:value="formData.author"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.author')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.isbn')">
          <Input
            v-model:value="formData.isbn"
            :maxlength="32"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.isbn')"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-book-manage.fields.category')"
          v-bind="validateInfos.category"
        >
          <Select
            v-model:value="formData.category"
            :options="categoryOptions"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.category')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.publisher')">
          <Input
            v-model:value="formData.publisher"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.publisher')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.publishDate')">
          <DatePicker
            v-model:value="formData.publishDate"
            class="w-full"
            :placeholder="
              $t('plugin.linapro-book-manage.placeholders.publishDate')
            "
            value-format="YYYY-MM-DD"
          />
        </FormItem>
      </div>
      <div class="grid lg:grid-cols-3 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-book-manage.fields.totalQuantity')"
          v-bind="validateInfos.totalQuantity"
        >
          <InputNumber
            v-model:value="formData.totalQuantity"
            class="w-full"
            :min="1"
            :precision="0"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.location')">
          <Input
            v-model:value="formData.location"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.location')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.status')">
          <RadioGroup
            v-model:value="formData.status"
            button-style="solid"
            option-type="button"
            :options="statusOptions"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-book-manage.fields.coverUrl')">
        <Input
          v-model:value="formData.coverUrl"
          :placeholder="$t('plugin.linapro-book-manage.placeholders.coverUrl')"
        />
      </FormItem>
      <FormItem :label="$t('plugin.linapro-book-manage.fields.remark')">
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
