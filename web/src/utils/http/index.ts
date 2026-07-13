import type {
  AxiosInstance,
  AxiosRequestConfig,
  CreateAxiosDefaults,
} from "axios";
import axios from "axios";
import { setupInterceptors } from "./interceptors";

export interface RequestInstance extends Omit<
  AxiosInstance,
  "request" | "get" | "delete" | "head" | "options" | "post" | "put" | "patch"
> {
  request<T = unknown>(config: AxiosRequestConfig): Promise<T>;
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
  head<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
  options<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>;
  post<T = unknown, D = unknown>(
    url: string,
    data?: D,
    config?: AxiosRequestConfig<D>,
  ): Promise<T>;
  put<T = unknown, D = unknown>(
    url: string,
    data?: D,
    config?: AxiosRequestConfig<D>,
  ): Promise<T>;
  patch<T = unknown, D = unknown>(
    url: string,
    data?: D,
    config?: AxiosRequestConfig<D>,
  ): Promise<T>;
}

export function createAxios(
  options: CreateAxiosDefaults = {},
): RequestInstance {
  const defaultOptions = {
    baseURL: import.meta.env.VITE_AXIOS_BASE_URL,
    timeout: 12000,
  };
  const service = axios.create({
    ...defaultOptions,
    ...options,
  });
  setupInterceptors(service);
  return service as RequestInstance;
}

export const request = createAxios();

export const mockRequest = createAxios({
  baseURL: "/mock-api",
});
