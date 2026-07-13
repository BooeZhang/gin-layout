import type { RouteRecordRaw } from 'vue-router'
import type { AccessButtonRecord } from './access'

export type LayoutName = 'full' | 'empty'

export type DynamicRouteRecord = Omit<RouteRecordRaw, 'component' | 'children'> & {
  name: string
  path: string
  component?: string | RouteRecordRaw['component']
  meta: AppRouteMeta
}

export interface AppRouteMeta {
  title?: string
  icon?: string
  layout?: LayoutName | string
  keepAlive?: boolean
  parentKey?: string
  btns?: AccessButtonRecord[]
  order?: number
  requiresAuth?: boolean
  accessCodes?: string[]
}

export interface TabRecord {
  name?: string
  path: string
  title?: string
  icon?: string
  keepAlive?: boolean
}
