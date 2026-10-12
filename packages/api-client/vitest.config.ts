import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    // 串行执行测试文件（原因见 packages/shared/vitest.config.ts 注释）：
    // 小规模用例下串行更可复现，同时规避受限临时目录下的模块落盘回退。
    fileParallelism: false,
  },
})
