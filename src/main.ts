import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

async function bootstrap() {
  try {
    const script = document.createElement('script')
    script.src = '/auth/config.js'
    await new Promise((resolve, reject) => {
      script.onload = resolve
      script.onerror = reject
      document.head.appendChild(script)
    })
  } catch {
    // no config served (tests / local dev without server)
  }
  createApp(App).use(router).mount('#app')
}

bootstrap()
