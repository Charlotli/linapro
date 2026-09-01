<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  Button,
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
import { formatRmbUppercase } from './oa-format';
import {
  flowInfo,
  requestSubmit,
  type FlowFieldInput,
} from './oa-client';

const emit = defineEmits<{ reload: [] }>();

interface AttachmentRow {
  url: string;
}

interface DetailRow {
  [key: string]: any;
}

interface FormData {
  flowId?: number;
  title: string;
  form: Record<string, any>;
  detailRows: Record<string, DetailRow[]>;
  attachmentRows: Record<string, AttachmentRow[]>;
}

const defaultValues: FormData = {
  flowId: undefined,
  title: '',
  form: {},
  detailRows: {},
  attachmentRows: {},
};

const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  $t('plugin.linapro-oa-approval.drawer.submitTitle'),
);

const activeFields = ref<FlowFieldInput[]>([]);

const formRules = reactive({
  flowId: [
    { message: $t('plugin.linapro-oa-approval.validation.flowId'), required: true },
  ],
  title: [
    { message: $t('plugin.linapro-oa-approval.validation.title'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const flowOptions = ref<{ id: number; flowName: string }[]>([]);

async function loadFlowOptions() {
  const res = await flowList({ pageSize: 100, status: 1 });
  flowOptions.value = res?.items ?? [];
}

async function handleFlowChange(flowId: number) {
  formData.value.form = {};
  formData.value.detailRows = {};
  formData.value.attachmentRows = {};
  if (!flowId) {
    activeFields.value = [];
    return;
  }
  const detail = await flowInfo(flowId);
  const fields = detail?.fields ?? [];
  if (fields.length === 0) {
    // Built-in standard form fallback: amount/content/attachments.
    activeFields.value = [
      {
        key: 'amount',
        label: $t('plugin.linapro-oa-approval.fields.amount'),
        type: 'number',
        required: false,
        asAmount: true,
        options: [],
        columns: [],
      },
      {
        key: 'content',
        label: $t('plugin.linapro-oa-approval.fields.content'),
        type: 'textarea',
        required: false,
        asAmount: false,
        options: [],
        columns: [],
      },
      {
        key: 'attachments',
        label: $t('plugin.linapro-oa-approval.fields.attachments'),
        type: 'attachment',
        required: false,
        asAmount: false,
        options: [],
        columns: [],
      },
    ];
  } else {
    activeFields.value = fields;
  }
}

function fieldRowsKey(key: string) {
  return `detail-${key}`;
}

function attachmentRowsKey(key: string) {
  return `attach-${key}`;
}

function detailColumnSum(field: FlowFieldInput, columnKey: string): string {
  const rows = formData.value.detailRows[fieldRowsKey(field.key)] ?? [];
  const total = rows.reduce((sum, row) => {
    const value = Number(row[columnKey]);
    return Number.isFinite(value) ? sum + value : sum;
  }, 0);
  const text = total.toLocaleString('zh-CN', { minimumFractionDigits: 2 });
  const uppercase = formatRmbUppercase(total);
  return uppercase ? `${text}（${uppercase}）` : text;
}

function addDetailRow(field: FlowFieldInput) {
  const storeKey = fieldRowsKey(field.key);
  formData.value.detailRows[storeKey] = formData.value.detailRows[storeKey] ?? [];
  if (formData.value.detailRows[storeKey].length >= 50) {
    message.warning($t('plugin.linapro-oa-approval.messages.detailRowLimit'));
    return;
  }
  formData.value.detailRows[storeKey].push({});
}

function removeDetailRow(field: FlowFieldInput, index: number) {
  formData.value.detailRows[fieldRowsKey(field.key)].splice(index, 1);
}

function addAttachmentRow(field: FlowFieldInput) {
  const storeKey = attachmentRowsKey(field.key);
  formData.value.attachmentRows[storeKey] =
    formData.value.attachmentRows[storeKey] ?? [];
  if (formData.value.attachmentRows[storeKey].length >= 10) {
    message.warning($t('plugin.linapro-oa-approval.messages.attachmentLimit'));
    return;
  }
  formData.value.attachmentRows[storeKey].push({ url: '' });
}

function removeAttachmentRow(field: FlowFieldInput, index: number) {
  formData.value.attachmentRows[attachmentRowsKey(field.key)].splice(index, 1);
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[760px]',
  fullscreenButton: true,
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    await loadFlowOptions();
    formData.value = { ...defaultValues };
    activeFields.value = [];
    resetFields();
  },
});

function collectFormValues(): {
  values: Record<string, any>;
  error: string | null;
} {
  const values: Record<string, any> = {};
  for (const field of activeFields.value) {
    if (field.type === 'detail') {
      const rows = (formData.value.detailRows[fieldRowsKey(field.key)] ?? []).map(
        (row) => {
          const clean: Record<string, any> = {};
          for (const column of field.columns) {
            const cell = row[column.key];
            if (cell === undefined || cell === null || cell === '') continue;
            clean[column.key] = cell;
          }
          return clean;
        },
      );
      if (rows.length > 0) {
        values[field.key] = rows;
      }
      continue;
    }
    if (field.type === 'attachment') {
      const urls = (formData.value.attachmentRows[attachmentRowsKey(field.key)] ?? [])
        .map((row) => row.url.trim())
        .filter((url) => url !== '');
      if (urls.length > 0) {
        values[field.key] = urls;
      }
      continue;
    }
    const value = formData.value.form[field.key];
    if (value === undefined || value === null || value === '') continue;
    values[field.key] = value;
  }
  for (const field of activeFields.value) {
    if (!field.required) continue;
    const value = values[field.key];
    const empty =
      value === undefined ||
      value === null ||
      (typeof value === 'string' && value.trim() === '') ||
      (Array.isArray(value) && value.length === 0);
    if (empty) {
      return {
        values,
        error: $t('plugin.linapro-oa-approval.messages.requiredFieldMissing', {
          field: field.label,
        }),
      };
    }
  }
  return { values, error: null };
}

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();

    if (!formData.value.flowId) {
      message.error($t('plugin.linapro-oa-approval.validation.flowId'));
      return;
    }
    const { values, error } = collectFormValues();
    if (error) {
      message.error(error);
      return;
    }

    await requestSubmit({
      flowId: Number(formData.value.flowId),
      title: formData.value.title,
      form: values,
    });
    message.success($t('plugin.linapro-oa-approval.messages.submitSuccess'));
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
          :label="$t('plugin.linapro-oa-approval.fields.flowId')"
          v-bind="validateInfos.flowId"
        >
          <Select
            v-model:value="formData.flowId"
            :options="
              flowOptions.map((flow) => ({
                label: flow.flowName,
                value: flow.id,
              }))
            "
            :placeholder="$t('plugin.linapro-oa-approval.placeholders.flowId')"
            @change="(value: any) => handleFlowChange(Number(value))"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-oa-approval.fields.title')"
          v-bind="validateInfos.title"
        >
          <Input
            v-model:value="formData.title"
            :maxlength="256"
            :placeholder="$t('plugin.linapro-oa-approval.placeholders.title')"
          />
        </FormItem>
      </div>

      <FormItem
        v-for="field in activeFields"
        :key="field.key"
        :label="field.label"
      >
        <Select
          v-if="field.type === 'select'"
          v-model:value="formData.form[field.key]"
          :options="
            field.options.map((option) => ({ label: option, value: option }))
          "
          :placeholder="$t('plugin.linapro-oa-approval.placeholders.selectInput')"
        />
        <InputNumber
          v-else-if="field.type === 'number'"
          v-model:value="formData.form[field.key]"
          class="w-full"
          :min="0"
          :precision="2"
          :formatter="
            (value: any) => {
              const text = String(value ?? '');
              return text === '' ? '' : `¥ ${text}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
            }
          "
          :parser="(value: any) => Number(String(value).replace(/[^\d.]/g, ''))"
        />
        <DatePicker
          v-else-if="field.type === 'date'"
          v-model:value="formData.form[field.key]"
          class="w-full"
          value-format="YYYY-MM-DD"
        />
        <Textarea
          v-else-if="field.type === 'textarea'"
          v-model:value="formData.form[field.key]"
          :rows="3"
        />
        <Input
          v-else-if="field.type === 'text'"
          v-model:value="formData.form[field.key]"
          :maxlength="512"
        />
        <template v-else-if="field.type === 'attachment'">
          <div class="flex flex-col gap-2">
            <div
              v-for="(row, index) in
                formData.attachmentRows[attachmentRowsKey(field.key)] ?? []"
              :key="index"
              class="flex items-center gap-2"
            >
              <Input
                v-model:value="row.url"
                class="flex-1"
                :placeholder="
                  $t('plugin.linapro-oa-approval.placeholders.attachmentUrl')
                "
              />
              <Button danger size="small" @click="removeAttachmentRow(field, index)">
                {{ $t('pages.common.delete') }}
              </Button>
            </div>
            <Button class="w-fit" size="small" type="dashed" @click="addAttachmentRow(field)">
              {{ $t('plugin.linapro-oa-approval.messages.addAttachment') }}
            </Button>
          </div>
        </template>
        <template v-else-if="field.type === 'detail'">
          <div class="flex flex-col gap-2">
            <table class="w-full border-collapse text-sm">
              <thead>
                <tr>
                  <th
                    v-for="column in field.columns"
                    :key="column.key"
                    class="border border-solid border-gray-200 px-2 py-1 text-left"
                  >
                    {{ column.label }}
                  </th>
                  <th class="w-16 border border-solid border-gray-200"></th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(row, rowIndex) in
                    formData.detailRows[fieldRowsKey(field.key)] ?? []"
                  :key="rowIndex"
                >
                  <td
                    v-for="column in field.columns"
                    :key="column.key"
                    class="border border-solid border-gray-200 p-1"
                  >
                    <InputNumber
                      v-if="column.type === 'number'"
                      v-model:value="row[column.key]"
                      class="w-full"
                      :min="0"
                      :precision="2"
                      size="small"
                    />
                    <DatePicker
                      v-else-if="column.type === 'date'"
                      v-model:value="row[column.key]"
                      class="w-full"
                      size="small"
                      value-format="YYYY-MM-DD"
                    />
                    <Input
                      v-else
                      v-model:value="row[column.key]"
                      size="small"
                    />
                  </td>
                  <td class="border border-solid border-gray-200 p-1 text-center">
                    <Button danger size="small" @click="removeDetailRow(field, rowIndex)">
                      {{ $t('pages.common.delete') }}
                    </Button>
                  </td>
                </tr>
                <tr
                  v-if="field.columns.some((column) => column.type === 'number')"
                >
                  <td
                    v-for="column in field.columns"
                    :key="`sum-${column.key}`"
                    class="border border-solid border-gray-200 px-2 py-1 font-medium"
                  >
                    <template v-if="column.type === 'number'">
                      {{ $t('plugin.linapro-oa-approval.messages.total') }}:
                      {{ detailColumnSum(field, column.key) }}
                    </template>
                  </td>
                  <td class="w-16 border border-solid border-gray-200"></td>
                </tr>
              </tbody>
            </table>
            <Button class="w-fit" size="small" type="dashed" @click="addDetailRow(field)">
              {{ $t('plugin.linapro-oa-approval.messages.addDetailRow') }}
            </Button>
          </div>
        </template>
      </FormItem>

      <div class="text-muted-foreground text-xs">
        {{ $t('plugin.linapro-oa-approval.messages.attachmentHint') }}
      </div>
    </Form>
  </Modal>
</template>
