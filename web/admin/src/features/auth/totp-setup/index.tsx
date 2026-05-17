import { Link } from '@tanstack/react-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useI18n } from '@/context/i18n-provider'
import { AuthLayout } from '../auth-layout'
import { TotpSetupForm } from './components/totp-setup-form'

export function TotpSetup() {
  const { t } = useI18n()

  return (
    <AuthLayout>
      <Card className='w-full border shadow-sm'>
        <CardHeader className='pb-4 pt-6'>
          <div className='space-y-1.5'>
            <CardTitle className='text-base font-medium'>
              {t('features.auth.totpSetup.title')}
            </CardTitle>
            <CardDescription>
              {t('features.auth.totpSetup.description')}
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <TotpSetupForm />
        </CardContent>
        <CardFooter className='border-t pt-6'>
          <p className='w-full text-center text-sm text-muted-foreground'>
            {t('features.auth.totpSetup.footerText')}{' '}
            <Link
              to='/sign-in'
              className='underline underline-offset-4 hover:text-primary'
            >
              {t('features.auth.totpSetup.backLink')}
            </Link>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
