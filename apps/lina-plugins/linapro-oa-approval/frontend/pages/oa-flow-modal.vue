<script setup lang="ts">
import { computed, reactive, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import {
  Button,
  Form,
  FormItem,
  Input,
  message,
  RadioGroup,
  Select,
} from 'ant-design-vue';

import { $t } from '#/locales';
import {
  flowAdd,
  flowInfo,
  flowUpdate,
  userOptions,
  type FlowColumnInput,
  type FlowFieldInput,
} from './oa-client';
import { useDictStore } from '#/store/dict';

const emit = defineEmits<{ reload: [] }>();

interface NodeRow {
  approverId?: number;
}

interface FieldRow {
  key: string;
  label: string;
  type: FlowFieldInput['type'];
  required: boolean;
  asAmount: boolean;
  optionsText: string;
  columns: FlowColumnInput[];
}

interface FormData {
  id?: number;
  flowType: number;
  flowName: string;
  description: string;
  status: number;
  nodeRows: NodeRow[];
  fieldRows: FieldRow[];
}

const fieldTypeOptions = computed(() =>
  ([
    ['text', 'fieldType.text'],
    ['textarea', 'fieldType.textarea'],
    ['number', 'fieldType.number'],
    ['date', 'fieldType.date'],
    ['select', 'fieldType.select'],
    ['attachment', 'fieldType.attachment'],
    ['detail', 'fieldType.detail'],
  ] as const).map(([value, key]) => ({
    label: $t(`plugin.linapro-oa-approval.${key}`),
    value,
  })),
);

const columnTypeOptions = computed(() =>
  (['text', 'number', 'date'] as const).map((value) => ({
    label: $t(`plugin.linapro-oa-approval.fieldType.${value}`),
    value,
  })),
);

const defaultValues: FormData = {
  id: undefined,
  flowType: 1,
  flowName: '',
  description: '',
  status: 1,
  nodeRows: [{ approverId: undefined }],
  fieldRows: [],
};

const isEdit = computed(() => !!formData.value.id);
const formData = ref<FormData>({ ...defaultValues });
const title = computed(() =>
  isEdit.value
    ? $t('plugin.linapro-oa-approval.drawer.flowEditTitle')
    : $t('plugin.linapro-oa-approval.drawer.flowCreateTitle'),
);

// Node-list and form-field checks run manually at confirm time: wiring them
// into Form.useForm re-validates on every row add/remove and shows errors
// before the user ever submits.
const formRules = reactive({
  flowType: [
    { message: $t('plugin.linapro-oa-approval.validation.flowType'), required: true },
  ],
  flowName: [
    { message: $t('plugin.linapro-oa-approval.validation.flowName'), required: true },
  ],
});

const { validate, validateInfos, resetFields } = Form.useForm(
  formData,
  formRules,
);

const nodeRowsError = ref('');

function validateNodeRows(): boolean {
  const hasApprover = formData.value.nodeRows.some((row) => !!row.approverId);
  nodeRowsError.value = hasApprover
    ? ''
    : $t('plugin.linapro-oa-approval.validation.nodes');
  return hasApprover;
}

const dictStore = useDictStore();
const flowTypeOptions = ref<{ label: string; value: number }[]>([]);
const userOptionsList = ref<{ id: number; name: string }[]>([]);

async function loadUserOptions(keyword?: string) {
  const res = await userOptions(keyword);
  userOptionsList.value = res?.list ?? [];
}

function addNodeRow() {
  formData.value.nodeRows.push({ approverId: undefined });
  validateNodeRows();
}

function removeNodeRow(index: number) {
  formData.value.nodeRows.splice(index, 1);
  validateNodeRows();
}

function moveNodeRow(index: number, offset: number) {
  const target = index + offset;
  if (target < 0 || target >= formData.value.nodeRows.length) {
    return;
  }
  const rows = formData.value.nodeRows;
  [rows[index], rows[target]] = [rows[target], rows[index]];
}

function nextFieldKey(): string {
  let index = formData.value.fieldRows.length + 1;
  let candidate = `field_${index}`;
  const used = new Set(formData.value.fieldRows.map((row) => row.key));
  while (used.has(candidate)) {
    index += 1;
    candidate = `field_${index}`;
  }
  return candidate;
}

function addFieldRow() {
  if (formData.value.fieldRows.length >= 30) {
    message.warning($t('plugin.linapro-oa-approval.messages.formFieldLimit'));
    return;
  }
  formData.value.fieldRows.push({
    key: nextFieldKey(),
    label: '',
    type: 'text',
    required: false,
    asAmount: false,
    optionsText: '',
    columns: [],
  });
}

function removeFieldRow(index: number) {
  formData.value.fieldRows.splice(index, 1);
}

function moveFieldRow(index: number, offset: number) {
  const target = index + offset;
  if (target < 0 || target >= formData.value.fieldRows.length) {
    return;
  }
  const rows = formData.value.fieldRows;
  [rows[index], rows[target]] = [rows[target], rows[index]];
}

function addColumnRow(field: FieldRow) {
  let index = field.columns.length + 1;
  let candidate = `col_${index}`;
  const used = new Set(field.columns.map((column) => column.key));
  while (used.has(candidate)) {
    index += 1;
    candidate = `col_${index}`;
  }
  field.columns.push({ key: candidate, label: '', type: 'text' });
}

function removeColumnRow(field: FieldRow, index: number) {
  field.columns.splice(index, 1);
}

const [Modal, modalApi] = useVbenModal({
  class: 'w-[880px]',
  fullscreenButton: true,
  onConfirm: handleConfirm,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const dicts = await dictStore.getDictOptionsAsync(
      'plugin_oa_approval_flow_type',
    );
    flowTypeOptions.value = dicts.map((item: any) => ({
      label: item.label,
      value: Number(item.value),
    }));
    await loadUserOptions();
    const data = modalApi.getData();
    if (data?.id) {
      modalApi.setState({ confirmLoading: true });
      try {
        const record = await flowInfo(data.id);
        formData.value = {
          id: record.id,
          flowType: record.flowType,
          flowName: record.flowName,
          description: record.description || '',
          status: record.status,
          nodeRows: (record.nodes ?? []).map((node) => ({
            approverId: node.approverId,
          })),
          fieldRows: (record.fields ?? []).map((field) => ({
            key: field.key,
            label: field.label,
            type: field.type,
            required: field.required,
            asAmount: field.asAmount,
            optionsText: (field.options ?? []).join(','),
            columns: (field.columns ?? []).map((column) => ({ ...column })),
          })),
        };
      } finally {
        modalApi.setState({ confirmLoading: false });
      }
    } else {
      formData.value = {
        ...defaultValues,
        nodeRows: [{ approverId: undefined }],
        fieldRows: [],
      };
      resetFields();
    }
  },
});

