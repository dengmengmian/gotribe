import { useEffect, useMemo, useRef, useState, lazy, Suspense } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { EditorErrorBoundary } from '@/components/editor-error-boundary'
import { useI18n } from '@/context/i18n-provider'
import type { Config, ConfigCreateParams, ConfigUpdateParams } from '../types/config'

const SlateEditor = lazy(() =>
  import('@/components/editor').then((m) => ({ default: m.SlateEditor }))
)

const JsonEditor = lazy(() =>
  import('@/components/json-editor').then((m) => ({ default: m.JsonEditor }))
)

const createConfigFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.content.config.form.validation.titleRequired')),
    description: z.string().min(1, t('features.content.config.form.validation.descriptionRequired')),
    project_id: z.number().min(1, t('features.content.config.form.validation.projectRequired')),
    alias: z.string().min(1, t('features.content.config.form.validation.aliasRequired')),
    type: z.union([z.literal(1), z.literal(2)]),
    md_content: z.string().optional(),
  })

type ConfigFormValues = z.infer<ReturnType<typeof createConfigFormSchema>>

type ConfigFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** 创建时提交 */
  onSubmit?: (data: ConfigCreateParams) => void
  /** 编辑时提交，ID 由传入的 config 带出 */
  onSubmitUpdate?: (id: number, data: ConfigUpdateParams) => void
  isLoading?: boolean
  projectList: { id: number; title: string }[]
  /** 编辑时的配置详情（由父组件拉取后传入），有值即为编辑模式 */
  editConfig?: Config | null
}

export function ConfigFormDialog({
  open,
  onOpenChange,
  onSubmit,
  onSubmitUpdate,
  isLoading = false,
  projectList,
  editConfig,
}: ConfigFormDialogProps) {
  const { t } = useI18n()
  const isEdit = Boolean(editConfig)
  const [editValuesApplied, setEditValuesApplied] = useState(false)
  const [jsonEditorMounted, setJsonEditorMounted] = useState(false)
  const configFormSchema = useMemo(() => createConfigFormSchema(t), [t])
  const form = useForm<ConfigFormValues>({
    resolver: zodResolver(configFormSchema),
    defaultValues: {
      title: '',
      description: '',
      project_id: undefined,
      alias: '',
      type: 1,
      md_content: '',
    },
  })

  const typeValue = form.watch('type')

  useEffect(() => {
    if (!open) {
      setEditValuesApplied(false)
      setJsonEditorMounted(false)
      return
    }
    if (editConfig) {
      form.reset({
        title: editConfig.title ?? '',
        description: editConfig.description ?? '',
        project_id: editConfig.project_id ?? undefined,
        alias: (editConfig.alias ?? '').trim(),
        type: (editConfig.type === 2 ? 2 : 1) as 1 | 2,
        md_content: editConfig.info ?? editConfig.md_content ?? '',
      })
      setEditValuesApplied(true)
      if (editConfig.type === 2) {
        const id = requestAnimationFrame(() => setJsonEditorMounted(true))
        return () => cancelAnimationFrame(id)
      }
      setJsonEditorMounted(true)
    } else {
      form.reset({
        title: '',
        description: '',
        project_id: undefined,
        alias: '',
        type: 1,
        md_content: '',
      })
      setEditValuesApplied(true)
      setJsonEditorMounted(true)
    }
  }, [open, editConfig, form])

  const prevTypeRef = useRef<1 | 2>(1)
  useEffect(() => {
    if (!open) return
    if (prevTypeRef.current === 1 && typeValue === 2) {
      form.setValue('md_content', '')
    }
    prevTypeRef.current = typeValue
  }, [open, typeValue, form])

  const handleSubmit = (values: ConfigFormValues) => {
    const mdContent = values.md_content ?? ''
    if (isEdit && editConfig && onSubmitUpdate) {
      onSubmitUpdate(editConfig.id, {
        title: values.title,
        description: values.description,
        project_id: values.project_id || undefined,
        info: mdContent,
        md_content: mdContent,
      })
    } else if (onSubmit) {
      onSubmit({
        title: values.title,
        description: values.description,
        project_id: values.project_id,
        alias: values.alias,
        type: values.type,
        md_content: mdContent,
        info: mdContent,
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{isEdit ? t('features.content.config.form.editTitle') : t('features.content.config.form.createTitle')}</DialogTitle>
          <DialogDescription>
            {isEdit ? t('features.content.config.form.editDescription') : t('features.content.config.form.createDescription')}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className='flex-1 overflow-y-auto pr-2 space-y-4 min-h-0'
          >
            <FormField
              control={form.control}
              name='title'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.title')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.titlePlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='description'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.description')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.descriptionPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='project_id'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.project')}</FormLabel>
                  <Select onValueChange={(v) => field.onChange(Number(v))} value={field.value != null ? String(field.value) : ''}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('features.content.config.form.projectPlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {projectList.map((p) => (
                        <SelectItem key={p.id} value={String(p.id)}>
                          {p.title}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='alias'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.alias')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.aliasPlaceholder')} {...field} disabled={isEdit} className={isEdit ? 'opacity-60' : ''} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='type'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.type')}</FormLabel>
                  <Select
                    onValueChange={(v) => field.onChange(Number(v) as 1 | 2)}
                    value={String(field.value)}
                    disabled={isEdit}
                  >
                    <FormControl>
                      <SelectTrigger className={isEdit ? 'opacity-60' : ''}>
                        <SelectValue placeholder={t('features.content.config.form.typePlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value='1'>{t('features.content.config.form.typeRichText')}</SelectItem>
                      <SelectItem value='2'>{t('features.content.config.form.typeJson')}</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='md_content'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.content')}</FormLabel>
                  <FormControl>
                    <div className='min-h-[200px]'>
                      {isEdit && !editValuesApplied ? (
                        <div className='min-h-[200px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                          {t('features.content.config.form.contentLoading')}
                        </div>
                      ) : typeValue === 2 && isEdit && !jsonEditorMounted ? (
                        <div className='min-h-[200px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                          {t('features.content.config.form.contentLoading')}
                        </div>
                      ) : typeValue === 1 ? (
                        <EditorErrorBoundary
                          fallback={
                            <Textarea
                              className='min-h-[200px] font-mono text-sm resize-none border-0 rounded-md'
                              value={field.value ?? ''}
                              onChange={(e) => field.onChange(e.target.value)}
                              placeholder={t('features.content.config.form.contentPlaceholder')}
                            />
                          }
                        >
                          <Suspense
                            fallback={
                              <div className='min-h-[200px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                                {t('features.content.config.form.contentLoading')}
                              </div>
                            }
                          >
                            <SlateEditor
                              value={field.value ?? ''}
                              onChange={field.onChange}
                              minHeight='min-h-[200px]'
                              outputMode='json'
                              autoHeight={true}
                            />
                          </Suspense>
                        </EditorErrorBoundary>
                      ) : (
                        <Suspense
                          fallback={
                            <div className='min-h-[200px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                              {t('features.content.config.form.contentLoading')}
                            </div>
                          }
                        >
                          <JsonEditor
                            key={isEdit && editConfig ? `edit-${editConfig.id}` : 'create'}
                            value={field.value ?? ''}
                            onChange={field.onChange}
                            minHeight='min-h-[200px]'
                            initialContent={isEdit && editConfig ? (editConfig.info ?? editConfig.md_content ?? '') : undefined}
                          />
                        </Suspense>
                      )}
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter className='shrink-0 pt-4 border-t mt-4'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                {t('features.content.config.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading ? t('features.content.config.form.submitting') : t('features.content.config.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
