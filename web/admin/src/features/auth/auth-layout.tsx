import { useI18n } from '@/context/i18n-provider'

type AuthLayoutProps = {
  children: React.ReactNode
}

export function AuthLayout({ children }: AuthLayoutProps) {
  const { t } = useI18n()

  return (
    <main className='relative flex min-h-svh items-center justify-center overflow-hidden bg-background px-4 py-8 text-foreground sm:px-6'>
      <div className='pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:44px_44px] opacity-20 dark:opacity-15' />
      <div className='pointer-events-none absolute inset-x-0 top-0 h-40 bg-[linear-gradient(180deg,var(--background),transparent)]' />
      <div className='pointer-events-none absolute inset-x-0 bottom-0 h-52 bg-[linear-gradient(0deg,var(--background),transparent)]' />

      <div className='relative flex w-full max-w-[560px] flex-col'>
        <div className='mb-7 text-center'>
          <p className='text-sm font-medium text-muted-foreground'>{t('features.auth.layout.description')}</p>
          <h2 className='mt-2 text-2xl font-semibold tracking-tight text-foreground sm:text-3xl'>
            {t('features.auth.layout.headline')}
          </h2>
        </div>
        {children}
      </div>
    </main>
  )
}
