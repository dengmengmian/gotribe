/** 文章（接口 /api/post 返回项） */
export interface Post {
  id: number
  slug: string
  title: string
  description: string
  category_id: number
  project_id: number
  user_id: number
  author: string
  content: string
  html_content: string
  ext: string
  icon: string
  tag: string
  type: number
  is_top: number
  is_passwd: number
  category: unknown | null
  tags: unknown | null
  project: unknown | null
  created_at: string
  status: number
  location: string
  people: string
  time: string
  images: string[] | null
  unit_price: number
  video: string
  password?: string
  show_time?: string
}

/** 文章列表查询参数 */
export interface PostListParams {
  id?: number
  title?: string
  status?: string
  project_id?: number
  page?: number
  per_page?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/** 创建/更新文章请求参数 */
export interface PostParams {
  title: string
  slug?: string
  description?: string
  author?: string
  user_id?: number
  content?: string
  status?: number
  icon?: string
  category_id?: number
  is_top?: number
  tag?: string
  // 新增字段
  type?: number
  project_id?: number
  is_passwd?: number
  video?: string
  images?: string[]
  password?: string
  html_content?: string
  show_time?: string
  /** 自定义字段 JSON 字符串，存 key-value 对象 */
  ext?: string
}

/** 文章列表响应 */
export interface PostListResponse {
  posts: Post[]
  total: number
}
