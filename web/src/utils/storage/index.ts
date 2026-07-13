import { createStorage } from './storage'

const prefixKey = 'vue-naive-admin_'

interface StorageFactoryOptions {
  prefixKey?: string
}

export function createLocalStorage(option: StorageFactoryOptions = {}) {
  return createStorage({
    prefixKey: option.prefixKey || '',
    storage: localStorage,
  })
}

export function createSessionStorage(option: StorageFactoryOptions = {}) {
  return createStorage({
    prefixKey: option.prefixKey || '',
    storage: sessionStorage,
  })
}

export const lStorage = createLocalStorage({ prefixKey })

export const sStorage = createSessionStorage({ prefixKey })
