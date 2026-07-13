import type { RouteRecordRaw, Router } from "vue-router";
import { useAccessStore, useAuthStore, useUserStore } from "@/store";
import { getUserAccessMenus, getUserInfo } from "@/store/helper";
import { resolveViewComponent } from "./route-component";

const WHITE_LIST = ["/login", "/404"];
export function createAccessGuard(router: Router): void {
  router.beforeEach(async (to) => {
    const authStore = useAuthStore();
    const token = authStore.accessToken;

    /** 没有token */
    if (!token) {
      if (WHITE_LIST.includes(to.path)) return true;
      return { path: "login", query: { ...to.query, redirect: to.path } };
    }

    // 有token的情况
    if (to.path === "/login") return { path: "/" };
    if (WHITE_LIST.includes(to.path)) return true;

    const userStore = useUserStore();
    const accessStore = useAccessStore();
    if (!userStore.userInfo) {
      const [user, accessMenus] = await Promise.all([
        getUserInfo(),
        getUserAccessMenus(),
      ]);
      console.log(user);
      userStore.setUser(user);
      accessStore.setAccessMenus(accessMenus);
      const routeComponents = import.meta.glob("/src/views/**/*.vue");
      accessStore.accessRoutes.forEach((route) => {
        if (typeof route.component === "string") {
          route.component = resolveViewComponent(
            route.component,
            routeComponents as Record<string, RouteRecordRaw["component"]>,
          );
        }
        !router.hasRoute(route.name) &&
          router.addRoute(route as RouteRecordRaw);
      });
      return { ...to, replace: true };
    }

    const routes = router.getRoutes();
    if (routes.find((route) => route.name === to.name)) return true;

    return { name: "404", query: { path: to.fullPath } };
  });
}
