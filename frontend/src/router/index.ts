import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/views/blog/BlogLayout.vue'),
      children: [
        { path: '', component: () => import('@/views/blog/Home.vue') },
        { path: 'article/:id', component: () => import('@/views/blog/ArticleDetail.vue') },
        { path: 'about', component: () => import('@/views/blog/About.vue') },
      ],
    },
    {
      path: '/admin/login',
      component: () => import('@/views/admin/Login.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/admin/dashboard' },
        { path: 'dashboard', component: () => import('@/views/admin/Dashboard.vue') },
        { path: 'articles', component: () => import('@/views/admin/ArticleList.vue') },
        { path: 'articles/new', component: () => import('@/views/admin/ArticleEdit.vue') },
        { path: 'articles/:id/edit', component: () => import('@/views/admin/ArticleEdit.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.isAuthed()) {
      return { path: '/admin/login', query: { redirect: to.fullPath } }
    }
  }
  return true
})

export default router
