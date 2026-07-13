/// <reference types="vite/client" />

import type {
  computed as vueComputed,
  h as vueH,
  nextTick as vueNextTick,
  ref as vueRef,
  unref as vueUnref,
  watch as vueWatch,
} from 'vue'
import type { useRoute as vueUseRoute, useRouter as vueUseRouter } from 'vue-router'

import type {
  DialogInstance,
  KeyedMessageOptions,
} from '@/utils/naiveTools'

export {}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'

  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, any>
  export default component
}

declare module '*.webp'
declare module '*.png'
declare module '*.svg'

declare module 'virtual:page-pathes' {
  const pagePathes: Record<string, () => Promise<unknown>>
  export default pagePathes
}

declare module 'virtual:icons' {
  const icons: Record<string, string>
  export default icons
}

declare module 'axios' {
  interface AxiosRequestConfig {
    needToken?: boolean
    needTip?: boolean
    authRetry?: boolean
    skipAuthRefresh?: boolean
  }
}

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    layout?: string
    keepAlive?: boolean
    parentKey?: string
    order?: number
  }
}

declare global {
  interface ImportMetaEnv {
    readonly VITE_TITLE: string
    readonly VITE_AXIOS_BASE_URL: string
    readonly VITE_PUBLIC_PATH: string
    readonly VITE_PROXY_TARGET: string
    readonly VITE_USE_HASH: string
  }

  interface Window {
    $loadingBar?: any
    $message?: GlobalMessageApi
    $dialog?: DialogInstance
    $notification?: any
  }

  const $loadingBar: any
  const $message: GlobalMessageApi
  const $dialog: DialogInstance
  const $notification: any
  const ref: typeof vueRef
  const computed: typeof vueComputed
  const watch: typeof vueWatch
  const nextTick: typeof vueNextTick
  const unref: typeof vueUnref
  const h: typeof vueH
  const useRoute: typeof vueUseRoute
  const useRouter: typeof vueUseRouter

  interface GlobalMessageApi {
    loading(content: string | string[], option?: KeyedMessageOptions): unknown
    success(content: string | string[], option?: KeyedMessageOptions): unknown
    error(content: string | string[], option?: KeyedMessageOptions): unknown
    info(content: string | string[], option?: KeyedMessageOptions): unknown
    warning(content: string | string[], option?: KeyedMessageOptions): unknown
    destroy(key: string, duration?: number): void
    destroyAll?: () => void
  }
}
