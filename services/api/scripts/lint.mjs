#!/usr/bin/env node
/**
 * 跨平台的 golangci-lint 启动器。
 *
 * 为什么需要它：golangci-lint 常被 `go install` 到 $(go env GOPATH)/bin 而
 * 未加入 PATH；Windows 上还多一个 .exe 后缀，且本机没有 make，无法借用
 * Makefile 的探测逻辑。这里按 PATH → GOPATH/bin 的顺序探测，找到即执行。
 *
 * 显式指定：设置环境变量 GOLANGCI_LINT=/path/to/golangci-lint
 */
import { spawnSync } from 'node:child_process'
import { existsSync } from 'node:fs'
import { dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const apiDir = resolve(dirname(fileURLToPath(import.meta.url)), '..')

/** PATH 中是否有可用的 golangci-lint */
function fromPath() {
  const probe = spawnSync('golangci-lint', ['--version'], {
    shell: true,
    stdio: 'ignore',
  })
  return probe.status === 0
}

/** $(go env GOPATH)/bin 下是否有 golangci-lint（兼容 Windows 的 .exe） */
function fromGopath() {
  const probe = spawnSync('go', ['env', 'GOPATH'], { encoding: 'utf8', shell: true })
  if (probe.status !== 0) return null

  const gopath = probe.stdout.trim()
  if (!gopath) return null

  for (const name of ['golangci-lint', 'golangci-lint.exe']) {
    const candidate = join(gopath, 'bin', name)
    if (existsSync(candidate)) return candidate
  }
  return null
}

function resolveBinary() {
  if (process.env.GOLANGCI_LINT) return process.env.GOLANGCI_LINT
  if (fromPath()) return 'golangci-lint'

  const gopathBinary = fromGopath()
  if (gopathBinary) return gopathBinary

  return null
}

const binary = resolveBinary()
if (!binary) {
  console.error(
    '未找到 golangci-lint。请先安装：\n' +
      '  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0\n' +
      '或设置环境变量 GOLANGCI_LINT 指向可执行文件。',
  )
  process.exit(1)
}

console.log(`golangci-lint: ${binary}`)

const result = spawnSync(binary, ['run', './...'], {
  cwd: apiDir,
  stdio: 'inherit',
  shell: process.platform === 'win32',
})

process.exit(result.status ?? 1)
