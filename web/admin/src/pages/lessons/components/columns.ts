import type { ColumnDef } from '@tanstack/vue-table'
import { h } from 'vue'
import DataTableColumnHeader from '@/components/data-table/column-header.vue'
import { SelectColumn } from '@/components/data-table/table-columns'
import Badge from '@/components/ui/badge/Badge.vue'
import type { Lesson } from '../data/schema'
import DataTableRowActions from './data-table-row-actions.vue'

export const columns: ColumnDef<Lesson>[] = [
  {
    ...SelectColumn,
    meta: {
      sticky: 'left',
    },
  } as ColumnDef<Lesson>,
  {
    accessorKey: 'title',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '课程标题' }),
    cell: ({ row }) => h('div', { class: 'font-medium max-w-[300px] truncate' }, row.getValue('title')),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'level',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '难度等级' }),
    cell: ({ row }) => {
      const level = row.getValue('level') as number
      const colors = ['', 'bg-green-100 text-green-700', 'bg-blue-100 text-blue-700', 'bg-yellow-100 text-yellow-700', 'bg-orange-100 text-orange-700', 'bg-red-100 text-red-700']
      return h(Badge, { 
        variant: 'secondary',
        class: colors[level] || ''
      }, () => `Level ${level}`)
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'tags',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '标签' }),
    cell: ({ row }) => {
      const tags = row.getValue('tags') as string[] || []
      if (tags.length === 0) return h('div', { class: 'text-muted-foreground' }, '-')
      return h('div', { class: 'flex flex-wrap gap-1' }, 
        tags.slice(0, 3).map((tag: string) => 
          h(Badge, { variant: 'outline', class: 'text-xs' }, () => tag)
        )
      )
    },
    enableSorting: false,
  },
  {
    accessorKey: 'duration',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '时长' }),
    cell: ({ row }) => {
      const duration = row.getValue('duration') as number
      if (!duration) return h('div', { class: 'text-muted-foreground' }, '-')
      return h('div', { class: 'text-sm' }, `${duration.toFixed(1)}s`)
    },
    enableSorting: false,
  },
  {
    accessorKey: 'is_public',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '公开状态' }),
    cell: ({ row }) => {
      const isPublic = row.getValue('is_public') as boolean
      return h(Badge, { 
        variant: isPublic ? 'default' : 'secondary',
        class: isPublic ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
      }, () => isPublic ? '公开' : '私密')
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: 'created_at',
    header: ({ column }) => h(DataTableColumnHeader<Lesson>, { column, title: '创建时间' }),
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
