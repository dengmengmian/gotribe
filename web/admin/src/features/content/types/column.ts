/** 专栏（列表项） */
export interface Column {
  id: number
  title: string
  info: string
  project_id: number
  icon: string
  description: string
  created_at: string
}

/** 列表查询参数（与接口 snake_case 一致） */
export interface ColumnListParams {
  project_id?: number
  title?: string
  page?: number
  per_page?: number
}

/** 列表接口返回 */
export interface ColumnListResponse {
  columns: Column[]
  total: number
}

/** 新增专栏请求体（均为必填） */
export interface ColumnCreateParams {
  title: string
  description: string
  info: string
  project_id: number
  icon: string
}

/** 更新专栏请求体 */
export interface ColumnUpdateParams {
  title?: string
  info?: string
  description?: string
  icon?: string
  project_id?: number
}
