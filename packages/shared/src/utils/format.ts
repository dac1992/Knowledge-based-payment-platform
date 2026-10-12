/**
 * 通用格式化工具。所有函数对异常输入都返回安全结果，不抛错。
 */

function toDate(input: string | number | Date | null | undefined): Date | null {
  if (input === null || input === undefined || input === '') return null
  const d = input instanceof Date ? input : new Date(input)
  return Number.isNaN(d.getTime()) ? null : d
}

function pad(n: number): string {
  return n < 10 ? `0${n}` : String(n)
}

/** 格式化为 YYYY-MM-DD。 */
export function formatDate(input: string | number | Date | null | undefined): string {
  const d = toDate(input)
  if (!d) return '-'
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

/** 格式化为 YYYY-MM-DD HH:mm。 */
export function formatDateTime(input: string | number | Date | null | undefined): string {
  const d = toDate(input)
  if (!d) return '-'
  return `${formatDate(d)} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/** 格式化为 YYYY-MM-DD HH:mm:ss。 */
export function formatDateTimeFull(input: string | number | Date | null | undefined): string {
  const d = toDate(input)
  if (!d) return '-'
  return `${formatDateTime(d)}:${pad(d.getSeconds())}`
}

/** 手机号脱敏：13800000001 → 138****0001。 */
export function maskPhone(phone: string | null | undefined): string {
  if (!phone) return '-'
  if (phone.length < 7) return phone
  return `${phone.slice(0, 3)}****${phone.slice(-4)}`
}

/** 相对时间：刚刚 / N 分钟前 / N 小时前 / N 天前 / 具体日期。 */
export function relativeTime(
  input: string | number | Date | null | undefined,
  now = Date.now(),
): string {
  const d = toDate(input)
  if (!d) return '-'

  const diff = now - d.getTime()
  if (diff < 0) return formatDateTime(d)
  if (diff < 60_000) return '刚刚'

  const minutes = Math.floor(diff / 60_000)
  if (minutes < 60) return `${minutes} 分钟前`

  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours} 小时前`

  const days = Math.floor(hours / 24)
  if (days < 30) return `${days} 天前`

  return formatDate(d)
}

/** 字节数转可读体积。 */
export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  const value = bytes / Math.pow(1024, i)
  return `${value >= 100 || i === 0 ? Math.round(value) : value.toFixed(1)} ${units[i]}`
}
