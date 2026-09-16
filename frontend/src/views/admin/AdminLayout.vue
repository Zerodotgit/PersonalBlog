<template>
  <el-container class="admin">
    <el-aside width="220px" class="aside">
      <div class="brand">博客后台</div>
      <el-menu
        :default-active="activeMenu"
        router
        class="menu"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/admin/articles">
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/articles/new">
          <el-icon><EditPen /></el-icon>
          <span>写文章</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="spacer"></div>
        <el-dropdown @command="onCommand">
          <span class="user">
            {{ auth.username || 'admin' }}
            <el-icon><CaretBottom /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="home">查看博客</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { DataLine, Document, EditPen, CaretBottom } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const activeMenu = computed(() => route.path)

function onCommand(cmd: string) {
  if (cmd === 'home') {
    router.push('/')
  } else if (cmd === 'logout') {
    auth.logout()
    router.replace('/admin/login')
  }
}
</script>

<style scoped>
.admin {
  height: 100vh;
}
.aside {
  background: #1f2937;
  color: #fff;
}
.brand {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.menu {
  border-right: none;
  background: transparent;
}
:deep(.el-menu) {
  background: transparent;
}
:deep(.el-menu-item) {
  color: #cbd5e1;
}
:deep(.el-menu-item.is-active),
:deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}
.header {
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 24px;
}
.spacer {
  flex: 1;
}
.user {
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--color-text);
}
.main {
  background: var(--color-bg);
  padding: 24px;
}
</style>
