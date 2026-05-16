export type Menu = {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  name: string
  title: string
  icon: string
  path: string
  redirect: string | null
  component: string | null
  sort: number
  status: number // 1: 正常/启用, 2: 禁用
  hidden: number // 1: 隐藏, 2: 显示
  no_cache: number // 1: 不缓存, 2: 缓存
  always_show: number // 1: 总是显示, 0: 不总是显示
  breadcrumb: number // 1: 显示面包屑, 0: 不显示
  active_menu: string | null
  parent_id: number
  creator: string
  children?: Menu[]
  roles?: unknown[] | null
}

export type MenuTreeResponse = {
  menu_tree: Menu[]
}

export type MenuListResponse = {
  menus: Menu[]
  total: number
}

export type CreateMenuParams = Partial<Menu>

export type UpdateMenuParams = Partial<Menu> & {
  id: number
}

export type BatchDeleteMenuParams = {
  menu_ids: number[]
}

export type UserMenuResponse = {
  menus: Menu[]
}

export type UserMenuTreeResponse = {
  menu_tree: Menu[]
}
