<script setup lang="ts" generic="T">
import type {
  Table as VueTable,
} from '@tanstack/vue-table'

import {
  FlexRender,
} from '@tanstack/vue-table'

import DataTableLoading from '@/components/data-table/table-loading.vue'
import DataTablePagination from '@/components/data-table/table-pagination.vue'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import type { DataTableProps } from './types'

defineProps<DataTableProps<T> & {
  table: VueTable<T>
}>()
</script>

<template>
  <div class="space-y-4">
    <slot name="toolbar" />

    <div class="border rounded-md overflow-auto">
      <Table>
        <TableHeader>
          <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
            <TableHead 
              v-for="header in headerGroup.headers" 
              :key="header.id"
              :class="[
                header.column.columnDef.meta?.sticky === 'left' && 'sticky left-0 z-10 bg-background',
                header.column.columnDef.meta?.sticky === 'right' && 'sticky right-0 z-10 bg-background',
                (header.column.columnDef.meta?.sticky === 'left' || header.column.columnDef.meta?.sticky === 'right') && 'shadow-[0_0_0_1px_hsl(var(--border))]',
              ]"
            >
              <FlexRender v-if="!header.isPlaceholder" :render="header.column.columnDef.header" :props="header.getContext()" />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody v-if="!loading">
          <template v-if="table.getRowModel().rows?.length">
            <TableRow
              v-for="row in table.getRowModel().rows"
              :key="row.id"
              :data-state="row.getIsSelected() && 'selected'"
            >
              <TableCell 
                v-for="cell in row.getVisibleCells()" 
                :key="cell.id"
                :class="[
                  cell.column.columnDef.meta?.sticky === 'left' && 'sticky left-0 z-10 bg-background',
                  cell.column.columnDef.meta?.sticky === 'right' && 'sticky right-0 z-10 bg-background',
                  (cell.column.columnDef.meta?.sticky === 'left' || cell.column.columnDef.meta?.sticky === 'right') && 'shadow-[0_0_0_1px_hsl(var(--border))]',
                ]"
              >
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>

          <TableRow v-else>
            <TableCell
              :colspan="columns.length"
              class="h-24 text-center"
            >
              No results.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <DataTableLoading v-if="loading" />
    </div>

    <DataTablePagination v-if="!loading" :table="table" :server-pagination="serverPagination" />
  </div>
</template>
