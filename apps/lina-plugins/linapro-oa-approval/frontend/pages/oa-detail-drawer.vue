<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import {
  Button,
  Descriptions,
  DescriptionsItem,
  Input,
  message,
  Modal,
  Popconfirm,
  Select,
  Textarea,
  Timeline,
  TimelineItem,
} from 'ant-design-vue';

import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

import { formatRmbUppercase, formatYuan } from './oa-format';
import { DictTag } from '#/components/dict';
import { useDictStore } from '#/store/dict';
import { useUserStore } from '@vben/stores';

import type { RequestFormField } from './oa-client';
import {
  requestApprove,
  requestAppender,
  requestComment,
  requestInfo,
  requestReject,
  requestResubmit,
  requestWithdraw,
  userOptions,
  type ApprovalRequest,
  type RequestDetail as RequestDetailData,
} from './oa-client';
import { buildApprovalPrintDocument } from './oa-print-document';

const emit = defineEmits<{ reload: [] }>();

const dictStore = useDictStore();
const userStore = useUserStore();
const flowTypeDicts = ref<any[]>([]);
const statusDicts = ref<any[]>([]);
const actionDicts = ref<any[]>([]);

const request = ref<ApprovalRequest | null>(null);
const nodes = ref<{ order: number; approverId: number; approverName: string }[]>([]);
const records = ref<any[]>([]);
const formFields = ref<RequestFormField[]>([]);
const formValues = ref<Record<string, any>>({});

const replyComment = ref('');
const appendApproverId = ref<number | undefined>(undefined);
const userOptionsList = ref<{ id: number; name: string }[]>([]);

const requestNo = computed(() =>
  request.value ? `OA-${String(request.value.id).padStart(8, '0')}` : '',
);

const currentUserId = computed(() => {
  const info: any = userStore.userInfo || {};
  return Number(info.userId ?? info.id ?? 0);
});

const isPending = computed(() => request.value?.status === 1);
const canResubmit = computed(
  () =>
    isApplicant.value &&
    (request.value?.status === 3 || request.value?.status === 4),
);
const isApplicant = computed(() => request.value?.applicantId === currentUserId.value);
const isCurrentApprover = computed(
  () => request.value?.status === 1 && request.value?.currentApproverId === currentUserId.value,
);

onMounted(async () => {
  [flowTypeDicts.value, statusDicts.value, actionDicts.value] = await Promise.all([
    dictStore.getDictOptionsAsync('plugin_oa_approval_flow_type'),
    dictStore.getDictOptionsAsync('plugin_oa_approval_status'),
    dictStore.getDictOptionsAsync('plugin_oa_approval_action'),
  ]);
});

function dictLabel(dicts: any[], value: number): string {
  const hit = dicts.find(
    (item: any) => String(item.value) === String(value),
  );
  return hit?.label ?? String(value);
}

function renderFieldValue(field: RequestFormField): string {
  const value = formValues.value[field.key];
  if (value === undefined || value === null || value === '') return '';
  if (Array.isArray(value)) {
    if (field.type === 'attachment') {
      return value.join('; ');
    }
    return value.length === 0 ? '' : JSON.stringify(value);
  }
  return String(value);
}

function hasNumberColumn(field: RequestFormField): boolean {
  return field.columns.some((column) => column.type === 'number');
}

function detailColumnSum(field: RequestFormField, columnKey: string): string {
  const rows = (formValues.value[field.key] as any[]) ?? [];
  const total = rows.reduce((sum, row) => {
    const value = Number(row[columnKey]);
    return Number.isFinite(value) ? sum + value : sum;
  }, 0);
  const text = total.toLocaleString('zh-CN', { minimumFractionDigits: 2 });
  const uppercase = formatRmbUppercase(total);
  return uppercase ? `${text}（${uppercase}）` : text;
}