function validateFieldRows(): string | null {
  const seen = new Set<string>();
  for (const field of formData.value.fieldRows) {
    // Blank storage keys are auto-filled from the row position; users never
    // need to invent one manually.
    if (!field.key.trim()) {
      field.key = nextFieldKey();
    }
    const key = field.key.trim();
    if (!/^[a-zA-Z][a-zA-Z0-9_]{0,63}$/.test(key)) {
      return $t('plugin.linapro-oa-approval.validation.fieldKey');
    }
    if (seen.has(key)) {
      return $t('plugin.linapro-oa-approval.validation.fieldKey');
    }
    seen.add(key);
    if (!field.label.trim()) {
      return $t('plugin.linapro-oa-approval.validation.fieldLabel');
    }
    if (field.type === 'select' && field.optionsText.trim() === '') {
      return $t('plugin.linapro-oa-approval.validation.fieldOptions');
    }
    if (field.type === 'detail' && field.columns.length === 0) {
      return $t('plugin.linapro-oa-approval.validation.fieldColumns');
    }
    for (const column of field.columns) {
      if (!column.key.trim()) {
        column.key = `col_${field.columns.indexOf(column) + 1}`;
      }
      if (!column.label.trim()) {
        return $t('plugin.linapro-oa-approval.validation.fieldColumns');
      }
    }
  }
  return null;
}

