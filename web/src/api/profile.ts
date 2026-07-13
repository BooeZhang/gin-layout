import type { PasswordPayload } from '@/types/auth'
import type { UserPayload } from '@/types/user'
import { request } from '@/utils'

export default {
  changePassword: (data: PasswordPayload) =>
    request.post<void, PasswordPayload>('/auth/password', data),
  updateProfile: (data: UserPayload) =>
    request.put<void, UserPayload>(`/v1/users/${data.id}`, toProfilePayload(data)),
}

function toProfilePayload(data: UserPayload): UserPayload {
  const payload = { ...data };
  if (payload.nickName !== undefined) {
    payload.nick_name = payload.nickName;
    delete payload.nickName;
  }
  return payload;
}
