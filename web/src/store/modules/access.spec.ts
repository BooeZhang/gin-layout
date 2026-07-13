import { createPinia, setActivePinia } from 'pinia'
import { useAccessStore } from './access'
import type { AccessMenuRecord } from '@/types/app'

function expect(condition: unknown, message: string) {
  if (!condition)
    throw new Error(message)
}

function findMenuByKey(items: Array<{ key?: string, children?: any[] }>, key: string): any {
  for (const item of items) {
    if (item.key === key)
      return item
    const found = findMenuByKey(item.children || [], key)
    if (found)
      return found
  }
  return null
}

setActivePinia(createPinia())

const accessStore = useAccessStore()
const accessMenus: AccessMenuRecord[] = [
  {
    id: 1,
    parentId: null,
    name: '系统管理',
    routeName: 'SysMgt',
    code: 'SysMgt',
    type: 'catalog',
    path: '/pms',
    component: '',
    icon: 'i-fe:grid',
    sort: 2,
    hidden: false,
    cache: true,
    enabled: true,
    children: [
      {
        id: 2,
        parentId: 1,
        name: '用户管理',
        routeName: 'UserMgt',
        code: 'UserMgt',
        type: 'menu',
        path: '/pms/user',
        component: '/src/views/pms/user/index.vue',
        icon: 'i-fe:user',
        sort: 1,
        hidden: false,
        cache: true,
        enabled: true,
        children: [
          {
            id: 3,
            parentId: 2,
            name: '创建新用户',
            routeName: 'AddUser',
            code: 'AddUser',
            type: 'button',
            method: 'POST',
            sort: 1,
            hidden: false,
            cache: false,
            enabled: true,
          },
        ],
      },
    ],
  },
  {
    id: 4,
    parentId: null,
    name: '个人资料',
    routeName: 'UserProfile',
    code: 'UserProfile',
    type: 'menu',
    path: '/profile',
    component: '/src/views/profile/index.vue',
    icon: 'i-fe:user',
    sort: 99,
    hidden: true,
    cache: true,
    enabled: true,
  },
]

accessStore.setAccessMenus(accessMenus)

expect(!!findMenuByKey(accessStore.sideMenus, 'SysMgt'), 'visible catalog menu should be generated')
expect(!!findMenuByKey(accessStore.sideMenus, 'UserMgt'), 'visible child menu should be generated')
expect(!findMenuByKey(accessStore.sideMenus, 'UserProfile'), 'hidden menu should not be shown in side menu')
expect(accessStore.accessRoutes.some(route => route.name === 'UserMgt'), 'visible menu route should be registered')
expect(accessStore.accessRoutes.some(route => route.name === 'UserProfile'), 'hidden menu route should still be registered')
expect(
  accessStore.accessRoutes.some(route => route.name === 'UserMgt' && route.meta.btns?.some(btn => btn.code === 'AddUser')),
  'button access codes should be attached to parent route',
)
