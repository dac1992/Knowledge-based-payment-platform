import * as ElementPlusIcons from '@element-plus/icons-vue'
import { setUnauthorizedHandler } from '@yanmai/api-client'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { createPinia } from 'pinia'
import { createApp } from 'vue'

import '@yanmai/ui/styles/tokens.css'
import 'element-plus/dist/index.css'
import '@yanmai/ui/styles/element-theme.css'
import './styles/main.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { locale: zhCn })

for (const [name, component] of Object.entries(ElementPlusIcons)) {
  app.component(name, component)
}

setUnauthorizedHandler(() => {
  const current = router.currentRoute.value
  if (current.name === 'login') return
  void router.replace({ name: 'login', query: { redirect: current.fullPath } })
})

app.mount('#app')
