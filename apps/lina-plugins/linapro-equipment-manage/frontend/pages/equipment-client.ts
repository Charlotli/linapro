import { pluginApiPath, requestClient } from '#/api/request';

const pluginID = 'linapro-equipment-manage';

function api(pathName: string) {
  return pluginApiPath(pluginID, pathName);
}

export interface Equipment {
  id: number;
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
  createdBy: number;
  createdByName: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface EquipmentListParams {
  pageNum?: number;
  pageSize?: number;
  name?: string;
  type?: number;
  status?: number;
}

export interface MaintenanceRecord {
  id: number;
  equipmentId: number;
  equipmentName: string;
  maintType: number;
  maintDate: string;
  maintainer: string;
  cost: number;
  content: string;
  result: string;
  remark: string;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface MaintenanceListParams {
  pageNum?: number;
  pageSize?: number;
  equipmentId?: number;
  maintType?: number;
  dateStart?: string;
  dateEnd?: string;
}

export async function equipmentList(params?: EquipmentListParams) {
  const res = await requestClient.get<{ list: Equipment[]; total: number }>(
    api('equipment'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export interface EquipmentOption {
  id: number;
  equipmentCode: string;
  equipmentName: string;
  status: number;
}

export function equipmentOptions(keyword?: string) {
  return requestClient.get<{ list: EquipmentOption[] }>(
    api('equipment/options'),
    { params: { keyword } },
  );
}

export function equipmentInfo(id: number) {
  return requestClient.get<Equipment>(api(`equipment/${id}`));
}

export function equipmentAdd(data: Partial<Equipment>) {
  return requestClient.post(api('equipment'), data);
}

export function equipmentUpdate(id: number, data: Partial<Equipment>) {
  return requestClient.put(api(`equipment/${id}`), data);
}

export function equipmentDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(api('equipment'), {
    params: { ids: list },
  });
}

export function maintenanceInfo(id: number) {
  return requestClient.get<MaintenanceRecord>(api(`maintenance/${id}`));
}

export async function maintenanceList(params?: MaintenanceListParams) {
  const res = await requestClient.get<{
    list: MaintenanceRecord[];
    total: number;
  }>(api('maintenance'), { params });
  return { items: res.list, total: res.total };
}

export function maintenanceAdd(data: Partial<MaintenanceRecord>) {
  return requestClient.post(api('maintenance'), data);
}

export function maintenanceUpdate(id: number, data: Partial<MaintenanceRecord>) {
  return requestClient.put(api(`maintenance/${id}`), data);
}

export function maintenanceDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(api('maintenance'), {
    params: { ids: list },
  });
}
