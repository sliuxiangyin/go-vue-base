<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { useCreateRoleMutation } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

const emit = defineEmits<{
  created: []
}>()

const open = ref(false)
const form = ref({
  name: '',
  display_name: '',
  description: '',
  status: 1,
})

const { mutate: createRole, isPending } = useCreateRoleMutation()

function handleSubmit() {
  if (!form.value.name || !form.value.display_name) {
    toast.error('请填写必填字段')
    return
  }

  createRole(form.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('创建成功')
        open.value = false
        resetForm()
        emit('created')
      } else {
        toast.error(response.message || '创建失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '创建失败')
    },
  })
}

function resetForm() {
  form.value = {
    name: '',
    display_name: '',
    description: '',
    status: 1,
  }
}

function handleOpenChange(newOpen: boolean) {
  open.value = newOpen
  if (!newOpen) {
    resetForm()
  }
}
</script>

<template>
  <UiDialog :open="open" @update:open="handleOpenChange">
    <UiDialogTrigger as-child>
      <UiButton>
        <Plus class="mr-2 h-4 w-4" />
        创建角色
      </UiButton>
    </UiDialogTrigger>
    <UiDialogContent class="sm:max-w-[500px]">
      <UiDialogHeader>
        <UiDialogTitle>创建角色</UiDialogTitle>
        <UiDialogDescription>
          创建新的系统角色
        </UiDialogDescription>
      </UiDialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid gap-2">
          <UiLabel for="name">角色标识 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="name"
            v-model="form.name"
            placeholder="例如: editor"
          />
          <p class="text-xs text-muted-foreground">英文标识，用于系统识别</p>
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
        <UiButton variant="outline" @click="open = false">
          取消
        </UiButton>
        <UiButton :disabled="isPending" @click="handleSubmit">
          {{ isPending ? '创建中...' : '创建' }}
        </UiButton>
      </UiDialogFooter>
    </UiDialogContent>
  </UiDialog>
</template>
