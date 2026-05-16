import { request } from '@/service'
import type { OperationLog } from '../types/log'

export const getLogList = async (params?: {
  current?: number
  username?: string
  ip?: string
  path?: string
  status?: string
  page?: number
  per_page?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}) => {
  return request.get<{ logs: OperationLog[]; total: number }>(
    '/api/operation/list',
    { params }
  )
}

export async function deleteLog(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>(
    '/api/operation/delete/batch',
    { data: { operation_log_ids: [id] } }
  )
}