function handlePrint() {
  if (!request.value) return;
  const html = buildApprovalPrintDocument({
    request: request.value,
    flowTypeLabel: dictLabel(flowTypeDicts.value, request.value.flowType),
    statusLabel: dictLabel(statusDicts.value, request.value.status),
    nodes: nodes.value,
    records: records.value.map((record: any) => ({
      id: record.id,
      actorName: record.actorName,
      comment: record.comment,
      createdAt: record.createdAt,
      actionLabel: dictLabel(actionDicts.value, record.action),
    })),
    formFields: formFields.value,
    form: formValues.value,
  });
  // Assign the document before insertion: appending an empty iframe first
  // fires an extra load event, which printed a blank page once and the real
  // document again (the duplicate print dialogs the user reported).
  const previous = document.getElementById('oa-approval-print-frame');
  previous?.remove();
  const frame = document.createElement('iframe');
  frame.id = 'oa-approval-print-frame';
  frame.style.position = 'fixed';
  frame.style.right = '0';
  frame.style.bottom = '0';
  frame.style.width = '0';
  frame.style.height = '0';
  frame.style.border = '0';
  frame.onload = () => {
    frame.contentWindow?.focus();
    frame.contentWindow?.print();
    window.setTimeout(() => frame.remove(), 1000);
  };
  frame.srcdoc = html;
  document.body.append(frame);
}

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[860px]',
  footer: false,
  onOpenChange: async (isOpen: boolean) => {
    if (!isOpen) return;
    const data = drawerApi.getData<{ id: number }>();
    if (!data?.id) return;
    const detail: RequestDetailData = await requestInfo(data.id);
    request.value = detail as unknown as ApprovalRequest;
    applyDetail(detail);
    nodes.value = detail.nodes ?? [];
    records.value = detail.records ?? [];
    replyComment.value = '';
    appendApproverId.value = undefined;
  },
});

function parseAttachments(raw: string): string[] {
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed.map(String) : [];
  } catch {
    return [];
  }
}

function applyDetail(detail: RequestDetailData) {
  request.value = detail as unknown as ApprovalRequest;
  nodes.value = detail.nodes ?? [];
  records.value = detail.records ?? [];
  formFields.value = detail.formFields ?? [];
  formValues.value = detail.form ?? {};
  replyComment.value = '';
  appendApproverId.value = undefined;
}

async function refreshDetail() {
  if (!request.value) return;
  const detail: RequestDetailData = await requestInfo(request.value.id);
  applyDetail(detail);
  nodes.value = detail.nodes ?? [];
  records.value = detail.records ?? [];
}

async function withReload(action: () => Promise<unknown>, successTip: string) {
  try {
    await action();
    message.success(successTip);
    await refreshDetail();
    emit('reload');
  } catch (error) {
    console.error(error);
  }
}

function handleApprove() {
  if (!request.value) return;
  withReload(
    () => requestApprove(request.value!.id, replyComment.value),
    $t('plugin.linapro-oa-approval.messages.approveSuccess'),
  );
}

function handleReject() {
  if (!request.value) return;
  Modal.confirm({
    title: $t('plugin.linapro-oa-approval.messages.rejectConfirmTitle'),
    content: $t('plugin.linapro-oa-approval.messages.rejectConfirmContent'),
    okType: 'danger',
    onOk: () =>
      withReload(
        () => requestReject(request.value!.id, replyComment.value || $t('plugin.linapro-oa-approval.messages.rejectDefaultComment')),
        $t('plugin.linapro-oa-approval.messages.rejectSuccess'),
      ),
  });
}

function handleComment() {
  if (!request.value || !replyComment.value.trim()) {
    message.warning($t('plugin.linapro-oa-approval.validation.comment'));
    return;
  }
  withReload(
    () => requestComment(request.value!.id, replyComment.value),
    $t('plugin.linapro-oa-approval.messages.commentSuccess'),
  );
}

function handleAppend() {
  if (!request.value || !appendApproverId.value) {
    message.warning($t('plugin.linapro-oa-approval.validation.appendApprover'));
    return;
  }
  withReload(
    () => requestAppender(request.value!.id, Number(appendApproverId.value)),
    $t('plugin.linapro-oa-approval.messages.appendSuccess'),
  );
}

function handleWithdraw() {
  if (!request.value) return;
  withReload(
    () => requestWithdraw(request.value!.id),
    $t('plugin.linapro-oa-approval.messages.withdrawSuccess'),
  );
}

