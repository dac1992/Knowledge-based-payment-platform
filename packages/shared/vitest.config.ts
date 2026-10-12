import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    // 串行执行测试文件：
    //   1) 当前用例规模很小，串行的耗时增加可忽略，但失败定位更清晰、结果更可复现；
    //   2) 并发 module runner 会触发 vitest「把已转换模块落盘为临时文件再 import」的
    //      回退路径，在只读/受限的系统临时目录下会直接失败。串行执行可规避该回退。
    fileParallelism: false,
  },
})
