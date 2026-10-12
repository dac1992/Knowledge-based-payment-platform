#!/usr/bin/env node
/**
 * 由后端 OpenAPI 文档生成类型定义。
 *
 * 使用前需先生成 swagger.json：
 *   cd services/api && make swagger     # 或： go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
 *
 * 然后：
 *   pnpm --filter @yanmai/api-client gen:openapi
 */
import { existsSync, mkdirSync, writeFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const swaggerPath = resolve(here, '../../../services/api/docs/swagger.json')
const outPath = resolve(here, '../src/generated/schema.d.ts')

if (!existsSync(swaggerPath)) {
  console.error(`[gen:openapi] 未找到 OpenAPI 文档：${swaggerPath}`)
  console.error('[gen:openapi] 请先执行：cd services/api && make swagger')
  process.exit(1)
}

let openapi
try {
  ;({ default: openapi } = await import('openapi-typescript'))
} catch {
  console.error('[gen:openapi] 缺少依赖 openapi-typescript，请先安装：')
  console.error('  pnpm --filter @yanmai/api-client add -D openapi-typescript')
  process.exit(1)
}

const ast = await openapi(swaggerPath)
mkdirSync(dirname(outPath), { recursive: true })
writeFileSync(outPath, `${String(ast)}\n`, 'utf8')
console.log(`[gen:openapi] 已生成 ${outPath}`)
