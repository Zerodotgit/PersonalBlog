<template>
  <div class="blog-app">
    <header class="nav">
      <div class="nav-inner">
        <router-link to="/" class="logo">个人博客</router-link>
        <nav class="links">
          <router-link to="/">首页</router-link>
          <router-link to="/about">关于</router-link>
          <router-link to="/admin/dashboard" class="admin-link">后台</router-link>
          <button class="theme-toggle" @click="toggleTheme" :aria-label="isDark ? '切换到浅色模式' : '切换到深色模式'">
            <svg v-if="isDark" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="4"/>
              <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"/>
            </svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
            </svg>
          </button>
        </nav>
      </div>
    </header>
    <main class="content">
      <router-view />
    </main>
    <footer class="footer">
      <span>© {{ year }} 个人博客 · Powered by Go + Vue</span>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'

const THEME_KEY = 'blog-theme'
const isDark = ref(false)
let userOverride = false
let mediaQuery: MediaQueryList | null = null

function applyTheme(dark: boolean) {
  isDark.value = dark
  if (dark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

function toggleTheme() {
  userOverride = true
  const next = !isDark.value
  applyTheme(next)
  localStorage.setItem(THEME_KEY, next ? 'dark' : 'light')
}

function initTheme() {
  const saved = localStorage.getItem(THEME_KEY)
  if (saved === 'dark') {
    userOverride = true
    applyTheme(true)
    return
  }
  if (saved === 'light') {
    userOverride = true
    applyTheme(false)
    return
  }
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  applyTheme(mediaQuery.matches)
  mediaQuery.addEventListener('change', (e) => {
    if (!userOverride) applyTheme(e.matches)
  })
}

onMounted(initTheme)

onBeforeUnmount(() => {
  if (mediaQuery) {
    mediaQuery.removeEventListener('change', () => {})
  }
})

const year = computed(() => new Date().getFullYear())
</script>

<style scoped>
.blog-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.nav {
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--shadow);
}
.nav-inner {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.logo {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text);
}
.links {
  display: flex;
  gap: 22px;
}
.links a {
  color: var(--color-text-light);
  font-weight: 500;
}
.links a.router-link-active,
.links a.router-link-exact-active,
.links a:hover {
  color: var(--color-primary);
}
.admin-link {
  color: var(--color-primary) !important;
}
.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  padding: 0;
  border-radius: 50%;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-light);
  cursor: pointer;
  transition: all 0.25s ease;
}
.theme-toggle:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
  transform: rotate(18deg);
}
.theme-toggle:active {
  transform: scale(0.92);
}
.content {
  flex: 1;
  width: 100%;
  max-width: 1080px;
  margin: 0 auto;
  padding: 32px 20px;
}
.footer {
  text-align: center;
  padding: 20px;
  color: var(--color-text-light);
  font-size: 13px;
  border-top: 1px solid var(--color-border);
}
</style>