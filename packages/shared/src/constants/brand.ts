import type { UserRole } from '../types/api'

/**
 * 品牌色板 —— 与《UI 设计规范》保持一致。
 * 语义色遵守国内股市惯例：涨/正向为红，跌/负向为绿。
 */
export const BRAND = {
  /** 主色：研脉蓝 */
  primary: '#1A56DB',
  primaryDark: '#0E3A99',
  primaryLight: '#EBF1FE',
  /** 辅色：学术金 */
  accent: '#C9A227',
  accentLight: '#FBF4DF',

  ink: '#1A1F2B',
  ink2: '#3D4759',
  ink3: '#6B7688',
  ink4: '#98A2B3',

  line: '#E4E8EF',
  lineLight: '#EEF1F6',
  bg: '#F4F6FA',
  card: '#FFFFFF',

  /** 正向（国内惯例：红涨） */
  up: '#DC2626',
  upLight: '#FDECEC',
  /** 负向（国内惯例：绿跌） */
  down: '#16A34A',
  downLight: '#E9F8EF',

  warn: '#D97706',
  warnLight: '#FEF5E7',
  info: '#2563EB',
  infoLight: '#EBF1FE',
} as const

/** 角色中文名。 */
export const ROLE_LABELS: Record<UserRole, string> = {
  admin: '平台管理员',
  lecturer: '讲师',
  student: '学员',
}

/** 角色对应的标签配色（Element Plus tag type）。 */
export const ROLE_TAG_TYPES: Record<UserRole, 'danger' | 'warning' | 'primary'> = {
  admin: 'danger',
  lecturer: 'warning',
  student: 'primary',
}

/** 账号状态中文名。 */
export const STATUS_LABELS: Record<number, string> = {
  0: '已禁用',
  1: '正常',
}

/** 本地存储键前缀。 */
export const STORAGE_PREFIX = 'yanmai'

/** 存储键名。 */
export const STORAGE_KEYS = {
  accessToken: 'accessToken',
  refreshToken: 'refreshToken',
  profile: 'profile',
} as const
