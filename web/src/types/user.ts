import type { Id } from "./common";
import type { RoleRecord } from "./role";

export interface UserInfo {
  id?: Id;
  account?: string;
  nickName?: string;
  nick_name?: string;
  avatar?: string;
  email?: string;
  phone?: string;
  enabled?: boolean;
  createdAt?: string;
  created_at?: string;
  updatedAt?: string;
  updated_at?: string;
  roles?: RoleRecord[];
}

export interface UserPayload extends Partial<UserInfo> {
  id?: Id;
  password?: string;
  roleIds?: Id[];
  [key: string]: unknown;
}
