/** 积分项（与接口返回一致） */
export interface PointItem {
  id?: number
  point: number
  user_id: number
  reason: string
  nickname: string
  created_at: string
  updated_at: string
}

/** 列表接口返回 data 结构 */
export interface PointListResponse {
  points?: PointItem[]
  total?: number
}

/** 列表查询参数 */
export interface PointListParams {
  page?: number
  per_page?: number
  user_id?: string
  nickname?: string
  project_id?: number
}

/** 创建积分请求参数 */
export interface PointCreateParams {
  user_id: number
  project_id: number
  point: number
}
