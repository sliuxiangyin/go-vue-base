<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { useCreateLessonMutation } from '@/services/api/lesson.api'
import { toast } from 'vue-sonner'
import type { LessonForm } from '../data/schema'

const emit = defineEmits<{
  created: []
}>()

const open = ref(false)
const form = ref<LessonForm>({
  title: '',
  content_en: '',
  content_zh: '',
  level: 1,
  is_public: false,
  tags: [],
})

const tagInput = ref('')

const { mutate: createLesson, isPending } = useCreateLessonMutation()

function handleSubmit() {
  if (!form.value.title || !form.value.content_en) {
    toast.error('请填写必填字段')
    return
  }

  createLesson(form.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('创建成功')
        open.value = false
        resetForm()
        emit('created')
      } else {
        toast.error(response.message || '创建失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '创建失败')
    },
  })
}

function resetForm() {
  form.value = {
    title: '',
    content_en: '',
    content_zh: '',
    level: 1,
    is_public: false,
    tags: [],
  }
  tagInput.value = ''
}

function handleOpenChange(newOpen: boolean) {
  open.value = newOpen
  if (!newOpen) {
    resetForm()
  }
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
  <UiDialog :open="open" @update:open="handleOpenChange">
    <UiDialogTrigger as-child>
      <UiButton>
        <Plus class="mr-2 h-4 w-4" />
        创建课程
      </UiButton>
    </UiDialogTrigger>
    <UiDialogContent class="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
      <UiDialogHeader>
        <UiDialogTitle>创建英文课程</UiDialogTitle>
        <UiDialogDescription>
          创建新的英语学习课程
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
        <UiButton variant="outline" @click="open = false">
          取消
        </UiButton>
        <UiButton :disabled="isPending" @click="handleSubmit">
          {{ isPending ? '创建中...' : '创建' }}
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
