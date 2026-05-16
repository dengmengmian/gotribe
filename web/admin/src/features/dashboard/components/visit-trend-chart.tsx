import { useMemo } from 'react'
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Skeleton } from '@/components/ui/skeleton'
import type { VisitPoint } from '../service'

interface VisitTrendChartProps {
  chartColors: string[]
  apiData?: VisitPoint[]
}

export default function VisitTrendChart({ chartColors, apiData }: VisitTrendChartProps) {
  const data = useMemo(() => {
    if (apiData && apiData.length > 0) {
      return apiData.map((p) => ({
        day: p.date,
        visits: p.visits,
        pageViews: p.page_views,
      }))
    }
    return []
  }, [apiData])

  const tooltipContent = useMemo(
    () => (tooltipProps: { payload?: ReadonlyArray<{ payload?: { day: string; visits: number; pageViews: number } }> }) => {
      const p = tooltipProps?.payload?.[0]?.payload
      if (!p) return null
      return (
        <div className="rounded-md border border-border bg-card px-3 py-2 text-sm shadow-sm">
          <div className="font-medium">{p.day}</div>
          <div className="space-y-1 text-muted-foreground">
            <div>访客: <span className="font-semibold text-foreground">{p.visits}</span></div>
            <div>浏览: <span className="font-semibold text-foreground">{p.pageViews}</span></div>
          </div>
        </div>
      )
    },
    [],
  )

  if (chartColors.some((c) => !c) || data.length === 0) {
    return <Skeleton className="h-[280px] w-full" />
  }

  return (
    <div className="h-[280px]">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="colorVisits" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor={chartColors[0]} stopOpacity={0.3} />
              <stop offset="50%" stopColor={chartColors[0]} stopOpacity={0.1} />
              <stop offset="100%" stopColor={chartColors[0]} stopOpacity={0.05} />
            </linearGradient>
            <linearGradient id="colorPageViews" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor={chartColors[1]} stopOpacity={0.3} />
              <stop offset="50%" stopColor={chartColors[1]} stopOpacity={0.1} />
              <stop offset="100%" stopColor={chartColors[1]} stopOpacity={0.05} />
            </linearGradient>
          </defs>
          <XAxis dataKey="day" tick={{ fontSize: 12, fill: 'hsl(var(--muted-foreground))' }} tickLine={false} axisLine={false} />
          <YAxis tick={{ fontSize: 12, fill: 'hsl(var(--muted-foreground))' }} tickLine={false} axisLine={false} width={40} />
          <Tooltip content={tooltipContent} />
          <Area type="monotone" dataKey="visits" name="访客数" stroke={chartColors[0]} strokeWidth={2} fill="url(#colorVisits)" activeDot={{ r: 4, fill: chartColors[0], stroke: 'hsl(var(--background))', strokeWidth: 2 }} />
          <Area type="monotone" dataKey="pageViews" name="浏览量" stroke={chartColors[1]} strokeWidth={2} fill="url(#colorPageViews)" activeDot={{ r: 4, fill: chartColors[1], stroke: 'hsl(var(--background))', strokeWidth: 2 }} />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  )
}
