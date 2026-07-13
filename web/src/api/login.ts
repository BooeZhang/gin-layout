import type { TokenPayload } from "@/types/auth";
import { request } from "@/utils";

interface ToggleRolePayload {
  roleId?: number | string;
  [key: string]: unknown;
}

export default {
  toggleRole: (data: ToggleRolePayload) =>
    request.post<TokenPayload, ToggleRolePayload>("/auth/role/toggle", data),
};
