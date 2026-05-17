import { useEffect, useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { AlertTriangle, Copy, Loader2, ShieldCheck } from 'lucide-react'
import { useI18n } from '@/context/i18n-provider'
import { useSetAccessToken, useSetAuthUser } from '@/stores/auth-store'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from '@/components/ui/input-otp'
import {
  confirmEnrollPendingTOTP,
  enrollPendingTOTP,
  type TOTPBindResult,
} from '../../totp/service'
import { getCurrentUser } from '@/service/user'

const STEP_TOKEN_KEY = 'totp_bind_step_token'
const REDIRECT_KEY = 'totp_bind_redirect'
const DEFAULT_REDIRECT = '/dashboard'

export function TotpSetupForm() {
  const { t } = useI18n()
  const navigate = useNavigate()
  const setAccessToken = useSetAccessToken()
  const setAuthUser = useSetAuthUser()
  const [stepToken] = useState<string | null>(() => sessionStorage.getItem(STEP_TOKEN_KEY))
  const [enrollResult, setEnrollResult] = useState<TOTPBindResult | null>(null)
  const [enrolling, setEnrolling] = useState(true)
  const [confirming, setConfirming] = useState(false)
  const [code, setCode] = useState('')

  // 没有 step_token 直接回登录页
  useEffect(() => {
    if (!stepToken) {
      navigate({ to: '/sign-in', replace: true })
    }
  }, [stepToken, navigate])

  // 进入页面立即发起 enroll，一次性产出 secret/QR/备份码
  useEffect(() => {
    if (!stepToken) return
    let cancelled = false
    enrollPendingTOTP(stepToken)
      .then((result) => {
        if (!cancelled) setEnrollResult(result)
      })
      .catch((error: Error) => {
        if (!cancelled) {
          toast.error(error.message || t('features.auth.totpSetup.enrollFailed'))
          navigate({ to: '/sign-in', replace: true })
        }
      })
      .finally(() => {
        if (!cancelled) setEnrolling(false)
      })
    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [stepToken])

  async function copy(value: string) {
    try {
      await navigator.clipboard.writeText(value)
      toast.success(t('features.auth.totpSetup.copied'))
    } catch {
      toast.error(t('features.auth.totpSetup.copyFailed'))
    }
  }

  async function handleConfirm() {
    if (!stepToken || code.length !== 6) return
    setConfirming(true)
    try {
      const result = await confirmEnrollPendingTOTP(stepToken, code)
      setAccessToken(result.token)
      let displayName = ''
      try {
        const userInfoResult = await getCurrentUser()
        if (userInfoResult && userInfoResult.admin) {
          const admin = userInfoResult.admin as { nickname?: string; username?: string }
          displayName = admin.nickname?.trim() || admin.username || ''
          setAuthUser({ ...userInfoResult.admin })
        } else {
          setAuthUser({})
        }
      } catch (error) {
        setAuthUser({})
        // eslint-disable-next-line no-console
        console.error('Failed to get user info:', error)
      }
      const redirect = sessionStorage.getItem(REDIRECT_KEY) || DEFAULT_REDIRECT
      sessionStorage.removeItem(STEP_TOKEN_KEY)
      sessionStorage.removeItem(REDIRECT_KEY)
      toast.success(t('features.auth.signIn.welcomeBack', { username: displayName || 'admin' }))
      navigate({ to: redirect || DEFAULT_REDIRECT, replace: true })
    } catch (error) {
      const message = error instanceof Error ? error.message : t('features.auth.totpSetup.confirmFailed')
      toast.error(message)
    } finally {
      setConfirming(false)
    }
  }

  if (enrolling || !enrollResult) {
    return (
      <div className='flex items-center justify-center py-12'>
        <Loader2 className='size-5 animate-spin text-muted-foreground' />
      </div>
    )
  }

  return (
    <div className='space-y-5'>
      <div className='flex flex-col items-center gap-3'>
        <div className='rounded-lg border border-border bg-white p-3 shadow-sm'>
          <QRCodeSVG value={enrollResult.otpauth_url} size={180} level='M' />
        </div>
        <div className='w-full space-y-1'>
          <Label className='text-xs text-muted-foreground'>
            {t('features.auth.totpSetup.manualSecret')}
          </Label>
          <div className='flex items-center gap-2'>
            <Input value={enrollResult.secret} readOnly className='font-mono text-xs' />
            <Button
              type='button'
              size='icon'
              variant='outline'
              onClick={() => copy(enrollResult.secret)}
              aria-label='copy'
            >
              <Copy className='size-4' />
            </Button>
          </div>
        </div>
      </div>

      <div className='space-y-2 rounded-lg border border-amber-300/70 bg-amber-50 p-3 text-xs text-amber-900 dark:border-amber-500/40 dark:bg-amber-500/10 dark:text-amber-200'>
        <div className='flex items-center gap-2 font-medium'>
          <AlertTriangle className='size-4' />
          {t('features.auth.totpSetup.recoveryWarningTitle')}
        </div>
        <p>{t('features.auth.totpSetup.recoveryWarningBody')}</p>
        <div className='grid grid-cols-2 gap-1 pt-1 font-mono text-xs'>
          {enrollResult.recovery_codes.map((c) => (
            <code key={c} className='rounded bg-amber-100 px-1.5 py-0.5 dark:bg-amber-500/20'>
              {c}
            </code>
          ))}
        </div>
        <Button
          type='button'
          variant='outline'
          size='sm'
          className='mt-1 w-full'
          onClick={() => copy(enrollResult.recovery_codes.join('\n'))}
        >
          <Copy className='mr-2 size-3' />
          {t('features.auth.totpSetup.copyAll')}
        </Button>
      </div>

      <div className='space-y-2'>
        <Label className='text-xs text-muted-foreground'>
          {t('features.auth.totpSetup.confirmCodeLabel')}
        </Label>
        <InputOTP
          maxLength={6}
          value={code}
          onChange={setCode}
          containerClassName='justify-between sm:[&>[data-slot="input-otp-group"]>div]:w-12'
        >
          <InputOTPGroup>
            <InputOTPSlot index={0} />
            <InputOTPSlot index={1} />
          </InputOTPGroup>
          <InputOTPSeparator />
          <InputOTPGroup>
            <InputOTPSlot index={2} />
            <InputOTPSlot index={3} />
          </InputOTPGroup>
          <InputOTPSeparator />
          <InputOTPGroup>
            <InputOTPSlot index={4} />
            <InputOTPSlot index={5} />
          </InputOTPGroup>
        </InputOTP>
      </div>

      <Button
        type='button'
        className='w-full'
        onClick={handleConfirm}
        disabled={code.length !== 6 || confirming}
      >
        {confirming ? <Loader2 className='mr-2 size-4 animate-spin' /> : <ShieldCheck className='mr-2 size-4' />}
        {t('features.auth.totpSetup.confirmAndContinue')}
      </Button>
    </div>
  )
}
