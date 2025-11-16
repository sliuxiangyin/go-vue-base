import type { ColumnDef } from '@tanstack/vue-table'
import { h } from 'vue'
import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { SelectColumn } from '@/components/data-table/table-columns'
import Badge from '@/components/ui/badge/Badge.vue'
import UiButton from '@/components/ui/button/Button.vue'
import type { Permission } from '@/services/api/rbac.api'
import DataTableRowActions from './data-table-row-actions.vue'
import { ChevronDown, ChevronRight } from 'lucide-vue-next'

export const columns: ColumnDef<Permission>[] = [
  SelectColumn as ColumnDef<Permission>,
  {
    accessorKey: 'type',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '类型' }),
    cell: ({ row }) => {
      const type = row.getValue('type') as string
      return h(Badge, { 
        default: type,
        variant: type === 'backend' ? 'default' : 'secondary',
        class: type === 'backend' ? 'bg-blue-100 text-blue-700' : 'bg-purple-100 text-purple-700'
      }, () => type === 'backend' ? '后端API' : '前端菜单')
    },
    enableSorting: false,
  },
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '权限标识' }),
    cell: ({ row }) => {
      const hasChildren = row.original.children && row.original.children.length > 0
      const canExpand = row.getCanExpand()
      
      return h('div', { class: 'flex items-center gap-2' }, [
        // 展开/折叠按钮（仅有子项时显示）
        canExpand && hasChildren
          ? h(UiButton, {
              variant: 'ghost',
              size: 'sm',
              class: 'h-6 w-6 p-0',
              onClick: () => row.toggleExpanded()
            }, () => h(row.getIsExpanded() ? ChevronDown : ChevronRight, { class: 'h-4 w-4' }))
          : h('div', { class: 'w-6' }),
        // 权限标识
        h('div', { 
          class: 'font-mono text-sm',
          style: { paddingLeft: `${row.depth * 1.5}rem` }
        }, row.getValue('name')),
      ])
    },
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'display_name',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '权限名称' }),
    cell: ({ row }) => h('div', {}, row.getValue('display_name')),
    enableSorting: false,
  },
  {
    accessorKey: 'category',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '分类' }),
    cell: ({ row }) => {
      const category = row.getValue('category') as string
      return h(Badge, { variant: 'secondary' }, () => category || '未分类')
    },
    enableSorting: false,
  },
  {
    accessorKey: 'path',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '路由路径' }),
    cell: ({ row }) => {
      const path = row.original.path
      const type = row.original.type
      if (type === 'frontend' && path) {
        return h('div', { class: 'font-mono text-xs text-blue-600' }, path)
      }
      return h('div', { class: 'text-muted-foreground' }, '-')
    },
    enableSorting: false,
  },
  {
    accessorKey: 'icon',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '图标' }),
    cell: ({ row }) => {
      const icon = row.original.icon
      const type = row.original.type
      if (type === 'frontend' && icon) {
        return h('div', { class: 'flex items-center gap-1' }, [
          h('span', { class: 'text-xs bg-gray-100 px-2 py-1 rounded' }, icon)
        ])
      }
      return h('div', { class: 'text-muted-foreground' }, '-')
    },
    enableSorting: false,
  },
  {
    accessorKey: 'description',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '描述' }),
    cell: ({ row }) => h('div', { class: 'max-w-[300px] truncate text-muted-foreground' }, row.getValue('description') || '-'),
    enableSorting: false,
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '状态' }),
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
    header: ({ column }) => h(DataTableColumnHeader<Permission>, { column, title: '创建时间' }),
    cell: ({ row }) => {
      const date = new Date(row.getValue('created_at'))
      return h('div', { class: 'text-sm text-muted-foreground' }, date.toLocaleString('zh-CN'))
    },
    enableSorting: false,
  },
  {
    id: 'actions',
    cell: ({ row }) => h(DataTableRowActions, { row }),
  },
]
