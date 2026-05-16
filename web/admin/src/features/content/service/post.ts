import { request } from '@/service'
import type { Post, PostListParams, PostListResponse, PostParams } from '../types/post'

type RawPost = Partial<Post> & {
  ID?: number
  post_id?: number | string
}

function normalizePost(rawPost: RawPost): Post {
  return {
    ...rawPost,
    id: rawPost.id ?? rawPost.ID ?? (rawPost.post_id ? Number(rawPost.post_id) : 0),
    slug: rawPost.slug ?? '',
  } as Post
}

/** 获取文章列表（分页、按 id / title 筛选） */
export const getPostList = async (params?: PostListParams): Promise<PostListResponse> => {
  const query = {
    ...params,
  }
  const data = await request.get<PostListResponse>('/api/post', { params: query })
  // 兼容后端字段名：优先 id，其次 ID，最后尝试将旧版 post_id 转为数字
  const rawPosts = data.posts as RawPost[] | undefined
  const posts = (rawPosts ?? []).map(normalizePost)
  return { posts, total: data.total ?? 0 }
}

/** 获取文章详情（编辑回显）；GET /api/post/:id */
export const getPostDetail = async (id: number): Promise<Post> => {
  const data = await request.get<{ post: RawPost }>(`/api/post/${id}`)
  return normalizePost(data.post)
}

/** 创建文章 */
export async function createPost(params: PostParams): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/post', params)
}

/** 更新文章（路径参数为 id） */
export async function updatePost(
  id: number,
  params: Partial<PostParams>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/post/${id}`, params)
}

/** 快速发布文章 */
export async function publishPost(id: number): Promise<{ success: boolean }> {
  return request.put<{ success: boolean }>(`/api/post/${id}`)
}

/** 删除文章（body: ids 为数组） */
export async function deletePost(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/post', {
    data: { post_ids: [id] },
  })
}
