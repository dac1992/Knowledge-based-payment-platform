import * as ElementPlusIcons from '@element-plus/icons-vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { createPinia } from 'pinia'
import { createApp } from 'vue'

import '@yanmai/ui/styles/tokens.css'
import 'element-plus/dist/index.css'
import '@yanmai/ui/styles/element-theme.css'
import './styles/main.css'

import { setUnauthorizedHandler } from '@yanmai/api-client'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })

// 统一注册 Element Plus 图标，模板中可直接用 <el-icon><User /></el-icon>
for (const [name, component] of Object.entries(ElementPlusIcons)) {
  app.component(name, component)
}

// 令牌彻底失效时统一跳登录页（避免各处重复处理）
setUnauthorizedHandler(() => {
  const current = router.currentRoute.value
  if (current.name === 'login') return
  void router.replace({ name: 'login', query: { redirect: current.fullPath } })
})

app.mount('#app')
