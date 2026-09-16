import http from './http'

export interface ArticleListItem {
  id: number
  title: string
  summary: string
  cover_image: string
  tags: string
  status: number
  view_count: number
  created_at: string
}

export interface ArticleDetail {
  id: number
  title: string
  summary: string
  content: string
  cover_image: string
  tags: string
  status: number
  view_count: number
  created_at: string
  updated_at: string
}

export interface ListResponse {
  total: number
  items: ArticleListItem[]
  page: number
  size: number
}

export interface ArticleInput {
  title: string
  summary: string
  content: string
  cover_image: string
  tags: string
  status: number
}

// 公开 API
export function fetchPublicArticles(page = 1, size = 10) {
  return http
    .get<ListResponse>('/articles', { params: { page, size } })
    .then((r) => r.data)
}

export function fetchArticle(id: number) {
  return http.get<ArticleDetail>(`/articles/${id}`).then((r) => r.data)
}

// 管理 API
export function fetchAdminArticles(page = 1, size = 10, status?: number) {
  return http
    .get<ListResponse>('/admin/articles', { params: { page, size, status } })
    .then((r) => r.data)
}

export function fetchAdminArticle(id: number) {
  return http.get<ArticleDetail>(`/admin/articles/${id}`).then((r) => r.data)
}

export function createArticle(input: ArticleInput) {
  return http.post<ArticleDetail>('/admin/articles', input).then((r) => r.data)
}

export function updateArticle(id: number, input: ArticleInput) {
  return http.put<ArticleDetail>(`/admin/articles/${id}`, input).then((r) => r.data)
}

export function deleteArticle(id: number) {
  return http.delete(`/admin/articles/${id}`).then((r) => r.data)
}

export interface StatsData {
  stats: {
    total_articles: number
    published_count: number
    draft_count: number
    total_views: number
  }
  recent: ArticleListItem[]
}

export function fetchStats() {
  return http.get<StatsData>('/admin/stats').then((r) => r.data)
}
