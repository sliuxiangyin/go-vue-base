<script setup lang="ts">
import type { DataTableProps } from '@/components/data-table/types'
import DataTable from '@/components/data-table/data-table.vue'
import { generateVueTable } from '@/components/data-table/use-generate-vue-table'
import type { Permission } from '@/services/api/rbac.api'
import DataTableToolbar from './data-table-toolbar.vue'
import { getExpandedRowModel } from '@tanstack/vue-table'

interface Props extends DataTableProps<Permission> {
  total: number
  page: number
  pageSize: number
  permissionType?: 'backend' | 'frontend'
}

const props = defineProps<Props>()
const emit = defineEmits<{
  search: [category: string, status?: number]
  pageChange: [page: number]
  refresh: []
}>()

// 扩展配置：仅在前端菜单权限时启用树形展示
const table = computed(() => {
  const extraConfig = props.permissionType === 'frontend'
    ? {
        getSubRows: (row: Permission) => row.children || [],
        getExpandedRowModel: getExpandedRowModel(),
      }
    : {}
  
  return generateVueTable<Permission>(props, extraConfig)
})
</script>

<template>
  <DataTable :columns :data :loading :table>
    <template #toolbar>
      <DataTableToolbar 
        :table 
        class="w-full overflow-x-auto"
        @search="emit('search', $event.category, $event.status)"
        @refresh="emit('refresh')"
      />
    </template>
  </DataTable>
</template>
