import type { Ref } from 'vue'
import type { RouteRecordNameGeneric } from 'vue-router'

const lastDataMap = new Map<string | symbol, unknown>()
function resolveAliveKey(key?: RouteRecordNameGeneric): string | symbol {
  return key ?? 'default'
}

export function useAliveData<T extends Record<string, unknown> = Record<string, unknown>>(
  initData = {} as T,
  key?: RouteRecordNameGeneric,
): { aliveData: Ref<T>, reset: () => void } {
  const aliveKey = resolveAliveKey(key ?? useRoute().name)
  const lastData = lastDataMap.get(aliveKey)
  const aliveData = ref((lastData || { ...initData }) as T) as Ref<T>

  watch(
    aliveData,
    (v) => {
      lastDataMap.set(aliveKey, v)
    },
    { deep: true },
  )

  return {
    aliveData,
    reset() {
      aliveData.value = { ...initData }
      lastDataMap.delete(aliveKey)
    },
  }
}
