import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'
import { initAnalytics } from './firebase'

async function bootstrap() {
  try {
    const script = document.createElement('script')
    script.src = '/auth/config.js'
    await Promise.race([
      new Promise((resolve, reject) => {
        script.onload = resolve
        script.onerror = reject
        document.head.appendChild(script)
      }),
      new Promise(resolve => setTimeout(resolve, 5000)),
    ])
  } catch {
    // no config served (tests / local dev without server)
  }
  try { initAnalytics() } catch { /* analytics must never block app mount */ }
  createApp(App).use(router).mount('#app')
}

bootstrap()
