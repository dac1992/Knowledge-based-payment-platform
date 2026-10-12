/**
 * lint-staged：只处理本次暂存的文件，保证 pre-commit 快速且不误伤。
 *
 * ESLint 9 的 flat config 不做目录级联，只认「当前工作目录」下的 eslint.config.js。
 * 在仓库根直接跑 eslint 会找不到配置，因此按 workspace 分组，
 * 用 `pnpm -C <包目录> exec eslint` 切进各包执行，让各自的配置与依赖生效。
 *
 * Prettier 统一使用根目录 .prettierrc.json（引用 @yanmai/config/prettier）。
 * Go 交给 gofmt；golangci-lint 较慢，留给 CI 与 `make lint`。
 */
import { relative, resolve } from 'node:path'

const repoRoot = process.cwd()

/** 统一转成「相对仓库根 + 正斜杠」的路径；lint-staged 传入的可能是绝对路径。 */
function toRepoRelative(file) {
  return relative(repoRoot, resolve(repoRoot, file)).replace(/\\/g, '/')
}

/** 提取所属 workspace（如 packages/shared）；根级文件返回 null。 */
function workspaceOf(repoRelativePath) {
  const matched = repoRelativePath.match(/^((?:apps|packages)\/[^/]+)\//)
  return matched ? matched[1] : null
}

const quote = (names) => names.map((n) => `"${n}"`).join(' ')

/**
 * 不参与自动格式化的路径。
 *
 * 两类必须排除：
 * 1. 生成物（OpenAPI 文档、契约类型）——被格式化会与 CI 重新生成的版本产生差异，
 *    直接导致契约校验失败。
 * 2. 包管理器锁文件 —— pnpm 有自己的写法，被 Prettier 改写后每次 install 都会
 *    产生无意义的 diff。
 */
const EXCLUDED_PATTERNS = [
  /^packages\/api-client\/src\/generated\//,
  /^services\/api\/docs\//,
  /^pnpm-lock\.yaml$/,
  /^docs\//,
]

const isExcluded = (repoRelativePath) =>
  EXCLUDED_PATTERNS.some((pattern) => pattern.test(repoRelativePath))

/** 按 workspace 分组，分别在其目录内执行 eslint --fix。 */
function eslintCommands(repoRelativeFiles) {
  const groups = new Map()
  for (const file of repoRelativeFiles) {
    const ws = workspaceOf(file)
    if (!ws) continue
    if (!groups.has(ws)) groups.set(ws, [])
    groups.get(ws).push(file)
  }

  const commands = []
  for (const [ws, files] of groups) {
    const insideWorkspace = files.map((f) => f.slice(ws.length + 1))
    commands.push(`pnpm -C ${ws} exec eslint --fix ${quote(insideWorkspace)}`)
  }
  return commands
}

export default {
  '*.{js,mjs,cjs,ts,vue}': (files) => {
    const relPaths = files.map(toRepoRelative).filter((f) => !isExcluded(f))
    if (relPaths.length === 0) return []
    // 根级配置文件（如 lint-staged.config.mjs）没有对应的 eslint 配置，
    // 只做格式化，不强行 lint，避免误报阻断提交。
    return [...eslintCommands(relPaths), `pnpm exec prettier --write ${quote(relPaths)}`]
  },
  '*.{json,md,yml,yaml,css,scss,html}': (files) => {
    const relPaths = files.map(toRepoRelative).filter((f) => !isExcluded(f))
    return relPaths.length > 0 ? [`pnpm exec prettier --write ${quote(relPaths)}`] : []
  },
  '*.go': ['gofmt -w'],
}
