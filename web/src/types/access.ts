import type { Component } from 'vue'
import type { Id } from './common'
import type { LayoutName } from './route'

export type AccessMenuType = 'catalog' | 'menu' | 'button'

export interface AccessMenuRecord {
  id?: Id
  parentId?: Id | null
  parent_id?: Id | null
  code: string
  name: string
  routeName?: string
  type: AccessMenuType
  path?: string
  redirect?: string
  icon?: string
  component?: string
  layout?: LayoutName | string
  activeMenu?: string
  active_menu?: string
  link?: string
  query?: string
  remark?: string
  level?: number
  hidden?: boolean
  cache?: boolean
  affix?: boolean
  breadcrumb?: boolean
  alwaysShow?: boolean
  always_show?: boolean
  external?: boolean
  iframe?: boolean
  enabled?: boolean
  method?: string
  apiPath?: string
  permCode?: string
  description?: string
  sort?: number
  children?: AccessMenuRecord[]
}

export interface AccessButtonRecord {
  code: string
  name: string
}

export interface SideMenuRecord {
  label?: string
  key?: string
  path?: string
  icon?: () => Component
  order: number
  children?: SideMenuRecord[]
}

export interface AccessMenuQueryParams {
  parentId?: Id
}
