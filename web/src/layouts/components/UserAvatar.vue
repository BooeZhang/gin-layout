<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div id="user-dropdown" class="flex cursor-pointer items-center">
      <n-avatar round :size="36" :src="userStore.avatar || defaultAvatar" />
      <div
        v-if="userStore.userInfo"
        class="ml-12 flex-col flex-shrink-0 items-center"
      >
        <span class="text-14">{{
          userStore.nickName ?? userStore.account
        }}</span>
      </div>
    </div>
  </n-dropdown>

  <RoleSelect ref="roleSelectRef" />
</template>

<script setup>
import api from "@/api";
import { RoleSelect } from "@/layouts/components";
import { useAccessStore, useAuthStore, useUserStore } from "@/store";
import defaultAvatar from "@/assets/images/avatar.png";

const router = useRouter();
const userStore = useUserStore();
const authStore = useAuthStore();
const accessStore = useAccessStore();

const options = reactive([
  {
    label: "个人资料",
    key: "profile",
    icon: () => h("i", { class: "i-material-symbols:person-outline text-14" }),
    show: computed(() =>
      accessStore.accessRoutes?.some((item) => item.path === "/profile"),
    ),
  },
  {
    label: "退出登录",
    key: "logout",
    icon: () => h("i", { class: "i-mdi:exit-to-app text-14" }),
  },
]);

function handleSelect(key) {
  switch (key) {
    case "profile":
      router.push("/profile");
      break;
    case "logout":
      $dialog.confirm({
        title: "提示",
        type: "info",
        content: "确认退出？",
        async confirm() {
          try {
            await api.logout();
          } catch (error) {
            console.error(error);
          }
          authStore.logout();
          $message.success("已退出登录");
        },
      });
      break;
  }
}
</script>
