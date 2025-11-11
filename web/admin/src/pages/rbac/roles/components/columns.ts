import type { ColumnDef } from '@tanstack/vue-table'
import { h } from 'vue'
import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { SelectColumn } from '@/components/data-table/table-columns'
import Badge from '@/components/ui/badge/Badge.vue'
import type { Role } from '@/services/api/rbac.api'
import DataTableRowActions from './data-table-row-actions.vue'

export const columns: ColumnDef<Role>[] = [
  {
    ...SelectColumn,
    meta: {
      sticky: 'left',
    },
  } as ColumnDef<Role>,
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '角色标识' }),
    cell: ({ row }) => h('div', { class: 'font-medium' }, row.getValue('name')),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'display_name',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '角色名称' }),
    cell: ({ row }) => h('div', {}, row.getValue('display_name')),
    enableSorting: false,
  },
  {
    accessorKey: 'description',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '描述' }),
    cell: ({ row }) => h('div', { class: 'max-w-[300px] truncate text-muted-foreground' }, row.getValue('description') || '-'),
    enableSorting: false,
  },
  {
    accessorKey: 'is_system',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '系统角色' }),
    cell: ({ row }) => {
      const isSystem = row.getValue('is_system') as boolean
      return h(Badge, { 
        variant: isSystem ? 'secondary' : 'outline',
        class: isSystem ? 'bg-blue-100 text-blue-700' : ''
      }, () => isSystem ? '系统角色' : '普通角色')
    },
    enableSorting: false,
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '状态' }),
    cell: ({ row }) => {
      const status = row.getValue('status') as number
      return h(Badge, { 
        variant: status === 1 ? 'default' : 'secondary',
        class: status === 1 ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
      }, () => status === 1 ? '启用' : '禁用')
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'created_at',
    header: ({ column }) => h(DataTableColumnHeader<Role>, { column, title: '创建时间' }),
    cell: ({ row }) => {
      const date = new Date(row.getValue('created_at'))
      return h('div', { class: 'text-sm text-muted-foreground' }, date.toLocaleString('zh-CN'))
    },
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => h(DataTableRowActions, { row }),
    meta: {
      sticky: 'right',
    },
  },
]
