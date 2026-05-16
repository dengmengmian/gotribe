import { pinyin } from 'pinyin-pro'

/**
 * 从标题生成 URL 友好的 slug
 * - 英文/ASCII 标题：转小写，空格替换为连字符
 * - 中文标题：转拼音后用连字符连接
 */
export function generateSlug(title: string): string {
  const trimmed = title.trim()
  if (!trimmed) return ''

  const source = /[\u4e00-\u9fff]/.test(trimmed)
    ? pinyin(trimmed, {
        toneType: 'none',
        type: 'array',
        nonZh: 'consecutive',
      }).join('-')
    : trimmed

  return source
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

/** 校验 slug 格式是否合法 */
export function isValidSlug(slug: string): boolean {
  if (!slug) return true // 空 slug 视为合法（可选字段）
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug)
}
