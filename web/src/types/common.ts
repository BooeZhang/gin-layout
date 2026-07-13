export type Id = number | string;

export type Dict<T = unknown> = Record<string, T>;

export interface ApiResult<T = unknown> {
  code: number;
  data?: T;
  message?: string;
}

export type PageParams<T extends Dict = Record<never, never>> = {
  page?: number;
  pageSize?: number;
} & T;

export interface PageResult<T = unknown> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}
