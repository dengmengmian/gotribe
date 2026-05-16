import { request } from '@/service'
import { type PointListResponse, type PointListParams, type PointCreateParams } from '../types/point'

/**
 * 获取积分列表
 * GET /api/point?page=1&per_page=10&user_id=&nickname=&project_id=
 */
export async function getPointList(
  params?: PointListParams
): Promise<PointListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    page: params?.page ?? 1,
    per_page: params?.per_page ?? 10,
  }

  // 只有当参数有值时才添加到请求参数中
  if (params?.user_id != null && params.user_id !== '') {
    requestParams.user_id = params.user_id
  }
  if (params?.nickname != null && params.nickname !== '') {
    requestParams.nickname = params.nickname
  }
  if (params?.project_id != null && params.project_id !== 0) {
    requestParams.project_id = params.project_id
  }

  const data = await request.get<PointListResponse>('/api/point', {
    params: requestParams,
  })
  return data as PointListResponse
}

/**
 * 创建积分
 * POST /api/point
 */
export async function createPoint(
  data: PointCreateParams,
): Promise<unknown> {
  const result = await request.post('/api/point', data)
  return result ?? {}
}

/** 批量删除积分 DELETE /api/point body: { ids: [...] } */
export async function deletePoints(ids: number[]): Promise<unknown> {
  return request.delete('/api/point', { data: { ids } })
}
