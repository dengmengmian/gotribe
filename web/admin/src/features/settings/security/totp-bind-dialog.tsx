import { useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import { toast } from 'sonner'
import { AlertTriangle, Copy, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from '@/components/ui/input-otp'
import { confirmTOTP, type TOTPBindResult } from '@/features/auth/totp/service'

type TOTPBindDialogProps = {
  open: boolean
  bindResult: TOTPBindResult | null
  onConfirmed: () => void
  onClose: () => void
}

export function TOTPBindDialog({ open, bindResult, onConfirmed, onClose }: TOTPBindDialogProps) {
  const { t } = useTranslation()
  const [code, setCode] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  async function handleConfirm() {
    if (code.length !== 6 || !bindResult) return
    setIsLoading(true)
    try {
      await confirmTOTP(code)
      toast.success(t('features.settings.security.bindSuccess'))
      setCode('')
      onConfirmed()
    } catch (error) {
      const message = error instanceof Error ? error.message : t('features.settings.security.bindFailed')
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  async function copy(value: string) {
    try {
      await navigator.clipboard.writeText(value)
      toast.success(t('features.settings.security.copied'))
    } catch {
      toast.error(t('features.settings.security.copyFailed'))
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) onClose()
      }}
    >
      <DialogContent className='max-w-md'>
        <DialogHeader>
          <DialogTitle>{t('features.settings.security.bindTitle')}</DialogTitle>
          <DialogDescription>{t('features.settings.security.bindDescription')}</DialogDescription>
        </DialogHeader>

        {bindResult ? (
          <div className='space-y-5'>
            <div className='flex flex-col items-center gap-3'>
              <div className='rounded-lg border border-border bg-white p-3 shadow-sm'>
                <QRCodeSVG value={bindResult.otpauth_url} size={180} level='M' />
              </div>
              <div className='w-full space-y-1'>
                <Label className='text-xs text-muted-foreground'>
                  {t('features.settings.security.manualSecret')}
                </Label>
                <div className='flex items-center gap-2'>
                  <Input value={bindResult.secret} readOnly className='font-mono text-xs' />
                  <Button
                    type='button'
                    size='icon'
                    variant='outline'
                    onClick={() => copy(bindResult.secret)}
                    aria-label='copy'
                  >
                    <Copy className='size-4' />
                  </Button>
                </div>
              </div>
            </div>

            <div className='space-y-2 rounded-lg border border-amber-300/70 bg-amber-50 p-3 text-xs text-amber-900 dark:border-amber-500/40 dark:bg-amber-500/10 dark:text-amber-200'>
              <div className='flex items-center gap-2 font-medium'>
                <AlertTriangle className='size-4' />
                {t('features.settings.security.recoveryWarningTitle')}
              </div>
              <p>{t('features.settings.security.recoveryWarningBody')}</p>
              <div className='grid grid-cols-2 gap-1 pt-1 font-mono text-xs'>
                {bindResult.recovery_codes.map((c) => (
                  <code key={c} className='rounded bg-amber-100 px-1.5 py-0.5 dark:bg-amber-500/20'>
                    {c}
                  </code>
                ))}
              </div>
              <Button
                type='button'
                variant='outline'
                size='sm'
                className='mt-1 w-full'
                onClick={() => copy(bindResult.recovery_codes.join('\n'))}
              >
                <Copy className='mr-2 size-3' />
                {t('features.settings.security.copyAll')}
              </Button>
            </div>

            <div className='space-y-2'>
              <Label className='text-xs text-muted-foreground'>
                {t('features.settings.security.confirmCodeLabel')}
              </Label>
              <InputOTP
                maxLength={6}
                value={code}
                onChange={setCode}
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
            </div>
          </div>
        ) : (
          <div className='flex items-center justify-center py-10'>
            <Loader2 className='size-5 animate-spin text-muted-foreground' />
          </div>
        )}

        <DialogFooter>
          <Button type='button' variant='outline' onClick={onClose} disabled={isLoading}>
            {t('features.settings.security.cancel')}
          </Button>
          <Button
            type='button'
            onClick={handleConfirm}
            disabled={code.length !== 6 || isLoading || !bindResult}
          >
            {isLoading ? <Loader2 className='mr-2 size-4 animate-spin' /> : null}
            {t('features.settings.security.confirmBind')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
