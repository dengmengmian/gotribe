import { request } from '@/service'
import type {
  Config,
  ConfigListParams,
  ConfigListResponse,
  ConfigCreateParams,
  ConfigUpdateParams,
} from '../types/config'

/**
 * 获取配置列表
 * GET /api/config?id=&alias=&title=&page=1&per_page=3&project_id=&type=
 */
export async function getConfigList(
  params?: ConfigListParams
): Promise<ConfigListResponse> {
  const data = await request.get<ConfigListResponse>('/api/config', {
    params: {
      id: params?.id ?? '',
      alias: params?.alias ?? '',
      title: params?.title ?? '',
      page: params?.page ?? 1,
      per_page: params?.per_page ?? 10,
      project_id: params?.project_id,
      type: params?.type,
    },
  })
  return data as ConfigListResponse
}

/**
 * 新增配置
 * POST /api/config  body: title, description, info, project_id, alias, type, md_content
 */
export async function createConfig(
  data: ConfigCreateParams
): Promise<{ id?: number }> {
  const result = await request.post<{ id?: number }>('/api/config', data)
  return (result ?? {}) as { id?: number }
}

/**
 * 获取配置详情
 * GET /api/config/:id
 */
export async function getConfig(configID: number): Promise<Config> {
  const id = String(configID)
  const data = await request.get<{ config: Config }>(`/api/config/${id}`)
  return (data as { config: Config }).config
}

/**
 * 更新配置
 * PATCH /api/config/:id  body: title, description, project_id, info, md_content
 */
export async function updateConfig(
  configID: number,
  data: ConfigUpdateParams
): Promise<void> {
  const id = String(configID)
  await request.patch(`/api/config/${id}`, data)
}

/**
 * 删除配置
 * DELETE /api/config  body: { ids: number[] }
 */
export async function deleteConfig(configIds: number): Promise<void> {
  await request.delete('/api/config', {
    data: { ids: [configIds] },
  })
}
