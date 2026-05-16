export type User = {
  id: number
  user_id: string
  username: string // 用户名
  nickname: string // 用户昵称
  email?: string
  phone?: string
  /** 仅请求体使用，API 不返回 */
  password?: string
  avatar_url: string // 头像URL
  sex: 'M' | 'F' | '' | 'U' // 性别：M-男，F-女，U-未知
  project_id: number // 项目ID
  status: number // 状态
  birthday: string // 生日
  point: number // 积分
  created_at: string // 创建时间
}

export type UserListResponse = {
  users: User[]
  total: number
}
