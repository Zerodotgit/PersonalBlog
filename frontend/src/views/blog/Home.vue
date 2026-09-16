<template>
  <div class="home">
    <div class="hero">
      <h1>记录与分享开发实践</h1>
      <p>这里是我写开发技术文章的地方，涉及后端、前端、数据库与工程化等话题。</p>
    </div>

    <el-skeleton v-if="loading" :rows="6" animated />

    <div v-else-if="articles.length === 0" class="empty">
      还没有文章，先去后台发布一篇吧。
    </div>

    <div v-else class="article-grid">
      <article
        v-for="a in articles"
        :key="a.id"
        class="card"
        @click="goDetail(a.id)"
      >
        <div
          class="cover"
          :style="coverStyle(a.cover_image)"
        ></div>
        <div class="body">
          <h2 class="title">{{ a.title }}</h2>
          <p class="summary">{{ a.summary || '暂无摘要' }}</p>
          <div class="meta">
            <span>{{ formatDate(a.created_at) }}</span>
            <span v-if="a.tags" class="tags">{{ a.tags }}</span>
            <span class="views">浏览 {{ a.view_count }}</span>
          </div>
        </div>
      </article>
    </div>

    <div v-if="total > size" class="pager">
      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="size"
        :current-page="page"
        @current-change="onPage"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchPublicArticles, type ArticleListItem } from '@/api/article'

const router = useRouter()
const loading = ref(true)
const articles = ref<ArticleListItem[]>([])
const total = ref(0)
const page = ref(1)
const size = 9

async function load() {
  loading.value = true
  try {
    const res = await fetchPublicArticles(page.value, size)
    articles.value = res.items || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function onPage(p: number) {
  page.value = p
  load()
}

function goDetail(id: number) {
  router.push(`/article/${id}`)
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('zh-CN')
}

function coverStyle(cover: string) {
  if (cover) {
    return { backgroundImage: `url(${cover})` }
  }
  return {}
}

onMounted(load)
</script>

<style scoped>
.hero {
  text-align: center;
  padding: 24px 0 40px;
}
.hero h1 {
  font-size: 30px;
  margin: 0 0 10px;
}
.hero p {
  color: var(--color-text-light);
  margin: 0;
}
.empty {
  text-align: center;
  color: var(--color-text-light);
  padding: 60px 0;
}
.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}
.card {
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
  display: flex;
  flex-direction: column;
}
.card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}
.cover {
  height: 150px;
  background-size: cover;
  background-position: center;
  background-image: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.body {
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}
.title {
  font-size: 18px;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.summary {
  color: var(--color-text-light);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}
.meta {
  display: flex;
  gap: 14px;
  font-size: 12px;
  color: var(--color-text-light);
  flex-wrap: wrap;
}
.tags {
  color: var(--color-primary);
}
.pager {
  margin-top: 36px;
  display: flex;
  justify-content: center;
}
</style>
