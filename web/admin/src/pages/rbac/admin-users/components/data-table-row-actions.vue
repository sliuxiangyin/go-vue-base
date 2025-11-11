<script setup lang="ts">
import { MoreHorizontal, Pencil, Key, Shield, Trash } from 'lucide-vue-next'
import type { AdminUser } from '@/services/api/rbac.api'
import { 
  useUpdateAdminUserMutation, 
  useUpdateAdminUserPasswordMutation,
  useDeleteAdminUserMutation,
  useGetAllRolesQuery,
  useGetUserRolesQuery,
  useAssignRolesToUserMutation
} from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'
import { useQueryClient } from '@tanstack/vue-query'

interface Props {
  row: AdminUser
}

const props = defineProps<Props>()

const queryClient = useQueryClient()
const showEditDialog = ref(false)
const showPasswordDialog = ref(false)
const showRolesDialog = ref(false)
const showDeleteDialog = ref(false)

// 编辑表单
const editForm = ref({
  email: props.row.email || '',
  nickname: props.row.nickname || '',
  avatar: props.row.avatar || '',
  status: props.row.status,
})

// 密码表单
const passwordForm = ref({
  password: '',
  confirmPassword: '',
})

// 角色分配
const { data: allRolesData } = useGetAllRolesQuery()
const { data: userRolesData } = useGetUserRolesQuery(props.row.id)
const selectedRoleIds = ref<number[]>([])

const allRoles = computed(() => allRolesData.value?.data || [])

// 初始化已选角色
watch(() => userRolesData.value, (data) => {
  if (data?.data) {
    selectedRoleIds.value = data.data.map(r => r.id)
  }
}, { immediate: true })

const { mutate: updateAdminUser, isPending: isUpdating } = useUpdateAdminUserMutation(props.row.id)
const { mutate: updatePassword, isPending: isUpdatingPassword } = useUpdateAdminUserPasswordMutation(props.row.id)
const { mutate: assignRoles, isPending: isAssigningRoles } = useAssignRolesToUserMutation(props.row.id)
const { mutate: deleteAdminUser, isPending: isDeleting } = useDeleteAdminUserMutation()

function handleUpdate() {
  updateAdminUser(editForm.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('更新成功')
        showEditDialog.value = false
        queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      } else {
        toast.error(response.message || '更新失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '更新失败')
    },
  })
}

function handleUpdatePassword() {
  if (!passwordForm.value.password) {
    toast.error('请输入新密码')
    return
  }
  if (passwordForm.value.password.length < 6) {
    toast.error('密码至少6位')
    return
  }
  if (passwordForm.value.password !== passwordForm.value.confirmPassword) {
    toast.error('两次密码输入不一致')
    return
  }

  updatePassword({ password: passwordForm.value.password }, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('密码重置成功')
        showPasswordDialog.value = false
        passwordForm.value = { password: '', confirmPassword: '' }
      } else {
        toast.error(response.message || '重置失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '重置失败')
    },
  })
}

function handleAssignRoles() {
  assignRoles({ role_ids: selectedRoleIds.value }, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('角色分配成功')
        showRolesDialog.value = false
        queryClient.invalidateQueries({ queryKey: ['admin-users'] })
        queryClient.invalidateQueries({ queryKey: ['user-roles', props.row.id] })
      } else {
        toast.error(response.message || '分配失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '分配失败')
    },
  })
}

function handleDelete() {
  deleteAdminUser(props.row.id, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('删除成功')
        showDeleteDialog.value = false
        queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      } else {
        toast.error(response.message || '删除失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '删除失败')
    },
  })
}
</script>

