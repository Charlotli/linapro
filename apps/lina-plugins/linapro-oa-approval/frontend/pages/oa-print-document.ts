/**
 * Builds the standalone A4 print document for one approval request. Shared by
 * callers that push the document into a hidden iframe and invoke the browser
 * print dialog, so the printed output is always identical to the prepared
 * document.
 */
import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

import { formatRmbUppercase, formatYuan } from './oa-format';

export interface PrintNode {
  order: number;
  approverId: number;
  approverName: string;
  statusLabel?: string;
}

export interface PrintRecord {
  id: number;
  actionLabel: string;
  actorName: string;
  comment: string;
  createdAt: number | null;
}

export interface PrintFormField {
  key: string;
  label: string;
  type: string;
  columns: { key: string; label: string; type: string }[];
}

export interface ApprovalPrintInput {
  request: {
    id: number;
    title: string;
    amount: number;
    content: string;
    attachments: string;
    status: number;
    currentNodeOrder: number;
    applicantName: string;
    createdAt: number | null;
  };
  flowTypeLabel: string;
  statusLabel: string;
  nodes: PrintNode[];
  records: PrintRecord[];
  formFields: PrintFormField[];
  form: Record<string, any>;
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}

function parseAttachmentUrls(raw: string): string[] {
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed.map(String) : [];
  } catch {
    return [];
  }
}

function chainStatusLabel(
  node: PrintNode,
  request: ApprovalPrintInput['request'],
): string {
  if (node.statusLabel) return escapeHtml(node.statusLabel);
  if (node.order < (request.currentNodeOrder ?? 0) || request.status === 2) {
    return escapeHtml($t('plugin.linapro-oa-approval.print.chainDone'));
  }
  if (node.order === request.currentNodeOrder && request.status === 1) {
    return escapeHtml($t('plugin.linapro-oa-approval.print.chainPending'));
  }
  return escapeHtml($t('plugin.linapro-oa-approval.print.chainWaiting'));
}

