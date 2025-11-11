import { z } from 'zod'

// 语义意群
export const semanticChunkSchema = z.object({
  text: z.string(),
  start: z.number(),
  end: z.number(),
})
export type SemanticChunk = z.infer<typeof semanticChunkSchema>

// 逐词时间戳
export const wordTimestampSchema = z.object({
  word: z.string(),
  start: z.number(),
  end: z.number(),
  confidence: z.number().optional(),
})
export type WordTimestamp = z.infer<typeof wordTimestampSchema>

// 发音信息
export const phoneticSchema = z.object({
  word: z.string(),
  ipa: z.string(),
})
export type Phonetic = z.infer<typeof phoneticSchema>

// 英文课程
export const lessonSchema = z.object({
  id: z.number(),
  title: z.string(),
  audio_url: z.string().optional(),
  duration: z.number().optional(),
  content_en: z.string(),
  content_zh: z.string().optional(),
  semantic_json: z.array(semanticChunkSchema).optional().nullable(),
  word_timestamp_json: z.array(wordTimestampSchema).optional().nullable(),
  tags: z.array(z.string()).optional().nullable(),
  level: z.number().min(1).max(5),
  is_public: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
})
export type Lesson = z.infer<typeof lessonSchema>

export const lessonListSchema = z.array(lessonSchema)

// 创建/更新课程的表单数据（不包含 audio_url，后端生成）
export const lessonFormSchema = z.object({
  title: z.string().min(1, '标题不能为空').max(255, '标题最多255个字符'),
  content_en: z.string().min(1, '英文内容不能为空'),
  content_zh: z.string().optional(),
  semantic_json: z.array(semanticChunkSchema).optional(),
  word_timestamp_json: z.array(wordTimestampSchema).optional(),
  tags: z.array(z.string()).optional(),
  level: z.number().min(1, '难度等级最小为1').max(5, '难度等级最大为5'),
  is_public: z.boolean().default(false),
})
export type LessonForm = z.infer<typeof lessonFormSchema>
