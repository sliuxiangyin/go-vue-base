<script setup lang="ts">
import { storeToRefs } from 'pinia'

import Loading from '@/components/loading.vue'
import { Toaster } from '@/components/ui/sonner'
import { THEMES } from '@/constants/themes'
import { useThemeStore } from '@/stores/theme'
import { useServerEvents } from '@/composables/use-server-events'

const themeStore = useThemeStore()
const { theme: t, radius } = storeToRefs(themeStore)

// 在应用级别建立 SSE 连接，避免路由切换时重复连接
useServerEvents()

watchEffect(() => {
  document.documentElement.classList.remove(...THEMES.map(theme => `theme-${theme}`))
  document.documentElement.classList.add(`theme-${t.value}`)
  document.documentElement.style.setProperty('--radius', `${radius.value}rem`)
})
</script>

<template>
  <Toaster />
  <!-- <VueQueryDevtools /> -->

  <Suspense>
    <router-view v-slot="{ Component, route }">
      <component :is="Component" :key="route" />
    </router-view>

    <template #fallback>
      <Loading />
    </template>
  </Suspense>
</template>
