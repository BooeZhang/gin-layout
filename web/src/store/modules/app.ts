import { generate, getRgbStr } from '@arco-design/color'
import { useDark } from '@vueuse/core'
import { defineStore } from 'pinia'
import type { RemovableRef } from '@vueuse/core'
import { defaultPrimaryColor, naiveThemeOverrides } from '@/settings'

interface AppState {
  collapsed: boolean
  isDark: RemovableRef<boolean>
  primaryColor: string
  naiveThemeOverrides: typeof naiveThemeOverrides
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    collapsed: false,
    isDark: useDark(),
    primaryColor: defaultPrimaryColor,
    naiveThemeOverrides,
  }),
  actions: {
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    setCollapsed(b: boolean) {
      this.collapsed = b
    },
    toggleDark() {
      this.isDark = !this.isDark
    },
    setPrimaryColor(color: string) {
      this.primaryColor = color
    },
    setThemeColor(color?: string, isDark?: boolean) {
      color = color ?? this.primaryColor
      isDark = isDark ?? this.isDark
      const colors = generate(color, {
        list: true,
        dark: isDark,
      })
      document.body.style.setProperty('--primary-color', getRgbStr(colors[5]))
      this.naiveThemeOverrides.common = {
        ...(this.naiveThemeOverrides.common || {}),
        primaryColor: colors[5],
        primaryColorHover: colors[4],
        primaryColorSuppl: colors[4],
        primaryColorPressed: colors[6],
      }
    },
  },
  persist: {
    pick: ['collapsed', 'primaryColor', 'naiveThemeOverrides'],
    storage: sessionStorage,
  },
})
