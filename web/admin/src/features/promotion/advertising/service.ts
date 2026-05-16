import { request } from '@/service'
import type {
  AdListParams,
  AdListResponse,
  AdCreateParams,
  AdUpdateParams,
} from './types'

/**
 * 获取广告列表
 * GET /api/ad?page=1&per_page=10&scene_id=&title&status
 */
export async function getAdList(
  params?: AdListParams
): Promise<AdListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    page: params?.page ?? 1,
    per_page: params?.per_page ?? 10,
  }

  // 只有当参数有值时才添加到请求参数中
  if (params?.scene_id != null && params.scene_id !== '') {
    requestParams.scene_id = params.scene_id
  }
  if (params?.title != null && params.title !== '') {
    requestParams.title = params.title
  }
  if (params?.status != null && params.status !== '') {
    requestParams.status = params.status
  }
  if (params?.sort_by) {
    requestParams.sort_by = params.sort_by
  }
  if (params?.sort_order) {
    requestParams.sort_order = params.sort_order
  }

  const data = await request.get<AdListResponse>('/api/ad', {
    params: requestParams,
  })
  return data as AdListResponse
}

/**
 * 创建广告
 * POST /api/ad?project_id
 */
export async function createAd(
  data: AdCreateParams
): Promise<unknown> {
  const requestData: Record<string, unknown> = {
    title: data.title,
    description: data.description,
    scene_id: data.scene_id,
    url: data.url,
    url_type: data.url_type,
    image: data.image,
    sort: data.sort,
    status: data.status,
    ext: data.ext,
  }
  // 只有当 video 有值时才添加
  if (data.video) {
    requestData.video = data.video
  }
  const result = await request.post('/api/ad', requestData)
  return result ?? {}
}

/**
 * 更新广告
 * PATCH /api/ad/:adID
 */
export async function updateAd(
  adID: string,
  data: AdUpdateParams
): Promise<void> {
  const id = adID.trim()
  await request.patch(`/api/ad/${id}`, data)
}

/**
 * 删除广告
 * DELETE /api/ad  body: { ids: number[] }
 */
export async function deleteAd(adIds: string): Promise<void> {
  await request.delete('/api/ad', {
    data: { ids: [Number(adIds.trim())] },
  })
}
