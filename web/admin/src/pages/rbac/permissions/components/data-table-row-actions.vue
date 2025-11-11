<script setup lang="ts">
import { MoreHorizontal, Pencil, Trash2 } from 'lucide-vue-next'
import type { Row } from '@tanstack/vue-table'
import type { Permission } from '@/services/api/rbac.api'
import { useDeletePermissionMutation, useUpdatePermissionMutation, useGetPermissionsQuery } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

interface Props {
  row: Row<Permission>
}

const props = defineProps<Props>()

const showEditDialog = ref(false)
const editForm = ref({
  display_name: '',
  description: '',
  category: '',
  type: 'backend' as 'backend' | 'frontend',
  path: '',
  icon: '',
  parent_id: undefined as number | undefined,
  sort: 0,
  status: 1,
})

const { mutate: deletePermission } = useDeletePermissionMutation()
const { mutate: updatePermission, isPending } = useUpdatePermissionMutation(props.row.original.id)

// 获取前端权限列表用于父级选择
const { data: frontendPermissionsData } = useGetPermissionsQuery({
  page: 1,
  page_size: 100,
  type: 'frontend',
  status: 1,
})

const frontendPermissions = computed(() => {
  const list = frontendPermissionsData.value?.data?.list || []
  // 排除当前权限，避免自己作为自己的父级
  return list.filter(p => p.id !== props.row.original.id)
})

function handleEdit() {
  editForm.value = {
    display_name: props.row.original.display_name,
    description: props.row.original.description || '',
    category: props.row.original.category,
    type: props.row.original.type || 'backend',
    path: props.row.original.path || '',
    icon: props.row.original.icon || '',
    parent_id: props.row.original.parent_id,
    sort: props.row.original.sort || 0,
    status: props.row.original.status,
  }
  showEditDialog.value = true
}

function handleUpdate() {
  // 处理 parent_id，0 转为 undefined
  const submitData = {
    ...editForm.value,
    parent_id: editForm.value.parent_id === 0 ? undefined : editForm.value.parent_id,
  }

  updatePermission(submitData, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('更新成功')
        showEditDialog.value = false
      } else {
        toast.error(response.message || '更新失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '更新失败')
    },
  })
}

function handleDelete() {
  if (confirm('确定要删除这个权限吗？此操作不可恢复。')) {
    deletePermission(props.row.original.id, {
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
  <div>
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
        <UiDropdownMenuItem @click="handleEdit">
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

    <UiDialog v-model:open="showEditDialog">
      <UiDialogContent class="sm:max-w-[500px] max-h-[90vh]">
        <UiDialogHeader>
          <UiDialogTitle>编辑权限</UiDialogTitle>
          <UiDialogDescription>
            编辑权限信息（权限标识不可修改）
          </UiDialogDescription>
        </UiDialogHeader>
        <div class="grid gap-4 py-4 overflow-y-auto max-h-[60vh]">
          <div class="grid gap-2">
            <UiLabel>权限标识</UiLabel>
            <UiInput :model-value="row.original.name" disabled />
          </div>
          <div class="grid gap-2">
            <UiLabel>权限类型</UiLabel>
            <UiInput :model-value="editForm.type === 'backend' ? '后端API权限' : '前端菜单权限'" disabled />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-display_name">权限名称 <span class="text-red-500">*</span></UiLabel>
            <UiInput
              id="edit-display_name"
              v-model="editForm.display_name"
              placeholder="例如: 创建用户"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-category">分类 <span class="text-red-500">*</span></UiLabel>
            <UiInput
              id="edit-category"
              v-model="editForm.category"
              :placeholder="editForm.type === 'backend' ? '例如: 用户管理' : '例如: 菜单'"
            />
          </div>
          <div v-if="editForm.type === 'frontend'" class="grid gap-2">
            <UiLabel for="edit-path">路由路径</UiLabel>
            <UiInput
              id="edit-path"
              v-model="editForm.path"
              placeholder="例如: /users"
            />
          </div>
          <div v-if="editForm.type === 'frontend'" class="grid gap-2">
            <UiLabel for="edit-parent_id">父级菜单</UiLabel>
            <UiSelect v-model="editForm.parent_id">
              <UiSelectTrigger id="edit-parent_id">
                <UiSelectValue placeholder="选择父级菜单（可选）" />
              </UiSelectTrigger>
              <UiSelectContent>
                <UiSelectItem :value="0">无（顶级菜单）</UiSelectItem>
                <UiSelectItem 
                  v-for="perm in frontendPermissions" 
                  :key="perm.id" 
                  :value="perm.id"
                >
                  {{ perm.display_name }} ({{ perm.name }})
                </UiSelectItem>
              </UiSelectContent>
            </UiSelect>
          </div>
          <div v-if="editForm.type === 'frontend'" class="grid gap-2">
            <UiLabel for="edit-icon">图标</UiLabel>
            <UiInput
              id="edit-icon"
              v-model="editForm.icon"
              placeholder="例如: Users"
            />
          </div>
          <div v-if="editForm.type === 'frontend'" class="grid gap-2">
            <UiLabel for="edit-sort">排序</UiLabel>
            <UiInput
              id="edit-sort"
              v-model.number="editForm.sort"
              type="number"
              placeholder="0"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-description">描述</UiLabel>
            <UiTextarea
              id="edit-description"
              v-model="editForm.description"
              placeholder="权限描述"
              rows="3"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-status">状态</UiLabel>
            <UiSelect v-model="editForm.status">
              <UiSelectTrigger>
                <UiSelectValue />
              </UiSelectTrigger>
              <UiSelectContent>
                <UiSelectItem :value="1">启用</UiSelectItem>
                <UiSelectItem :value="0">禁用</UiSelectItem>
              </UiSelectContent>
            </UiSelect>
          </div>
        </div>
        <UiDialogFooter>
          <UiButton variant="outline" @click="showEditDialog = false">
            取消
          </UiButton>
          <UiButton :disabled="isPending" @click="handleUpdate">
            {{ isPending ? '保存中...' : '保存' }}
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>
  </div>
</template>
