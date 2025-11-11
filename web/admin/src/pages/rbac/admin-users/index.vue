<script setup lang="ts">
import Page from '@/components/global-layout/basic-page.vue'
import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import AdminUserCreate from './components/admin-user-create.vue'
import { useGetAdminUsersQuery } from '@/services/api/rbac.api'

const page = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const searchStatus = ref<number | string>('-1')

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  keyword: searchKeyword.value,
  status: searchStatus.value === '-1' ? undefined : Number(searchStatus.value),
}))

const { data, isLoading, refetch } = useGetAdminUsersQuery(queryParams)

const adminUsers = computed(() => data.value?.data?.list || [])
const total = computed(() => data.value?.data?.total || 0)

function handleSearch(keyword: string, status?: number) {
  searchKeyword.value = keyword
  searchStatus.value = status === undefined ? '-1' : status
  page.value = 1
}

function handlePageChange(newPage: number) {
  page.value = newPage
}
</script>

<template>
  <Page
    title="管理员管理"
    description="管理系统管理员账号和权限"
    sticky
  >
    <template #actions>
      <AdminUserCreate @created="refetch" />
    </template>

    <div class="overflow-x-auto">
      <DataTable 
        :loading="isLoading" 
        :data="adminUsers" 
        :columns="columns"
        :total="total"
        :page="page"
        :page-size="pageSize"
        @search="handleSearch"
        @page-change="handlePageChange"
        @refresh="refetch"
      />
    </div>
  </Page>
</template>
