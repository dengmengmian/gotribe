/** 评论（列表项，与接口返回一致） */
export interface Comment {
  id?: number
  comment_id?: string
  project_id?: number
  status?: number
  user_id?: number
  object_id?: string
  object_type?: number
  comment?: string
  html_content?: string
  nickname?: string
  ip?: string
  country?: string
  region_name?: string
  city?: string
  created_at?: string
  updated_at?: string
  [key: string]: unknown
}

/** 列表查询参数 */
export interface CommentListParams {
  page?: number
  per_page?: number
  status?: string | number
  nickname?: string
  project_id?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/** 列表接口返回 data 结构 */
export interface CommentListResponse {
  comments?: Comment[]
  total?: number
}
