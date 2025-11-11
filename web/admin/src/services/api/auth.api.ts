import type { AxiosResponse } from 'axios'
import { useAxios } from '@/composables/use-axios'

const { axiosInstance } = useAxios()

// 登录请求参数
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应数据
export interface LoginResponse {
  code: number
  message: string
  data: {
    token: string
    user: {
      id: number
      username: string
      email: string
      nickname: string
      avatar: string
    }
  }
}

// 用户信息响应
export interface UserInfoResponse {
  code: number
  message: string
  data: {
    id: number
    username: string
    email: string
    nickname: string
    avatar: string
    last_login: string
  }
}

// 修改密码请求参数
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 通用响应
export interface CommonResponse {
  code: number
  message: string
}

/**
 * 登录
 */
export function login(data: LoginRequest): Promise<AxiosResponse<LoginResponse>> {
  return axiosInstance.post('/admin/auth/login', data)
}

/**
 * 登出
 */
export function logout(): Promise<AxiosResponse<CommonResponse>> {
  return axiosInstance.post('/admin/auth/logout')
}

/**
 * 获取用户信息
 */
export function getUserInfo(): Promise<AxiosResponse<UserInfoResponse>> {
  return axiosInstance.get('/admin/auth/info')
}

/**
 * 修改密码
 */
export function changePassword(data: ChangePasswordRequest): Promise<AxiosResponse<CommonResponse>> {
  return axiosInstance.post('/admin/auth/change-password', data)
}
