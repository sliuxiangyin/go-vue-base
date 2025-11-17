<script setup lang="ts">
import { Upload, RefreshCw, Trash2 } from 'lucide-vue-next'
import Page from '@/components/global-layout/basic-page.vue'
import SortSelect from '@/components/sort-select/index.vue'
import type { TSort } from '@/components/sort-select/types'
import FileCard from './components/file-card.vue'
import FileEditDialog from './components/file-edit-dialog.vue'
import { 
  useGetFilesQuery, 
  useDeleteFileMutation, 
  useBatchDeleteFilesMutation,
  useDownloadFile,
  type File,
  type FileType,
  type FileSource
} from '@/services/api/file.api'
import { useFileUpload } from '@/composables/use-file-upload'
import { toast } from 'vue-sonner'

const page = ref(1)
const pageSize = ref(20)
const searchTerm = ref('')
const fileType = ref<FileType>('')
const fileSource = ref<FileSource>('')
const sort = ref<TSort>('desc')

const fileTypes: { label: string; value: FileType }[] = [
  { label: '全部类型', value: 'undefined' },
  { label: '图片', value: 'image' },
  { label: '音频', value: 'audio' },
  { label: '视频', value: 'video' },
  { label: '文档', value: 'document' },
  { label: '其他', value: 'other' },
]

const fileSources: { label: string; value: FileSource }[] = [
  { label: '全部来源', value: 'undefined' },
  { label: '上传', value: 'upload' },
  { label: '下载', value: 'download' },
  { label: '生成', value: 'generate' },
]

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  file_type: fileType.value,
  source: fileSource.value,
  keyword: searchTerm.value,
}))

const { data, isLoading, refetch } = useGetFilesQuery(queryParams)
const { mutate: deleteFile } = useDeleteFileMutation()
const { mutate: batchDelete } = useBatchDeleteFilesMutation()
const { mutate: downloadFile } = useDownloadFile()
const { uploading, progress, uploadFile } = useFileUpload()

const files = computed(() => {
  let list = data.value?.data?.list || []
  
  // 客户端排序
  if (sort.value === 'asc') {
    list = [...list].sort((a, b) => a.name.localeCompare(b.name))
  } else {
    list = [...list].sort((a, b) => b.name.localeCompare(a.name))
  }
  
  return list
})

const total = computed(() => data.value?.data?.total || 0)

// 文件编辑
const editingFile = ref<File | null>(null)
const showEditDialog = ref(false)

function handleEdit(file: File) {
  editingFile.value = file
  showEditDialog.value = true
}

// 文件删除
function handleDelete(id: number) {
  if (!confirm('确定要删除这个文件吗？')) return
  
  deleteFile({ id, delete_physical: true }, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('删除成功')
        refetch()
      } else {
        toast.error(response.message || '删除失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '删除失败')
    },
  })
}

// 文件下载
function handleDownload(id: number) {
  downloadFile(id, {
    onSuccess: (response) => {
      if (response.code === 0 && response.data) {
        const link = document.createElement('a')
        link.href = `/${response.data.file_path}`
        link.download = response.data.file_name
        link.click()
      } else {
        toast.error(response.message || '下载失败')
      }
    },
    onError: () => {
      toast.error('下载失败')
    },
  })
}

// 文件上传
const fileInputRef = ref<HTMLInputElement | null>(null)

function triggerFileUpload() {
  fileInputRef.value?.click()
}

async function handleFileUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return
  
  try {
    await uploadFile(file, (prog) => {
      console.log(`上传进度: ${prog.percentage}%`)
    })
    toast.success('上传成功')
    refetch()
    // 重置文件输入
    if (target) target.value = ''
  } catch (error: any) {
    toast.error(`上传失败: ${error.message}`)
  }
}

// 分页
function handlePageChange(newPage: number) {
  page.value = newPage
}
</script>