<template>
  <div class="flex items-center gap-2">
    <UiDropdownMenu>
      <UiDropdownMenuTrigger as-child>
        <UiButton variant="ghost" class="h-8 w-8 p-0">
          <span class="sr-only">打开菜单</span>
          <MoreHorizontal class="h-4 w-4" />
        </UiButton>
      </UiDropdownMenuTrigger>
      <UiDropdownMenuContent align="end">
        <UiDropdownMenuItem @click="showEditDialog = true">
          <Pencil class="mr-2 h-4 w-4" />
          编辑
        </UiDropdownMenuItem>
        <UiDropdownMenuItem @click="showPasswordDialog = true">
          <Key class="mr-2 h-4 w-4" />
          重置密码
        </UiDropdownMenuItem>
        <UiDropdownMenuItem @click="showRolesDialog = true">
          <Shield class="mr-2 h-4 w-4" />
          分配角色
        </UiDropdownMenuItem>
        <UiDropdownMenuSeparator />
        <UiDropdownMenuItem class="text-red-600" @click="showDeleteDialog = true">
          <Trash class="mr-2 h-4 w-4" />
          删除
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>

    <!-- 编辑对话框 -->
    <UiDialog v-model:open="showEditDialog">
      <UiDialogContent class="sm:max-w-[500px] max-h-[90vh]">
        <UiDialogHeader>
          <UiDialogTitle>编辑管理员</UiDialogTitle>
          <UiDialogDescription>
            编辑管理员信息（用户名不可修改）
          </UiDialogDescription>
        </UiDialogHeader>
        <div class="grid gap-4 py-4 overflow-y-auto max-h-[60vh]">
          <div class="grid gap-2">
            <UiLabel>用户名</UiLabel>
            <UiInput :model-value="row.username" disabled />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-email">邮箱</UiLabel>
            <UiInput
              id="edit-email"
              v-model="editForm.email"
              type="email"
              placeholder="请输入邮箱"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-nickname">昵称</UiLabel>
            <UiInput
              id="edit-nickname"
              v-model="editForm.nickname"
              placeholder="请输入昵称"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="edit-avatar">头像URL</UiLabel>
            <UiInput
              id="edit-avatar"
              v-model="editForm.avatar"
              placeholder="请输入头像URL"
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
          <UiButton :disabled="isUpdating" @click="handleUpdate">
            {{ isUpdating ? '保存中...' : '保存' }}
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>

    <!-- 重置密码对话框 -->
    <UiDialog v-model:open="showPasswordDialog">
      <UiDialogContent class="sm:max-w-[400px]">
        <UiDialogHeader>
          <UiDialogTitle>重置密码</UiDialogTitle>
          <UiDialogDescription>
            为 {{ row.username }} 设置新密码
          </UiDialogDescription>
        </UiDialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <UiLabel for="new-password">新密码</UiLabel>
            <UiInput
              id="new-password"
              v-model="passwordForm.password"
              type="password"
              placeholder="请输入新密码（至少6位）"
            />
          </div>
          <div class="grid gap-2">
            <UiLabel for="confirm-password">确认密码</UiLabel>
            <UiInput
              id="confirm-password"
              v-model="passwordForm.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
            />
          </div>
        </div>
        <UiDialogFooter>
          <UiButton variant="outline" @click="showPasswordDialog = false">
            取消
          </UiButton>
          <UiButton :disabled="isUpdatingPassword" @click="handleUpdatePassword">
            {{ isUpdatingPassword ? '重置中...' : '确认重置' }}
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>

    <!-- 分配角色对话框 -->
    <UiDialog v-model:open="showRolesDialog">
      <UiDialogContent class="sm:max-w-[600px] max-h-[90vh]">
        <UiDialogHeader>
          <UiDialogTitle>分配角色 - {{ row.username }}</UiDialogTitle>
          <UiDialogDescription>
            为管理员分配系统角色
          </UiDialogDescription>
        </UiDialogHeader>
        <div class="overflow-y-auto max-h-[60vh] py-4">
          <div class="grid gap-3">
            <div 
              v-for="role in allRoles" 
              :key="role.id"
              class="flex items-start gap-2"
            >
              <UiCheckbox
                :id="`role-${role.id}`"
                :model-value="selectedRoleIds.includes(role.id)"
                @update:model-value="(checked) => {
                  if (checked === true) {
                    selectedRoleIds = [...selectedRoleIds, role.id]
                  } else {
                    selectedRoleIds = selectedRoleIds.filter(id => id !== role.id)
                  }
                }"
              />
              <div class="flex flex-col">
                <label 
                  :for="`role-${role.id}`" 
                  class="text-sm font-medium cursor-pointer"
                >
                  {{ role.display_name }}
                </label>
                <span class="text-xs text-muted-foreground">{{ role.description || role.name }}</span>
              </div>
            </div>
          </div>
        </div>
        <UiDialogFooter>
          <UiButton variant="outline" @click="showRolesDialog = false">
            取消
          </UiButton>
          <UiButton :disabled="isAssigningRoles" @click="handleAssignRoles">
            {{ isAssigningRoles ? '保存中...' : '保存' }}
          </UiButton>
        </UiDialogFooter>
      </UiDialogContent>
    </UiDialog>

    <!-- 删除确认对话框 -->
    <UiAlertDialog v-model:open="showDeleteDialog">
      <UiAlertDialogContent>
        <UiAlertDialogHeader>
          <UiAlertDialogTitle>确认删除</UiAlertDialogTitle>
          <UiAlertDialogDescription>
            确定要删除管理员 <span class="font-semibold text-foreground">{{ row.username }}</span> 吗？此操作无法撤销。
          </UiAlertDialogDescription>
        </UiAlertDialogHeader>
        <UiAlertDialogFooter>
          <UiAlertDialogCancel>取消</UiAlertDialogCancel>
          <UiAlertDialogAction :disabled="isDeleting" @click="handleDelete">
            {{ isDeleting ? '删除中...' : '确认删除' }}
          </UiAlertDialogAction>
        </UiAlertDialogFooter>
      </UiAlertDialogContent>
    </UiAlertDialog>
  </div>
</template>
