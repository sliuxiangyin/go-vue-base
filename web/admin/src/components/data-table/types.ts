import type { ColumnDef } from '@tanstack/vue-table'

// 扩展 TanStack Table 的 ColumnMeta 类型
declare module '@tanstack/vue-table' {
  interface ColumnMeta<TData, TValue> {
    sticky?: 'left' | 'right'
  }
}

export interface FacetedFilterOption {
  label: string
  value: string
  icon?: Component
}

export interface ServerPagination {
  page: number
  pageSize: number
  total: number
  onPageChange: (page: number) => void
  onPageSizeChange: (pageSize: number) => void
}

export interface DataTableProps<T> {
  loading?: boolean
  columns: ColumnDef<T, any>[]
  data: T[]
  serverPagination?: ServerPagination
}
