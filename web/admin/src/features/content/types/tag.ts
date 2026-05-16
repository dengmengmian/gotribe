export type Tag = {
  id: number
  title: string
  slug: string
  color: string
  description: string
  sort?: number
  status?: number
  count?: number
  created_at: string
}

export type TagListResponse = {
  tags: Tag[]
  total: number
}
