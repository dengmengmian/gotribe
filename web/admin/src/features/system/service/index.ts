import { request } from '@/service'
import type { Admin } from '../types/admin'

/** 获取当前管理员信息；GET /api/admin/info，用于 Profile 页展示 */
export interface AdminInfoParams {
  project_id?: string
  title?: string
  page?: number
  per_page?: number
}

export async function getAdminInfo(
  params?: AdminInfoParams
): Promise<Admin | null> {
  const data = await request.get<{ admin: Admin }>(
    '/api/admin/info',
    { params: params ?? { page: 1, per_page: 3 } }
  )
  return data?.admin ?? null
}

/** 修改当前管理员密码；PUT /api/admin/changePwd */
export async function changePassword(params: {
  old_password: string
  new_password: string
}): Promise<void> {
  await request.put<void>('/api/admin/changePwd', params)
}

/** 系统管理员列表 /api/admin/list */
export const getAdminList = async (params: {
  current: number
  username: string
  nickname: string
  status: string
  mobile: string
  page: number
  per_page: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}) => {
  return request.get<{ admins: Admin[]; total: number }>('/api/admin/list', { params })
}

/** 新增管理员 */
export async function createAdmin(
  params: Partial<Admin>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/admin/create', params)
}

/** 更新管理员 */
export async function updateAdmin(
  params: Partial<Admin>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/admin/update/${params.id}`, params)
}

/** 删除管理员 */
export async function deleteAdmin(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/admin/delete/batch', {
    data: { user_ids: [id] },
  })
}
