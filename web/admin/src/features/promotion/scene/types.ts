/** 广告场景（列表项，与接口返回一致） */
export interface Scene {
  id: number
  title: string
  description?: string
  project_id: number
  projectTitle?: string
  created_at?: string
  updated_at?: string
}

/** 列表查询参数 */
export interface SceneListParams {
  page?: number
  per_page?: number
  project_id?: number
}

/** 列表接口返回 data 结构 */
export interface SceneListResponse {
  adScenes: Scene[]
  total: number
}

/** 新建广告场景请求参数 */
export interface SceneCreateParams {
  title: string
  description: string
  project_id: number
}

/** 更新广告场景请求参数（PATCH body） */
export interface SceneUpdateParams {
  title?: string
  description?: string
  project_id?: number
}
