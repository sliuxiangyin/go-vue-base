<template>
  <div class="home">
    <h1 class="text-3xl font-bold mb-4">Home Page</h1>
    <p class="mb-4">Welcome to the Database AI application!</p>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-2">Key-Value Store</h2>
        <p>Store and retrieve key-value pairs using our API.</p>
        <router-link to="/api-test" class="mt-2 inline-block text-blue-500 hover:text-blue-700">
          访问测试页面 →
        </router-link>
      </div>
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-2">Database Account Management</h2>
        <p>Manage your database accounts with ease.</p>
        <router-link to="/db-account" class="mt-2 inline-block text-blue-500 hover:text-blue-700">
          管理数据库账号 →
        </router-link>
      </div>
    </div>
    <Button
        variant="outline" @click="() => {
     toast.warning('wqeqe')
    }"
    >
      Add to calendar
    </Button>
  </div>
</template>

<script setup lang="ts">
import { toast } from 'vue-sonner'
// Home page component
import {Button} from "@/components/ui/button";
</script>