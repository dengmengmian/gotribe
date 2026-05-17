import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { PageSkeleton } from '@/components/page-skeleton'

const TotpSetup = lazy(() =>
  import('@/features/auth/totp-setup').then((m) => ({
    default: m.TotpSetup,
  }))
)

export const Route = createFileRoute('/(auth)/totp-setup')({
  component: () => (
    <Suspense fallback={<PageSkeleton />}>
      <TotpSetup />
    </Suspense>
  ),
})
