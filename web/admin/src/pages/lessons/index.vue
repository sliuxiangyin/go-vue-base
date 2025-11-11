<script setup lang="ts">
import Page from '@/components/global-layout/basic-page.vue'
import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import LessonCreate from './components/lesson-create.vue'
import { useGetLessonsQuery } from '@/services/api/lesson.api'

const page = ref(1)
const pageSize = ref(10)
const searchLevel = ref<number | undefined>(undefined)
const searchIsPublic = ref<boolean | undefined>(undefined)

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  level: searchLevel.value,
  is_public: searchIsPublic.value,
}))

const { data, isLoading, refetch } = useGetLessonsQuery(queryParams)

const lessons = computed(() => data.value?.data?.list || [])
const total = computed(() => data.value?.data?.total || 0)

function handleSearch(level?: number, isPublic?: boolean) {
  searchLevel.value = level
  searchIsPublic.value = isPublic
  page.value = 1
}

function handlePageChange(newPage: number) {
  page.value = newPage
}
</script>

<template>
  <Page
    title="英文课程管理"
    description="管理英语学习课程内容"
    sticky
  >
    <template #actions>
      <LessonCreate @created="refetch" />
    </template>
    <div class="overflow-x-auto">
      <DataTable 
        :loading="isLoading" 
        :data="lessons" 
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
