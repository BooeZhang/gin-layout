<template>
  <router-link class="side-logo h-60 f-c-c" to="/">
    <TheLogo />
    <h2
      v-show="!appStore.collapsed"
      class="ml-10 max-w-140 flex-shrink-0 truncate text-16 font-bold"
    >
      {{ title }}
    </h2>
  </router-link>
</template>

<script setup lang="ts">
import { useAppStore } from '@/store'

const title = import.meta.env.VITE_TITLE

const appStore = useAppStore()
</script>

<style scoped>
.side-logo {
  color: var(--admin-text);
  border-bottom: 1px solid var(--admin-border);
}

.side-logo h2 {
  letter-spacing: 0;
}
</style>
