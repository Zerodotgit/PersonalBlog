<template>
  <div class="md-editor">
    <div class="toolbar">
      <span class="tip">支持 Markdown 语法</span>
    </div>
    <div class="panes">
      <textarea
        class="input"
        :value="modelValue"
        :placeholder="placeholder"
        spellcheck="false"
        @input="onInput"
      ></textarea>
      <div class="preview markdown-body" v-html="preview"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, nextTick, watch } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

const cache = ref('')

const preview = computed(() => {
  const raw = marked.parse(props.modelValue || '', { breaks: true }) as string
  const clean = DOMPurify.sanitize(raw)
  cache.value = clean
  nextTick(highlightAll)
  return clean
})

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}

function highlightAll() {
  const blocks = document.querySelectorAll('.md-editor .preview pre code')
  blocks.forEach((el) => {
    if (!el.getAttribute('data-highlighted')) {
      hljs.highlightElement(el as HTMLElement)
      el.setAttribute('data-highlighted', 'yes')
    }
  })
}

watch(
  () => props.modelValue,
  () => {
    nextTick(highlightAll)
  },
)

onMounted(() => {
  nextTick(highlightAll)
})
</script>

<style scoped>
.md-editor {
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--color-surface);
}
.toolbar {
  padding: 8px 14px;
  border-bottom: 1px solid var(--color-border);
  background: #fafafa;
}
.tip {
  font-size: 12px;
  color: var(--color-text-light);
}
.panes {
  display: grid;
  grid-template-columns: 1fr 1fr;
  height: 420px;
}
.input {
  border: none;
  border-right: 1px solid var(--color-border);
  resize: none;
  outline: none;
  padding: 14px 16px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 14px;
  line-height: 1.6;
  background: #fff;
}
.preview {
  padding: 14px 16px;
  overflow: auto;
  font-size: 14px;
}
</style>
