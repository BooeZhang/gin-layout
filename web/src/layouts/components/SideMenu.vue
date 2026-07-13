<template>
  <n-menu
    ref="menu"
    class="side-menu"
    accordion
    :indent="18"
    :collapsed-icon-size="22"
    :collapsed-width="64"
    :collapsed="appStore.collapsed"
    :options="accessStore.sideMenus"
    :value="activeKey"
    @update:value="handleMenuSelect"
  />
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAccessStore, useAppStore } from '@/store'

interface MenuOption {
  path?: string
  [key: string]: unknown
}

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const accessStore = useAccessStore()

const activeKey = computed(() => route.meta?.parentKey || route.name)

const menu = ref<{ showOption: () => void } | null>(null)
watch(route, async () => {
  await nextTick()
  menu.value?.showOption()
})

function handleMenuSelect(_key: string, item: MenuOption) {
  if (!item.path)
    return
  router.push(item.path)
}
</script>

<style>
.side-menu {
  --n-item-text-color: var(--admin-muted) !important;
  --n-item-text-color-hover: var(--admin-text) !important;
}

.side-menu:not(.n-menu--collapsed) {
  .n-menu-item-content {
    height: 40px;
    margin: 2px 8px;
    border-radius: 10px;
    transition: background-color 0.2s ease, color 0.2s ease;
    &::before {
      left: 0;
      right: 0;
      border-radius: 10px;
      background-color: transparent;
    }
    &:hover::before {
      background-color: var(--admin-hover);
    }
    &.n-menu-item-content--selected::before {
      background-color: var(--admin-active);
      border-left: 3px solid rgb(var(--primary-color));
    }
  }
}
</style>
