<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 p-4">
    <div class="w-full max-w-md">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8 space-y-8">
        <div class="flex justify-center">
          <div class="bg-gray-200 border-2 border-dashed rounded-xl w-16 h-16 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-gray-600">
              <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
          </div>
        </div>
        
        <LoginForm />
      </div>
      
      <p class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
        通过登录即表示您同意我们的
        <a href="#" class="underline underline-offset-4 hover:text-primary">
          服务条款
        </a>
        和
        <a href="#" class="underline underline-offset-4 hover:text-primary">
          隐私政策
        </a>
        。
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
</script>