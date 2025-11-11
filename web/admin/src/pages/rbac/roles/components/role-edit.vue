<script setup lang="ts">
import type { Role } from '@/services/api/rbac.api'
import { useUpdateRoleMutation } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

interface Props {
  role: Role
  open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const form = ref({
  display_name: props.role.display_name,
  description: props.role.description || '',
  status: props.role.status,
})

const { mutate: updateRole, isPending } = useUpdateRoleMutation(props.role.id)

function handleSubmit() {
  if (!form.value.display_name) {
    toast.error('请填写角色名称')
    return
  }

  updateRole(form.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('更新成功')
        emit('update:open', false)
      } else {
        toast.error(response.message || '更新失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '更新失败')
    },
  })
}
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="sm:max-w-[500px]">
      <UiDialogHeader>
        <UiDialogTitle>编辑角色</UiDialogTitle>
        <UiDialogDescription>
          编辑角色信息（角色标识不可修改）
        </UiDialogDescription>
      </UiDialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid gap-2">
          <UiLabel>角色标识</UiLabel>
          <UiInput :model-value="role.name" disabled />
        </div>
        <div class="grid gap-2">
          <UiLabel for="display_name">角色名称 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="display_name"
            v-model="form.display_name"
            placeholder="例如: 编辑者"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="description">描述</UiLabel>
          <UiTextarea
            id="description"
            v-model="form.description"
            placeholder="角色描述"
            rows="3"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="status">状态</UiLabel>
          <UiSelect v-model="form.status">
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
        <UiButton variant="outline" @click="emit('update:open', false)">
          取消
        </UiButton>
        <UiButton :disabled="isPending" @click="handleSubmit">
          {{ isPending ? '保存中...' : '保存' }}
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
