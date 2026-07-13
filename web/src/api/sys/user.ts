import type { PasswordPayload } from "@/types/auth";
import type { Dict, Id, PageParams, PageResult } from "@/types/common";
import type { RoleRecord } from "@/types/role";
import type { UserInfo, UserPayload } from "@/types/user";
import { request } from "@/utils";
import type { AccessMenuRecord } from "@/types/access";

export default {
  // 创建用户
  create: (data: UserPayload) =>
    request.post<UserInfo, UserPayload>("/v1/users", toUserPayload(data)),

  list: (params: PageParams = {}) =>
    request
      .get<PageResult<UserInfo>>("/v1/users", { params: toUserQuery(params) })
      .then((data) => ({
        ...data,
        items: data.items.map(normalizeUser),
      })),
  // 更新用户
  update: (data: UserPayload) =>
    request.put<void, UserPayload>(`/v1/users/${data.id}`, toUserPayload(data)),
  // 删除用户
  delete: (id: Id) => request.delete<void>(`/v1/users/${id}`),
  // getUser: (id: Id) => request.get<UserInfo>("/user/${id}"),
  // 获取用户信息
  getUser: () => request.get<UserInfo>("/v1/users/details").then(normalizeUser),
  // 获取用户菜单
  getUserMenu: () => request.get<AccessMenuRecord[]>("/v1/users/menus"),

  resetPwd: (id: Id, data: PasswordPayload) =>
    request.patch<void, PasswordPayload>(`/user/password/reset/${id}`, data),

  // 获取所有角色
  getAllRoles: () => request.get<RoleRecord[]>("/v1/roles/all"),
};

function toUserPayload(data: UserPayload): UserPayload {
  const payload = { ...data };
  if (payload.nickName !== undefined) {
    payload.nick_name = payload.nickName;
    delete payload.nickName;
  }
  return payload;
}

function toUserQuery(params: PageParams): PageParams<Dict> {
  const query: PageParams<Dict> = { ...params };
  if (query.nickName !== undefined) {
    query.nick_name = query.nickName;
    delete query.nickName;
  }
  if (query.enabled === 0) query.enabled = false;
  if (query.enabled === 1) query.enabled = true;
  return query;
}

function normalizeUser(user: UserInfo): UserInfo {
  return {
    ...user,
    nickName: user.nickName ?? user.nick_name,
    createdAt: user.createdAt ?? user.created_at,
    updatedAt: user.updatedAt ?? user.updated_at,
  };
}
