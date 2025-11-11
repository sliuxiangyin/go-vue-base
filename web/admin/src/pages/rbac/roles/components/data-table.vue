<script setup lang="ts">
import type { DataTableProps } from '@/components/data-table/types'
import DataTable from '@/components/data-table/data-table.vue'
import { generateVueTable } from '@/components/data-table/use-generate-vue-table'
import type { Role } from '@/services/api/rbac.api'
import DataTableToolbar from './data-table-toolbar.vue'

interface Props extends DataTableProps<Role> {
  total: number
  page: number
  pageSize: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  search: [name: string, status?: number]
  pageChange: [page: number]
  refresh: []
}>()

const table = generateVueTable<Role>(props)
</script>

<template>
  <DataTable :columns :data :loading :table>
    <template #toolbar>
      <DataTableToolbar 
        :table 
        class="w-full overflow-x-auto"
        @search="emit('search', $event.name, $event.status)"
        @refresh="emit('refresh')"
      />
    </template>
  </DataTable>
</template>
