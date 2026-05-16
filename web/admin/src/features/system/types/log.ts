export interface OperationLog {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  username: string
  ip: string
  ip_location: string
  method: string
  path: string
  desc: string
  status: number
  start_time: string
  time_cost: number
  user_agent: string
}
