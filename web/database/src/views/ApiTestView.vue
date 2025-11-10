<template>
  <div class="max-w-4xl mx-auto p-6">
    <h1 class="text-3xl font-bold mb-6">API 测试页面</h1>
    
    <!-- 设置键值对 -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold mb-4">设置键值对</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">键 (Key)</label>
          <input 
            v-model="newEntry.key" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="输入键名"
          >
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">值 (Value)</label>
          <input 
            v-model="newEntry.value" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="输入值"
          >
        </div>
      </div>
      <button 
        @click="setValue" 
        :disabled="!newEntry.key || !newEntry.value"
        class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        保存
      </button>
    </div>

    <!-- 获取值 -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-xl font-semibold mb-4">获取值</h2>
      <div class="flex gap-4">
        <div class="flex-grow">
          <label class="block text-sm font-medium text-gray-700 mb-1">键 (Key)</label>
          <input 
            v-model="getKey" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="输入要查询的键名"
          >
        </div>
        <div class="self-end">
          <button 
            @click="getValue" 
            :disabled="!getKey"
            class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed h-[42px]"
          >
            查询
          </button>
        </div>
      </div>
      <div v-if="getResult" class="mt-4 p-4 bg-gray-50 rounded-md">
        <p class="font-medium">查询结果:</p>
        <p>键: {{ getResult.key }}</p>
        <p>值: {{ getResult.value }}</p>
      </div>
      <div v-if="getError" class="mt-4 p-4 bg-red-50 text-red-700 rounded-md">
        <p>{{ getError }}</p>
      </div>
    </div>

    <!-- 响应状态 -->
    <div v-if="responseMessage" class="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-md">
      {{ responseMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// 状态变量
const newEntry = ref({ key: '', value: '' })
const getKey = ref('')
const getResult = ref<{key: string, value: string} | null>(null)
const getError = ref('')
const responseMessage = ref('')

// 设置键值对
const setValue = async () => {
  try {
    const response = await fetch('/api/test', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        key: newEntry.value.key,
        value: newEntry.value.value
      })
    })

    if (response.ok) {
      responseMessage.value = '键值对保存成功!'
      // 清空表单
      newEntry.value.key = ''
      newEntry.value.value = ''
      
      // 3秒后清除消息
      setTimeout(() => {
        responseMessage.value = ''
      }, 3000)
    } else {
      responseMessage.value = `保存失败: ${response.status} ${response.statusText}`
    }
  } catch (error) {
    responseMessage.value = `请求失败: ${(error as Error).message}`
  }
}

// 获取值
const getValue = async () => {
  try {
    // 清除之前的结果
    getResult.value = null
    getError.value = ''
    
    const response = await fetch(`/api/test/${getKey.value}`)
    
    if (response.ok) {
      const data = await response.json()
      getResult.value = data
    } else if (response.status === 404) {
      getError.value = '未找到指定的键'
    } else {
      getError.value = `查询失败: ${response.status} ${response.statusText}`
    }
  } catch (error) {
    getError.value = `请求失败: ${(error as Error).message}`
  }
}
</script>