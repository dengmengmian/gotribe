import { request } from '@/service'

export interface TOTPStatus {
  bound: boolean
  enabled: boolean
  last_used_at?: string
  remaining_recovery_codes: number
}

export interface TOTPBindResult {
  secret: string
  otpauth_url: string
  recovery_codes: string[]
}

export interface TOTPVerifyResult {
  stage: 'ok'
  token: string
  expires: string
}

/** 登录后两步校验，提交 step_token + 6 位 code（或备份码） */
export function verifyTOTP(stepToken: string, code: string) {
  return request.post<TOTPVerifyResult>('/api/base/totp/verify', {
    step_token: stepToken,
    code,
  })
}

/** 查询当前账户的 TOTP 状态（登录态） */
export function getTOTPStatus() {
  return request.get<TOTPStatus>('/api/base/totp/status')
}

/** 发起绑定（登录态）；返回 secret/QR/备份码，仅本次返回 */
export function bindTOTP() {
  return request.post<TOTPBindResult>('/api/base/totp/bind')
}

/** 用一次 6 位码确认并激活绑定 */
export function confirmTOTP(code: string) {
  return request.post<null>('/api/base/totp/confirm', { code })
}

/** 自助解绑，需提供当前 6 位码 */
export function deleteTOTP(code: string) {
  return request.delete<null>('/api/base/totp', { data: { code } })
}

/** 超管强制重置他人 TOTP */
export function adminResetTOTP(adminId: number) {
  return request.post<null>(`/api/base/admin/${adminId}/totp/reset`)
}
