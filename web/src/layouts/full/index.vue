<template>
  <div class="admin-shell wh-full flex">
    <aside
      class="admin-sidebar flex-col flex-shrink-0 transition-width-300"
      :class="appStore.collapsed ? 'w-64' : 'w-220'"
    >
      <SideBar />
    </aside>

    <article class="admin-content w-0 flex-col flex-1">
      <AppHeader class="admin-header h-60 flex-shrink-0" />
      <div class="full-layout-tabs p-12">
        <AppTab class="flex-shrink-0" />
      </div>
      <slot />
    </article>
  </div>
</template>

<script setup lang="ts">
import { AppTab } from '@/layouts/components'
import { useAppStore } from '@/store'
import AppHeader from './header/index.vue'
import SideBar from './sidebar/index.vue'

const appStore = useAppStore()
</script>

<style>
.collapsed {
  width: 64px;
}

.full-layout-tabs {
  border-bottom: 1px solid var(--admin-border);
  background: var(--admin-surface);
}
</style>
