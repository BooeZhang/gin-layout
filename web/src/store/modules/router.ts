import { defineStore } from 'pinia'
import type { RouteRecordNameGeneric } from 'vue-router'

export const useRouterStore = defineStore('router', () => {
  const router = useRouter()
  const route = useRoute()

  function resetRouter(accessRoutes: Array<{ name?: RouteRecordNameGeneric }>) {
    accessRoutes.forEach((item) => {
      if (item.name && router.hasRoute(item.name))
        router.removeRoute(item.name)
    })
  }

  return {
    router,
    route,
    resetRouter,
  }
})
