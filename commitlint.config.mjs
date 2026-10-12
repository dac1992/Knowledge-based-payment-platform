/**
 * 提交信息规范（Conventional Commits）。
 *
 * 面向中文团队做了取舍：类型前缀仍走国际约定（便于生成 CHANGELOG 与语义化版本），
 * 但标题正文用中文书写，因此不做英文向的大小写、句号等强制。
 */
export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat', // 新功能
        'fix', // 缺陷修复
        'docs', // 文档
        'style', // 格式调整（不影响行为）
        'refactor', // 重构
        'perf', // 性能优化
        'test', // 测试
        'build', // 构建系统或依赖
        'ci', // CI 配置
        'chore', // 杂项
        'revert', // 回滚
      ],
    ],
    // 标题为中文，关闭大小写与结尾句号限制
    'subject-case': [0],
    'subject-full-stop': [0],
    // 标题限长，正文与脚注可自由换行写中文说明
    'header-max-length': [2, 'always', 100],
    'body-max-line-length': [0],
    'footer-max-line-length': [0],
  },
}
