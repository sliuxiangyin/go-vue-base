<script setup lang="ts">
import type { DataTableProps } from '@/components/data-table/types'
import DataTable from '@/components/data-table/data-table.vue'
import { generateVueTable } from '@/components/data-table/use-generate-vue-table'
import type { AdminUser } from '@/services/api/rbac.api'
import DataTableToolbar from './data-table-toolbar.vue'

interface Props extends DataTableProps<AdminUser> {
  total: number
  page: number
  pageSize: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  search: [keyword: string, status?: number]
  pageChange: [page: number]
  refresh: []
}>()

const table = generateVueTable<AdminUser>(props)
</script>

<template>
  <DataTable :columns :data :loading :table>
    <template #toolbar>
      <DataTableToolbar 
        :table 
        class="w-full overflow-x-auto"
        @search="emit('search', $event.keyword, $event.status)"
        @refresh="emit('refresh')"
      />
    </template>
  </DataTable>
</template>
