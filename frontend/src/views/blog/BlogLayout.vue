<template>
  <div class="blog-app">
    <header class="nav">
      <div class="nav-inner">
        <router-link to="/" class="logo">个人博客</router-link>
        <nav class="links">
          <router-link to="/">首页</router-link>
          <router-link to="/about">关于</router-link>
          <router-link to="/admin/dashboard" class="admin-link">后台</router-link>
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
import { computed } from 'vue'
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
.links a.router-link-exact-active,
.links a:hover {
  color: var(--color-primary);
}
.admin-link {
  color: var(--color-primary) !important;
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
