<template>
  <div class="edit">
    <div class="head">
      <h2 class="page-title">{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="head-actions">
        <el-button @click="$router.back()">返回</el-button>
      </div>
    </div>

    <el-skeleton v-if="loading" :rows="8" animated />

    <el-form v-else :model="form" label-position="top" class="form">
      <el-form-item label="标题" required>
        <el-input v-model="form.title" placeholder="请输入文章标题" maxlength="255" show-word-limit />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="16">
          <el-form-item label="摘要">
            <el-input
              v-model="form.summary"
              type="textarea"
              :rows="2"
              placeholder="一句话简介，留空则前台显示“暂无摘要”"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="标签">
            <el-input v-model="form.tags" placeholder="逗号分隔，如：Go,Vue" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="封面图 URL">
        <el-input v-model="form.cover_image" placeholder="可选，填写图片链接" />
      </el-form-item>

      <el-form-item label="正文（Markdown）" required>
        <MarkdownEditor v-model="form.content" placeholder="在此使用 Markdown 编写正文……" />
      </el-form-item>

      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">发布</el-radio>
          <el-radio :value="0">草稿</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="saving" @click="onSave">
          {{ isEdit ? '保存修改' : '创建文章' }}
        </el-button>
        <el-button @click="$router.back()">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import {
  fetchAdminArticle,
  createArticle,
  updateArticle,
  type ArticleInput,
} from '@/api/article'

const route = useRoute()
const router = useRouter()

const id = computed(() => (route.params.id ? Number(route.params.id) : 0))
const isEdit = computed(() => id.value > 0)

const loading = ref(false)
const saving = ref(false)
const form = reactive<ArticleInput>({
  title: '',
  summary: '',
  content: '',
  cover_image: '',
  tags: '',
  status: 1,
})

async function loadArticle() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const a = await fetchAdminArticle(id.value)
    form.title = a.title
    form.summary = a.summary
    form.content = a.content
    form.cover_image = a.cover_image
    form.tags = a.tags
    form.status = a.status
  } finally {
    loading.value = false
  }
}

async function onSave() {
  if (!form.title.trim()) {
    ElMessage.warning('请填写标题')
    return
  }
  if (!form.content.trim()) {
    ElMessage.warning('请填写正文')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateArticle(id.value, { ...form })
      ElMessage.success('已保存')
    } else {
      await createArticle({ ...form })
      ElMessage.success('已创建')
    }
    router.push('/admin/articles')
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadArticle)
</script>

<style scoped>
.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.page-title {
  margin: 0;
  font-size: 22px;
}
.form {
  background: var(--color-surface);
  border-radius: var(--radius);
  padding: 24px;
  box-shadow: var(--shadow);
}
</style>
