import { Suspense, lazy } from 'react'
import type { UseFormReturn } from 'react-hook-form'
import { useI18n } from '@/context/i18n-provider'
import {
  FormField,
  FormItem,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
import { EditorErrorBoundary } from '@/components/editor-error-boundary'
import type { ArticleFormValues } from '../article-form-page'

const SlateEditor = lazy(() =>
  import('@/components/editor').then((m) => ({ default: m.SlateEditor }))
)

interface ArticleEditorProps {
  form: UseFormReturn<ArticleFormValues>
  isEdit?: boolean
}

export function ArticleEditor({ form, isEdit = false }: ArticleEditorProps) {
  const { t } = useI18n()

  return (
    <div className='mx-auto w-full max-w-[52rem] px-5 pt-6 pb-24 sm:px-8'>
      {/* 标题输入 */}
      <div className='mx-auto w-full max-w-[52rem]'>
        <FormField
          control={form.control}
          name='title'
          render={({ field }) => (
            <FormItem className='space-y-0'>
              <FormControl>
                <Textarea
                  placeholder={t(
                    'features.content.article.form.titlePlaceholder'
                  )}
                  className='min-h-[2.5rem] resize-none overflow-hidden border-none px-0 py-0 text-3xl leading-snug font-bold tracking-tight shadow-none placeholder:text-muted-foreground/50 focus-visible:ring-0 sm:text-4xl'
                  rows={1}
                  onInput={(e) => {
                    const target = e.target as HTMLTextAreaElement
                    target.style.height = 'auto'
                    target.style.height = `${target.scrollHeight}px`
                  }}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </div>

      {/* Slate 编辑器 */}
      <FormField
        control={form.control}
        name='content'
        render={({ field }) => (
          <FormItem className='flex min-h-0 flex-col'>
            <FormControl>
              <EditorErrorBoundary
                fallback={
                  <div className='flex min-h-[320px] flex-1 flex-col items-center justify-center rounded-md border border-input bg-background text-muted-foreground'>
                    <p className='mb-2 text-sm'>
                      {t('features.content.article.form.editorLoadError')}
                    </p>
                    <Textarea
                      className='min-h-[120px] w-full max-w-2xl resize-none font-mono text-xs'
                      value={field.value ?? ''}
                      onChange={field.onChange}
                      placeholder={t(
                        'features.content.article.form.contentPlaceholder'
                      )}
                      readOnly={false}
                    />
                  </div>
                }
              >
                <Suspense
                  fallback={
                    <div className='flex min-h-[320px] flex-col items-center justify-center text-muted-foreground'>
                      <p className='text-sm'>
                        {t('features.content.article.form.editorLoading')}
                      </p>
                    </div>
                  }
                >
                  <SlateEditor
                    value={field.value ?? ''}
                    onChange={field.onChange}
                    minHeight='min-h-[52vh]'
                    outputMode='json'
                    autoHeight={true}
                    autoFocus={!isEdit}
                  />
                </Suspense>
              </EditorErrorBoundary>
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}
