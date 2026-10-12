/**
 * 业务模块清单 —— 与《开发计划与技术选型》的 14 个阶段、以及后端
 * internal/modules 下实际挂载的模块保持一致。
 *
 * 管理后台工作台直接渲染这张表，用于直观呈现各业务域交付进度。
 */
export interface ModuleDescriptor {
  /** 模块 key，与后端路由前缀一致 */
  key: string
  /** 中文名 */
  name: string
  /** 计划交付阶段 */
  phase: number
  /** 是否已交付可用 */
  delivered: boolean
  /** 能力摘要 */
  summary: string
}

export const MODULES: ModuleDescriptor[] = [
  { key: 'auth', name: '认证', phase: 1, delivered: true, summary: '注册 / 登录 / 双 Token / 鉴权守卫' },
  { key: 'user', name: '用户', phase: 1, delivered: true, summary: '个人资料 / 用户列表 / 角色权限' },
  { key: 'course', name: '课程内容中心', phase: 3, delivered: false, summary: '分类体系 / 课程 SPU·SKU / 章节课时 / 审核上下架' },
  { key: 'order', name: '订单', phase: 5, delivered: false, summary: '购物车 / 下单 / 支付回调 / 权益开通 / 退款' },
  { key: 'payment', name: '支付渠道', phase: 5, delivered: false, summary: '微信支付 / 支付宝 / 验签 / 对账' },
  { key: 'learning', name: '学习中心', phase: 6, delivered: false, summary: '云点播 VOD / 播放器 / 进度 / 笔记 / 打卡' },
  { key: 'marketing', name: '营销中心', phase: 7, delivered: false, summary: '优惠券 / 拼团 / 秒杀 / 分销推广' },
  { key: 'member', name: '会员与积分', phase: 8, delivered: false, summary: '会员卡 / 权益 / 积分 / 成长等级' },
  { key: 'live', name: '直播', phase: 9, delivered: false, summary: '直播间 / 弹幕打赏问答 / 回放 / 直播带货' },
  { key: 'circle', name: '圈子社群', phase: 11, delivered: false, summary: '圈子 / 话题 / 问答 / 打卡 / 资料库 / 私信' },
  { key: 'finance', name: '财务结算', phase: 12, delivered: false, summary: '讲师分成 / 佣金结算 / 提现审批' },
  { key: 'stats', name: '数据看板', phase: 12, delivered: false, summary: '经营看板 / 转化漏斗 / 留存分析' },
]

/** 已交付模块数量。 */
export const DELIVERED_MODULE_COUNT = MODULES.filter((m) => m.delivered).length

/** 六端矩阵（用于工作台呈现整体蓝图）。 */
export const CLIENTS = [
  { key: 'admin-web', name: '管理后台 Web', stack: 'Vue3 + Vite + Element Plus', batch: 1, ready: true },
  { key: 'lecturer-web', name: '讲师工作台', stack: 'Vue3 + Vite + Element Plus', batch: 1, ready: true },
  { key: 'user-uniapp', name: '小程序 / H5 / APP', stack: 'uni-app (Vue3)', batch: 1, ready: true },
  { key: 'user-web', name: '用户 PC 端', stack: 'Nuxt 3 (SSR)', batch: 2, ready: false },
  { key: 'desktop', name: '桌面客户端', stack: 'Tauri 2 + Vue3', batch: 3, ready: false },
] as const
