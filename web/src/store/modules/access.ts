import { defineStore } from 'pinia'
import { h } from 'vue'
import type { AccessMenuRecord, DynamicRouteRecord, SideMenuRecord } from '@/types/app'

const HOME_SIDE_MENU: SideMenuRecord = {
  label: '首页',
  key: 'Home',
  path: '/',
  icon: () => h('i', { class: 'i-fe:home text-16' }),
  order: Number.MIN_SAFE_INTEGER,
}

export const useAccessStore = defineStore('access', {
  state: (): {
    accessRoutes: DynamicRouteRecord[]
    accessMenus: AccessMenuRecord[]
    sideMenus: SideMenuRecord[]
  } => ({
    accessRoutes: [],
    accessMenus: [],
    sideMenus: [],
  }),
  actions: {
    setAccessMenus(accessMenus: AccessMenuRecord[]) {
      this.accessRoutes = []
      this.accessMenus = accessMenus
      const sideMenus = this.accessMenus
        .filter(item => this.isMenuItem(item))
        .map(item => this.getSideMenuItem(item))
        .filter((item): item is SideMenuRecord => !!item)
        .sort((a, b) => a.order - b.order)
      this.sideMenus = [HOME_SIDE_MENU, ...sideMenus.filter(item => item.key !== HOME_SIDE_MENU.key)]
    },
    getSideMenuItem(item: AccessMenuRecord, parent?: SideMenuRecord): SideMenuRecord | null {
      const route = this.generateRoute(item, item.hidden ? parent?.key : null)
      if (item.enabled !== false && route.path && !route.path.startsWith('http'))
        this.accessRoutes.push(route)
      const sideMenuItem: SideMenuRecord = {
        label: route.meta.title,
        key: route.name.toString(),
        path: route.path,
        icon: () => h('i', { class: `${route.meta.icon} text-16` }),
        order: item.sort ?? 0,
      }
      const children = item.children?.filter(item => this.isMenuItem(item)) || []
      if (children.length) {
        sideMenuItem.children = children
          .map(child => this.getSideMenuItem(child, sideMenuItem))
          .filter((item): item is SideMenuRecord => !!item)
          .sort((a, b) => a.order - b.order)
        if (!sideMenuItem.children.length)
          delete sideMenuItem.children
      }
      if (item.hidden)
        return null
      return sideMenuItem
    },
    generateRoute(item: AccessMenuRecord, parentKey?: string | null): DynamicRouteRecord {
      return {
        name: item.routeName || item.code,
        path: item.path || '',
        redirect: item.redirect,
        component: item.component,
        meta: {
          icon: item.icon,
          title: item.name,
          keepAlive: !!item.cache,
          parentKey: parentKey ?? undefined,
          btns: item.children
            ?.filter(item => item.type === 'button')
            .map(item => ({ code: item.code, name: item.name })),
        },
      }
    },
    isMenuItem(item: AccessMenuRecord): boolean {
      return item.type === 'catalog' || item.type === 'menu'
    },
    resetAccess() {
      this.$reset()
    },
  },
})
