import { cn } from '@/lib/utils'

export function Logo({ className, ...props }: { className?: string, props?: React.HTMLAttributes<HTMLImageElement> }) {
  return (
    <img className={cn('size-6', className)} src="/images/gotribe.png" alt="GoTribe Logo" {...props} />
  )
}
