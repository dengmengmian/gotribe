import { request } from '@/service'

export interface Role {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  name: string
  keyword: string
  desc: string
  status: number
  sort: number
  creator: string
  admins: unknown[]
  menus: unknown[]
}

export interface RoleListData {
  roles: Role[]
  total: number
}

export type MenuItem = {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  name: string;
  title: string;
  icon: string;
  path: string;
  redirect: string | null;
  component: string | null;
  sort: number;
  status: number;
  hidden: number;
  no_cache: number;
  always_show: number;
  breadcrumb: number;
  active_menu: string | null;
  parent_id: number;
  creator: string;
  children: MenuItem[];
  roles: unknown[] | null;
};

export type MenuList = {
  menu_tree?: MenuItem[];
};

export async function getRoleList(): Promise<RoleListData> {
  return request.get<RoleListData>('/api/role/list')
}

export async function getMenuAccessTree(parentId: number): Promise<MenuList> {
  return request.get<MenuList>(`/api/menu/access/tree/${parentId}`)
}