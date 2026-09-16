<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>

    <el-skeleton v-if="loading" :rows="4" animated />

    <template v-else-if="data">
      <div class="stat-cards">
        <div class="stat-card">
          <div class="num">{{ data.stats.total_articles }}</div>
          <div class="label">文章总数</div>
        </div>
        <div class="stat-card">
          <div class="num">{{ data.stats.published_count }}</div>
          <div class="label">已发布</div>
        </div>
        <div class="stat-card">
          <div class="num">{{ data.stats.draft_count }}</div>
          <div class="label">草稿</div>
        </div>
        <div class="stat-card accent">
          <div class="num">{{ data.stats.total_views }}</div>
          <div class="label">总浏览量</div>
        </div>
      </div>

      <div class="recent">
        <div class="section-head">
          <h3>最近文章</h3>
          <el-button text type="primary" @click="$router.push('/admin/articles')">
            查看全部
          </el-button>
        </div>
        <el-table :data="data.recent" stripe>
          <el-table-column prop="title" label="标题" min-width="200" />
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag v-if="row.status === 1" type="success" size="small">已发布</el-tag>
              <el-tag v-else type="info" size="small">草稿</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" label="浏览" width="90" />
          <el-table-column label="创建时间" width="140">
            <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
          </el-table-column>
        </el-table>
        <div v-if="data.recent.length === 0" class="empty">还没有文章，去写一篇吧。</div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchStats, type StatsData } from '@/api/article'

const loading = ref(true)
const data = ref<StatsData | null>(null)

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('zh-CN')
}

onMounted(async () => {
  try {
    data.value = await fetchStats()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-title {
  margin: 0 0 24px;
  font-size: 22px;
}
.stat-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}
.stat-card {
  background: var(--color-surface);
  border-radius: var(--radius);
  padding: 24px;
  box-shadow: var(--shadow);
  text-align: center;
}
.stat-card.accent {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}
.stat-card.accent .label {
  color: rgba(255, 255, 255, 0.85);
}
.num {
  font-size: 32px;
  font-weight: 700;
}
.label {
  color: var(--color-text-light);
  font-size: 14px;
  margin-top: 6px;
}
.recent {
  background: var(--color-surface);
  border-radius: var(--radius);
  padding: 20px 24px;
  box-shadow: var(--shadow);
}
.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.section-head h3 {
  margin: 0;
  font-size: 16px;
}
.empty {
  text-align: center;
  color: var(--color-text-light);
  padding: 24px 0;
}
</style>
