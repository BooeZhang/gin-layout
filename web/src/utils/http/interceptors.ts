import type {
  AxiosError,
  AxiosInstance,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from "axios";
import type { TokenPayload } from "@/types/app";
import type { ApiResult } from "@/types/common";
import { useAuthStore } from "@/store";
import { resolveResError } from "./helpers";

interface ApiErrorPayload extends Omit<Partial<ApiResult>, "code"> {
  code?: number | string;
  message?: string;
  [key: string]: unknown;
}

// 认证过期的code
const AUTH_EXPIRED_CODES = [401, 4011, 4012, 4013, 4014, 4015, 4016];
let refreshTokenPromise: Promise<TokenPayload> | null = null;

function normalizeTokenPayload(payload: TokenPayload): TokenPayload {
  return {
    ...payload,
    accessToken: payload.accessToken ?? payload.access_token,
    refreshToken: payload.refreshToken ?? payload.refresh_token,
  };
}

export function setupInterceptors(axiosInstance: AxiosInstance): void {
  const SUCCESS_CODES = [0, 200];

  async function resResolve(response: AxiosResponse<ApiErrorPayload>) {
    const { data, status, config, statusText, headers } = response;
    if (String(headers["content-type"] ?? "").includes("json")) {
      if (data?.code == null) {
        return Promise.resolve(data);
      }

      if (SUCCESS_CODES.includes(Number(data.code))) {
        return Promise.resolve(data.data);
      }
      const code = data?.code ?? status;

      const refreshResult = await tryRefreshAndReplay(
        axiosInstance,
        config,
        code,
      );
      if (refreshResult) return refreshResult;

      const needTip = config?.needTip !== false;

      // 根据code处理对应的操作，并返回处理后的message
      const message = resolveResError(
        code,
        data?.message ?? statusText,
        needTip,
      );

      return Promise.reject({ code, message, error: data ?? response });
    }
    return Promise.resolve(data ?? response);
  }

  async function resReject(error: AxiosError<ApiErrorPayload>) {
    if (!error || !error.response) {
      const code = error?.code;
      /** 根据code处理对应的操作，并返回处理后的message */
      const message = resolveResError(code, error.message);
      return Promise.reject({ code, message, error });
    }

    const { data, status, config } = error.response;
    const code = data?.code ?? status;

    const refreshResult = await tryRefreshAndReplay(
      axiosInstance,
      config,
      code,
    );
    if (refreshResult) return refreshResult;

    const needTip = config?.needTip !== false;
    const message = resolveResError(
      code,
      data?.message ?? error.message,
      needTip,
    );
    return Promise.reject({
      code,
      message,
      error: error.response?.data || error.response,
    });
  }

  axiosInstance.interceptors.request.use(reqResolve, reqReject);
  axiosInstance.interceptors.response.use(
    resResolve as unknown as (response: AxiosResponse) => AxiosResponse,
    resReject,
  );
}

function reqResolve(config: InternalAxiosRequestConfig) {
  // 处理不需要token的请求
  if (config.needToken === false) {
    return config;
  }

  const { accessToken } = useAuthStore();
  if (accessToken) {
    // token: Bearer + xxx
    config.headers.Authorization = `Bearer ${accessToken}`;
  }

  return config;
}

function reqReject(error: AxiosError) {
  return Promise.reject(error);
}

function isAuthExpiredCode(code: number | string | undefined): boolean {
  return AUTH_EXPIRED_CODES.includes(Number(code));
}

async function tryRefreshAndReplay(
  axiosInstance: AxiosInstance | undefined,
  config: InternalAxiosRequestConfig | undefined,
  code: number | string | undefined,
) {
  if (
    !axiosInstance ||
    !config ||
    !isAuthExpiredCode(code) ||
    config.authRetry ||
    config.skipAuthRefresh
  ) {
    return undefined;
  }

  const authStore = useAuthStore();
  if (!authStore.refreshToken) {
    return logoutAndReject();
  }

  config.authRetry = true;

  let tokenPayload: TokenPayload;
  try {
    tokenPayload = normalizeTokenPayload(
      await refreshTokenOnce(axiosInstance, authStore.refreshToken),
    );
  } catch {
    return logoutAndReject();
  }

  if (!tokenPayload.accessToken || !tokenPayload.refreshToken) {
    return logoutAndReject();
  }

  authStore.setToken(tokenPayload);
  config.headers.Authorization = `Bearer ${tokenPayload.accessToken}`;
  return axiosInstance(config);
}

async function logoutAndReject() {
  await useAuthStore().logout();
  return Promise.reject({
    code: 401,
    message: "登录已过期，请重新登录",
  });
}

async function refreshTokenOnce(
  axiosInstance: AxiosInstance,
  refreshToken: string,
): Promise<TokenPayload> {
  if (!refreshTokenPromise) {
    refreshTokenPromise = axiosInstance
      .post<
        unknown,
        TokenPayload
      >("/auth/refresh-token", { refresh_token: refreshToken }, { needToken: false, needTip: false, skipAuthRefresh: true })
      .finally(() => {
        refreshTokenPromise = null;
      });
  }
  return refreshTokenPromise;
}
