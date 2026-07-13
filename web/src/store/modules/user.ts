import { defineStore } from "pinia";
import type { Id, RoleRecord, UserInfo } from "@/types/app";

export const useUserStore = defineStore("user", {
  state: (): { userInfo: UserInfo | null } => ({
    userInfo: null,
  }),
  getters: {
    userId: (state) => state.userInfo?.id,
    account: (state) => state.userInfo?.account,
    username: (state) => state.userInfo?.account,
    nickName: (state) => state.userInfo?.nickName,
    avatar: (state) => state.userInfo?.avatar,
    roles: (state) => state.userInfo?.roles || ([] as RoleRecord[]),
    currentRole: (state) => state.userInfo?.roles?.[0],
  },
  actions: {
    setUser(user: UserInfo) {
      this.userInfo = user;
    },
    setUserId(id: Id) {
      this.userInfo = { ...(this.userInfo || {}), id };
    },
    resetUser() {
      this.$reset();
    },
  },
});
