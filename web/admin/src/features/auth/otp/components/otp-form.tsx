import { useEffect, useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { Loader2, ShieldCheck } from 'lucide-react'
import { useI18n } from '@/context/i18n-provider'
import { useSetAccessToken, useSetAuthUser } from '@/stores/auth-store'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from '@/components/ui/input-otp'
import { verifyTOTP } from '../../totp/service'
import { getCurrentUser } from '@/service/user'

type OtpFormProps = React.HTMLAttributes<HTMLFormElement>

const STEP_TOKEN_KEY = 'totp_step_token'
const STEP_REDIRECT_KEY = 'totp_step_redirect'
const DEFAULT_REDIRECT = '/dashboard'

export function OtpForm({ className, ...props }: OtpFormProps) {
  const navigate = useNavigate()
  const setAccessToken = useSetAccessToken()
  const setAuthUser = useSetAuthUser()
  const { t } = useI18n()
  const [isLoading, setIsLoading] = useState(false)
  const [useRecoveryCode, setUseRecoveryCode] = useState(false)
  const [stepToken] = useState<string | null>(() => sessionStorage.getItem(STEP_TOKEN_KEY))

  // 没有 step_token 直接跳回登录页
  useEffect(() => {
    if (!stepToken) {
      navigate({ to: '/sign-in', replace: true })
    }
  }, [stepToken, navigate])

  const otpSchema = z.object({
    otp: z.string().min(6, t('features.auth.otp.validation.otpLength')).max(6, t('features.auth.otp.validation.otpLength')),
  })
  const recoverySchema = z.object({
    otp: z.string().min(4, t('features.auth.otp.validation.recoveryRequired')),
  })

  const form = useForm<{ otp: string }>({
    resolver: zodResolver(useRecoveryCode ? recoverySchema : otpSchema),
    defaultValues: { otp: '' },
  })

  // eslint-disable-next-line react-hooks/incompatible-library
  const otp = form.watch('otp')

  async function completeLogin(token: string) {
    setAccessToken(token)
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

    const redirect = sessionStorage.getItem(STEP_REDIRECT_KEY) || DEFAULT_REDIRECT
    sessionStorage.removeItem(STEP_TOKEN_KEY)
    sessionStorage.removeItem(STEP_REDIRECT_KEY)
    navigate({ to: redirect || DEFAULT_REDIRECT, replace: true })
    return t('features.auth.signIn.welcomeBack', { username: displayName || 'admin' })
  }

  function onSubmit(data: { otp: string }) {
    if (!stepToken) return
    setIsLoading(true)
    toast.promise(
      verifyTOTP(stepToken, data.otp).then((result) => completeLogin(result.token)),
      {
        loading: t('features.auth.otp.verifying'),
        success: (message) => {
          setIsLoading(false)
          return message
        },
        error: (error) => {
          setIsLoading(false)
          return error?.message || t('features.auth.otp.error')
        },
      }
    )
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-3', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='otp'
          render={({ field }) => (
            <FormItem>
              <FormLabel className='sr-only'>{t('features.auth.otp.otpLabel')}</FormLabel>
              <FormControl>
                {useRecoveryCode ? (
                  <Input
                    autoComplete='one-time-code'
                    placeholder={t('features.auth.otp.recoveryPlaceholder')}
                    className='h-10 font-mono tracking-widest'
                    {...field}
                  />
                ) : (
                  <InputOTP
                    maxLength={6}
                    {...field}
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
                )}
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          className='mt-2'
          disabled={isLoading || (useRecoveryCode ? otp.length < 4 : otp.length < 6)}
        >
          {isLoading ? <Loader2 className='animate-spin' /> : <ShieldCheck />}
          {t('features.auth.otp.verify')}
        </Button>
        <Button
          type='button'
          variant='link'
          className='h-auto px-0 text-xs text-muted-foreground'
          onClick={() => {
            form.reset({ otp: '' })
            setUseRecoveryCode((v) => !v)
          }}
        >
          {useRecoveryCode
            ? t('features.auth.otp.switchToCode')
            : t('features.auth.otp.switchToRecovery')}
        </Button>
      </form>
    </Form>
  )
}
