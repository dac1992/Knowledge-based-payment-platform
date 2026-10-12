import { describe, expect, it } from 'vitest'

import { formatBytes, formatDate, formatDateTime, maskPhone, relativeTime } from './format'

describe('formatDate / formatDateTime', () => {
  it('格式化合法日期', () => {
    const d = new Date(2026, 9, 12, 15, 5, 8)
    expect(formatDate(d)).toBe('2026-10-12')
    expect(formatDateTime(d)).toBe('2026-10-12 15:05')
  })

  it('对非法输入返回占位符而不是抛错', () => {
    expect(formatDate('')).toBe('-')
    expect(formatDate(null)).toBe('-')
    expect(formatDate(undefined)).toBe('-')
    expect(formatDateTime('not-a-date')).toBe('-')
  })
})

describe('maskPhone', () => {
  it('隐藏中间四位', () => {
    expect(maskPhone('13800000001')).toBe('138****0001')
  })

  it('过短或空值原样返回', () => {
    expect(maskPhone('123')).toBe('123')
    expect(maskPhone('')).toBe('-')
    expect(maskPhone(undefined)).toBe('-')
  })
})

describe('relativeTime', () => {
  const now = new Date(2026, 9, 12, 12, 0, 0).getTime()

  it('按区间降级展示', () => {
    expect(relativeTime(new Date(now - 30_000), now)).toBe('刚刚')
    expect(relativeTime(new Date(now - 5 * 60_000), now)).toBe('5 分钟前')
    expect(relativeTime(new Date(now - 3 * 3_600_000), now)).toBe('3 小时前')
    expect(relativeTime(new Date(now - 2 * 86_400_000), now)).toBe('2 天前')
  })

  it('超过 30 天回退为具体日期', () => {
    expect(relativeTime(new Date(2026, 7, 1), now)).toBe('2026-08-01')
  })

  it('未来时间直接展示日期', () => {
    expect(relativeTime(new Date(now + 86_400_000), now)).toBe('2026-10-13 12:00')
  })
})

describe('formatBytes', () => {
  it('换算单位', () => {
    expect(formatBytes(0)).toBe('0 B')
    expect(formatBytes(512)).toBe('512 B')
    expect(formatBytes(1024)).toBe('1.0 KB')
    expect(formatBytes(1024 * 1024 * 3)).toBe('3.0 MB')
  })
})
