import { pluginApiPath, requestClient } from '#/api/request';

const pluginID = 'linapro-live-manage';

function liveApi(pathName: string) {
  return pluginApiPath(pluginID, pathName);
}

export interface LiveRoom {
  id: number;
  roomCode: string;
  roomName: string;
  roomType: number;
  status: number;
  description: string;
  createdBy: number;
  createdByName: string;
  updatedBy: number;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface LiveRoomOption {
  id: number;
  roomCode: string;
  roomName: string;
  status: number;
}

export interface LiveRoomListParams {
  pageNum?: number;
  pageSize?: number;
  roomName?: string;
  roomType?: number;
  status?: number;
}

export interface LiveSong {
  name: string;
  singer: string;
  order: number;
}

export interface LiveContent {
  id: number;
  roomId: number;
  roomName: string;
  title: string;
  content: string;
  liveDate: string;
  pageUrl: string;
  pushUrl: string;
  liveUrl: string;
  coverUrl: string;
  songName: string;
  songList: string;
  leadSinger: string;
  accompaniment: string;
  host: string;
  sermonTitle: string;
  preacher: string;
  preacherIdentity: string;
  scriptureRef: string;
  scriptureContent: string;
  outline: string;
  deviceInfo: string;
  reception: string;
  state: number;
  isPublic: number;
  startTime: number | null;
  createdBy: number;
  createdByName: string;
  updatedBy: number;
  createdAt: number | null;
  updatedAt: number | null;
}

export interface LiveContentListParams {
  pageNum?: number;
  pageSize?: number;
  title?: string;
  roomId?: number;
  state?: number;
  isPublic?: number;
  liveDateStart?: string;
  liveDateEnd?: string;
}

export async function liveroomList(params?: LiveRoomListParams) {
  const res = await requestClient.get<{ list: LiveRoom[]; total: number }>(
    liveApi('liveroom'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export function liveroomOptions(keyword?: string) {
  return requestClient.get<{ list: LiveRoomOption[] }>(
    liveApi('liveroom/options'),
    { params: { keyword } },
  );
}

export function liveroomAdd(data: Partial<LiveRoom>) {
  return requestClient.post(liveApi('liveroom'), data);
}

export function liveroomUpdate(id: number, data: Partial<LiveRoom>) {
  return requestClient.put(liveApi(`liveroom/${id}`), data);
}

export function liveroomDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(liveApi('liveroom'), {
    params: { ids: list },
  });
}

export function liveroomInfo(id: number) {
  return requestClient.get<LiveRoom>(liveApi(`liveroom/${id}`));
}

export async function liveList(params?: LiveContentListParams) {
  const res = await requestClient.get<{ list: LiveContent[]; total: number }>(
    liveApi('live'),
    { params },
  );
  return { items: res.list, total: res.total };
}

export function liveAdd(data: Partial<LiveContent>) {
  return requestClient.post(liveApi('live'), data);
}

export function liveUpdate(id: number, data: Partial<LiveContent>) {
  return requestClient.put(liveApi(`live/${id}`), data);
}

export function liveDelete(ids: number[] | string) {
  const list =
    typeof ids === 'string'
      ? ids
          .split(',')
          .map((part) => Number(part.trim()))
          .filter((id) => Number.isFinite(id) && id > 0)
      : ids;
  return requestClient.delete(liveApi('live'), {
    params: { ids: list },
  });
}

export function liveInfo(id: number) {
  return requestClient.get<LiveContent>(liveApi(`live/${id}`));
}

export function liveStart(id: number) {
  return requestClient.put(liveApi(`live/${id}/start`));
}

export function liveStop(id: number) {
  return requestClient.put(liveApi(`live/${id}/stop`));
}
