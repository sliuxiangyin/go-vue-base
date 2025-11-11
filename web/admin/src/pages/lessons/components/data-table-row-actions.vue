<script setup lang="ts">
import { MoreHorizontal, Pencil, Trash2, Eye } from 'lucide-vue-next'
import type { Row } from '@tanstack/vue-table'
import type { Lesson } from '../data/schema'
import { useDeleteLessonMutation } from '@/services/api/lesson.api'
import { toast } from 'vue-sonner'
import LessonEdit from './lesson-edit.vue'
import LessonView from './lesson-view.vue'

interface Props {
  row: Row<Lesson>
}

const props = defineProps<Props>()

const showEditDialog = ref(false)
const showViewDialog = ref(false)

const { mutate: deleteLesson } = useDeleteLessonMutation()

function handleDelete() {
  if (confirm('确定要删除这个课程吗？此操作不可恢复。')) {
    deleteLesson(props.row.original.id, {
      onSuccess: (response) => {
        if (response.code === 0) {
          toast.success('删除成功')
        } else {
          toast.error(response.message || '删除失败')
        }
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.message || '删除失败')
      },
    })
  }
}
</script>

<template>
  <div class="flex items-center gap-2">
    <UiButton
      variant="ghost"
      size="sm"
      @click="showViewDialog = true"
    >
      <Eye class="h-4 w-4" />
      查看
    </UiButton>

    <UiDropdownMenu>
      <UiDropdownMenuTrigger as-child>
        <UiButton
          variant="ghost"
          class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
        >
          <MoreHorizontal class="h-4 w-4" />
          <span class="sr-only">打开菜单</span>
        </UiButton>
      </UiDropdownMenuTrigger>
      <UiDropdownMenuContent align="end" class="w-[160px]">
        <UiDropdownMenuItem @click="showEditDialog = true">
          <Pencil class="mr-2 h-4 w-4" />
          编辑
        </UiDropdownMenuItem>
        <UiDropdownMenuSeparator />
        <UiDropdownMenuItem
          class="text-red-600"
          @click="handleDelete"
        >
          <Trash2 class="mr-2 h-4 w-4" />
          删除
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>

    <LessonEdit
      v-if="showEditDialog"
      :lesson="row.original"
      :open="showEditDialog"
      @update:open="showEditDialog = $event"
    />

    <LessonView
      v-if="showViewDialog"
      :lesson="row.original"
      :open="showViewDialog"
      @update:open="showViewDialog = $event"
    />
  </div>
</template>
