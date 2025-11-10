<template>
  <div class="flex flex-col h-full p-2">
    <h2 class="text-lg font-semibold mb-3 p-2 shrink-0">数据库表</h2>

    <!-- 搜索框 -->
    <div class="mb-3 flex items-center">
      <Button 
        @click="toggleSearchMode" 
        variant="outline" 
        size="sm" 
        class="mr-2"
      >
        <Funnel v-if="searchMode === 'filter'" class="h-4 w-4" />
        <Highlighter v-else class="h-4 w-4" />
      </Button>
      <Input
          v-model="searchQuery"
          placeholder="搜索数据表... (支持多标签搜索，用逗号、空格分隔)"
          class="w-full"
      />
    </div>
    <!-- 搜索提示 -->
    <div v-if="searchTags.length > 1" class="mb-2 text-sm text-muted-foreground">
      当前搜索标签: {{ searchTags.join(', ') }}
    </div>
    <!-- 表列表 -->
    <ScrollArea class="flex-1 min-h-0">
      <div class="space-y-2 items-top flex gap-x-2 "
           v-for="table in displayedTables" :key="table.name">
        <Checkbox :id="table.name" :model-value="table.selected"  @update:modelValue="(val) => table.selected! = val"    />
        <div class="grid gap-1.5 leading-none">
          <label
              :for="table.name"
              class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              :class="isHighlighted(table.name) || isHighlighted(table.comment) ? 'bg-yellow-200 dark:bg-yellow-800' : ''"
          >
            {{ table.name }}
          </label>
          <label 
            :for="table.name" 
            class="text-xs text-muted-foreground"
            :class="isHighlighted(table.comment) ? 'bg-yellow-200 dark:bg-yellow-800' : ''"
          >
            {{ table.comment }}
          </label>
        </div>
      </div>
      <!-- 无搜索结果提示 -->
      <div v-if="displayedTables.length === 0 && searchQuery" class="text-center text-muted-foreground py-4">
        未找到匹配的数据表
      </div>
    </ScrollArea>
    <!-- 底部固定工具栏 -->
    <div class="border-t p-2 flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <Checkbox id="select-all" :checked="isAllSelected" @update:modelValue="toggleSelectAll" />
        <label for="select-all" class="text-sm font-medium leading-none">
          {{ isAllSelected ? '取消全选' : '全选' }}
        </label>
      </div>
      <div class="flex items-center space-x-2">
        <!-- 其它扩展功能可以放在这里 -->
        <span class="text-sm text-muted-foreground">{{ selectedCount }} 项已选择</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, computed, watch} from 'vue'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Funnel,Highlighter } from "lucide-vue-next"
interface TableInfo {
  name: string
  comment: string
  selected?: boolean
}

const props = defineProps<{
  tables: TableInfo[]
}>()

const emit = defineEmits<{
  (e: 'selection-change', selectedTables: TableInfo[]): void
}>()

const tables = ref<TableInfo[]>([])

// 搜索模式：'filter' 为筛选模式，'highlight' 为高亮模式
const searchMode = ref<'filter' | 'highlight'>('highlight')

watch(
    () => props.tables,
    (newTables) => {
      tables.value = newTables.map((table) => ({
        ...table,
        selected: table.selected ?? false, // 防止重复初始化
      }))
    },
    { immediate: true, deep: true }
)

const searchQuery = ref('')

// 切换搜索模式
const toggleSearchMode = () => {
  searchMode.value = searchMode.value === 'filter' ? 'highlight' : 'filter'
}

// 解析搜索标签
const searchTags = computed(() => {
  if (!searchQuery.value) return []
  return searchQuery.value
    .split(/[,，\s]+/) // 支持逗号、中文逗号和空格分隔
    .map(tag => tag.trim())
    .filter(tag => tag.length > 0)
})

// 检查文本是否应该高亮
const isHighlighted = (text: string) => {
  if (!searchQuery.value || searchMode.value !== 'highlight') return false
  
  const tags = searchTags.value
  if (tags.length === 0) return false
  const lowerText = text.toLowerCase()
  return tags.some(tag => lowerText.includes(tag.toLowerCase()))
}

// 检查表是否匹配所有标签
const isTableMatch = (table: TableInfo) => {
  const tags = searchTags.value
  if (tags.length === 0) return true
  const name = table.name.toLowerCase()
  const comment = table.comment.toLowerCase()
  
  // 每个标签都必须在表名或注释中找到
  return tags.every(tag => {
    const lowerTag = tag.toLowerCase()
    return name.includes(lowerTag) || (comment && comment.includes(lowerTag))
  })
}

// 计算表的匹配度得分（匹配的标签数量）
const getMatchScore = (table: TableInfo): number => {
  const tags = searchTags.value
  if (tags.length === 0) return 0
  
  const name = table.name.toLowerCase()
  const comment = table.comment.toLowerCase()
  
  let score = 0
  tags.forEach(tag => {
    const lowerTag = tag.toLowerCase()
    if (name.includes(lowerTag)) score++
    if (comment && comment.includes(lowerTag)) score++
  })
  
  return score
}

// 显示的表列表（根据搜索模式不同而变化）
const displayedTables = computed(() => {
  if (!searchQuery.value) return tables.value
  
  if (searchMode.value === 'filter') {
    // 筛选模式：只显示匹配所有标签的表
    return tables.value.filter(isTableMatch)
  } else {
    // 高亮模式：显示所有表，但根据匹配度排序（匹配的表排在前面）
    return [...tables.value].sort((a, b) => {
      const scoreA = getMatchScore(a)
      const scoreB = getMatchScore(b)
      // 匹配度高的排在前面，匹配度相同的保持原有顺序
      return scoreB - scoreA
    })
  }
})

const selectedCount = computed(() => {
  return tables.value.filter((t) => t.selected).length
})

const isAllSelected = computed(() => {
  return tables.value.length > 0 && tables.value.every((t) => t.selected)
})

const toggleSelectAll = (checked: boolean) => {
  tables.value.forEach((t) => (t.selected = checked))
}

// 监听选中数据变化并触发事件
watch(
    () => tables.value.filter(t => t.selected),
    (selectedTables) => {
      emit('selection-change', selectedTables)
    },
    { deep: true }
)
</script>
