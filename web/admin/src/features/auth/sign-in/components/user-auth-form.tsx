import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn, ShieldCheck } from 'lucide-react'
import { toast } from 'sonner'
import { useSetAccessToken, useSetAuthUser } from '@/stores/auth-store'
import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'
import { login, type LoginResult } from '../service'
import { getCurrentUser } from '@/service/user'

interface UserAuthFormProps extends React.HTMLAttributes<HTMLFormElement> {
  redirectTo?: string
}

const AUTH_ROUTES = ['/sign-in', '/otp', '/totp-setup']
const DEFAULT_SIGN_IN_REDIRECT = '/dashboard'

function resolveSignInRedirect(redirectTo?: string): string {
  if (!redirectTo) return DEFAULT_SIGN_IN_REDIRECT

  let targetPath = redirectTo
  let pathname = redirectTo
  try {
    const url = new URL(redirectTo, window.location.origin)
    if (url.origin !== window.location.origin) return DEFAULT_SIGN_IN_REDIRECT
    pathname = url.pathname
    targetPath = `${url.pathname}${url.search}${url.hash}`
  } catch {
    pathname = redirectTo
  }

  if (!targetPath.startsWith('/')) return DEFAULT_SIGN_IN_REDIRECT

  if (AUTH_ROUTES.some(route => pathname.startsWith(route))) {
    return DEFAULT_SIGN_IN_REDIRECT
  }

  return targetPath
}

export function UserAuthForm({
  className,
  redirectTo,
  ...props
}: UserAuthFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const setAccessToken = useSetAccessToken()
  const setAuthUser = useSetAuthUser()
  const { t } = useI18n()

  // Create form schema with translations
  const formSchema = z.object({
    username: z.string().min(1, t('features.auth.validation.usernameRequired')),
    password: z.string().min(1, t('features.auth.validation.passwordRequired')),
  })

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  async function completeLogin(token: string, username: string, showMfaReminder: boolean) {
    setAccessToken(token)
    let displayName = username
    try {
      const userInfoResult = await getCurrentUser()
      if (userInfoResult && userInfoResult.admin) {
        const admin = userInfoResult.admin as { nickname?: string; username?: string }
        displayName = admin.nickname?.trim() || admin.username || username
        setAuthUser({ ...userInfoResult.admin })
      } else {
        setAuthUser({})
      }
    } catch (error) {
      setAuthUser({})
      // eslint-disable-next-line no-console
      console.error('Failed to get user info:', error)
    }

    const targetPath = resolveSignInRedirect(redirectTo)
    navigate({ to: targetPath, replace: true })

    if (showMfaReminder) {
      // 登录成功但未绑 TOTP：提示一次，并提供「去绑定」入口；可关闭忽略
      toast.warning(t('features.auth.mfa.reminderTitle'), {
        description: t('features.auth.mfa.reminderDescription'),
        action: {
          label: t('features.auth.mfa.reminderAction'),
          onClick: () => navigate({ to: '/personal-center' }),
        },
        duration: 8000,
      })
    }

    return t('features.auth.signIn.welcomeBack', { username: displayName })
  }

  function onSubmit(data: z.infer<typeof formSchema>) {
    setIsLoading(true)

    toast.promise(
      login({ username: data.username, password: data.password }).then(async (result: LoginResult) => {
        if (result.stage === 'totp_required') {
          if (!result.step_token) {
            throw new Error(t('features.auth.signIn.error'))
          }
          sessionStorage.setItem('totp_step_token', result.step_token)
          sessionStorage.setItem('totp_step_redirect', redirectTo ?? '')
          navigate({ to: '/otp', replace: true })
          return t('features.auth.signIn.totpRequired')
        }
        if (result.stage === 'bind_required') {
          if (!result.step_token) {
            throw new Error(t('features.auth.signIn.error'))
          }
          sessionStorage.setItem('totp_bind_step_token', result.step_token)
          sessionStorage.setItem('totp_bind_redirect', redirectTo ?? '')
          navigate({ to: '/totp-setup', replace: true })
          return t('features.auth.signIn.bindRequired')
        }
        if (!result.token) {
          throw new Error(t('features.auth.signIn.error'))
        }
        return completeLogin(result.token, data.username, Boolean(result.mfa_reminder))
      }),
      {
        loading: t('features.auth.signIn.loggingIn'),
        success: (message) => {
          setIsLoading(false)
          return message
        },
        error: (error) => {
          setIsLoading(false)
          return error?.message || t('features.auth.signIn.error')
        },
      }
    )
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-4', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='username'
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('features.auth.signIn.account')}</FormLabel>
              <FormControl>
                <Input
                  autoComplete='username'
                  className='h-10'
                  placeholder={t('features.auth.signIn.accountPlaceholder')}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem className='relative'>
              <FormLabel>{t('features.auth.signIn.password')}</FormLabel>
              <FormControl>
                <PasswordInput
                  autoComplete='current-password'
                  placeholder={t('features.auth.signIn.passwordPlaceholder')}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type='submit' size='lg' className='mt-1 w-full' disabled={isLoading}>
          {isLoading ? <Loader2 className='animate-spin' /> : <LogIn />}
          {t('features.auth.signIn.login')}
        </Button>
        <p className='flex items-start gap-2 rounded-md border border-border/70 bg-muted/35 px-3 py-2 text-xs leading-5 text-muted-foreground'>
          <ShieldCheck className='mt-0.5 size-4 shrink-0 text-foreground/70' />
          <span>
          {redirectTo ? t('features.auth.signIn.redirectHint') : t('features.auth.signIn.securityHint')}
          </span>
        </p>
      </form>
    </Form>
  )
}
