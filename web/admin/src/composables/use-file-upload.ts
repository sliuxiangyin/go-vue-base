import { ref } from 'vue'
import { useAxios } from './use-axios'

const { axiosInstance } = useAxios()

export interface UploadProgress {
  loaded: number
  total: number
  percentage: number
}

export interface ChunkUploadResult {
  success: boolean
  message?: string
  file_path?: string
}

/**
 * 文件分段上传 Composable
 * @param chunkSize 每个分片的大小（字节），默认 2MB
 */
export function useFileUpload(chunkSize = 2 * 1024 * 1024) {
  const uploading = ref(false)
  const progress = ref<UploadProgress>({
    loaded: 0,
    total: 0,
    percentage: 0,
  })

  /**
   * 上传单个文件（支持分段上传）
   * @param file 要上传的文件
   * @param onProgress 进度回调
   * @returns 文件相对路径
   */
  async function uploadFile(
    file: File,
    onProgress?: (progress: UploadProgress) => void
  ): Promise<string> {
    if (!file) {
      throw new Error('请选择要上传的文件')
    }

    uploading.value = true
    progress.value = {
      loaded: 0,
      total: file.size,
      percentage: 0,
    }

    try {
      // 如果文件小于分片大小，直接上传
      if (file.size <= chunkSize) {
        return await uploadSmallFile(file, onProgress)
      }

      // 大文件分片上传
      return await uploadLargeFile(file, onProgress)
    } finally {
      uploading.value = false
    }
  }

  /**
   * 小文件直接上传
   */
  async function uploadSmallFile(
    file: File,
    onProgress?: (progress: UploadProgress) => void
  ): Promise<string> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await axiosInstance.post<{
      code: number
      message: string
      data: { file_path: string }
    }>('/admin/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const loaded = progressEvent.loaded
          const total = progressEvent.total
          const percentage = Math.round((loaded / total) * 100)

          progress.value = { loaded, total, percentage }
          onProgress?.({ loaded, total, percentage })
        }
      },
    })

    if (response.data.code !== 0) {
      throw new Error(response.data.message || '上传失败')
    }

    return response.data.data.file_path
  }

  /**
   * 大文件分片上传
   */
  async function uploadLargeFile(
    file: File,
    onProgress?: (progress: UploadProgress) => void
  ): Promise<string> {
    const totalChunks = Math.ceil(file.size / chunkSize)
    const fileId = `${Date.now()}_${file.name}`
    let uploadedSize = 0

    // 上传每个分片
    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
      const start = chunkIndex * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      const chunk = file.slice(start, end)

      const formData = new FormData()
      formData.append('file', chunk, file.name)
      formData.append('chunk_index', String(chunkIndex))
      formData.append('total_chunks', String(totalChunks))
      formData.append('file_id', fileId)
      formData.append('original_name', file.name)

      const response = await axiosInstance.post<{
        code: number
        message: string
        data: { file_path?: string; chunk_uploaded?: boolean }
      }>('/admin/files/upload/chunk', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      if (response.data.code !== 0) {
        throw new Error(response.data.message || `分片 ${chunkIndex + 1} 上传失败`)
      }

      // 更新进度
      uploadedSize += chunk.size
      const percentage = Math.round((uploadedSize / file.size) * 100)
      progress.value = {
        loaded: uploadedSize,
        total: file.size,
        percentage,
      }
      onProgress?.({
        loaded: uploadedSize,
        total: file.size,
        percentage,
      })

      // 如果是最后一个分片，返回文件路径
      if (chunkIndex === totalChunks - 1 && response.data.data.file_path) {
        return response.data.data.file_path
      }
    }

    throw new Error('文件上传完成但未返回文件路径')
  }

  return {
    uploading,
    progress,
    uploadFile,
  }
}
