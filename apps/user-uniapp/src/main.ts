import { createPinia } from 'pinia'
import { createSSRApp } from 'vue'

import App from './App.vue'

// uni-app 约定：运行时要求导出 createApp 并返回 { app }，此处类型由 uni 运行时决定。
export function createApp(): { app: unknown } {
  const app = createSSRApp(App)
  app.use(createPinia())
  return { app }
}
