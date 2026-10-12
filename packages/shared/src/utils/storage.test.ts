import { afterEach, describe, expect, it, vi } from 'vitest'

import { createStorage } from './storage'

/** 最小 localStorage 实现，供 Node 环境下的单测使用。 */
function installFakeLocalStorage(): void {
  const map = new Map<string, string>()
  const fake: Storage = {
    get length() {
      return map.size
    },
    key: (i: number) => Array.from(map.keys())[i] ?? null,
    getItem: (k: string) => (map.has(k) ? (map.get(k) as string) : null),
    setItem: (k: string, v: string) => void map.set(k, v),
    removeItem: (k: string) => void map.delete(k),
    clear: () => map.clear(),
  }
  vi.stubGlobal('localStorage', fake)
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('createStorage', () => {
  it('按前缀读写并支持 JSON 结构', () => {
    installFakeLocalStorage()
    const s = createStorage('test')

    s.set('profile', { id: 1, nickname: '管理员' })
    expect(s.get('profile')).toEqual({ id: 1, nickname: '管理员' })
    expect(globalThis.localStorage.getItem('test:profile')).toBe('{"id":1,"nickname":"管理员"}')
  })

  it('缺失键返回兜底值', () => {
    installFakeLocalStorage()
    const s = createStorage('test')

    expect(s.get('missing')).toBeUndefined()
    expect(s.get('missing', 'fallback')).toBe('fallback')
  })

  it('remove 只删除目标键', () => {
    installFakeLocalStorage()
    const s = createStorage('test')

    s.set('a', 1)
    s.set('b', 2)
    s.remove('a')

    expect(s.get('a')).toBeUndefined()
    expect(s.get('b')).toBe(2)
  })

  it('clear 只清理自身命名空间', () => {
    installFakeLocalStorage()
    const a = createStorage('app-a')
    const b = createStorage('app-b')

    a.set('token', 'x')
    b.set('token', 'y')
    a.clear()

    expect(a.get('token')).toBeUndefined()
    expect(b.get('token')).toBe('y')
  })

  it('localStorage 不可用时静默降级，不抛错', () => {
    vi.stubGlobal('localStorage', undefined)
    const s = createStorage('test')

    expect(() => s.set('k', 'v')).not.toThrow()
    expect(s.get('k', 'safe')).toBe('safe')
    expect(() => s.remove('k')).not.toThrow()
    expect(() => s.clear()).not.toThrow()
  })

  it('存储内容被污染时不抛错', () => {
    installFakeLocalStorage()
    const s = createStorage('test')

    globalThis.localStorage.setItem('test:broken', '{not-json')
    expect(s.get('broken', 'fallback')).toBe('fallback')
  })
})
