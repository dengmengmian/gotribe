import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { FormPageSkeleton } from '@/components/page-skeleton'

const SettingsSecurity = lazy(() =>
  import('@/features/settings/security').then((m) => ({
    default: m.SettingsSecurity,
  }))
)

export const Route = createFileRoute('/_authenticated/personal-center/security')({
  component: () => (
    <Suspense fallback={<FormPageSkeleton />}>
      <SettingsSecurity />
    </Suspense>
  ),
})
