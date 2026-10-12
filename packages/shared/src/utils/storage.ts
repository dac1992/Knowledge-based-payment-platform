import { STORAGE_PREFIX } from '../constants/brand'

interface UniStorage {
  getStorageSync: (key: string) => unknown
  setStorageSync: (key: string, value: unknown) => void
  removeStorageSync: (key: string) => void
  clearStorageSync: () => void
}

/** 存储抽象：Web 走 localStorage，uni-app 走 uni.*StorageSync。 */
export interface StorageLike {
  get<T>(key: string, fallback?: T): T | undefined
  set(key: string, value: unknown): void
  remove(key: string): void
  clear(): void
}

function getUniStorage(): UniStorage | null {
  const g = globalThis as { uni?: UniStorage }
  const uni = g.uni
  if (uni && typeof uni.getStorageSync === 'function' && typeof uni.setStorageSync === 'function') {
    return uni
  }
  return null
}

function getLocalStorage(): Storage | null {
  try {
    const g = globalThis as { localStorage?: Storage }
    const ls = g.localStorage
    if (!ls) return null
    // 隐私模式 / 禁用 Cookie 时 setItem 会抛错，提前探测避免运行期崩溃
    const probe = '__yanmai_probe__'
    ls.setItem(probe, '1')
    ls.removeItem(probe)
    return ls
  } catch {
    return null
  }
}

/**
 * 创建带命名空间的存储实例。
 *
 * @param prefix 键前缀，默认 `yanmai`
 */
export function createStorage(prefix: string = STORAGE_PREFIX): StorageLike {
  const uni = getUniStorage()
  const ls = uni ? null : getLocalStorage()
  const full = (key: string) => `${prefix}:${key}`

  return {
    get<T>(key: string, fallback?: T): T | undefined {
      try {
        if (uni) {
          const v = uni.getStorageSync(full(key))
          return v === '' || v === null || v === undefined ? fallback : (v as T)
        }
        if (!ls) return fallback
        const raw = ls.getItem(full(key))
        if (raw === null || raw === '') return fallback
        return JSON.parse(raw) as T
      } catch {
        return fallback
      }
    },

    set(key: string, value: unknown): void {
      try {
        if (uni) {
          uni.setStorageSync(full(key), value)
          return
        }
        ls?.setItem(full(key), JSON.stringify(value))
      } catch {
        // 存储失败不应中断业务流程
      }
    },

    remove(key: string): void {
      try {
        if (uni) {
          uni.removeStorageSync(full(key))
          return
        }
        ls?.removeItem(full(key))
      } catch {
        // 忽略
      }
    },

    clear(): void {
      try {
        if (uni) {
          uni.clearStorageSync()
          return
        }
        if (!ls) return
        const keys: string[] = []
        for (let i = 0; i < ls.length; i++) {
          const k = ls.key(i)
          if (k && k.startsWith(`${prefix}:`)) keys.push(k)
        }
        keys.forEach((k) => ls.removeItem(k))
      } catch {
        // 忽略
      }
    },
  }
}

/** 默认存储实例。 */
export const storage = createStorage()
