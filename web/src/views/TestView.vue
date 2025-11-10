<script setup lang="ts">
import {RouterLink, RouterView} from "vue-router";
</script>
<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Navigation -->

    <nav class="bg-white shadow-md">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold">Database AI</h1>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <RouterLink
                  to="/test/home"
                  class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  active-class="border-indigo-500 text-gray-900"
              >
                Home
              </RouterLink>
              <RouterLink
                  to="/test/about"
                  class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  active-class="border-indigo-500 text-gray-900"
              >
                About
              </RouterLink>
              <RouterLink
                  to="/test/api-test"
                  class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  active-class="border-indigo-500 text-gray-900"
              >
                API测试
              </RouterLink>
              <RouterLink
                  to="/test/db-account"
                  class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  active-class="border-indigo-500 text-gray-900"
              >
                数据库账号
              </RouterLink>
              <RouterLink
                  to="/test/login"
                  class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  active-class="border-indigo-500 text-gray-900"
              >
                登录
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </nav>
    <!-- Page Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <RouterView />
      </div>
    </main>
  </div>
</template>

<style scoped>

</style>