<script setup lang="ts">
import type { DataTableProps } from '@/components/data-table/types'
import DataTable from '@/components/data-table/data-table.vue'
import { generateVueTable } from '@/components/data-table/use-generate-vue-table'
import type { Lesson } from '../data/schema'
import DataTableToolbar from './data-table-toolbar.vue'

interface Props extends DataTableProps<Lesson> {
  total: number
  page: number
  pageSize: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  search: [level?: number, isPublic?: boolean]
  pageChange: [page: number]
  refresh: []
}>()

const table = generateVueTable<Lesson>(props)
</script>

<template>
  <DataTable :columns :data :loading :table>
    <template #toolbar>
      <DataTableToolbar 
        :table 
        class="w-full overflow-x-auto"
        @search="emit('search', $event.level, $event.isPublic)"
        @refresh="emit('refresh')"
      />
    </template>
  </DataTable>
</template>
