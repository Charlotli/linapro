<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  DatePicker,
  Form,
  FormItem,
  Input,
  message,
  Select,
  Textarea,
} from 'ant-design-vue';

import { $t } from '#/locales';
import { bookOptions, borrowAdd } from './book-client';

const emit = defineEmits<{ reload: [] }>();

interface FormData {
  bookId?: number;
  borrower: string;
  borrowDate: string;
  dueDate: string;
  remark: string;
}

const defaultValues: FormData = {
  bookId: undefined,
  borrower: '',
  borrowDate: '',
  dueDate: '',
  remark: '',
};

const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  $t('plugin.linapro-book-manage.drawer.borrowTitle'),
);

const formRules = reactive({
  bookId: [
    { message: $t('plugin.linapro-book-manage.validation.book'), required: true },
  ],
  borrower: [
    { message: $t('plugin.linapro-book-manage.validation.borrower'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const availableBooks = ref<{
  id: number;
  title: string;
  availableQuantity: number;
}[]>([]);

async function loadAvailableBooks() {
  const res = await bookOptions();
  availableBooks.value = (res?.list ?? []).filter(
    (book) => book.availableQuantity > 0,
  );
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[600px]',
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    await loadAvailableBooks();
    formData.value = { ...defaultValues };
    resetFields();
  },
});

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();

    await borrowAdd({
      bookId: Number(formData.value.bookId),
      borrower: formData.value.borrower,
      borrowDate: formData.value.borrowDate || undefined,
      dueDate: formData.value.dueDate || undefined,
      remark: formData.value.remark,
    });
    message.success($t('plugin.linapro-book-manage.messages.borrowSuccess'));
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
        :label="$t('plugin.linapro-book-manage.fields.book')"
        v-bind="validateInfos.bookId"
      >
        <Select
          v-model:value="formData.bookId"
          :options="
            availableBooks.map((book) => ({
              label: `${book.title}（${$t(
                'plugin.linapro-book-manage.fields.available',
              )} ${book.availableQuantity}）`,
              value: book.id,
            }))
          "
          :placeholder="$t('plugin.linapro-book-manage.placeholders.book')"
          show-search
          option-filter-prop="label"
        />
      </FormItem>
      <div class="grid lg:grid-cols-2 sm:grid-cols-1">
        <FormItem
          :label="$t('plugin.linapro-book-manage.fields.borrower')"
          v-bind="validateInfos.borrower"
        >
          <Input
            v-model:value="formData.borrower"
            :maxlength="64"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.borrower')"
          />
        </FormItem>
        <FormItem :label="$t('plugin.linapro-book-manage.fields.dueDate')">
          <DatePicker
            v-model:value="formData.dueDate"
            class="w-full"
            :placeholder="$t('plugin.linapro-book-manage.placeholders.dueDate')"
            value-format="YYYY-MM-DD"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-book-manage.fields.remark')">
        <Textarea v-model:value="formData.remark" :maxlength="512" :rows="2" />
      </FormItem>
    </Form>
  </Modal>
</template>
