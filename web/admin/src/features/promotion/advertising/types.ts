/** 广告（列表项，与接口返回一致） */
export interface Ad {
  id?: number
  title?: string
  description?: string
  scene_id?: number
  SceneTitle?: string
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  url_type?: number
  ext?: string
  created_at?: string
  updated_at?: string
  [key: string]: unknown
}

/** 列表查询参数 */
export interface AdListParams {
  page?: number
  per_page?: number
  scene_id?: string
  title?: string
  status?: string | number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

/** 列表接口返回 data 结构 */
export interface AdListResponse {
  ads?: Ad[]
  total?: number
}

/** 创建广告请求参数 */
export interface AdCreateParams {
  title: string
  description: string
  scene_id: number
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  url_type?: number
  ext?: string
}

/** 更新广告请求参数 */
export interface AdUpdateParams {
  title?: string
  description?: string
  scene_id?: number
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  url_type?: number
  ext?: string
}
