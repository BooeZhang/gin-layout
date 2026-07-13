import type { Id, PageParams, PageResult } from "@/types/common";
import type { AccessMenuRecord } from "@/types/access";
import type { RolePayload, RoleRecord } from "@/types/role";
import type { UserInfo } from "@/types/user";
import { request } from "@/utils";

interface RoleUsersPayload {
  userIds: Id[];
}

export default {
  // 创建角色
  create: (data: RolePayload) =>
    request.post<RoleRecord, RolePayload>("/v1/roles", data),
  // 获取角色列表
  list: (params: PageParams = {}) =>
    request.get<PageResult<RoleRecord>>("/v1/roles", { params }),
  getOne: (id: Id) => request.get<RoleRecord>(`/v1/roles/${id}`),
  // 更新角色
  update: (data: RolePayload) =>
    request.put<void, RolePayload>(`/v1/roles/${data.id}`, data),
  // 删除角色
  delete: (id: Id) => request.delete<void>(`/v1/roles/${id}`),
  // 添加角色用户
  addRoleUsers: (roleId: Id, data: RoleUsersPayload) =>
    request.put<void, RoleUsersPayload>(`/v1/roles/user-add/${roleId}`, data),
  getAllAccessMenuTree: () => request.get<AccessMenuRecord[]>("/v1/menus"),
  getRoleMenuTree: () => request.get<AccessMenuRecord[]>("/v1/menus"),
  getAllUsers: (params: PageParams = {}) =>
    request.get<PageResult<UserInfo>>("/v1/users", { params }),

  removeRoleUsers: (roleId: Id, data: RoleUsersPayload) =>
    request.put<void, RoleUsersPayload>(`/v1/roles/user-remove/${roleId}`, data),
};
