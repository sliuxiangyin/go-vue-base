<script setup lang="ts">
import { useUpdateLessonMutation } from '@/services/api/lesson.api'
import { toast } from 'vue-sonner'
import type { Lesson, LessonForm } from '../data/schema'

interface Props {
  lesson: Lesson
  open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const form = ref<LessonForm>({
  title: props.lesson.title,
  content_en: props.lesson.content_en,
  content_zh: props.lesson.content_zh || '',
  level: props.lesson.level,
  is_public: props.lesson.is_public,
  tags: props.lesson.tags || [],
  semantic_json: props.lesson.semantic_json || undefined,
  word_timestamp_json: props.lesson.word_timestamp_json || undefined,
  phonetic_json: props.lesson.phonetic_json || undefined,
})

const tagInput = ref('')

const { mutate: updateLesson, isPending } = useUpdateLessonMutation(props.lesson.id)

function handleSubmit() {
  if (!form.value.title || !form.value.content_en) {
    toast.error('请填写必填字段')
    return
  }

  updateLesson(form.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('更新成功')
        emit('update:open', false)
      } else {
        toast.error(response.message || '更新失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '更新失败')
    },
  })
}

function addTag() {
  const tag = tagInput.value.trim()
  if (tag && !form.value.tags?.includes(tag)) {
    if (!form.value.tags) {
      form.value.tags = []
    }
    form.value.tags.push(tag)
    tagInput.value = ''
  }
}

function removeTag(index: number) {
  form.value.tags?.splice(index, 1)
}
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
      <UiDialogHeader>
        <UiDialogTitle>编辑课程</UiDialogTitle>
        <UiDialogDescription>
          编辑课程信息（注意：AudioURL 由后端自动生成）
        </UiDialogDescription>
      </UiDialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid gap-2">
          <UiLabel for="title">课程标题 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="title"
            v-model="form.title"
            placeholder="例如: Forrest Gump - Opening Narration"
          />
        </div>
        
        <div class="grid gap-2">
          <UiLabel for="content_en">英文内容 <span class="text-red-500">*</span></UiLabel>
          <UiTextarea
            id="content_en"
            v-model="form.content_en"
            placeholder="英文原文内容"
            rows="4"
          />
        </div>

        <div class="grid gap-2">
          <UiLabel for="content_zh">中文翻译</UiLabel>
          <UiTextarea
            id="content_zh"
            v-model="form.content_zh"
            placeholder="中文翻译（可选）"
            rows="3"
          />
        </div>

        <div class="grid gap-2">
          <UiLabel for="level">难度等级 <span class="text-red-500">*</span></UiLabel>
          <UiSelect v-model="form.level">
            <UiSelectTrigger>
              <UiSelectValue />
            </UiSelectTrigger>
            <UiSelectContent>
              <UiSelectItem :value="1">Level 1</UiSelectItem>
              <UiSelectItem :value="2">Level 2</UiSelectItem>
              <UiSelectItem :value="3">Level 3</UiSelectItem>
              <UiSelectItem :value="4">Level 4</UiSelectItem>
              <UiSelectItem :value="5">Level 5</UiSelectItem>
            </UiSelectContent>
          </UiSelect>
        </div>

        <div class="grid gap-2">
          <UiLabel>标签</UiLabel>
          <div class="flex gap-2">
            <UiInput
              v-model="tagInput"
              placeholder="输入标签后按回车"
              @keyup.enter="addTag"
            />
            <UiButton type="button" @click="addTag">添加</UiButton>
          </div>
          <div v-if="form.tags && form.tags.length > 0" class="flex flex-wrap gap-2 mt-2">
            <UiBadge
              v-for="(tag, index) in form.tags"
              :key="index"
              variant="secondary"
              class="cursor-pointer"
              @click="removeTag(index)"
            >
              {{ tag }} ×
            </UiBadge>
          </div>
        </div>

        <div class="flex items-center space-x-2">
          <UiSwitch
            id="is_public"
            v-model:checked="form.is_public"
          />
          <UiLabel for="is_public" class="cursor-pointer">公开课程</UiLabel>
        </div>
      </div>
      <UiDialogFooter>
        <UiButton variant="outline" @click="emit('update:open', false)">
          取消
        </UiButton>
        <UiButton :disabled="isPending" @click="handleSubmit">
          {{ isPending ? '更新中...' : '更新' }}
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