async function handleConfirm() {
  try {
    modalApi.lock(true);
    await validate();

    if (!validateNodeRows()) {
      message.error(nodeRowsError.value);
      return;
    }

    const fieldError = validateFieldRows();
    if (fieldError) {
      message.error(fieldError);
      return;
    }

    const nodes = formData.value.nodeRows
      .map((row) => row.approverId)
      .filter((id) => Number.isFinite(id) && Number(id) > 0)
      .map((id) => ({ approverId: Number(id) }));
    if (nodes.length === 0) {
      message.error($t('plugin.linapro-oa-approval.validation.nodes'));
      return;
    }

    const fields: FlowFieldInput[] = formData.value.fieldRows.map((row) => ({
      key: row.key.trim(),
      label: row.label.trim(),
      type: row.type,
      required: row.required,
      asAmount: row.asAmount,
      options:
        row.type === 'select'
          ? row.optionsText
              .split(',')
              .map((option) => option.trim())
              .filter((option) => option !== '')
          : [],
      columns:
        row.type === 'detail'
          ? row.columns.map((column) => ({
              key: column.key.trim(),
              label: column.label.trim(),
              type: column.type,
            }))
          : [],
    }));

    const { id, nodeRows: _nodeRows, fieldRows: _fieldRows, ...values } = formData.value;
    if (isEdit.value && id) {
      await flowUpdate(id, { ...values, nodes, fields });
      message.success($t('pages.common.updateSuccess'));
    } else {
      await flowAdd({ ...values, nodes, fields });
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
          :label="$t('plugin.linapro-oa-approval.fields.flowType')"
          v-bind="validateInfos.flowType"
        >
          <Select
            v-model:value="formData.flowType"
            :options="flowTypeOptions"
            :placeholder="$t('plugin.linapro-oa-approval.placeholders.flowType')"
          />
        </FormItem>
        <FormItem
          :label="$t('plugin.linapro-oa-approval.fields.flowName')"
          v-bind="validateInfos.flowName"
        >
          <Input
            v-model:value="formData.flowName"
            :maxlength="128"
            :placeholder="$t('plugin.linapro-oa-approval.placeholders.flowName')"
          />
        </FormItem>
      </div>
      <FormItem :label="$t('plugin.linapro-oa-approval.fields.status')">
        <RadioGroup
          v-model:value="formData.status"
          button-style="solid"
          option-type="button"
          :options="[
            {
              label: $t('plugin.linapro-oa-approval.status.enabled'),
              value: 1,
            },
            {
              label: $t('plugin.linapro-oa-approval.status.disabled'),
              value: 0,
            },
          ]"
        />
      </FormItem>
      <FormItem :label="$t('plugin.linapro-oa-approval.fields.description')">
        <Input.TextArea
          v-model:value="formData.description"
          :maxlength="512"
          :placeholder="$t('plugin.linapro-oa-approval.placeholders.description')"
          :rows="2"
        />
      </FormItem>
      <FormItem
        :label="$t('plugin.linapro-oa-approval.fields.nodes')"
        :validate-status="nodeRowsError ? 'error' : ''"
        :help="nodeRowsError"
      >
        <div class="flex flex-col gap-2">
          <div
            v-for="(row, index) in formData.nodeRows"
            :key="index"
            class="flex items-center gap-2"
          >
            <span class="w-16 shrink-0 text-center">
              {{ $t('plugin.linapro-oa-approval.messages.nodeLabel', { order: index + 1 }) }}
            </span>
            <Select
              v-model:value="row.approverId"
              class="flex-1"
              :filterOption="false"
              :options="
                userOptionsList.map((user) => ({
                  label: user.name,
                  value: user.id,
                }))
              "
              :placeholder="$t('plugin.linapro-oa-approval.placeholders.approver')"
              show-search
              @search="loadUserOptions"
            />
            <Button
              :disabled="index === 0"
              size="small"
              @click="moveNodeRow(index, -1)"
            >
              ↑
            </Button>
            <Button
              :disabled="index === formData.nodeRows.length - 1"
              size="small"
              @click="moveNodeRow(index, 1)"
            >
              ↓
            </Button>
            <Button
              :disabled="formData.nodeRows.length === 1"
              danger
              size="small"
              @click="removeNodeRow(index)"
            >
              {{ $t('pages.common.delete') }}
            </Button>
          </div>
          <Button class="w-fit" size="small" type="dashed" @click="addNodeRow">
            {{ $t('plugin.linapro-oa-approval.messages.addApprover') }}
          </Button>
        </div>
      </FormItem>

      <FormItem :label="$t('plugin.linapro-oa-approval.fields.formFields')">
        <div class="flex flex-col gap-3">
          <div
            v-for="(field, index) in formData.fieldRows"
            :key="index"
            class="rounded border border-solid border-gray-200 p-2"
          >
            <div class="flex flex-wrap items-center gap-2">
              <span class="w-8 shrink-0 text-center">{{ index + 1 }}</span>
              <Input
                v-model:value="field.label"
                class="w-48"
                :maxlength="64"
                :placeholder="$t('plugin.linapro-oa-approval.placeholders.fieldLabel')"
              />
              <Select
                v-model:value="field.type"
                class="w-36"
                :options="fieldTypeOptions"
              />
              <label class="flex items-center gap-1 text-sm">
                <input v-model="field.required" type="checkbox" />
                {{ $t('plugin.linapro-oa-approval.fields.required') }}
              </label>
              <label
                v-if="field.type === 'number'"
                class="flex items-center gap-1 text-sm"
              >
                <input v-model="field.asAmount" type="checkbox" />
                {{ $t('plugin.linapro-oa-approval.fields.asAmount') }}
              </label>
              <Input
                v-model:value="field.key"
                class="ml-auto w-32"
                :maxlength="64"
                :placeholder="$t('plugin.linapro-oa-approval.placeholders.fieldKeyAuto')"
              />
              <Button
                :disabled="index === 0"
                size="small"
                @click="moveFieldRow(index, -1)"
              >
                ↑
              </Button>
              <Button
                :disabled="index === formData.fieldRows.length - 1"
                size="small"
                @click="moveFieldRow(index, 1)"
              >
                ↓
              </Button>
              <Button
                danger
                size="small"
                @click="removeFieldRow(index)"
              >
                {{ $t('pages.common.delete') }}
              </Button>
            </div>
            <div v-if="field.type === 'select'" class="mt-2">
              <Input
                v-model:value="field.optionsText"
                :placeholder="$t('plugin.linapro-oa-approval.placeholders.fieldOptions')"
              />
            </div>
            <div v-if="field.type === 'detail'" class="mt-2 flex flex-col gap-2">
              <div
                v-for="(column, columnIndex) in field.columns"
                :key="columnIndex"
                class="flex items-center gap-2"
              >
                <Input
                  v-model:value="column.key"
                  class="w-36"
                  :maxlength="64"
                  :placeholder="$t('plugin.linapro-oa-approval.placeholders.columnKeyAuto')"
                />
                <Input
                  v-model:value="column.label"
                  class="w-44"
                  :maxlength="64"
                  :placeholder="$t('plugin.linapro-oa-approval.placeholders.columnLabel')"
                />
                <Select
                  v-model:value="column.type"
                  class="w-32"
                  :options="columnTypeOptions"
                />
                <Button
                  danger
                  size="small"
                  @click="removeColumnRow(field, columnIndex)"
                >
                  {{ $t('pages.common.delete') }}
                </Button>
              </div>
              <Button
                class="w-fit"
                size="small"
                type="dashed"
                @click="addColumnRow(field)"
              >
                {{ $t('plugin.linapro-oa-approval.messages.addColumn') }}
              </Button>
            </div>
          </div>
          <Button class="w-fit" size="small" type="dashed" @click="addFieldRow">
            {{ $t('plugin.linapro-oa-approval.messages.addFormField') }}
          </Button>
          <div class="text-muted-foreground text-xs">
            {{ $t('plugin.linapro-oa-approval.messages.formFieldHint') }}
          </div>
        </div>
      </FormItem>
    </Form>
  </Modal>
</template>