export function buildApprovalPrintDocument(input: ApprovalPrintInput): string {
  const { request } = input;
  const requestNo = `OA-${String(request.id).padStart(8, '0')}`;

  const rows: string[] = [
    `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.title'))}</th>
  <td>${escapeHtml(request.title)}</td>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.flowType'))}</th>
  <td>${escapeHtml(input.flowTypeLabel)}</td>
</tr>`,
    `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.amount'))}</th>
  <td>${escapeHtml(formatYuan(request.amount))}</td>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.applicant'))}</th>
  <td>${escapeHtml(request.applicantName)}</td>
</tr>`,
    `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.status'))}</th>
  <td>${escapeHtml(input.statusLabel)}</td>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.requestDate'))}</th>
  <td>${escapeHtml(formatTimestamp(request.createdAt))}</td>
</tr>`,
  ];

  const contentRow = `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.content'))}</th>
  <td colspan="3">${escapeHtml(request.content ?? '')}</td>
</tr>`;

  const attachmentUrls = parseAttachmentUrls(request.attachments);
  const attachmentRows =
    attachmentUrls.length > 0
      ? `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.fields.attachments'))}</th>
  <td colspan="3">${attachmentUrls
    .map((url, index) => `${index + 1}. ${escapeHtml(url)}`)
    .join('<br/>')}</td>
</tr>`
      : '';

  const amountUppercaseRow = `<tr>
  <th>${escapeHtml($t('plugin.linapro-oa-approval.print.amountInWords'))}</th>
  <td colspan="3">${escapeHtml(formatRmbUppercase(request.amount))}</td>
</tr>`;

  const dynamicRows = input.formFields
    .map((field) => {
      if (field.type === 'detail') {
        const rows = (Array.isArray(input.form[field.key])
          ? input.form[field.key]
          : []) as Record<string, any>[];
        if (rows.length === 0) return '';
        const headers = field.columns
          .map((column) => `<th>${escapeHtml(column.label)}</th>`)
          .join('');
        const body = rows
          .map(
            (row) =>
              `<tr>${field.columns
                .map(
                  (column) =>
                    `<td>${escapeHtml(String(row[column.key] ?? '-'))}</td>`,
                )
                .join('')}</tr>`,
          )
          .join('');
        const totals = field.columns.some((column) => column.type === 'number')
          ? `<tr>${field.columns
              .map((column) => {
                if (column.type !== 'number') return '<td></td>';
                const total = rows.reduce((sum, row) => {
                  const cell = Number(row[column.key]);
                  return Number.isFinite(cell) ? sum + cell : sum;
                }, 0);
                return `<td><b>${escapeHtml($t('plugin.linapro-oa-approval.messages.total'))}: ${formatYuan(total)}</b>（${escapeHtml(formatRmbUppercase(total))}）</td>`;
              })
              .join('')}</tr>`
          : '';
        return `<div class="section">${escapeHtml(field.label)}</div>
  <table>
    <thead><tr>${headers}</tr></thead>
    <tbody>${body}${totals}</tbody>
  </table>`;
      }
      const raw = input.form[field.key];
      let text = '';
      if (Array.isArray(raw)) {
        text = raw.map((item) => escapeHtml(String(item))).join('<br/>');
      } else if (raw !== undefined && raw !== null && raw !== '') {
        text = escapeHtml(String(raw));
      }
      if (text === '') return '';
      return `<tr>
  <th>${escapeHtml(field.label)}</th>
  <td colspan="3">${text}</td>
</tr>`;
    })
    .join('');

  const chainHeaders = [
    escapeHtml($t('plugin.linapro-oa-approval.print.chainOrder')),
    escapeHtml($t('plugin.linapro-oa-approval.fields.currentApprover')),
    escapeHtml($t('plugin.linapro-oa-approval.print.chainStatus')),
    escapeHtml($t('plugin.linapro-oa-approval.print.signTitle')),
    escapeHtml($t('plugin.linapro-oa-approval.print.chainDate')),
  ];
  const chainRows = input.nodes
    .map(
      (node) => `<tr>
  <td class="center">${escapeHtml($t('plugin.linapro-oa-approval.messages.nodeLabel', { order: node.order }))}</td>
  <td>${escapeHtml(node.approverName)}</td>
  <td>${chainStatusLabel(node, request)}</td>
  <td class="sign"></td>
  <td></td>
</tr>`,
    )
    .join('');

  const recordHeaders = [
    '#',
    escapeHtml($t('plugin.linapro-oa-approval.print.timelineAction')),
    escapeHtml($t('plugin.linapro-oa-approval.print.timelineActor')),
    escapeHtml($t('plugin.linapro-oa-approval.print.timelineComment')),
    escapeHtml($t('plugin.linapro-oa-approval.print.chainDate')),
  ];
  const recordRows = input.records
    .map(
      (record, index) => `<tr>
  <td class="center">${index + 1}</td>
  <td>${escapeHtml(record.actionLabel)}</td>
  <td>${escapeHtml(record.actorName)}</td>
  <td>${escapeHtml(record.comment ?? '')}</td>
  <td>${escapeHtml(formatTimestamp(record.createdAt))}</td>
</tr>`,
    )
    .join('');

  return `<!doctype html>
<html>
<head>
<meta charset="utf-8"/>
<title>${escapeHtml($t('plugin.linapro-oa-approval.print.docTitle'))} ${requestNo}</title>
<style>
  @page { size: A4; margin: 14mm 12mm; }
  * { box-sizing: border-box; }
  body { font-family: "Microsoft YaHei", "PingFang SC", sans-serif; color: #000; margin: 0; }
  h1 { font-size: 22px; text-align: center; margin: 0 0 4px; }
  .no { text-align: right; font-size: 12px; margin-bottom: 10px; }
  .section { font-size: 14px; font-weight: 700; margin: 14px 0 6px; }
  table { width: 100%; border-collapse: collapse; font-size: 12px; }
  th, td { border: 1px solid #000; padding: 6px 8px; text-align: left; vertical-align: top; }
  th { background: #f2f2f2; font-weight: 600; white-space: nowrap; }
  td.center, th.center { text-align: center; }
  td.sign { width: 34mm; }
  .footer { display: flex; justify-content: space-between; margin-top: 16px; font-size: 11px; }
</style>
</head>
<body>
  <h1>${escapeHtml($t('plugin.linapro-oa-approval.print.docTitle'))}</h1>
  <div class="no">${escapeHtml($t('plugin.linapro-oa-approval.print.requestNo'))}: ${requestNo}</div>
  <table><tbody>
    ${rows.join('')}
    ${amountUppercaseRow}
    ${contentRow}
    ${attachmentRows}
  </tbody></table>
  <div class="section">${escapeHtml($t('plugin.linapro-oa-approval.print.chainTitle'))}</div>
  <table>
    <thead><tr>${chainHeaders.map((header) => `<th>${header}</th>`).join('')}</tr></thead>
    <tbody>${chainRows}</tbody>
  </table>
  <div class="section">${escapeHtml($t('plugin.linapro-oa-approval.print.timelineTitle'))}</div>
  <table>
    <thead><tr>${recordHeaders.map((header) => `<th>${header}</th>`).join('')}</tr></thead>
    <tbody>${recordRows}</tbody>
  </table>
  <div class="footer">
    <span>${escapeHtml($t('plugin.linapro-oa-approval.print.printTime'))}: ${escapeHtml(formatTimestamp(Date.now()))}</span>
    <span>${escapeHtml($t('plugin.linapro-oa-approval.print.archiveNote'))}</span>
  </div>
</body>
</html>`;
}
