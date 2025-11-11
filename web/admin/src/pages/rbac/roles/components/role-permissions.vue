<script setup lang="ts">
import type { Role } from '@/services/api/rbac.api'
import { 
  useGetPermissionsQuery, 
  useGetRolePermissionsQuery,
  useAssignPermissionsToRoleMutation 
} from '@/services/api/rbac.api'
import { toast } from 'vue-sonner'

interface Props {
  role: Role
  open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

// Tab 切换状态
const activeTab = ref<'backend' | 'frontend'>('backend')

// 获取后端权限列表
const backendParams = computed(() => ({
  page: 1,
  page_size: 1000,
  type: 'backend' as const,
  status: 1,
}))
const { data: backendPermissionsData } = useGetPermissionsQuery(backendParams)

// 获取前端权限列表
const frontendParams = computed(() => ({
  page: 1,
  page_size: 1000,
  type: 'frontend' as const,
  status: 1,
}))
const { data: frontendPermissionsData } = useGetPermissionsQuery(frontendParams)

// 获取角色已有权限
const { data: rolePermissionsData, refetch: refetchRolePermissions } = useGetRolePermissionsQuery(props.role.id)
const { mutate: assignPermissions, isPending } = useAssignPermissionsToRoleMutation(props.role.id)

// 当前显示的权限列表
const currentPermissions = computed(() => {
  if (activeTab.value === 'backend') {
    return backendPermissionsData.value?.data?.list || []
  } else {
    return frontendPermissionsData.value?.data?.list || []
  }
})

const selectedPermissions = ref<number[]>([])

// 按分类分组权限
const permissionsByCategory = computed(() => {
  const grouped: Record<string, typeof currentPermissions.value> = {}
  currentPermissions.value.forEach(perm => {
    const category = perm.category || '其他'
    if (!grouped[category]) {
      grouped[category] = []
    }
    grouped[category].push(perm)
  })
  return grouped
})

// 初始化已选权限
watch(() => rolePermissionsData.value, (data) => {
  if (data?.data) {
    selectedPermissions.value = data.data.map(p => p.id)
  }
}, { immediate: true })

function handleSubmit() {
  assignPermissions({ permission_ids: selectedPermissions.value }, {
    onSuccess: (response) => {
      if (response.code === 0) {
        toast.success('权限分配成功')
        refetchRolePermissions()
        emit('update:open', false)
      } else {
        toast.error(response.message || '分配失败')
      }
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || '分配失败')
    },
  })
}

function toggleCategory(category: string) {
  const categoryPermIds = permissionsByCategory.value[category].map(p => p.id)
  const allSelected = categoryPermIds.every(id => selectedPermissions.value.includes(id))
  
  if (allSelected) {
    selectedPermissions.value = selectedPermissions.value.filter(id => !categoryPermIds.includes(id))
  } else {
    selectedPermissions.value = [...new Set([...selectedPermissions.value, ...categoryPermIds])]
  }
}

function isCategorySelected(category: string) {
  const categoryPermIds = permissionsByCategory.value[category].map(p => p.id)
  return categoryPermIds.every(id => selectedPermissions.value.includes(id))
}

function handleTabChange(tab: string | number) {
  if (typeof tab === 'string') {
    activeTab.value = tab as 'backend' | 'frontend'
  }
}
</script>

<template>
  <UiDialog :open="open" @update:open="emit('update:open', $event)">
    <UiDialogContent class="sm:max-w-[800px] max-h-[90vh]">
      <UiDialogHeader>
        <UiDialogTitle>分配权限 - {{ role.display_name }}</UiDialogTitle>
        <UiDialogDescription>
          为角色分配系统权限（后端API权限和前端菜单权限）
        </UiDialogDescription>
      </UiDialogHeader>
      
      <!-- Tab 切换 -->
      <div class="border-b">
        <UiTabs :model-value="activeTab" @update:model-value="handleTabChange">
          <UiTabsList class="grid w-full grid-cols-2">
            <UiTabsTrigger value="backend">
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>
                后端API权限
              </span>
            </UiTabsTrigger>
            <UiTabsTrigger value="frontend">
              <span class="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="7" x="3" y="3" rx="1"/><rect width="7" height="7" x="14" y="3" rx="1"/><rect width="7" height="7" x="14" y="14" rx="1"/><rect width="7" height="7" x="3" y="14" rx="1"/></svg>
                前端菜单权限
              </span>
            </UiTabsTrigger>
          </UiTabsList>
        </UiTabs>
      </div>
      
      <div class="overflow-y-auto max-h-[55vh] py-4">
        <div v-if="Object.keys(permissionsByCategory).length === 0" class="text-center text-muted-foreground py-8">
          暂无权限数据
        </div>
        <div v-for="(permissions, category) in permissionsByCategory" :key="category" class="mb-6">
          <div class="mb-3 flex items-center gap-2">
            <UiCheckbox
              :id="`category-${category}`"
              :checked="isCategorySelected(category)"
               @update:model-value="toggleCategory(category)"
            />
            <label 
              :for="`category-${category}`" 
              class="text-sm font-semibold cursor-pointer"
            >
              {{ category }}
            </label>
          </div>
          <div class="grid grid-cols-2 gap-3 ml-6">
            <div 
              v-for="permission in permissions" 
              :key="permission.id"
              class="flex items-start gap-2"
            >
              <UiCheckbox
                :id="`perm-${permission.id}`"
                :checked="selectedPermissions.includes(permission.id)"
                @update:model-value="(checked: boolean) => {
                  if (checked) {
                    selectedPermissions.push(permission.id)
                  } else {
                    selectedPermissions = selectedPermissions.filter(id => id !== permission.id)
                  }
                }"
              />
              <div class="flex flex-col">
                <label 
                  :for="`perm-${permission.id}`" 
                  class="text-sm font-medium cursor-pointer"
                >
                  {{ permission.display_name }}
                </label>
                <span class="text-xs text-muted-foreground">{{ permission.name }}</span>
                <span v-if="activeTab === 'frontend' && permission.path" class="text-xs text-blue-500">{{ permission.path }}</span>
              </div>
            </div>
          </div>
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
