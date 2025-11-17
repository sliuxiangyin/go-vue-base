<script setup lang="ts">
import type { File } from '@/services/api/file.api'
import { useUpdateFileMutation } from '@/services/api/file.api'
import { toast } from 'vue-sonner'

interface Props {
  file: File | null
  open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  updated: []
}>()

const description = ref('')

watch(() => props.file, (newFile) => {
  if (newFile) {
    description.value = newFile.description || ''
  }
})

const { mutate: updateFile, isPending } = useUpdateFileMutation(props.file?.id || 0)

function handleSubmit() {
  if (!props.file) return

  updateFile({ description: description.value }, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('更新成功')
        emit('update:open', false)
        emit('updated')
      } else {
        toast.error(response.message || '更新失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '更新失败')
    },
  })
}
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="sm:max-w-[500px]">
      <UiDialogHeader>
        <UiDialogTitle>编辑文件信息</UiDialogTitle>
        <UiDialogDescription>
          更新文件描述信息
        </UiDialogDescription>
      </UiDialogHeader>
      <div v-if="file" class="grid gap-4 py-4">
        <div class="grid gap-2">
          <UiLabel>文件名</UiLabel>
          <UiInput :model-value="file.name" disabled />
        </div>
        <div class="grid gap-2">
          <UiLabel>文件路径</UiLabel>
          <UiInput :model-value="file.path" disabled class="font-mono text-xs" />
        </div>
        <div class="grid gap-2">
          <UiLabel for="description">描述</UiLabel>
          <UiTextarea
            id="description"
            v-model="description"
            placeholder="添加文件描述..."
            rows="4"
          />
        </div>
      </div>
      <UiDialogFooter>
        <UiButton variant="outline" @click="emit('update:open', false)">
          取消
        </UiButton>
        <UiButton :disabled="isPending" @click="handleSubmit">
          {{ isPending ? '保存中...' : '保存' }}
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
