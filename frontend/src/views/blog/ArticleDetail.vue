<template>
  <div class="detail">
    <el-skeleton v-if="loading" :rows="10" animated />

    <template v-else-if="article">
      <div class="header">
        <h1 class="title">{{ article.title }}</h1>
        <div class="meta">
          <span>{{ formatDate(article.created_at) }}</span>
          <span class="views">浏览 {{ article.view_count }}</span>
          <span v-if="article.tags" class="tags">{{ article.tags }}</span>
        </div>
      </div>

      <div class="markdown-body" v-html="rendered"></div>

      <div class="back">
        <router-link to="/">← 返回首页</router-link>
      </div>
    </template>

    <el-empty v-else description="文章不存在或尚未发布" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import { fetchArticle, type ArticleDetail } from '@/api/article'

const route = useRoute()
const loading = ref(true)
const article = ref<ArticleDetail | null>(null)
const rendered = ref('')

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('zh-CN')
}

function renderMarkdown(content: string) {
  const raw = marked.parse(content, { breaks: true }) as string
  return DOMPurify.sanitize(raw)
}

async function load() {
  const id = Number(route.params.id)
  if (!id) return
  loading.value = true
  try {
    article.value = await fetchArticle(id)
    rendered.value = renderMarkdown(article.value.content)
    await nextTick()
    highlightAll()
  } finally {
    loading.value = false
  }
}

function highlightAll() {
  const blocks = document.querySelectorAll('.markdown-body pre code')
  blocks.forEach((el) => {
    if (!el.getAttribute('data-highlighted')) {
      hljs.highlightElement(el as HTMLElement)
      el.setAttribute('data-highlighted', 'yes')
    }
  })
}

watch(() => route.params.id, load)

onMounted(load)
</script>

<style scoped>
.detail {
  max-width: var(--content-max);
  margin: 0 auto;
}
.header {
  margin-bottom: 28px;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 16px;
}
.title {
  font-size: 30px;
  margin: 0 0 12px;
  line-height: 1.3;
}
.meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: var(--color-text-light);
  flex-wrap: wrap;
}
.tags {
  color: var(--color-primary);
}
.back {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid var(--color-border);
}
</style>
