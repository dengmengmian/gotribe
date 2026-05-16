import { request } from '@/service'
import type {
  SceneCreateParams,
  SceneListParams,
  SceneListResponse,
  SceneUpdateParams,
} from './types'

/**
 * 获取广告场景列表
 * GET /api/ad/scene?page=1&per_page=10&project_id=
 */
export async function getSceneList(
  params?: SceneListParams
): Promise<SceneListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    page: params?.page ?? 1,
    per_page: params?.per_page ?? 10,
  }
  if (params?.project_id != null && params.project_id !== 0) {
    requestParams.project_id = params.project_id
  }
  const data = await request.get<SceneListResponse>('/api/ad/scene', {
    params: requestParams,
  })
  return data as SceneListResponse
}

/**
 * 新建广告场景
 * POST /api/ad/scene  body: title, description, project_id
 */
export async function createScene(
  data: SceneCreateParams
): Promise<unknown> {
  const result = await request.post('/api/ad/scene', data)
  return result ?? {}
}

/**
 * 更新广告场景
 * PATCH /api/ad/scene/:adSceneID  body: { title, description }
 */
export async function updateScene(
  adSceneID: string,
  data: SceneUpdateParams
): Promise<void> {
  const id = adSceneID.trim()
  await request.patch(`/api/ad/scene/${id}`, data)
}

/**
 * 删除广告场景
 * DELETE /api/ad/scene  body: { ids: number[] }
 */
export async function deleteScene(adSceneIds: string): Promise<void> {
  await request.delete('/api/ad/scene', {
    data: { ids: [Number(adSceneIds.trim())] },
  })
}
