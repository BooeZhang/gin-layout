<template>
  <div id="top-tab">
    <n-tabs
      :value="tabStore.activeTab"
      :closable="tabStore.tabs.length > 1"
      type="card"
      @close="(path: string) => tabStore.removeTab(path)"
    >
      <n-tab
        v-for="item in tabStore.tabs"
        :key="item.path"
        :name="item.path"
        @click="handleItemClick(item.path)"
        @contextmenu.prevent="handleContextMenu($event, item)"
      >
        {{ item.title }}
      </n-tab>
    </n-tabs>

    <ContextMenu
      v-if="contextMenuOption.show"
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </div>
</template>

<script setup lang="ts">
import { nextTick, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useTabStore } from '@/store'
import type { TabRecord } from '@/types/app'
import ContextMenu from './ContextMenu.vue'

const router = useRouter()
const tabStore = useTabStore()

const contextMenuOption = reactive({
  show: false,
  x: 0,
  y: 0,
  currentPath: '',
})

function handleItemClick(path: string) {
  tabStore.setActiveTab(path)
  router.push(path)
}

function showContextMenu() {
  contextMenuOption.show = true
}
function hideContextMenu() {
  contextMenuOption.show = false
}
function setContextMenu(x: number, y: number, currentPath: string) {
  Object.assign(contextMenuOption, { x, y, currentPath })
}

// 右击菜单
async function handleContextMenu(e: MouseEvent, tagItem: TabRecord) {
  const { clientX, clientY } = e
  hideContextMenu()
  setContextMenu(clientX, clientY, tagItem.path)
  await nextTick()
  showContextMenu()
}
</script>

<style scoped>
:deep(.n-tabs) {
  .n-tabs-tab {
    align-items: center;
    padding-left: 16px;
    height: 36px;
    background: var(--admin-surface) !important;
    border: 1px solid transparent !important;
    border-radius: 10px !important;
    margin-right: 6px;
    color: var(--admin-muted);
    transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease;
    &:hover {
      color: var(--admin-text);
      border-color: rgba(var(--primary-color), 0.45) !important;
      background: var(--admin-hover) !important;
    }
  }
  .n-tabs-tab--active {
    color: rgb(var(--primary-color));
    border-color: rgba(var(--primary-color), 0.55) !important;
    background-color: var(--admin-active) !important;
  }
  .n-tabs-pad,
  .n-tabs-tab-pad,
  .n-tabs-scroll-padding {
    border: none !important;
  }
}
</style>
