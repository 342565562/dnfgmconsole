import { createApp } from 'vue'
import App from './App.vue'

import router from './router'
import store from './store'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import locale from 'element-plus/lib/locale/lang/zh-cn'

import * as ElIcon from '@element-plus/icons-vue'

import SvgIcon from './components/SvgIcon/index.vue'
import './icons'
import { setupPermissionRouter } from './permission'

// 导入公共样式
import 'normalize.css/normalize.css'
import './styles/index.scss'

// init
;(async () => {
  const app = createApp(App)

  // load vue router
  app.use(router)

  //  load pinia store
  app.use(store)

  // load element plus ui
  app.use(ElementPlus, {
    size: 'default',
    locale
  })
  for (const iconName in ElIcon) {
    app.component(iconName, ElIcon[iconName])
  }

  // load 路由守卫
  setupPermissionRouter(router)

  // load svg组件
  app.component('svg-icon', SvgIcon)

  // 读取站点标题(失败不阻塞启动)
  try {
    const { getSiteConfig } = await import('./api/site')
    const site = await getSiteConfig()
    if (site?.title) document.title = site.title
  } catch (e) {
    // 忽略，使用默认标题
  }

  app.mount('#app', true)
})()
