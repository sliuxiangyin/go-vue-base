<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { useCreatePermissionMutation, useGetPermissionsQuery } from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

const emit = defineEmits<{
  created: []
}>()

const open = ref(false)
const form = ref({
  name: '',
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

// 获取前端权限列表用于父级选择
const { data: frontendPermissionsData } = useGetPermissionsQuery({
  page: 1,
  page_size: 100,
  type: 'frontend',
  status: 1,
})

const frontendPermissions = computed(() => frontendPermissionsData.value?.data?.list || [])

const { mutate: createPermission, isPending } = useCreatePermissionMutation()

function handleSubmit() {
  if (!form.value.name || !form.value.display_name || !form.value.category) {
    toast.error('请填写必填字段')
    return
  }

  // 处理 parent_id，0 转为 undefined
  const submitData = {
    ...form.value,
    parent_id: form.value.parent_id === 0 ? undefined : form.value.parent_id,
  }

  createPermission(submitData, {
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
    category: '',
    type: 'backend',
    path: '',
    icon: '',
    parent_id: undefined,
    sort: 0,
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
        创建权限
      </UiButton>
    </UiDialogTrigger>
    <UiDialogContent class="sm:max-w-[500px] max-h-[90vh]">
      <UiDialogHeader>
        <UiDialogTitle>创建权限</UiDialogTitle>
        <UiDialogDescription>
          创建新的系统权限
        </UiDialogDescription>
      </UiDialogHeader>
      <div class="grid gap-4 py-4 overflow-y-auto max-h-[60vh]">
        <div class="grid gap-2">
          <UiLabel for="type">权限类型 <span class="text-red-500">*</span></UiLabel>
          <UiSelect v-model="form.type">
            <UiSelectTrigger id="type">
              <UiSelectValue />
            </UiSelectTrigger>
            <UiSelectContent>
              <UiSelectItem value="backend">后端API权限</UiSelectItem>
              <UiSelectItem value="frontend">前端菜单权限</UiSelectItem>
            </UiSelectContent>
          </UiSelect>
          <p class="text-xs text-muted-foreground">
            后端API权限用于接口鉴权，前端菜单权限用于菜单显示
          </p>
        </div>
        <div class="grid gap-2">
          <UiLabel for="name">权限标识 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="name"
            v-model="form.name"
            :placeholder="form.type === 'backend' ? '例如: user.create' : '例如: menu.users'"
          />
          <p class="text-xs text-muted-foreground">
            {{ form.type === 'backend' ? '通常为路由名称，用于系统鉴权' : '菜单唯一标识' }}
          </p>
        </div>
        <div class="grid gap-2">
          <UiLabel for="display_name">权限名称 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="display_name"
            v-model="form.display_name"
            placeholder="例如: 创建用户"
          />
        </div>
        <div class="grid gap-2">
          <UiLabel for="category">分类 <span class="text-red-500">*</span></UiLabel>
          <UiInput
            id="category"
            v-model="form.category"
            :placeholder="form.type === 'backend' ? '例如: 用户管理' : '例如: 菜单'"
          />
        </div>
        <div v-if="form.type === 'frontend'" class="grid gap-2">
          <UiLabel for="path">路由路径</UiLabel>
          <UiInput
            id="path"
            v-model="form.path"
            placeholder="例如: /users"
          />
          <p class="text-xs text-muted-foreground">前端路由路径</p>
        </div>
        <div v-if="form.type === 'frontend'" class="grid gap-2">
          <UiLabel for="parent_id">父级菜单</UiLabel>
          <UiSelect v-model="form.parent_id">
            <UiSelectTrigger id="parent_id">
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
          <p class="text-xs text-muted-foreground">选择父级菜单可创建嵌套菜单</p>
        </div>
        <div v-if="form.type === 'frontend'" class="grid gap-2">
          <UiLabel for="icon">图标</UiLabel>
          <UiInput
            id="icon"
            v-model="form.icon"
            placeholder="例如: Users"
          />
          <p class="text-xs text-muted-foreground">Lucide 图标名称</p>
        </div>
        <div v-if="form.type === 'frontend'" class="grid gap-2">
          <UiLabel for="sort">排序</UiLabel>
          <UiInput
            id="sort"
            v-model.number="form.sort"
            type="number"
            placeholder="0"
          />
          <p class="text-xs text-muted-foreground">数字越小越靠前</p>
        </div>
        <div class="grid gap-2">
          <UiLabel for="description">描述</UiLabel>
          <UiTextarea
            id="description"
            v-model="form.description"
            placeholder="权限描述"
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
