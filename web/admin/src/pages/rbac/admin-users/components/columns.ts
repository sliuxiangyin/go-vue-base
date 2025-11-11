import type { ColumnDef } from '@tanstack/vue-table'
import type { AdminUser } from '@/services/api/rbac.api'
import { h } from 'vue'
import DataTableRowActions from './data-table-row-actions.vue'
import Badge from '@/components/ui/badge/Badge.vue'
export const columns: ColumnDef<AdminUser>[] = [
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }) => h('div', { class: 'w-20' }, row.getValue('id')),
    meta: {
      sticky: 'left',
    },
  },
  {
    accessorKey: 'username',
    header: '用户名',
    cell: ({ row }) => h('div', { class: 'font-medium' }, row.getValue('username')),
  },
  {
    accessorKey: 'nickname',
    header: '昵称',
    cell: ({ row }) => {
      const nickname = row.getValue('nickname') as string
      return h('div', nickname || '-')
    },
  },
  {
    accessorKey: 'email',
    header: '邮箱',
    cell: ({ row }) => {
      const email = row.getValue('email') as string
      return h('div', email || '-')
    },
  },
  {
    accessorKey: 'roles',
    header: '角色',
    cell: ({ row }) => {
      const roles = row.getValue('roles') as AdminUser['roles']
      if (!roles || roles.length === 0) {
        return h('div', { class: 'text-muted-foreground' }, '未分配')
      }
      return h('div', { class: 'flex flex-wrap gap-1' }, roles.map(role =>
        h(Badge, { variant: 'secondary', key: role.id }, () => role.display_name)
      ))
    },
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => {
      const status = row.getValue('status') as number
      return h(Badge, {
        variant: status === 1 ? 'default' : 'destructive',
      }, () => status === 1 ? '启用' : '禁用')
    },
  },
  {
    accessorKey: 'last_login',
    header: '最后登录',
    cell: ({ row }) => {
      const lastLogin = row.getValue('last_login') as string | null
      if (!lastLogin) {
        return h('div', { class: 'text-muted-foreground' }, '从未登录')
      }
      return h('div', { class: 'text-xs' }, new Date(lastLogin).toLocaleString('zh-CN'))
    },
  },
  {
    accessorKey: 'created_at',
    header: '创建时间',
    cell: ({ row }) => {
      const createdAt = row.getValue('created_at') as string
      return h('div', { class: 'text-xs' }, new Date(createdAt).toLocaleString('zh-CN'))
    },
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => h(DataTableRowActions, { row: row.original }),
    meta: {
      sticky: 'right',
    },
  },
]
