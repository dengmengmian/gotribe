import { createFileRoute, Outlet, useLocation } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentArticle = lazy(() =>
  import('@/features/content/article').then((m) => ({
    default: m.ContentArticle,
  }))
)

export const Route = createFileRoute('/_authenticated/content/article')({
  component: ArticleLayout,
})

function ArticleLayout() {
  const pathname = useLocation({ select: (location) => location.pathname })
  const isListRoute = pathname === '/content/article' || pathname === '/content/article/'

  if (isListRoute) {
    return (
      <Suspense fallback={<TablePageSkeleton />}>
        <ContentArticle />
      </Suspense>
    )
  }

  return <Outlet />
}
