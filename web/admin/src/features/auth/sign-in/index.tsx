import { useSearch } from '@tanstack/react-router'
import { Card, CardContent } from '@/components/ui/card'
import { LanguageSwitcher } from '@/components/language-switcher'
import { AuthLayout } from '../auth-layout'
import { UserAuthForm } from './components/user-auth-form'

export function SignIn() {
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })

  return (
    <AuthLayout>
      <div className='mb-5 flex justify-end'>
        <LanguageSwitcher />
      </div>
      <Card className='w-full rounded-lg border border-border/80 bg-card shadow-xl shadow-foreground/5'>
        <CardContent className='px-6 pb-6 pt-6 sm:px-8 sm:pb-8 sm:pt-8'>
          <UserAuthForm redirectTo={redirect} />
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
