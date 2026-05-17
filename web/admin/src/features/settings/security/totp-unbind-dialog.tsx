import { useState } from 'react'
import { toast } from 'sonner'
import { Loader2 } from 'lucide-react'
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
import { Label } from '@/components/ui/label'
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from '@/components/ui/input-otp'
import { deleteTOTP } from '@/features/auth/totp/service'

type TOTPUnbindDialogProps = {
  open: boolean
  onUnbound: () => void
  onClose: () => void
}

export function TOTPUnbindDialog({ open, onUnbound, onClose }: TOTPUnbindDialogProps) {
  const { t } = useTranslation()
  const [code, setCode] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  async function handleUnbind() {
    if (code.length !== 6) return
    setIsLoading(true)
    try {
      await deleteTOTP(code)
      toast.success(t('features.settings.security.unbindSuccess'))
      setCode('')
      onUnbound()
    } catch (error) {
      const message = error instanceof Error ? error.message : t('features.settings.security.unbindFailed')
      toast.error(message)
    } finally {
      setIsLoading(false)
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
          <DialogTitle>{t('features.settings.security.unbindTitle')}</DialogTitle>
          <DialogDescription>{t('features.settings.security.unbindDescription')}</DialogDescription>
        </DialogHeader>

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

        <DialogFooter>
          <Button type='button' variant='outline' onClick={onClose} disabled={isLoading}>
            {t('features.settings.security.cancel')}
          </Button>
          <Button
            type='button'
            variant='destructive'
            onClick={handleUnbind}
            disabled={code.length !== 6 || isLoading}
          >
            {isLoading ? <Loader2 className='mr-2 size-4 animate-spin' /> : null}
            {t('features.settings.security.confirmUnbind')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
