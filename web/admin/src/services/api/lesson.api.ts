import type { AxiosError } from 'axios'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { toValue, type MaybeRefOrGetter } from 'vue'
import { useAxios } from '@/composables/use-axios'
import type { IResponse } from '../types/response.type'
import type { Lesson, LessonForm } from '@/pages/lessons/data/schema'

const { axiosInstance } = useAxios()

export interface ListResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface LessonQueryParams {
  page: number
  page_size: number
  level?: number
  is_public?: boolean
}

// ========== 课程管理 API ==========

export function useGetLessonsQuery(params: MaybeRefOrGetter<LessonQueryParams>) {
  return useQuery<IResponse<ListResponse<Lesson>>, AxiosError>({
    queryKey: ['lessons', params],
    queryFn: async () => {
      const resolvedParams = toValue(params)
      const response = await axiosInstance.get('/admin/lessons', { params: resolvedParams })
      return response.data
    },
  })
}

export function useGetLessonQuery(id: number) {
  return useQuery<IResponse<Lesson>, AxiosError>({
    queryKey: ['lesson', id],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/lessons/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreateLessonMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Lesson>, AxiosError, LessonForm>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post('/admin/lessons', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['lessons'] })
    },
  })
}

export function useUpdateLessonMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Lesson>, AxiosError, LessonForm>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/lessons/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['lessons'] })
      queryClient.invalidateQueries({ queryKey: ['lesson', id] })
    },
  })
}

export function useDeleteLessonMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, number>({
    mutationFn: async (id) => {
      const response = await axiosInstance.delete(`/admin/lessons/${id}`)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['lessons'] })
    },
  })
}

export function useSearchLessonsByTagsQuery(tags: MaybeRefOrGetter<string[]>) {
  return useQuery<IResponse<{ list: Lesson[] }>, AxiosError>({
    queryKey: ['lessons-by-tags', tags],
    queryFn: async () => {
      const resolvedTags = toValue(tags)
      const response = await axiosInstance.get('/admin/lessons/search/tags', {
        params: { tags: resolvedTags.join(',') },
      })
      return response.data
    },
    enabled: computed(() => toValue(tags).length > 0),
  })
}
