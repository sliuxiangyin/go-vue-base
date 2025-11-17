<script setup lang="ts">
import { FileText, Image, Music, Video, File as FileIcon, Download, Trash2, Edit } from 'lucide-vue-next'
import type { File } from '@/services/api/file.api'
import { formatBytes } from '@/lib/utils'

interface Props {
  file: File
}

const props = defineProps<Props>()
const emit = defineEmits<{
  delete: [id: number]
  edit: [file: File]
  download: [id: number]
}>()

// 根据文件类型选择图标
const getFileIcon = (fileType: string) => {
  switch (fileType) {
    case 'image':
      return Image
    case 'audio':
      return Music
    case 'video':
      return Video
    case 'document':
      return FileText
    default:
      return FileIcon
  }
}

// 根据文件类型选择颜色
const getFileColor = (fileType: string) => {
  switch (fileType) {
    case 'image':
      return 'text-green-600 bg-green-100'
    case 'audio':
      return 'text-purple-600 bg-purple-100'
    case 'video':
      return 'text-red-600 bg-red-100'
    case 'document':
      return 'text-blue-600 bg-blue-100'
    default:
      return 'text-gray-600 bg-gray-100'
  }
}

// 格式化日期
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>

<template>
  <UiCard class="overflow-hidden hover:shadow-lg transition-shadow">
    <UiCardHeader class="p-4">
      <div class="flex items-start justify-between">
        <div class="flex items-center gap-3 flex-1 min-w-0">
          <div :class="['p-3 rounded-lg', getFileColor(file.file_type)]">
            <component :is="getFileIcon(file.file_type)" class="h-6 w-6" />
          </div>
          <div class="flex-1 min-w-0">
            <UiCardTitle class="text-base truncate" :title="file.name">
              {{ file.name }}
            </UiCardTitle>
            <p class="text-xs text-muted-foreground mt-1">
              {{ formatBytes(file.size) }} · {{ file.extension }}
            </p>
          </div>
        </div>
      </div>
    </UiCardHeader>

    <UiCardContent class="p-4 pt-0">
      <div class="space-y-2 text-sm">
        <div class="flex items-center justify-between">
          <span class="text-muted-foreground">类型</span>
          <UiBadge variant="secondary">{{ file.file_type }}</UiBadge>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-muted-foreground">来源</span>
          <UiBadge variant="outline">{{ file.source }}</UiBadge>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-muted-foreground">下载次数</span>
          <span>{{ file.download_count }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-muted-foreground">上传时间</span>
          <span>{{ formatDate(file.created_at) }}</span>
        </div>
        <div v-if="file.description" class="pt-2 border-t">
          <p class="text-xs text-muted-foreground line-clamp-2" :title="file.description">
            {{ file.description }}
          </p>
        </div>
      </div>
    </UiCardContent>

    <UiCardFooter class="p-4 pt-0 flex gap-2">
      <UiButton
        variant="outline"
        size="sm"
        class="flex-1"
        @click="emit('download', file.id)"
      >
        <Download class="h-4 w-4 mr-1" />
        下载
      </UiButton>
      <UiButton
        variant="outline"
        size="sm"
        @click="emit('edit', file)"
      >
        <Edit class="h-4 w-4" />
      </UiButton>
      <UiButton
        variant="outline"
        size="sm"
        @click="emit('delete', file.id)"
      >
        <Trash2 class="h-4 w-4" />
      </UiButton>
    </UiCardFooter>
  </UiCard>
</template>
