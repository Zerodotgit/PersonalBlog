<template>
  <div class="detail">
    <el-skeleton v-if="loading" :rows="10" animated />

    <template v-else-if="article">
      <div class="header">
        <h1 class="title">{{ article.title }}</h1>
        <div class="meta">
          <span class="meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
            {{ formatDate(article.created_at) }}
          </span>
          <span class="meta-item views">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
            浏览 {{ article.view_count }}
          </span>
          <span v-if="article.tags" class="meta-item tags">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 01-2.83 0L2 12V2h10l8.59 8.59a2 2 0 010 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
            {{ article.tags }}
          </span>
        </div>
      </div>

      <div class="markdown-body" v-html="rendered" ref="contentRef"></div>

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
const contentRef = ref<HTMLElement | null>(null)

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('zh-CN')
}

const slugify = (text: string): string => {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s\u4e00-\u9fa5-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

function setupMarked() {
  const renderer = new marked.Renderer()
  const headingIdMap = new Map<string, number>()

  renderer.heading = function ({ tokens, depth }) {
    const text = this.parser.parseInline(tokens)
    let id = slugify(text)
    const count = headingIdMap.get(id) || 0
    if (count > 0) id = `${id}-${count}`
    headingIdMap.set(id, count + 1)

    return `<h${depth} id="${id}"><a class="header-anchor" href="#${id}" aria-label="锚点">#</a>${text}</h${depth}>`
  }

  renderer.code = function ({ text, lang }) {
    const language = lang && hljs.getLanguage(lang) ? lang : ''
    const highlighted = language
      ? hljs.highlight(text, { language }).value
      : hljs.highlightAuto(text).value
    const displayLang = language || (hljs.highlightAuto(text).language || 'plaintext')

    return `
      <pre><div class="code-header"><span class="lang-tag">${displayLang}</span><button class="copy-btn" type="button" data-code="${encodeURIComponent(text)}">复制</button></div><code class="hljs">${highlighted}</code></pre>
    `
  }

  marked.setOptions({ breaks: true, gfm: true, renderer })
}

function renderMarkdown(content: string) {
  setupMarked()
  const raw = marked.parse(content) as string
  return DOMPurify.sanitize(raw, { ADD_ATTR: ['target'] })
}

async function load() {
  const id = Number(route.params.id)
  if (!id) return
  loading.value = true
  try {
    article.value = await fetchArticle(id)
    rendered.value = renderMarkdown(article.value.content)
    await nextTick()
    setupCopyButtons()
    highlightAll()
  } finally {
    loading.value = false
  }
}

function setupCopyButtons() {
  const content = contentRef.value
  if (!content) return

  content.querySelectorAll('.copy-btn').forEach((btn) => {
    btn.addEventListener('click', () => {
      const target = btn as HTMLButtonElement
      const code = decodeURIComponent(target.dataset.code || '')
      navigator.clipboard.writeText(code).then(() => {
        const originalText = target.textContent
        target.textContent = '已复制'
        target.classList.add('copied')
        setTimeout(() => {
          target.textContent = originalText
          target.classList.remove('copied')
        }, 1800)
      })
    })
  })
}

function highlightAll() {
  const blocks = document.querySelectorAll('.markdown-body pre code:not(.hljs)')
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
  margin-bottom: 32px;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 20px;
}
.title {
  font-size: 30px;
  margin: 0 0 14px;
  line-height: 1.3;
  font-weight: 700;
}
.meta {
  display: flex;
  gap: 20px;
  font-size: 13px;
  color: var(--color-text-light);
  flex-wrap: wrap;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
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