function handleResubmit() {
  if (!request.value) return;
  withReload(
    () => requestResubmit(request.value!.id),
    $t('plugin.linapro-oa-approval.messages.resubmitSuccess'),
  );
}

async function loadUserOptions(keyword?: string) {
  const res = await userOptions(keyword);
  userOptionsList.value = res?.list ?? [];
}
</script>

<template>
  <Drawer class="w-[860px]" :title="$t('plugin.linapro-oa-approval.drawer.detailTitle')">
    <div
      v-if="request"
      class="mb-4 rounded-lg border border-solid border-gray-100 bg-gradient-to-r from-blue-50 to-indigo-50 p-4"
    >
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <span class="truncate text-lg font-semibold">
              {{ request.title }}
            </span>
            <DictTag :dicts="statusDicts" :value="String(request.status)" />
          </div>
          <div class="text-muted-foreground mt-1 text-xs">
            {{ $t('plugin.linapro-oa-approval.print.requestNo') }}: {{ requestNo }}
            · {{ request.applicantName }}
            · {{ formatTimestamp(request.createdAt) }}
          </div>
        </div>
        <div class="shrink-0 text-right">
          <div
            class="text-2xl font-bold tabular-nums text-red-600"
            style="font-variant-numeric: tabular-nums"
          >
            {{ formatYuan(request.amount) }}
          </div>
          <div class="text-muted-foreground text-xs">
            {{ formatRmbUppercase(request.amount) }}
          </div>
        </div>
      </div>
    </div>

    <Descriptions v-if="request" :column="3" bordered size="small">
      <DescriptionsItem :label="$t('plugin.linapro-oa-approval.fields.flowType')">
        <DictTag :dicts="flowTypeDicts" :value="String(request.flowType)" />
      </DescriptionsItem>
      <DescriptionsItem
        v-if="request.status === 1"
        :label="$t('plugin.linapro-oa-approval.fields.currentApprover')"
      >
        {{ request.currentApproverName }}
      </DescriptionsItem>
      <DescriptionsItem :label="$t('plugin.linapro-oa-approval.fields.status')">
        <DictTag :dicts="statusDicts" :value="String(request.status)" />
      </DescriptionsItem>
      <DescriptionsItem :label="$t('plugin.linapro-oa-approval.fields.content')" :span="3">
        {{ request.content }}
      </DescriptionsItem>
      <DescriptionsItem
        v-if="parseAttachments(request.attachments).length > 0"
        :label="$t('plugin.linapro-oa-approval.fields.attachments')"
        :span="3"
      >
        <div class="flex flex-col gap-1">
          <a
            v-for="(url, index) in parseAttachments(request.attachments)"
            :key="index"
            :href="url"
            rel="noopener"
            target="_blank"
          >
            {{ url }}
          </a>
        </div>
      </DescriptionsItem>
    </Descriptions>

    <div v-if="formFields.length > 0" class="mt-4">
      <div class="mb-2 font-medium">
        {{ $t('plugin.linapro-oa-approval.messages.formContentTitle') }}
      </div>
      <table class="w-full border-collapse text-sm">
        <tbody>
          <tr
            v-for="field in formFields"
            :key="field.key"
            v-show="renderFieldValue(field) !== ''"
          >
            <td class="w-40 border border-solid border-gray-200 bg-gray-50 px-2 py-1 font-medium">
              {{ field.label }}
            </td>
            <td class="border border-solid border-gray-200 px-2 py-1">
              <template v-if="field.type === 'detail'">
                <table class="w-full border-collapse text-xs">
                  <thead>
                    <tr>
                      <th
                        v-for="column in field.columns"
                        :key="column.key"
                        class="border border-solid border-gray-200 px-2 py-1 text-left"
                      >
                        {{ column.label }}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="(row, rowIndex) in
                        (formValues[field.key] as any[]) ?? []"
                      :key="rowIndex"
                    >
                      <td
                        v-for="column in field.columns"
                        :key="column.key"
                        class="border border-solid border-gray-200 px-2 py-1"
                      >
                        {{ row[column.key] ?? '-' }}
                      </td>
                    </tr>
                    <tr v-if="hasNumberColumn(field)">
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
                    </tr>
                  </tbody>
                </table>
              </template>
              <template v-else>
                {{ renderFieldValue(field) }}
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="mt-4">
      <div class="mb-2 font-medium">
        {{ $t('plugin.linapro-oa-approval.messages.progressTitle') }}
      </div>
      <Timeline>
        <TimelineItem
          v-for="node in nodes"
          :key="node.order"
          :color="
            node.order < (request?.currentNodeOrder ?? 0) || request?.status === 2
              ? 'green'
              : node.order === request?.currentNodeOrder && request?.status === 1
                ? 'blue'
                : 'gray'
          "
        >
          {{ $t('plugin.linapro-oa-approval.messages.nodeLabel', { order: node.order }) }}
          · {{ node.approverName }}
          <span
            v-if="node.order === request?.currentNodeOrder && request?.status === 1"
            class="ml-1"
          >
            ({{ $t('plugin.linapro-oa-approval.messages.currentNode') }})
          </span>
        </TimelineItem>
      </Timeline>
    </div>

    <div class="mt-4">
      <div class="mb-2 font-medium">
        {{ $t('plugin.linapro-oa-approval.messages.timelineTitle') }}
      </div>
      <Timeline>
        <TimelineItem v-for="record in records" :key="record.id">
          <DictTag :dicts="actionDicts" :value="String(record.action)" />
          · {{ record.actorName }}
          <span class="text-muted-foreground ml-1 text-xs">
            {{ formatTimestamp(record.createdAt) }}
          </span>
          <div
            v-if="record.comment"
            class="text-muted-foreground mt-1 text-xs"
          >
            {{ record.comment }}
          </div>
        </TimelineItem>
      </Timeline>
    </div>

    <div class="mt-4 text-right">
      <Button @click="handlePrint">
        {{ $t('plugin.linapro-oa-approval.print.printButton') }}
      </Button>
    </div>

    <div v-if="isPending" class="mt-4 flex flex-col gap-2">
      <Textarea
        v-model:value="replyComment"
        :maxlength="2000"
        :placeholder="$t('plugin.linapro-oa-approval.placeholders.reply')"
        :rows="3"
        show-count
      />
      <div class="flex flex-wrap items-center gap-2">
        <Button
          v-if="isCurrentApprover"
          type="primary"
          @click="handleApprove"
        >
          {{ $t('plugin.linapro-oa-approval.actions.approve') }}
        </Button>
        <Button
          v-if="isCurrentApprover"
          danger
          @click="handleReject"
        >
          {{ $t('plugin.linapro-oa-approval.actions.reject') }}
        </Button>
        <Button
          v-if="isApplicant || isCurrentApprover"
          @click="handleComment"
        >
          {{ $t('plugin.linapro-oa-approval.actions.comment') }}
        </Button>
        <Popconfirm
          v-if="isApplicant && isPending"
          :title="$t('plugin.linapro-oa-approval.messages.withdrawConfirm')"
          @confirm="handleWithdraw"
        >
          <Button danger>
            {{ $t('plugin.linapro-oa-approval.actions.withdraw') }}
          </Button>
        </Popconfirm>
        <Popconfirm
          v-if="canResubmit"
          :title="$t('plugin.linapro-oa-approval.messages.resubmitConfirm')"
          @confirm="handleResubmit"
        >
          <Button type="primary">
            {{ $t('plugin.linapro-oa-approval.actions.resubmit') }}
          </Button>
        </Popconfirm>
      </div>
      <div v-if="isCurrentApprover" class="flex items-center gap-2">
        <span class="shrink-0">
          {{ $t('plugin.linapro-oa-approval.messages.appendLabel') }}
        </span>
        <Select
          v-model:value="appendApproverId"
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
          @focus="loadUserOptions()"
          @search="loadUserOptions"
        />
        <Button @click="handleAppend">
          {{ $t('plugin.linapro-oa-approval.actions.append') }}
        </Button>
      </div>
    </div>
  </Drawer>
</template>
