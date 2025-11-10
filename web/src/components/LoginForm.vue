<template>
  <div class="w-full max-w-sm mx-auto space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-3xl font-bold">登录</h1>
      <p class="text-muted-foreground">请输入您的凭据以访问您的账户</p>
    </div>
    
    <form @submit.prevent="onSubmit" class="space-y-4">
      <div class="space-y-2">
        <label for="email" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
          邮箱
        </label>
        <input
          id="email"
          v-model="form.email"
          type="email"
          placeholder="m@example.com"
          class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="isLoading"
        />
      </div>
      
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
            密码
          </label>
          <a href="#" class="text-sm font-medium text-primary underline-offset-4 hover:underline">
            忘记密码?
          </a>
        </div>
        <input
          id="password"
          v-model="form.password"
          type="password"
          class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="isLoading"
        />
      </div>
      
      <button
        type="submit"
        :disabled="isLoading"
        class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2 w-full"
      >
        <span v-if="isLoading">登录中...</span>
        <span v-else>登录</span>
      </button>
    </form>
    
    <div class="relative">
      <div class="absolute inset-0 flex items-center">
        <span class="w-full border-t"></span>
      </div>
      <div class="relative flex justify-center text-xs uppercase">
        <span class="bg-background px-2 text-muted-foreground">或者继续</span>
      </div>
    </div>
    
    <div class="flex gap-2">
      <button
        type="button"
        class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 w-full"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2 h-4 w-4">
          <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path>
          <path d="M9 18c-4.51 2-5-2-7-2"></path>
        </svg>
        Github
      </button>
      <button
        type="button"
        class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-10 w-full"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2 h-4 w-4">
          <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"></path>
        </svg>
        Twitter
      </button>
    </div>
    
    <p class="px-8 text-center text-sm text-muted-foreground">
      还没有账户?
      <a href="#" class="underline underline-offset-4 hover:text-primary">
        注册
      </a>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Form {
  email: string
  password: string
}

const form = ref<Form>({
  email: '',
  password: ''
})

const isLoading = ref(false)

const onSubmit = () => {
  isLoading.value = true
  
  // 模拟登录请求
  setTimeout(() => {
    console.log('登录表单提交:', form.value)
    isLoading.value = false
    // 这里可以添加实际的登录逻辑
  }, 1500)
}
</script>