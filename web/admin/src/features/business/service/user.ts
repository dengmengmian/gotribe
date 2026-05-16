import { request } from '@/service'
import type { User, UserListResponse } from '../types/user'

function mapUser(raw: Omit<User, 'user_id'>): User {
  return {
    ...raw,
    user_id: String(raw.id),
  }
}

// 获取用户详情（单条）
export async function getUserDetail(id: number): Promise<User> {
  const data = await request.get<{ user: User }>(`/api/user/${id}`)
  return mapUser(data.user)
}

// 获取用户列表（支持分页和筛选）
export const getUserList = async (params?: {
  current?: number
  user_id?: string
  project_id?: number
  page?: number
  per_page?: number
}) => {
  const query = {
    current: params?.current,
    user_id: params?.user_id,
    project_id: params?.project_id,
    page: params?.page,
    per_page: params?.per_page,
  }
  return request
    .get<UserListResponse>('/api/user', { params: query })
    .then((data) => ({
      ...data,
      users: (data.users ?? []).map((user) => mapUser(user)),
    }))
}

// 创建用户
export async function createUser(
  params: Partial<User>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/user', params)
}

// 更新用户
export async function updateUser(
  params: Partial<User>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(
    `/api/user/${params.id}`,
    params,
  )
}
