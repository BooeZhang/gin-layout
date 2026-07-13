import type { LoginPayload, TokenPayload } from "@/types/auth";
import { request } from "@/utils";

export default {
  login: (data: LoginPayload) =>
    request.post<TokenPayload, LoginPayload>("/auth/login", data, {
      needToken: false,
    }),

  // 刷新token
  refreshToken: (refreshToken: string) =>
    request.post<TokenPayload, { refresh_token: string }>(
      "/auth/refresh-token",
      { refresh_token: refreshToken },
      { needToken: false, needTip: false, skipAuthRefresh: true },
    ),

  // 登出
  logout: () => request.post("/auth/logout", {}, { needTip: false }),
};
