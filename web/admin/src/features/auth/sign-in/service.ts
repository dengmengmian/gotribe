import { request } from '@/service'

export interface LoginParams {
  username: string
  password: string
}

export type LoginStage = 'ok' | 'totp_required' | 'bind_required'

/** 登录响应：根据 stage 字段不同填充不同字段。 */
export interface LoginResult {
  stage: LoginStage
  /** stage='ok' 时存在：可用 access_token */
  token?: string
  /** stage='ok' 时存在：过期时间字符串 */
  expires?: string
  /** stage='ok' 时存在：未绑 TOTP 的提示标志 */
  mfa_reminder?: boolean
  /** stage='totp_required' 时存在：5 分钟内有效的 step token */
  step_token?: string
  /** stage='totp_required' 时存在：step token 过期时间字符串 */
  step_expires?: string
}

/** 登录接口 POST /api/base/login */
export async function login(body: LoginParams): Promise<LoginResult> {
  return request.post<LoginResult>('/api/base/login', body)
}
