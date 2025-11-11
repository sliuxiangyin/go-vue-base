<script setup lang="ts">
import Page from '@/components/global-layout/basic-page.vue'
import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import PermissionCreate from './components/permission-create.vue'
import { useGetPermissionsQuery } from '@/services/api/rbac.api'

const page = ref(1)
const pageSize = ref(10)
const searchCategory = ref('')
const searchStatus = ref<number | undefined>(undefined)
const activeTab = ref<'backend' | 'frontend'>('backend')

// 使用 computed 创建响应式查询参数
const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  category: searchCategory.value,
  type: activeTab.value,
  status: searchStatus.value,
}))

const { data, isLoading, refetch } = useGetPermissionsQuery(queryParams)

const permissions = computed(() => data.value?.data?.list || [])
const total = computed(() => data.value?.data?.total || 0)

function handleSearch(category: string, status?: number) {
  searchCategory.value = category
  searchStatus.value = status
  page.value = 1
  refetch()
}

function handlePageChange(newPage: number) {
  page.value = newPage
  refetch()
}

function handleTabChange(tab: string | number) {
    console.log(tab);
  if (typeof tab === 'string') {
    activeTab.value = tab as 'backend' | 'frontend'
    page.value = 1
    refetch()
  }
}
</script>

<template>
  <Page
    title="权限管理"
    description="管理系统权限配置"
    sticky
  >
    <template #actions>
      <PermissionCreate @created="refetch" />
    </template>
    
    <!-- Tab 切换 -->
    <div class="mb-4">
      <UiTabs :model-value="activeTab" @update:model-value="handleTabChange">
        <UiTabsList>
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

    <div class="overflow-x-auto">
      <DataTable 
        :loading="isLoading" 
        :data="permissions" 
        :columns="columns"
        :total="total"
        :page="page"
        :page-size="pageSize"
        :permission-type="activeTab"
        @search="handleSearch"
        @page-change="handlePageChange"
        @refresh="refetch"
      />
    </div>
  </Page>
</template>
