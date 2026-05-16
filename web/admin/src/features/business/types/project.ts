export type Project = {
  id: number
  title: string // 项目名称
  description: string // 项目描述
  name: string // 项目别名
  meta_title?: string // Meta标题
  meta_description?: string // Meta描述
  keywords: string // Meta关键词
  domain: string // 项目域名
  post_url: string // 内容链接
  icp: string // icp备案号
  author: string // 项目归属者
  baidu_analytics: string // 第三方js
  favicon: string // 网站图标
  public_security: string // 公安备案号
  nav_image: string // Nav图标
  info?: string // 额外信息
  push_token?: string
  created_at?: string
  updated_at?: string
}

export type ProjectListResponse = {
  projects: Project[]
  total: number
}
