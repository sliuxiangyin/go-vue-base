<script setup lang="ts">
import type { Table } from '@tanstack/vue-table'
import { RefreshCw } from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import type { Lesson } from '../data/schema'

interface Props {
  table: Table<Lesson>
}

defineProps<Props>()

const emit = defineEmits<{
  search: [params: { level?: number, isPublic?: boolean }]
  refresh: []
}>()

const selectedLevel = ref<number | undefined>(undefined)
const selectedIsPublic = ref<boolean | undefined>(undefined)

function handleLevelChange(value: any) {
  if (!value) return
  const strValue = String(value)
  selectedLevel.value = strValue === 'all' ? undefined : parseInt(strValue)
  emitSearch()
}

function handlePublicChange(value: any) {
  if (!value) return
  const strValue = String(value)
  if (strValue === 'all') {
    selectedIsPublic.value = undefined
  } else {
    selectedIsPublic.value = strValue === 'true'
  }
  emitSearch()
}

function emitSearch() {
  emit('search', {
    level: selectedLevel.value,
    isPublic: selectedIsPublic.value,
  })
}

function handleRefresh() {
  emit('refresh')
}
</script>

<template>
  <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
    <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
      <!-- 难度等级筛选 -->
      <Select @update:model-value="handleLevelChange">
        <SelectTrigger class="w-full sm:w-[180px]">
          <SelectValue placeholder="选择难度等级" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem value="all">全部等级</SelectItem>
            <SelectItem value="1">Level 1</SelectItem>
            <SelectItem value="2">Level 2</SelectItem>
            <SelectItem value="3">Level 3</SelectItem>
            <SelectItem value="4">Level 4</SelectItem>
            <SelectItem value="5">Level 5</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>

      <!-- 公开状态筛选 -->
      <Select @update:model-value="handlePublicChange">
        <SelectTrigger class="w-full sm:w-[180px]">
          <SelectValue placeholder="选择公开状态" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem value="all">全部状态</SelectItem>
            <SelectItem value="true">公开</SelectItem>
            <SelectItem value="false">私密</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>

    <!-- 刷新按钮 -->
    <Button
      variant="outline"
      size="sm"
      @click="handleRefresh"
    >
      <RefreshCw class="mr-2 size-4" />
      刷新
    </Button>
  </div>
</template>
