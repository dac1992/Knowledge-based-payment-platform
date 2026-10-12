#!/usr/bin/env node
/**
 * 由后端 OpenAPI 文档生成前端类型定义。
 *
 * 链路：
 *   swag（Swagger 2.0）──swagger2openapi──▶ OpenAPI 3.0 ──openapi-typescript──▶ schema.d.ts
 *
 * 为什么需要中间转换：swag 目前只能输出 Swagger 2.0（无 3.x 选项），
 * 而 openapi-typescript v7 只接受 openapi: 3.x，因此必须做一次规范转换。
 *
 * 使用：
 *   cd services/api && make sdk                      # 推荐：一步完成文档生成 + 类型同步
 *   # 或分步：
 *   cd services/api && make swagger
 *   pnpm --filter @yanmai/api-client gen:openapi
 */
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs'
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

/** 动态加载依赖，缺失时给出可执行的安装指引。 */
async function load(name) {
  try {
    return await import(name)
  } catch {
    console.error(`[gen:openapi] 缺少依赖 ${name}，请先安装：`)
    console.error(`  pnpm --filter @yanmai/api-client add -D ${name}`)
    process.exit(1)
  }
}

const { default: Swagger2OpenAPI } = await load('swagger2openapi')
const { default: openapiTS, astToString } = await load('openapi-typescript')

const swagger2 = JSON.parse(readFileSync(swaggerPath, 'utf8'))
if (swagger2.swagger !== '2.0') {
  console.error(
    `[gen:openapi] 预期 Swagger 2.0，实际为 ${swagger2.swagger || '(未声明)'}；` +
      '若 swag 未来直接输出 OpenAPI 3.x，请调整本脚本的前置转换逻辑。',
  )
  process.exit(1)
}

// patch：对不规范的文档做补全；warnOnly：转换告警不中断，避免小瑕疵阻断构建
const { openapi: openapi3 } = await Swagger2OpenAPI.convertObj(swagger2, {
  patch: true,
  warnOnly: true,
})

// alphabetize：字段名排序，保证同一份契约每次生成结果完全一致（CI 契约校验依赖此稳定性）
const ast = await openapiTS(openapi3, { alphabetize: true })

const banner = `/**
 * 本文件由脚本自动生成，请勿手动修改。
 *
 * 来源：services/api/docs/swagger.json
 * 生成：cd services/api && make sdk
 */

`

mkdirSync(dirname(outPath), { recursive: true })
writeFileSync(outPath, `${banner}${astToString(ast)}\n`, 'utf8')

console.log(`[gen:openapi] 已生成 ${outPath}`)
console.log(
  `[gen:openapi] 路径 ${Object.keys(openapi3.paths ?? {}).length} 个，` +
    `模型 ${Object.keys(openapi3.components?.schemas ?? {}).length} 个`,
)
