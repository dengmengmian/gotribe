import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { CheckCircle2, ShieldOff, Loader2, KeyRound } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ContentSection } from '../components/content-section'
import {
  bindTOTP,
  getTOTPStatus,
  type TOTPBindResult,
  type TOTPStatus,
} from '@/features/auth/totp/service'
import { TOTPBindDialog } from './totp-bind-dialog'
import { TOTPUnbindDialog } from './totp-unbind-dialog'

export function SettingsSecurity() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [bindResult, setBindResult] = useState<TOTPBindResult | null>(null)
  const [bindDialogOpen, setBindDialogOpen] = useState(false)
  const [unbindDialogOpen, setUnbindDialogOpen] = useState(false)

  const { data: status, isLoading } = useQuery({
    queryKey: ['admin', 'totp', 'status'],
    queryFn: () => getTOTPStatus(),
  })

  const bindMutation = useMutation({
    mutationFn: bindTOTP,
    onSuccess: (result) => {
      setBindResult(result)
      setBindDialogOpen(true)
    },
    onError: (error: Error) => {
      toast.error(error.message || t('features.settings.security.bindFailed'))
    },
  })

  function refreshStatus() {
    queryClient.invalidateQueries({ queryKey: ['admin', 'totp', 'status'] })
  }

  function handleStartBind() {
    bindMutation.mutate()
  }

  function handleBindClosed() {
    setBindDialogOpen(false)
    setBindResult(null)
  }

  function handleBindConfirmed() {
    handleBindClosed()
    refreshStatus()
  }

  function handleUnbound() {
    setUnbindDialogOpen(false)
    refreshStatus()
  }

  const isBound = Boolean(status?.bound && status?.enabled)

  return (
    <ContentSection
      title={t('features.settings.security.title')}
      desc={t('features.settings.security.desc')}
    >
      <div className='space-y-5'>
        {isLoading ? (
          <div className='flex items-center justify-center py-8'>
            <Loader2 className='size-5 animate-spin text-muted-foreground' />
          </div>
        ) : (
          <TOTPStatusPanel
            status={status ?? { bound: false, enabled: false, remaining_recovery_codes: 0 }}
            isBound={isBound}
            bindLoading={bindMutation.isPending}
            onStartBind={handleStartBind}
            onStartUnbind={() => setUnbindDialogOpen(true)}
          />
        )}
        <TOTPBindDialog
          open={bindDialogOpen}
          bindResult={bindResult}
          onConfirmed={handleBindConfirmed}
          onClose={handleBindClosed}
        />
        <TOTPUnbindDialog
          open={unbindDialogOpen}
          onUnbound={handleUnbound}
          onClose={() => setUnbindDialogOpen(false)}
        />
      </div>
    </ContentSection>
  )
}

type StatusPanelProps = {
  status: TOTPStatus
  isBound: boolean
  bindLoading: boolean
  onStartBind: () => void
  onStartUnbind: () => void
}

function TOTPStatusPanel({ status, isBound, bindLoading, onStartBind, onStartUnbind }: StatusPanelProps) {
  const { t } = useTranslation()
  return (
    <div className='rounded-lg border border-border/70 p-5'>
      <div className='flex items-start justify-between gap-4'>
        <div className='flex items-start gap-3'>
          <div className='rounded-md bg-muted p-2 text-muted-foreground'>
            <KeyRound className='size-5' />
          </div>
          <div>
            <div className='flex items-center gap-2'>
              <h4 className='text-base font-semibold'>{t('features.settings.security.totpTitle')}</h4>
              {isBound ? (
                <Badge className='bg-emerald-500 hover:bg-emerald-600'>
                  <CheckCircle2 className='mr-1 size-3' />
                  {t('features.settings.security.statusEnabled')}
                </Badge>
              ) : (
                <Badge variant='outline' className='text-muted-foreground'>
                  <ShieldOff className='mr-1 size-3' />
                  {t('features.settings.security.statusDisabled')}
                </Badge>
              )}
            </div>
            <p className='mt-1 text-sm text-muted-foreground'>
              {isBound
                ? t('features.settings.security.totpEnabledDesc', {
                    count: status.remaining_recovery_codes,
                  })
                : t('features.settings.security.totpDisabledDesc')}
            </p>
          </div>
        </div>
        <div>
          {isBound ? (
            <Button variant='outline' onClick={onStartUnbind}>
              {t('features.settings.security.unbindButton')}
            </Button>
          ) : (
            <Button onClick={onStartBind} disabled={bindLoading}>
              {bindLoading ? <Loader2 className='mr-2 size-4 animate-spin' /> : null}
              {t('features.settings.security.bindButton')}
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
