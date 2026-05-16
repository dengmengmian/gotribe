import { request } from '@/service'

export type AIGenerateTask = 'post_metadata' | 'post_slug' | 'post_description'

export interface AIGenerateResponse {
  result: Record<string, unknown>
}

export async function generateAI(params: {
  task: AIGenerateTask
  language?: string
  input: Record<string, unknown>
}): Promise<AIGenerateResponse> {
  return request.post<AIGenerateResponse>('/api/ai/generate', params)
}
