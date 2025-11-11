---
trigger: model_decision
description: web 增加菜单的位置web\admin\src\composables\use-sidebar.ts
---
import { useSidebar } from '@/composables/use-sidebar'
export const { sidebar, openSidebar, closeSidebar } = useSidebar()
export default defineComponent({
  name: 'Sidebar',
  setup() {
    return {
      sidebar,
      openSidebar,
      closeSidebar,
    }
  },
})
<style scoped>
@import './sidebar.css'
</style>
<template>
  <div class="sidebar" :class="{ 'sidebar-open': sidebar }">
    <div class="sidebar-content">
      <div class="sidebar-logo">
        <img src="/logo.svg" alt="Logo" />
        <span>Admin</span>
      </div>
      <div class="sidebar-menu">
        <a href="#" class="sidebar-menu-item" @click="openSidebar">
          <i class="ri-menu-line"></i>
          <span>Menu</span>
        </a>
        <a href="#" class="sidebar-menu-item" @click="closeSidebar">
          <i class="ri-close-line"></i>
          <span>Close</span>
        </a>
      </div>
    </div>
  </div>
</template>
