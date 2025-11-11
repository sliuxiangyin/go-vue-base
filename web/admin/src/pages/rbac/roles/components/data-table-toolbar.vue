<script setup lang="ts">
import { RefreshCw, X } from 'lucide-vue-next'
import type { Table } from '@tanstack/vue-table'
import type { Role } from '@/services/api/rbac.api'

interface Props {
  table: Table<Role>
}

defineProps<Props>()
const emit = defineEmits<{
  search: [params: { name: string, status?: number }]
  refresh: []
}>()

const searchName = ref('')
const searchStatus = ref<number | string>('-1')

function handleSearch() {
  emit('search', { 
    name: searchName.value, 
    status: searchStatus.value === '-1' ? undefined : Number(searchStatus.value)
  })
}

function handleReset() {
  searchName.value = ''
  searchStatus.value = '-1'
  emit('search', { name: '', status: undefined })
}
</script>

<template>
  <div class="flex items-center justify-between gap-2">
    <div class="flex flex-1 items-center space-x-2">
      <UiInput
        v-model="searchName"
        placeholder="搜索角色名称..."
        class="h-8 w-[150px] lg:w-[250px]"
        @keyup.enter="handleSearch"
      />
      <UiSelect v-model="searchStatus">
        <UiSelectTrigger class="h-8 w-[120px]">
          <UiSelectValue placeholder="状态" />
        </UiSelectTrigger>
        <UiSelectContent>
          <UiSelectItem value="-1">全部</UiSelectItem>
          <UiSelectItem :value="1">启用</UiSelectItem>
          <UiSelectItem :value="0">禁用</UiSelectItem>
        </UiSelectContent>
      </UiSelect>
      <UiButton variant="outline" size="sm" @click="handleSearch">
        搜索
      </UiButton>
      <UiButton 
        v-if="searchName || searchStatus !== '-1'" 
        variant="ghost" 
        size="sm" 
        @click="handleReset"
      >
        <X class="h-4 w-4" />
        重置
      </UiButton>
    </div>
    <div class="flex items-center space-x-2">
      <UiButton variant="outline" size="sm" @click="emit('refresh')">
        <RefreshCw class="h-4 w-4" />
        刷新
      </UiButton>
    </div>
  </div>
</template>
