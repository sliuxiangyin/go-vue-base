import { defineStore } from 'pinia'
import Cookies from 'universal-cookie'
import type { LoginRequest } from '@/services/api/auth.api'
import { getUserInfo, login, logout } from '@/services/api/auth.api'

const cookies = new Cookies()
const TOKEN_KEY = 'admin_token'

interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(cookies.get(TOKEN_KEY) || '')
  const userInfo = ref<UserInfo | null>(null)
  const isLogin = computed(() => !!token.value)

  /**
   * 登录
   */
  async function handleLogin(params: LoginRequest) {
    try {
      const { data } = await login(params)
      if (data.code === 0) {
        token.value = data.data.token
        userInfo.value = data.data.user
        // 保存 token 到 cookie，有效期 7 天
        cookies.set(TOKEN_KEY, data.data.token, { path: '/', maxAge: 7 * 24 * 60 * 60 })
        return { success: true, message: data.message }
      }
      else {
        return { success: false, message: data.message }
      }
    }
    catch (error: any) {
      const message = error.response?.data?.message || '登录失败，请稍后重试'
      return { success: false, message }
    }
  }

  /**
   * 登出
   */
  async function handleLogout() {
    try {
      await logout()
    }
    catch (error) {
      console.error('登出失败:', error)
    }
    finally {
      token.value = ''
      userInfo.value = null
      cookies.remove(TOKEN_KEY, { path: '/' })
    }
  }

  /**
   * 获取用户信息
   */
  async function fetchUserInfo() {
    try {
      const { data } = await getUserInfo()
      if (data.code === 0) {
        userInfo.value = data.data
        return true
      }
      return false
    }
    catch (error) {
      console.error('获取用户信息失败:', error)
      return false
    }
  }

  /**
   * 初始化认证状态
   */
  async function initAuth() {
    if (token.value) {
      await fetchUserInfo()
    }
  }

  return {
    token,
    userInfo,
    isLogin,
    handleLogin,
    handleLogout,
    fetchUserInfo,
    initAuth,
  }
})
