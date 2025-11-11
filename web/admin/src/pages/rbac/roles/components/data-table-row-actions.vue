<script setup lang="ts">
import { MoreHorizontal, Pencil, Trash2, Shield } from 'lucide-vue-next'
import type { Row } from '@tanstack/vue-table'
import type { Role } from '@/services/api/rbac.api'
import { useDeleteRoleMutation } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'
import RoleEdit from './role-edit.vue'
import RolePermissions from './role-permissions.vue'

interface Props {
  row: Row<Role>
}

const props = defineProps<Props>()

const showEditDialog = ref(false)
const showPermissionsDialog = ref(false)

const { mutate: deleteRole } = useDeleteRoleMutation()

function handleDelete() {
  if (props.row.original.is_system) {
    toast.error('系统角色不能删除')
    return
  }

  if (confirm('确定要删除这个角色吗？此操作不可恢复。')) {
    deleteRole(props.row.original.id, {
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
  <div class="flex items-center gap-2">
    <UiButton
      variant="ghost"
      size="sm"
      @click="showPermissionsDialog = true"
    >
      <Shield class="h-4 w-4" />
      权限
    </UiButton>

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
        <UiDropdownMenuItem @click="showEditDialog = true">
          <Pencil class="mr-2 h-4 w-4" />
          编辑
        </UiDropdownMenuItem>
        <UiDropdownMenuSeparator />
        <UiDropdownMenuItem
          class="text-red-600"
          :disabled="row.original.is_system"
          @click="handleDelete"
        >
          <Trash2 class="mr-2 h-4 w-4" />
          删除
        </UiDropdownMenuItem>
      </UiDropdownMenuContent>
    </UiDropdownMenu>

    <RoleEdit
      v-if="showEditDialog"
      :role="row.original"
      :open="showEditDialog"
      @update:open="showEditDialog = $event"
    />

    <RolePermissions
      v-if="showPermissionsDialog"
      :role="row.original"
      :open="showPermissionsDialog"
      @update:open="showPermissionsDialog = $event"
    />
  </div>
</template>
