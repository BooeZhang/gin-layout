import { isNullOrUndef } from '@/utils'

interface StorageOptions {
  prefixKey?: string
  storage?: globalThis.Storage
}

interface StorageItem<T = unknown> {
  value: T
  time: number
  expire?: number | null
}

class StorageWrapper {
  private storage: globalThis.Storage
  private prefixKey: string

  constructor(option: Required<StorageOptions>) {
    this.storage = option.storage
    this.prefixKey = option.prefixKey
  }

  getKey(key: string): string {
    return `${this.prefixKey}${key}`.toLowerCase()
  }

  set<T>(key: string, value: T, expire?: number): void {
    const stringData = JSON.stringify({
      value,
      time: Date.now(),
      expire: !isNullOrUndef(expire) ? new Date().getTime() + expire * 1000 : null,
    })
    this.storage.setItem(this.getKey(key), stringData)
  }

  get<T = unknown>(key: string): T | undefined {
    return this.getItem<T>(key)?.value
  }

  getItem<T = unknown>(key: string, def: StorageItem<T> | null = null): StorageItem<T> | null {
    const val = this.storage.getItem(this.getKey(key))
    if (!val)
      return def
    try {
      const data = JSON.parse(val) as StorageItem<T>
      const { value, time, expire } = data
      if (isNullOrUndef(expire) || expire > new Date().getTime()) {
        return { value, time }
      }
      this.remove(key)
      return def
    }
    catch (error) {
      console.error(error)
      this.remove(key)
      return def
    }
  }

  remove(key: string): void {
    this.storage.removeItem(this.getKey(key))
  }

  clear(): void {
    this.storage.clear()
  }
}

export function createStorage({ prefixKey = '', storage = sessionStorage }: StorageOptions = {}) {
  return new StorageWrapper({ prefixKey, storage })
}
