import { useState, useEffect, lazy, Suspense } from 'react'
import { Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Search,
  Plus,
  FileText,
  FileEdit,
  MessageCircle,
  Eye,

  Image,
  Tags,
  FolderOpen,
  AlertCircle,
  CheckCircle2,
  Clock,
  ChevronRight,
  Zap,
  Shield,
  AlertTriangle,
} from 'lucide-react'

import { useAuthUser } from '@/stores/auth-store'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'
import { getDashboard, clearCache, type DashboardData } from './service'
import { useCssVariables } from '@/hooks/use-theme-change'

const VisitTrendChart = lazy(() => import('./components/visit-trend-chart'))

const CHART_CSS_VARS = ['--chart-1', '--chart-2', '--chart-3', '--chart-4', '--chart-5']

const QUICK_ACTIONS = [
  { label: '新建文章', icon: Plus, path: '/content/article/new', color: 'bg-blue-500 hover:bg-blue-600' },
  { label: '上传图片', icon: Image, path: '/content/resource', color: 'bg-violet-500 hover:bg-violet-600' },
  { label: '管理标签', icon: Tags, path: '/content/tag', color: 'bg-amber-500 hover:bg-amber-600' },
  { label: '查看分类', icon: FolderOpen, path: '/content/category', color: 'bg-emerald-500 hover:bg-emerald-600' },
] as const

const SYSTEM_STATUS_STATIC = [
  { label: '系统运行时间', icon: Clock },
  { label: '数据库连接', icon: CheckCircle2 },
  { label: '缓存状态', icon: Zap },
  { label: '备份状态', icon: Shield },
] as const

function getColorClasses(color: string) {
  const colors: Record<string, { bg: string; text: string; border: string }> = {
    blue: { bg: 'bg-blue-500/10', text: 'text-blue-600', border: 'border-blue-200/60' },
    emerald: { bg: 'bg-emerald-500/10', text: 'text-emerald-600', border: 'border-emerald-200/60' },
    amber: { bg: 'bg-amber-500/10', text: 'text-amber-600', border: 'border-amber-200/60' },
    rose: { bg: 'bg-rose-500/10', text: 'text-rose-600', border: 'border-rose-200/60' },
  }
  return colors[color] || colors.blue
}

function getPriorityColor(priority: string) {
  switch (priority) {
    case 'high': return 'bg-rose-500 text-white'
    case 'medium': return 'bg-amber-500 text-white'
    default: return 'bg-slate-400 text-white'
  }
}

function getStatusBadge(status: number) {
  if (status === 2) return <Badge variant="default" className="bg-emerald-500 hover:bg-emerald-600">已发布</Badge>
  if (status === 1) return <Badge variant="secondary">草稿</Badge>
  return <Badge variant="outline">审核中</Badge>
}

function CurrentTime() {
  const [time, setTime] = useState(new Date())
  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000)
    return () => clearInterval(timer)
  }, [])
  const timeStr = time.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  const dateStr = time.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', weekday: 'short' })
  return (
    <span className="text-muted-foreground/80 ml-2 text-xs tabular-nums">
      · {dateStr} {timeStr}
    </span>
  )
}

function formatRelativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return '刚刚'
  if (diffMin < 60) return `${diffMin}分钟前`
  const diffHour = Math.floor(diffMin / 60)
  if (diffHour < 24) return `${diffHour}小时前`
  const diffDay = Math.floor(diffHour / 24)
  if (diffDay === 1) return '昨天'
  return `${diffDay}天前`
}

