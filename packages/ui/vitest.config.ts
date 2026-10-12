import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    include: ['src/**/*.spec.ts'],
    // 串行执行测试文件（原因见 packages/shared/vitest.config.ts 注释）。
    fileParallelism: false,
  },
})
