import { pluginApiPath, requestClient } from '#/api/request';

const pluginID = 'linapro-oa-approval';

function oaApi(pathName: string) {
  return pluginApiPath(pluginID, pathName);
}

export interface ApprovalFlow {
  id: number;
  flowType: number;
  flowName: string;
  description: string;
  status: number;
  nodeCount: number;
  createdBy: number;
  createdByName: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface FlowNodeItem {
  order: number;
  approverId: number;
  approverName: string;
}

export interface FlowNodeInput {
  approverId: number;
}

export interface FlowColumnInput {
  key: string;
  label: string;
  type: 'text' | 'number' | 'date';
}

export interface FlowFieldInput {
  key: string;
  label: string;
  type:
    | 'text'
    | 'textarea'
    | 'number'
    | 'date'
    | 'select'
    | 'attachment'
    | 'detail';
  required: boolean;
  asAmount: boolean;
  options: string[];
  columns: FlowColumnInput[];
}

export interface FlowFieldItemFull extends FlowFieldInput {
  order?: number;
}

export interface FlowListParams {
  pageNum?: number;
  pageSize?: number;
  flowName?: string;
  flowType?: number;
  status?: number;
}

export interface UserOption {
  id: number;
  name: string;
}

export interface ApprovalRequest {
  id: number;
  flowId: number;
  flowType: number;
  title: string;
  amount: number;
  content: string;
  attachments: string;
  status: number;
  currentNodeOrder: number;
  currentApproverId: number;
  currentApproverName: string;
  applicantId: number;
  applicantName: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface ApprovalRecordItem {
  id: number;
  nodeOrder: number;
  action: number;
  actorId: number;
  actorName: string;
  comment: string;
  createdAt: number | null;
}

export interface RequestListParams {
  pageNum?: number;
  pageSize?: number;
  scope?: 'mine' | 'pending';
  title?: string;
  flowType?: number;
  status?: number;
}

export interface RequestFormFieldColumn {
  key: string;
  label: string;
  type: 'text' | 'number' | 'date';
}

export interface RequestFormField {
  key: string;
  label: string;
  type:
    | 'text'
    | 'textarea'
    | 'number'
    | 'date'
    | 'select'
    | 'attachment'
    | 'detail';
  required: boolean;
  asAmount: boolean;
  options: string[];
  columns: RequestFormFieldColumn[];
}

export interface RequestDetail {
  nodes: FlowNodeItem[];
  records: ApprovalRecordItem[];
  formFields: RequestFormField[];
  form: Record<string, any>;
  [key: string]: any;
}

export async function flowList(params?: FlowListParams) {
  const res = await requestClient.get<{ list: ApprovalFlow[]; total: number }>(
    oaApi('approval/flow'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export function flowInfo(id: number) {
  return requestClient.get<
    ApprovalFlow & {
      nodes: FlowNodeItem[];
      fields: FlowFieldItemFull[];
    }
  >(oaApi(`approval/flow/${id}`));
}

export function flowAdd(
  data: Partial<ApprovalFlow> & {
    nodes: FlowNodeInput[];
    fields?: FlowFieldInput[];
  },
) {
  return requestClient.post(oaApi('approval/flow'), data);
}

export function flowUpdate(
  id: number,
  data: Partial<ApprovalFlow> & {
    nodes?: FlowNodeInput[];
    fields?: FlowFieldInput[];
  },
) {
  return requestClient.put(oaApi(`approval/flow/${id}`), data);
}

export function flowDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(oaApi('approval/flow'), {
    params: { ids: list },
  });
}

export function userOptions(keyword?: string) {
  return requestClient.get<{ list: UserOption[] }>(
    oaApi('approval/user/options'),
    { params: { keyword } },
  );
}

export async function requestList(params?: RequestListParams) {
  const res = await requestClient.get<{ list: ApprovalRequest[]; total: number }>(
    oaApi('approval/request'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export function requestInfo(id: number) {
  return requestClient.get<RequestDetail>(oaApi(`approval/request/${id}`));
}

export function requestSubmit(data: {
  flowId: number;
  title: string;
  form?: Record<string, any>;
}) {
  return requestClient.post(oaApi('approval/request'), data);
}

export function requestApprove(id: number, comment?: string) {
  return requestClient.put(oaApi(`approval/request/${id}/approve`), {
    comment,
  });
}

export function requestReject(id: number, comment: string) {
  return requestClient.put(oaApi(`approval/request/${id}/reject`), {
    comment,
  });
}

export function requestComment(id: number, comment: string) {
  return requestClient.put(oaApi(`approval/request/${id}/comment`), {
    comment,
  });
}

export function requestAppender(id: number, approverId: number) {
  return requestClient.put(oaApi(`approval/request/${id}/appender`), {
    approverId,
  });
}

export function requestWithdraw(id: number) {
  return requestClient.put(oaApi(`approval/request/${id}/withdraw`));
}

export function requestResubmit(id: number) {
  return requestClient.put(oaApi(`approval/request/${id}/resubmit`));
}

export function pendingCount() {
  return requestClient.get<{ count: number }>(
    oaApi('approval/pending-count'),
  );
}
