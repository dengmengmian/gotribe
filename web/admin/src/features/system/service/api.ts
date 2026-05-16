import { request } from '@/service'
import type { Api } from '../types/api'

export const getApiList = async (params?: {
  current?: number
  path?: string
  category?: string
  method?: string
  creator?: string
  page?: number
  per_page?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}) => {
  return request.get<{ apis: Api[]; total: number }>('/api/api/list', { params })
}

export async function createApi(
  params: Partial<Api>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/api/create', params)
}

export async function updateApi(
  params: Partial<Api>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/api/update/${params.id}`, params)
}

export async function deleteApi(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/api/delete/batch', {
    data: { api_ids: [id] },
  })
}