<template>
  <Page
    title="文件管理"
    description="管理所有上传和下载的文件"
    sticky
  >
    <!-- 工具栏 -->
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-col gap-4 sm:flex-row sm:flex-1">
        <!-- 搜索 -->
        <UiInput
          v-model:model-value="searchTerm"
          placeholder="搜索文件名..."
          class="h-9 w-full sm:w-[250px]"
        />

        <!-- 文件类型筛选 -->
        <UiSelect v-model:model-value="fileType">
          <UiSelectTrigger class="h-9 w-full sm:w-36">
            <UiSelectValue />
          </UiSelectTrigger>
          <UiSelectContent>
            <UiSelectItem
              v-for="type in fileTypes"
              :key="type.value"
              :value="type.value"
            >
              {{ type.label }}
            </UiSelectItem>
          </UiSelectContent>
        </UiSelect>

        <!-- 来源筛选 -->
        <UiSelect v-model:model-value="fileSource">
          <UiSelectTrigger class="h-9 w-full sm:w-36">
            <UiSelectValue />
          </UiSelectTrigger>
          <UiSelectContent>
            <UiSelectItem
              v-for="source in fileSources"
              :key="source.value"
              :value="source.value"
            >
              {{ source.label }}
            </UiSelectItem>
          </UiSelectContent>
        </UiSelect>
      </div>

      <!-- 操作按钮 -->
      <div class="flex items-center gap-2">
        <UiButton
          variant="outline"
          size="sm"
          :disabled="isLoading"
          @click="refetch"
        >
          <RefreshCw class="h-4 w-4 mr-1" :class="{ 'animate-spin': isLoading }" />
          刷新
        </UiButton>
        
        <input
          ref="fileInputRef"
          type="file"
          class="hidden"
          @change="handleFileUpload"
        />
        
        <UiButton
          size="sm"
          :disabled="uploading"
          @click="triggerFileUpload"
        >
          <Upload class="h-4 w-4 mr-1" />
          {{ uploading ? `上传中 ${progress.percentage}%` : '上传文件' }}
        </UiButton>
        
        <SortSelect v-model:sort="sort" />
      </div>
    </div>

    <!-- 上传进度 -->
    <div v-if="uploading" class="mt-4">
      <div class="w-full bg-secondary rounded-full h-2">
        <div 
          class="bg-primary h-2 rounded-full transition-all duration-300"
          :style="{ width: `${progress.percentage}%` }"
        />
      </div>
      <p class="text-sm text-muted-foreground mt-1">
        上传中... {{ progress.percentage }}%
      </p>
    </div>

    <!-- 文件网格 -->
    <main class="grid grid-cols-1 gap-4 mt-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <FileCard
        v-for="file in files"
        :key="file.id"
        :file="file"
        @delete="handleDelete"
        @edit="handleEdit"
        @download="handleDownload"
      />
    </main>

    <!-- 空状态 -->
    <div v-if="!isLoading && files.length === 0" class="flex flex-col items-center justify-center py-12">
      <Upload class="h-16 w-16 text-muted-foreground mb-4" />
      <h3 class="text-lg font-semibold">暂无文件</h3>
      <p class="text-sm text-muted-foreground mt-1">
        上传您的第一个文件
      </p>
      <UiButton class="mt-4" @click="triggerFileUpload">
        <Upload class="h-4 w-4 mr-2" />
        上传文件
      </UiButton>
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="flex items-center justify-center mt-6">
      <div class="flex items-center gap-2">
        <UiButton
          variant="outline"
          size="sm"
          :disabled="page === 1"
          @click="handlePageChange(page - 1)"
        >
          上一页
        </UiButton>
        <span class="text-sm text-muted-foreground">
          第 {{ page }} 页，共 {{ Math.ceil(total / pageSize) }} 页
        </span>
        <UiButton
          variant="outline"
          size="sm"
          :disabled="page >= Math.ceil(total / pageSize)"
          @click="handlePageChange(page + 1)"
        >
          下一页
        </UiButton>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <FileEditDialog
      v-model:open="showEditDialog"
      :file="editingFile"
      @updated="refetch"
    />
  </Page>
</template>

<route lang="yaml">
meta:
  auth: true
</route>