export function DashboardPage() {
  const user = useAuthUser()
  const chartColors = useCssVariables(CHART_CSS_VARS)
  const [mounted, setMounted] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    const timer = setTimeout(() => setMounted(true), 100)
    return () => clearTimeout(timer)
  }, [])

  const { data: dashboardResp, isLoading, error } = useQuery({
    queryKey: ['dashboard'],
    queryFn: () => getDashboard(),
  })

  const data: DashboardData | undefined = dashboardResp?.indexDate

  const queryClient = useQueryClient()
  const clearCacheMutation = useMutation({
    mutationFn: clearCache,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    },
  })

  if (!mounted || isLoading) {
    return <DashboardSkeleton />
  }

  if (error || !data) {
    return (
      <div className="min-h-full p-6 flex items-center justify-center">
        <Card className="p-8 text-center">
          <AlertCircle className="size-8 text-rose-500 mx-auto mb-4" />
          <p className="text-muted-foreground">加载仪表板失败，请刷新重试</p>
          <Button className="mt-4" onClick={() => window.location.reload()}>刷新页面</Button>
        </Card>
      </div>
    )
  }

  const statsCards = [
    { label: '总文章数', value: data.stats.total_posts, trend: 0, icon: FileText, color: 'blue', path: '/content/article' },
    { label: '草稿数', value: data.stats.draft_posts, trend: 0, icon: FileEdit, color: 'amber', path: '/content/article' },
    { label: '待审核评论', value: data.stats.pending_comments, trend: 0, icon: MessageCircle, color: 'rose', path: '/operation/comment' },
    { label: '近7日访问量', value: data.stats.week_visits, trend: 0, icon: Eye, color: 'emerald', path: '/dashboard' },
  ]

  const pendingItems = [
    { label: '待审核文章', count: data.pending.pending_review_posts ?? 0, path: '/content/article', priority: 'high' },
    { label: '待审核评论', count: data.pending.pending_review_comments ?? 0, path: '/operation/comment', priority: 'high' },
  ]

  const systemStatusItems = SYSTEM_STATUS_STATIC.map((item, i) => {
    let value = ''
    if (i === 0) value = '运行中'
    else if (i === 1) value = data.system_status.db_status === 'ok' ? '正常' : '异常'
    else if (i === 2) value = data.system_status.redis_status === 'ok' ? '已启用' : '异常'
    else if (i === 3) value = '手动'
    const ok = (i === 0 || value === '正常' || value === '已启用' || value === '运行中')
    return { ...item, value, status: ok ? ('normal' as const) : ('warning' as const) }
  })

  return (
    <div className="min-h-full p-6 space-y-6">
      <header className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">控制台</h1>
          <p className="text-sm text-muted-foreground">
            欢迎回来，{user?.nickname || user?.username || '管理员'}
            <CurrentTime />
          </p>
        </div>
        <div className="flex items-center gap-3">
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              placeholder="搜索文章、标签..."
              className="pl-9"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          <Link to="/content/article/new">
            <Button className="gap-2">
              <Plus className="size-4" />
              <span className="hidden sm:inline">新建文章</span>
            </Button>
          </Link>
          <Avatar className="size-9 cursor-pointer">
            <AvatarImage src={user?.avatar as string} />
            <AvatarFallback>{(user?.nickname || user?.username || 'A').charAt(0)}</AvatarFallback>
          </Avatar>
        </div>
      </header>

      {/* 统计卡片 */}
      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {statsCards.map(({ label, value, icon: Icon, color, path }) => {
          const colors = getColorClasses(color)
          return (
            <Link key={label} to={path}>
              <Card className="transition-all hover:shadow-md hover:-translate-y-0.5">
                <CardContent className="p-4">
                  <div className="flex items-start justify-between">
                    <div className="space-y-2">
                      <p className="text-xs text-muted-foreground">{label}</p>
                      <p className="text-2xl font-bold">{value.toLocaleString()}</p>
                    </div>
                    <div className={cn('rounded-lg p-2', colors.bg, colors.text)}>
                      <Icon className="size-5" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          )
        })}
      </section>

      {/* 访问趋势 + 快捷操作/待处理 */}
      <section className="grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div>
              <CardTitle className="text-base font-medium">访问趋势</CardTitle>
              <CardDescription>近7日网站访问统计</CardDescription>
            </div>
            <Badge variant="secondary">本周</Badge>
          </CardHeader>
          <CardContent>
            <Suspense fallback={<Skeleton className="h-[280px] w-full" />}>
              <VisitTrendChart chartColors={chartColors} apiData={data.visit_trend} />
            </Suspense>
          </CardContent>
        </Card>

        <div className="space-y-4">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">快捷操作</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2">
              {QUICK_ACTIONS.map(({ label, icon: Icon, path, color }) => (
                <Link key={label} to={path}>
                  <Button variant="secondary" className={cn('w-full justify-start gap-2 text-white', color)}>
                    <Icon className="size-4" />
                    {label}
                  </Button>
                </Link>
              ))}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">待处理提醒</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {pendingItems.map(({ label, count, path, priority }) => (
                <Link
                  key={label}
                  to={path}
                  className="flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-muted/50"
                >
                  <div className="flex items-center gap-3">
                    <span className={cn('flex size-5 items-center justify-center rounded-full text-[10px] font-bold', getPriorityColor(priority))}>
                      {count}
                    </span>
                    <span className="text-sm">{label}</span>
                  </div>
                  <ChevronRight className="size-4 text-muted-foreground" />
                </Link>
              ))}
            </CardContent>
          </Card>
        </div>
      </section>

      {/* 最近文章 + 最近评论/热门文章 */}
      <section className="grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div>
              <CardTitle className="text-base font-medium">最近文章</CardTitle>
            </div>
            <Link to="/content/article">
              <Button variant="ghost" size="sm" className="gap-1">
                查看全部 <ChevronRight className="size-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {data.recent_posts.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-8">暂无文章</p>
              ) : (
                data.recent_posts.map((post) => (
                  <div key={post.id} className="flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-muted/30">
                    <div className="flex items-center gap-3">
                      {getStatusBadge(post.status)}
                      <span className="text-sm font-medium">{post.title}</span>
                    </div>
                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      <span className="flex items-center gap-1"><Eye className="size-3" />{post.view}</span>
                      <span>{formatRelativeTime(post.created_at)}</span>
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>

        <div className="space-y-4">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">最近评论</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {data.recent_comments.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-4">暂无评论</p>
              ) : (
                data.recent_comments.map((comment) => (
                  <div key={comment.id} className="space-y-1 rounded-lg border p-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">{comment.nickname}</span>
                      <span className="text-xs text-muted-foreground">{formatRelativeTime(comment.created_at)}</span>
                    </div>
                    <p className="text-xs text-muted-foreground line-clamp-2">{comment.content}</p>
                    <p className="text-xs text-blue-600 truncate">《{comment.post_title}》</p>
                  </div>
                ))
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">热门文章</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {data.popular_posts.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-4">暂无数据</p>
              ) : (
                data.popular_posts.map((post, index) => (
                  <div key={post.id} className="flex items-center justify-between rounded-lg p-2 transition-colors hover:bg-muted/30">
                    <div className="flex items-center gap-2">
                      <span className={cn(
                        'flex size-5 items-center justify-center rounded text-[10px] font-bold',
                        index === 0 ? 'bg-amber-500 text-white' :
                          index === 1 ? 'bg-slate-400 text-white' :
                            index === 2 ? 'bg-orange-400 text-white' : 'bg-slate-200 text-slate-600',
                      )}>
                        {index + 1}
                      </span>
                      <span className="text-sm truncate max-w-[140px]">{post.title}</span>
                    </div>
                    <span className="text-xs text-muted-foreground">{post.view.toLocaleString()}</span>
                  </div>
                ))
              )}
            </CardContent>
          </Card>
        </div>
      </section>

      {/* 系统状态 / 缓存 / SEO */}
      <section className="grid gap-4 lg:grid-cols-3">
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">系统状态</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {systemStatusItems.map((item) => (
              <div key={item.label} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <item.icon className={cn('size-4', item.status === 'normal' ? 'text-emerald-500' : 'text-amber-500')} />
                  <span className="text-sm">{item.label}</span>
                </div>
                <span className={cn('text-sm font-medium', item.status === 'normal' ? 'text-emerald-600' : 'text-amber-600')}>
                  {item.value}
                </span>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">缓存管理</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">已使用缓存</span>
                <span className="font-medium">{data.cache_status.used_memory || '未知'}</span>
              </div>
              <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                <div
                  className="h-full rounded-full bg-blue-500 transition-all duration-500"
                  style={{ width: `${data.cache_status.used_percent || 0}%` }}
                />
              </div>
            </div>
            <div className="grid grid-cols-1 gap-2">
              <Button
                variant="outline"
                size="sm"
                className="gap-2"
                onClick={() => clearCacheMutation.mutate()}
                disabled={clearCacheMutation.isPending}
              >
                <Zap className="size-4" />
                {clearCacheMutation.isPending ? '清理中...' : '清理缓存'}
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">SEO 提醒</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {data.seo_alerts.map((alert, index) => (
              <div key={index} className={cn(
                'flex items-start gap-2 rounded-lg p-2',
                alert.type === 'warning' ? 'bg-amber-500/10' : alert.type === 'success' ? 'bg-emerald-500/10' : 'bg-blue-500/10',
              )}>
                {alert.type === 'warning' ? <AlertTriangle className="size-4 text-amber-500 mt-0.5" /> :
                  alert.type === 'success' ? <CheckCircle2 className="size-4 text-emerald-500 mt-0.5" /> :
                    <AlertCircle className="size-4 text-blue-500 mt-0.5" />}
                <span className="text-xs">{alert.message}</span>
              </div>
            ))}
          </CardContent>
          <CardFooter className="border-t pt-3">
            <Link to="/system/config" className="w-full">
              <Button variant="outline" size="sm" className="w-full">SEO 设置</Button>
            </Link>
          </CardFooter>
        </Card>
      </section>
    </div>
  )
}

export function DashboardSkeleton() {
  return (
    <div className="min-h-full p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-2">
          <Skeleton className="h-8 w-32" />
          <Skeleton className="h-4 w-48" />
        </div>
        <div className="flex items-center gap-3">
          <Skeleton className="h-10 w-80" />
          <Skeleton className="h-10 w-28" />
          <Skeleton className="h-9 w-9 rounded-full" />
        </div>
      </div>
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (<Skeleton key={i} className="h-24 w-full rounded-lg" />))}
      </div>
      <div className="grid gap-4 lg:grid-cols-3">
        <Skeleton className="h-[360px] w-full rounded-lg lg:col-span-2" />
        <div className="space-y-4">
          <Skeleton className="h-48 w-full rounded-lg" />
          <Skeleton className="h-48 w-full rounded-lg" />
        </div>
      </div>
      <div className="grid gap-4 lg:grid-cols-3">
        <Skeleton className="h-[400px] w-full rounded-lg lg:col-span-2" />
        <div className="space-y-4">
          <Skeleton className="h-48 w-full rounded-lg" />
          <Skeleton className="h-48 w-full rounded-lg" />
        </div>
      </div>
    </div>
  )
}
