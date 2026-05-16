/** 配置项（列表/详情） */
export interface Config {
  id: number
  alias: string
  title: string
  description: string
  type?: number
  info?: string
  md_content?: string
  project_id?: number
  project?: {
    id: number
    title?: string
    name?: string
  }
  created_at?: string
  updated_at?: string
}

/** 列表查询参数 */
export interface ConfigListParams {
  id?: number
  alias?: string
  title?: string
  page?: number
  per_page?: number
  project_id?: number
  type?: number
}

/** 列表接口返回 */
export interface ConfigListResponse {
  configs: Config[]
  total: number
}

/** 新增/创建配置请求体 */
export interface ConfigCreateParams {
  title: string
  description: string
  info?: string
  project_id: number
  alias: string
  type: number
  md_content?: string
}

/** 更新配置请求体（编辑时：不含 alias、type） */
export interface ConfigUpdateParams {
  title: string
  description: string
  project_id?: number
  info?: string
  md_content?: string
}
