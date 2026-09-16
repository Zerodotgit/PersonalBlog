<template>
  <div class="list">
    <div class="head">
      <h2 class="page-title">文章管理</h2>
      <div class="actions">
        <el-select v-model="statusFilter" placeholder="全部状态" clearable style="width: 140px" @change="onFilter">
          <el-option label="已发布" :value="1" />
          <el-option label="草稿" :value="0" />
        </el-select>
        <el-button type="primary" @click="$router.push('/admin/articles/new')">写文章</el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="data?.items || []" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="220" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="success" size="small">已发布</el-tag>
          <el-tag v-else type="info" size="small">草稿</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览" width="90" />
      <el-table-column label="创建时间" width="150">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button text type="primary" @click="edit(row.id)">编辑</el-button>
          <el-button text type="danger" @click="remove(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="data?.total || 0"
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
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchAdminArticles,
  deleteArticle,
  type ListResponse,
} from '@/api/article'

const router = useRouter()
const loading = ref(false)
const data = ref<ListResponse | null>(null)
const page = ref(1)
const size = 10
const statusFilter = ref<number | undefined>(undefined)

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('zh-CN')
}

async function load() {
  loading.value = true
  try {
    data.value = await fetchAdminArticles(page.value, size, statusFilter.value)
  } finally {
    loading.value = false
  }
}

function onPage(p: number) {
  page.value = p
  load()
}

function onFilter() {
  page.value = 1
  load()
}

function edit(id: number) {
  router.push(`/admin/articles/${id}/edit`)
}

async function remove(id: number) {
  try {
    await ElMessageBox.confirm('确定删除这篇文章吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteArticle(id)
    ElMessage.success('已删除')
    load()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '删除失败')
  }
}

onMounted(load)
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
.actions {
  display: flex;
  gap: 12px;
}
.pager {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
