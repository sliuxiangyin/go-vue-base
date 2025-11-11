<script setup lang="ts">
import type { Lesson } from '../data/schema'

interface Props {
  lesson: Lesson
  open: boolean
}

defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
      <UiDialogHeader>
        <UiDialogTitle>课程详情</UiDialogTitle>
      </UiDialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid gap-2">
          <div class="font-semibold">课程标题</div>
          <div class="text-sm">{{ lesson.title }}</div>
        </div>

        <div v-if="lesson.audio_url" class="grid gap-2">
          <div class="font-semibold">音频文件</div>
          <div class="text-sm text-muted-foreground">{{ lesson.audio_url }}</div>
        </div>

        <div v-if="lesson.duration" class="grid gap-2">
          <div class="font-semibold">时长</div>
          <div class="text-sm">{{ lesson.duration.toFixed(1) }}秒</div>
        </div>

        <div class="grid gap-2">
          <div class="font-semibold">英文内容</div>
          <div class="text-sm whitespace-pre-wrap border rounded-md p-3 bg-muted/50">{{ lesson.content_en }}</div>
        </div>

        <div v-if="lesson.content_zh" class="grid gap-2">
          <div class="font-semibold">中文翻译</div>
          <div class="text-sm whitespace-pre-wrap border rounded-md p-3 bg-muted/50">{{ lesson.content_zh }}</div>
        </div>

        <div class="grid gap-2">
          <div class="font-semibold">难度等级</div>
          <UiBadge variant="secondary" class="w-fit">Level {{ lesson.level }}</UiBadge>
        </div>

        <div v-if="lesson.tags && lesson.tags.length > 0" class="grid gap-2">
          <div class="font-semibold">标签</div>
          <div class="flex flex-wrap gap-2">
            <UiBadge
              v-for="(tag, index) in lesson.tags"
              :key="index"
              variant="outline"
            >
              {{ tag }}
            </UiBadge>
          </div>
        </div>

        <div class="grid gap-2">
          <div class="font-semibold">公开状态</div>
          <UiBadge :variant="lesson.is_public ? 'default' : 'secondary'" class="w-fit">
            {{ lesson.is_public ? '公开' : '私密' }}
          </UiBadge>
        </div>

        <div v-if="lesson.semantic_json && lesson.semantic_json.length > 0" class="grid gap-2">
          <div class="font-semibold">语义意群</div>
          <div class="space-y-1 text-sm border rounded-md p-3 bg-muted/50">
            <div v-for="(chunk, index) in lesson.semantic_json" :key="index" class="flex gap-2">
              <span class="text-muted-foreground">[{{ chunk.start.toFixed(2)}}s - {{ chunk.end.toFixed(2) }}s]</span>
              <span>{{ chunk.text }}</span>
            </div>
          </div>
        </div>

        <div class="grid gap-2">
          <div class="font-semibold">创建时间</div>
          <div class="text-sm text-muted-foreground">{{ new Date(lesson.created_at).toLocaleString('zh-CN') }}</div>
        </div>
      </div>
      <UiDialogFooter>
        <UiButton @click="emit('update:open', false)">
          关闭
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
