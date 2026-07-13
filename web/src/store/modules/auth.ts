import { defineStore } from "pinia";
import type { TokenPayload } from "@/types/app";
import {
  useAccessStore,
  useRouterStore,
  useTabStore,
  useUserStore,
} from "@/store";

export const useAuthStore = defineStore("auth", {
  state: (): TokenPayload => ({
    accessToken: undefined,
    refreshToken: undefined,
  }),
  actions: {
    setToken({ accessToken, refreshToken }: TokenPayload) {
      const nextAccessToken = accessToken;
      const nextRefreshToken = refreshToken;
      if (nextAccessToken) this.accessToken = nextAccessToken;
      if (nextRefreshToken) this.refreshToken = nextRefreshToken;
    },
    resetToken() {
      this.$reset();
    },
    toLogin() {
      const { router, route } = useRouterStore();
      router.replace({
        path: "/login",
        query: route.query,
      });
    },
    async switchCurrentRole(data: TokenPayload) {
      this.resetLoginState();
      await nextTick();
      this.setToken(data);
    },
    resetLoginState() {
      const { resetUser } = useUserStore();
      const { resetRouter } = useRouterStore();
      const { resetAccess, accessRoutes } = useAccessStore();
      const { resetTabs } = useTabStore();
      // 重置路由
      resetRouter(accessRoutes);
      // 重置用户
      resetUser();
      // 重置访问菜单和动态路由
      resetAccess();
      // 重置Tabs
      resetTabs();
      // 重置token
      this.resetToken();
    },
    async logout() {
      this.resetLoginState();
      this.toLogin();
    },
  },
  persist: {
    key: "vue-naivue-admin_auth",
  },
});
