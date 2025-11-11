<script setup lang="ts">
import { useCreateAdminUserMutation } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

const emit = defineEmits<{
  created: []
}>()

const open = ref(false)
const form = ref({
  username: '',
  password: '',
  email: '',
  nickname: '',
  avatar: '',
  status: 1,
})

const { mutate: createAdminUser, isPending } = useCreateAdminUserMutation()

function handleSubmit() {
  if (!form.value.username || !form.value.password) {
    toast.error('请填写用户名和密码')
    return
  }

  if (form.value.password.length < 6) {
    toast.error('密码至少6位')
    return
  }

  createAdminUser(form.value, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('创建成功')
        open.value = false
        emit('created')
        // 重置表单
        form.value = {
          username: '',
          password: '',
          email: '',
          nickname: '',
          avatar: '',
          status: 1,
        }
      } else {
        toast.error(response.message || '创建失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '创建失败')
    },
  })
}
</script>

<template>
  <UiDialog v-model:open="open">
    <UiDialogTrigger as-child>
      <UiButton>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" x2="19" y1="8" y2="14"/><line x1="22" x2="16" y1="11" y2="11"/></svg>
        创建管理员
      </UiButton>
    </UiDialogTrigger>
    <UiDialogContent class="sm:max-w-[500px] max-h-[90vh]">
      <UiDialogHeader>
        <UiDialogTitle>创建管理员</UiDialogTitle>
        <UiDialogDescription>
          创建新的系统管理员账号
        </UiDialogDescription>
      </UiDialogHeader>
      <div class="grid gap-4 py-4 overflow-y-auto max-h-[60vh]">
        <div class="grid gap-2">
          <UiLabel for="username">用户名 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="username"
            v-model="form.username"
            placeholder="请输入用户名"
            :disabled="isPending"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="password">密码 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            :disabled="isPending"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="email">邮箱</UiLabel>
          <UiInput
            id="email"
            v-model="form.email"
            type="email"
            placeholder="请输入邮箱"
            :disabled="isPending"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="nickname">昵称</UiLabel>
          <UiInput
            id="nickname"
            v-model="form.nickname"
            placeholder="请输入昵称"
            :disabled="isPending"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="avatar">头像URL</UiLabel>
          <UiInput
            id="avatar"
            v-model="form.avatar"
            placeholder="请输入头像URL"
            :disabled="isPending"
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
