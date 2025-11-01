<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface ApiResponse {
  message: string
}

const apiData = ref<ApiResponse | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await fetch('/api/hello')
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    apiData.value = await response.json()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unknown error occurred'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <main class="p-6">
    <div class="max-w-4xl mx-auto">
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">Database AI</h1>
        <p class="text-lg text-gray-600">GoFiber + Vue 3 + Shadcn-vue 全栈应用示例</p>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-2xl font-semibold text-gray-800 mb-4">API 测试</h2>
        
        <div v-if="loading" class="text-center py-4">
          <p class="text-gray-600">加载中...</p>
        </div>
        
        <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
          <p class="text-red-800">错误: {{ error }}</p>
        </div>
        
        <div v-else-if="apiData" class="bg-green-50 border border-green-200 rounded-md p-4">
          <p class="text-green-800">后端响应: {{ apiData.message }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-2">后端技术</h3>
          <ul class="list-disc list-inside text-gray-600">
            <li>GoFiber</li>
            <li>嵌入式静态文件</li>
            <li>RESTful API</li>
          </ul>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-2">前端技术</h3>
          <ul class="list-disc list-inside text-gray-600">
            <li>Vue 3 + TypeScript</li>
            <li>Shadcn-vue 组件库</li>
            <li>Tailwind CSS</li>
          </ul>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-xl font-semibold text-gray-800 mb-2">特性</h3>
          <ul class="list-disc list-inside text-gray-600">
            <li>前后端一体化</li>
            <li>自动构建打包</li>
            <li>开发热重载</li>
          </ul>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
main {
  min-height: 100vh;
  background-color: #f8fafc;
}
</style>