/** 文件类型（与后端 FILE_TYPE_* 一致） */
export const FILE_TYPE = {
  IMAGE: 1,
  VIDEO: 2,
  AUDIO: 3,
  ARCHIVE: 4,
  DOCUMENT: 5,
  FONT: 6,
  APP: 7,
  UNKNOWN: 8,
} as const

/** 接口 /api/resource 返回的单条资源（原始结构） */
export interface ResourceApiItem {
  id: number
  title: string
  description?: string
  url: string
  path: string
  file_type: number
  file_extension?: string
  size?: number
  created_at?: string
}

/** 上传接口返回的原始结构 */
export interface UploadResourceApiItem {
  file_ext: string
  key: string
  domain: string
  file_type: number
}

/** 前端统一使用的资源项（由 API 项映射） */
export interface ResourceItem {
  id: number
  url: string
  name: string
  type: number
  size?: number
  created_at?: string
}

/** 资源列表查询参数 */
export interface ResourceListParams {
  type?: number
  page?: number
  per_page?: number
}

/** 资源列表响应 */
export interface ResourceListResponse {
  list: ResourceItem[]
  total: number
}
