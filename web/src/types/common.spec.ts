import type { PageParams, PageResult } from './common'

type Expect<T extends true> = T
type HasKey<T, K extends PropertyKey> = K extends keyof T ? true : false
type Not<T extends boolean> = T extends true ? false : true

type _PageParamsUseBackendFields = Expect<
  HasKey<PageParams, 'page'> extends true
    ? HasKey<PageParams, 'pageSize'> extends true
      ? true
      : false
    : false
>
type _PageParamsDoNotExposePageNo = Expect<Not<HasKey<PageParams, 'pageNo'>>>

type _PageResultUseBackendFields = Expect<
  HasKey<PageResult, 'items'> extends true
    ? HasKey<PageResult, 'total'> extends true
      ? HasKey<PageResult, 'page'> extends true
        ? HasKey<PageResult, 'pageSize'> extends true
          ? true
          : false
        : false
      : false
    : false
>
type _PageResultDoesNotExposeLegacyCollections = Expect<
  Not<
    HasKey<PageResult, 'pageData'> extends true
      ? true
      : HasKey<PageResult, 'list'> extends true
        ? true
        : HasKey<PageResult, 'records'>
  >
>

const requestParams: PageParams = { page: 1, pageSize: 20 }
const responsePage: PageResult<string> = {
  items: ['alice'],
  total: 1,
  page: 1,
  pageSize: 20,
}

void requestParams
void responsePage
