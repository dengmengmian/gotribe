export type Api = {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  method: string
  path: string
  category: string
  desc: string
  creator: string,
  children: Api[]
}

export type ApiListResponse = {
  apis: Api[]
  total: number
}