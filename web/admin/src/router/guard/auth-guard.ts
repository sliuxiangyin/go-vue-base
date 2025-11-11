import type { Router } from 'vue-router'

import { storeToRefs } from 'pinia'

import pinia from '@/plugins/pinia/setup'
import { useAuthStore } from '@/stores/auth'

export function authGuard(router: Router) {
  const authStore = useAuthStore(pinia)
  const { isLogin } = storeToRefs(authStore)

  // 初始化认证状态
  authStore.initAuth()

  router.beforeEach((to, _from, next) => {
    console.log(to.meta.auth,unref(isLogin))
    // 如果路由需要认证且未登录，跳转到登录页
    if (!unref(isLogin) && to.path !== '/auth/sign-in') {
      next('/auth/sign-in')
      return
    }

    // 如果已登录且访问登录页，跳转到首页
    if (unref(isLogin) && to.path === '/auth/sign-in') {
      next('/dashboard')
      return
    }

    next()
  })
}
