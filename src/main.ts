import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

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
  console.log('[firebase-config]', JSON.stringify(window.__firebaseConfig))
  createApp(App).use(router).mount('#app')
}

bootstrap()
