import { defineStore } from 'pinia'
import type { TabRecord } from '@/types/app'
import { useRouterStore } from './router'

export const useTabStore = defineStore('tab', {
  state: (): { tabs: TabRecord[], activeTab: string, reloading: boolean } => ({
    tabs: [],
    activeTab: '',
    reloading: false,
  }),
  getters: {
    activeIndex: state => state.tabs.findIndex(item => item.path === state.activeTab),
  },
  actions: {
    async setActiveTab(path: string) {
      await nextTick() // tab栏dom更新完再设置激活，让tab栏定位到新增的tab上生效
      this.activeTab = path
    },
    setTabs(tabs: TabRecord[]) {
      this.tabs = tabs
    },
    addTab(tab: TabRecord) {
      const findIndex = this.tabs.findIndex(item => item.path === tab.path)
      if (findIndex !== -1) {
        this.tabs.splice(findIndex, 1, tab)
      }
      else {
        this.setTabs([...this.tabs, tab])
      }
      this.setActiveTab(tab.path)
    },
    async reloadTab(path: string, keepAlive?: boolean) {
      const findItem = this.tabs.find(item => item.path === path)
      if (!findItem)
        return
      // 更新key可让keepAlive失效
      if (keepAlive)
        findItem.keepAlive = false
      $loadingBar.start()
      this.reloading = true
      await nextTick()
      this.reloading = false
      findItem.keepAlive = !!keepAlive
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        $loadingBar.finish()
      }, 100)
    },
    async removeTab(path: string) {
      this.setTabs(this.tabs.filter(tab => tab.path !== path))
      if (path === this.activeTab) {
        const latestTab = this.tabs[this.tabs.length - 1]
        latestTab && useRouterStore().router?.push(latestTab.path)
      }
    },
    removeOther(curPath?: string) {
      curPath = curPath ?? this.activeTab
      this.setTabs(this.tabs.filter(tab => tab.path === curPath))
      if (curPath !== this.activeTab) {
        const latestTab = this.tabs[this.tabs.length - 1]
        latestTab && useRouterStore().router?.push(latestTab.path)
      }
    },
    removeLeft(curPath: string) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index >= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab)) {
        const latestTab = filterTabs[filterTabs.length - 1]
        latestTab && useRouterStore().router?.push(latestTab.path)
      }
    },
    removeRight(curPath: string) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index <= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab)) {
        const latestTab = filterTabs[filterTabs.length - 1]
        latestTab && useRouterStore().router?.push(latestTab.path)
      }
    },
    resetTabs() {
      this.$reset()
    },
  },
  persist: {
    pick: ['tabs'],
    storage: sessionStorage,
  },
})
