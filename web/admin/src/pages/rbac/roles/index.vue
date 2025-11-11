<script setup lang="ts">
import Page from '@/components/global-layout/basic-page.vue'
import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import RoleCreate from './components/role-create.vue'
import { useGetRolesQuery } from '@/services/api/rbac.api'

const page = ref(1)
const pageSize = ref(10)
const searchName = ref('')
const searchStatus = ref<number | undefined>(undefined)

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  name: searchName.value,
  status: searchStatus.value,
}))

const { data, isLoading, refetch } = useGetRolesQuery(queryParams)

const roles = computed(() => data.value?.data?.list || [])
const total = computed(() => data.value?.data?.total || 0)

function handleSearch(name: string, status?: number) {
  searchName.value = name
  searchStatus.value = status
  page.value = 1
}

function handlePageChange(newPage: number) {
  page.value = newPage
}
</script>

<template>
  <Page
    title="角色管理"
    description="管理系统角色和权限分配"
    sticky
  >
    <template #actions>
      <RoleCreate @created="refetch" />
    </template>
    <div class="overflow-x-auto">
      <DataTable 
        :loading="isLoading" 
        :data="roles" 
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
