import type { Id } from "@/types/common";
import type {
  AccessButtonRecord,
  AccessMenuQueryParams,
  AccessMenuRecord,
} from "@/types/access";
import axios from "axios";
import { request } from "@/utils";

export default {
  // 获取菜单树
  getMenuTree: () => request.get<AccessMenuRecord[]>("/v1/menus"),

  getComponents: () =>
    axios.get(`${import.meta.env.VITE_PUBLIC_PATH}components.json`),

  // 新增菜单
  addMenu: (data: Partial<AccessMenuRecord>) =>
    request.post<AccessMenuRecord, Partial<AccessMenuRecord>>(
      "/v1/menus",
      toMenuPayload(data),
    ),
  // 更新菜单
  updateMenu: (id: Id, data: Partial<AccessMenuRecord>) =>
    request.put<void, Partial<AccessMenuRecord>>(`/v1/menus/${id}`, toMenuPayload(data)),

  // 删除菜单
  deleteMenu: (id: Id) => request.delete<void>(`/v1/menus/${id}`),
};

function toMenuPayload(data: Partial<AccessMenuRecord>): Partial<AccessMenuRecord> {
  const payload = { ...data };
  if (payload.parentId !== undefined) {
    payload.parent_id = payload.parentId;
    delete payload.parentId;
  }
  if (payload.activeMenu !== undefined) {
    payload.active_menu = payload.activeMenu;
    delete payload.activeMenu;
  }
  if (payload.alwaysShow !== undefined) {
    payload.always_show = payload.alwaysShow;
    delete payload.alwaysShow;
  }
  return payload;
}
