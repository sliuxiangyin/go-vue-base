import type { AxiosError } from 'axios'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { toValue, type MaybeRefOrGetter } from 'vue'
import { useAxios } from '@/composables/use-axios'
import type { IResponse } from '../types/response.type'

const { axiosInstance } = useAxios()

// ========== 类型定义 ==========

export type FileType = 'image' | 'audio' | 'video' | 'document' | 'other' | 'undefined'
export type FileSource = 'upload' | 'download' | 'generate' | 'undefined'

export interface File {
  id: number
  created_at: string
  updated_at: string
  name: string
  path: string
  size: number
  mime_type: string
  extension: string
  file_type: FileType
  source: FileSource
  hash: string
  description?: string
  uploaded_by?: number
  reference_by?: string
  reference_id?: number
  download_url?: string
  download_count: number
  storage_type: string
}

export interface ListResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface FileQueryParams {
  page: number
  page_size: number
  file_type?: FileType
  source?: FileSource
  keyword?: string
}

export interface StorageStats {
  total_files: number
  total_size: number
  by_type: Array<{
    file_type: FileType
    count: number
    size: number
  }>
}

// ========== API 函数 ==========

export function useGetFilesQuery(params: MaybeRefOrGetter<FileQueryParams>) {
  return useQuery<IResponse<ListResponse<File>>, AxiosError>({
    queryKey: ['files', params],
    queryFn: async () => {
      const resolvedParams = toValue(params)
      const response = await axiosInstance.get('/admin/files', { params: resolvedParams })
      return response.data
    },
  })
}

export function useGetFileQuery(id: number) {
  return useQuery<IResponse<File>, AxiosError>({
    queryKey: ['file', id],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/files/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

export function useUpdateFileMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, { description: string }>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/files/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['files'] })
      queryClient.invalidateQueries({ queryKey: ['file', id] })
    },
  })
}

export function useDeleteFileMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, { id: number; delete_physical?: boolean }>({
    mutationFn: async ({ id, delete_physical = true }) => {
      const response = await axiosInstance.delete(`/admin/files/${id}`, {
        params: { delete_physical },
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['files'] })
    },
  })
}

export function useBatchDeleteFilesMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, { ids: number[]; delete_physical?: boolean }>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post('/admin/files/batch-delete', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['files'] })
    },
  })
}

export function useGetStorageStatsQuery() {
  return useQuery<IResponse<StorageStats>, AxiosError>({
    queryKey: ['storage-stats'],
    queryFn: async () => {
      const response = await axiosInstance.get('/admin/files/stats')
      return response.data
    },
  })
}

export function useDownloadFile() {
  return useMutation<IResponse<{ file_path: string; file_name: string }>, AxiosError, number>({
    mutationFn: async (id) => {
      const response = await axiosInstance.get(`/admin/files/${id}/download`)
      return response.data
    },
  })
}
