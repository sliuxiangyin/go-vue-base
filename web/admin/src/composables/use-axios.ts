import type { AxiosError } from 'axios'

import axios from 'axios'
import Cookies from 'universal-cookie'
import { toast } from 'vue-sonner'

import env from '@/utils/env'

const cookies = new Cookies()
const TOKEN_KEY = 'admin_token'

export function useAxios() {
  const axiosInstance = axios.create({
    baseURL: env.VITE_SERVER_API_URL + env.VITE_SERVER_API_PREFIX,
    timeout: env.VITE_SERVER_API_TIMEOUT,
  })

  axiosInstance.interceptors.request.use((config) => {
    // 添加 token 到请求头
    const token = cookies.get(TOKEN_KEY)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  }, (error) => {
    return Promise.reject(error)
  })

  axiosInstance.interceptors.response.use((response) => {
    return response
  }, (error: AxiosError<any>) => {
    // 处理错误响应
    if (error.response) {
      const { status, data } = error.response
      
      // 401 未授权，清除 token 并跳转登录
      if (status === 401) {
        cookies.remove(TOKEN_KEY, { path: '/' })
        toast.error(data?.message || '认证失败，请重新登录')
        // 跳转到登录页
        if (window.location.pathname !== '/auth/sign-in') {
          window.location.href = '/auth/sign-in'
        }
      }
      // 其他错误显示提示
      else if (data?.message) {
        toast.error(data.message)
      }
    }
    else if (error.request) {
      toast.error('网络错误，请检查网络连接')
    }
    else {
      toast.error('请求失败，请稍后重试')
    }
    
    return Promise.reject(error)
  })

  return {
    axiosInstance,
  }
}
