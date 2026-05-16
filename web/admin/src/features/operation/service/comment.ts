import { request } from '@/service'
import type { CommentListParams, CommentListResponse } from '../types/comment'

/**
 * 获取评论列表
 * GET /api/comment?page=1&per_page=10&status=&nickname=&project_id=
 */
export async function getCommentList(
  params?: CommentListParams
): Promise<CommentListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    page: params?.page ?? 1,
    per_page: params?.per_page ?? 10,
  }

  if (params?.status != null && params.status !== '') {
    requestParams.status = params.status
  }
  if (params?.nickname != null && params.nickname !== '') {
    requestParams.nickname = params.nickname
  }
  if (params?.project_id != null && params.project_id !== 0) {
    requestParams.project_id = params.project_id
  }
  if (params?.sort_by) {
    requestParams.sort_by = params.sort_by
  }
  if (params?.sort_order) {
    requestParams.sort_order = params.sort_order
  }

  const data = await request.get<CommentListResponse>('/api/comment', {
    params: requestParams,
  })
  return data as CommentListResponse
}

/**
 * 审核评论（通过）
 * PATCH /api/comment/:commentID  body: {}
 */
export async function approveComment(commentID: string): Promise<void> {
  const id = commentID.trim()
  await request.patch(`/api/comment/${id}`, {})
}

/**
 * 评论设为不通过
 * PATCH /api/comment/:commentID  body: { status: 1 }
 */
export async function rejectComment(commentID: string): Promise<void> {
  const id = commentID.trim()
  await request.patch(`/api/comment/${id}`, { status: 1 })
}

/** 批量删除评论 POST /api/comment body: { ids: [...] } */
export async function deleteComments(ids: number[]): Promise<unknown> {
  return request.delete('/api/comment', { data: { ids } })
}
