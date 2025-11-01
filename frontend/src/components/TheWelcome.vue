<script setup lang="ts">
import WelcomeItem from './WelcomeItem.vue'
import DatabaseIcon from './icons/IconDatabase.vue'
import CommunityIcon from './icons/IconCommunity.vue'
import SupportIcon from './icons/IconSupport.vue'
</script>

<template>
  <WelcomeItem>
    <template #icon>
      <DatabaseIcon />
    </template>
    <template #heading>Database AI</template>

    欢迎使用 Database AI 系统。这是一个结合了 GoFiber 后端和 Vue 前端的全栈应用。
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <CommunityIcon />
    </template>
    <template #heading>社区支持</template>

    Vue 的社区非常活跃，你可以在论坛、Discord 或 GitHub 上获得帮助。
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <SupportIcon />
    </template>
    <template #heading>技术支持</template>

    我们的团队会为您提供专业的技术支持服务。
  </WelcomeItem>
</template>