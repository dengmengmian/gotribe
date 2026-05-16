import { request } from '@/service'

export interface DashboardStats {
  total_posts: number
  draft_posts: number
  pending_comments: number
  week_visits: number
}

export interface VisitPoint {
  date: string
  visits: number
  page_views: number
}

export interface PostSummary {
  id: number
  title: string
  status: number
  created_at: string
  view: number
}

export interface CommentSummary {
  id: number
  nickname: string
  content: string
  post_title: string
  created_at: string
}

export interface PendingCounts {
  pending_review_posts: number
  pending_review_comments: number
}

export interface SystemStatus {
  db_status: string
  redis_status: string
}

export interface CacheStatus {
  used_memory: string
  used_percent: number
}

export interface SeoAlert {
  type: 'warning' | 'info' | 'success'
  message: string
}

export interface DashboardData {
  stats: DashboardStats
  visit_trend: VisitPoint[]
  recent_posts: PostSummary[]
  recent_comments: CommentSummary[]
  popular_posts: PostSummary[]
  pending: PendingCounts
  system_status: SystemStatus
  cache_status: CacheStatus
  seo_alerts: SeoAlert[]
}

export interface DashboardResponse {
  indexDate: DashboardData
}

export async function getDashboard(params?: { project_id?: string }) {
  return request.get<DashboardResponse>('/api/index/dashboard', { params })
}

export async function clearCache() {
  return request.post<{ message: string }>('/api/index/cache/clear')
}
