import type { AxiosError } from 'axios'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { toValue, type MaybeRefOrGetter } from 'vue'
import { useAxios } from '@/composables/use-axios'
import type { IResponse } from '../types/response.type'

const { axiosInstance } = useAxios()

// ========== 类型定义 ==========

export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  status: number
  is_system: boolean
  created_at: string
  updated_at: string
  permissions?: Permission[]
}

export interface Permission {
  id: number
  name: string
  display_name: string
  description: string
  category: string
  type: 'backend' | 'frontend'
  path: string
  icon: string
  parent_id?: number
  sort: number
  status: number
  created_at: string
  updated_at: string
  children?: Permission[]
}

export interface AdminUser {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  status: number
  last_login: string | null
  created_at: string
  updated_at: string
  roles?: Role[]
}

export interface ListResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface CreateRoleRequest {
  name: string
  display_name: string
  description?: string
  status: number
}

export interface UpdateRoleRequest {
  display_name: string
  description?: string
  status: number
}

export interface CreatePermissionRequest {
  name: string
  display_name: string
  description?: string
  category: string
  type: 'backend' | 'frontend'
  path?: string
  icon?: string
  parent_id?: number
  sort?: number
  status: number
}

export interface UpdatePermissionRequest {
  display_name: string
  description?: string
  category: string
  type: 'backend' | 'frontend'
  path?: string
  icon?: string
  parent_id?: number
  sort?: number
  status: number
}

export interface AssignPermissionsRequest {
  permission_ids: number[]
}

export interface AssignRolesRequest {
  role_ids: number[]
}

export interface CreateAdminUserRequest {
  username: string
  password: string
  email?: string
  nickname?: string
  avatar?: string
  status: number
}

export interface UpdateAdminUserRequest {
  email?: string
  nickname?: string
  avatar?: string
  status: number
}

export interface UpdateAdminUserPasswordRequest {
  password: string
}

// ========== 角色管理 API ==========

export function useGetRolesQuery(params: MaybeRefOrGetter<{ page: number, page_size: number, name?: string, status?: number }>) {
  return useQuery<IResponse<ListResponse<Role>>, AxiosError>({
    queryKey: ['roles', params],
    queryFn: async () => {
      const resolvedParams = toValue(params)
      const response = await axiosInstance.get('/admin/roles', { params: resolvedParams })
      return response.data
    },
  })
}

export function useGetRoleQuery(id: number) {
  return useQuery<IResponse<Role>, AxiosError>({
    queryKey: ['role', id],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/roles/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreateRoleMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Role>, AxiosError, CreateRoleRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post('/admin/roles', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
    },
  })
}

export function useUpdateRoleMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Role>, AxiosError, UpdateRoleRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/roles/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
      queryClient.invalidateQueries({ queryKey: ['role', id] })
    },
  })
}

export function useDeleteRoleMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, number>({
    mutationFn: async (id) => {
      const response = await axiosInstance.delete(`/admin/roles/${id}`)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['roles'] })
    },
  })
}

// ========== 权限管理 API ==========

export function useGetPermissionsQuery(params: MaybeRefOrGetter<{ page: number, page_size: number, category?: string, type?: string, status?: number }>) {
  return useQuery<IResponse<ListResponse<Permission>>, AxiosError>({
    queryKey: ['permissions', params],
    queryFn: async () => {
      const resolvedParams = toValue(params)
      const response = await axiosInstance.get('/admin/permissions', { params: resolvedParams })
      return response.data
    },
  })
}

export function useGetAllPermissionsQuery() {
  return useQuery<IResponse<Permission[]>, AxiosError>({
    queryKey: ['all-permissions'],
    queryFn: async () => {
      const response = await axiosInstance.get('/admin/permissions/all')
      return response.data
    },
  })
}

export function useGetPermissionQuery(id: number) {
  return useQuery<IResponse<Permission>, AxiosError>({
    queryKey: ['permission', id],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/permissions/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreatePermissionMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Permission>, AxiosError, CreatePermissionRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post('/admin/permissions', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      queryClient.invalidateQueries({ queryKey: ['all-permissions'] })
    },
  })
}

export function useUpdatePermissionMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<Permission>, AxiosError, UpdatePermissionRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/permissions/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      queryClient.invalidateQueries({ queryKey: ['permission', id] })
      queryClient.invalidateQueries({ queryKey: ['all-permissions'] })
    },
  })
}

export function useDeletePermissionMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, number>({
    mutationFn: async (id) => {
      const response = await axiosInstance.delete(`/admin/permissions/${id}`)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['permissions'] })
      queryClient.invalidateQueries({ queryKey: ['all-permissions'] })
    },
  })
}

// ========== 角色权限关联 API ==========

export function useGetRolePermissionsQuery(roleId: number) {
  return useQuery<IResponse<Permission[]>, AxiosError>({
    queryKey: ['role-permissions', roleId],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/roles/${roleId}/permissions`)
      return response.data
    },
    enabled: !!roleId,
  })
}

export function useAssignPermissionsToRoleMutation(roleId: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, AssignPermissionsRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post(`/admin/roles/${roleId}/permissions`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['role-permissions', roleId] })
      queryClient.invalidateQueries({ queryKey: ['role', roleId] })
    },
  })
}

// ========== 用户角色关联 API ==========

export function useGetUserRolesQuery(userId: number) {
  return useQuery<IResponse<Role[]>, AxiosError>({
    queryKey: ['user-roles', userId],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/users/${userId}/roles`)
      return response.data
    },
    enabled: !!userId,
  })
}

export function useAssignRolesToUserMutation(userId: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, AssignRolesRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post(`/admin/users/${userId}/roles`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-roles', userId] })
    },
  })
}

export function useGetUserPermissionsQuery(userId: number) {
  return useQuery<IResponse<Permission[]>, AxiosError>({
    queryKey: ['user-permissions', userId],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/users/${userId}/permissions`)
      return response.data
    },
    enabled: !!userId,
  })
}

// 获取当前用户的前端菜单权限
export function useGetCurrentUserFrontendPermissionsQuery() {
  const userIdRef = ref<number | null>(null)
  
  // 从 localStorage 或 store 获取用户ID
  if (typeof window !== 'undefined') {
    const authStore = localStorage.getItem('auth')
    if (authStore) {
      try {
        const parsed = JSON.parse(authStore)
        userIdRef.value = parsed.user?.id || null
      } catch (e) {
        console.error('Failed to parse auth store:', e)
      }
    }
  }

  return useQuery<IResponse<Permission[]>, AxiosError>({
    queryKey: ['current-user-frontend-permissions'],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/users/${userIdRef.value}/permissions?type=frontend`)
      return response.data
    },
    enabled: !!userIdRef.value,
  })
}

// ========== 管理员管理 API ==========

export function useGetAdminUsersQuery(params: MaybeRefOrGetter<{ page: number, page_size: number, keyword?: string, status?: number }>) {
  return useQuery<IResponse<ListResponse<AdminUser>>, AxiosError>({
    queryKey: ['admin-users', params],
    queryFn: async () => {
      const resolvedParams = toValue(params)
      const response = await axiosInstance.get('/admin/users', { params: resolvedParams })
      return response.data
    },
  })
}

export function useGetAdminUserQuery(id: number) {
  return useQuery<IResponse<AdminUser>, AxiosError>({
    queryKey: ['admin-user', id],
    queryFn: async () => {
      const response = await axiosInstance.get(`/admin/users/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

export function useCreateAdminUserMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<AdminUser>, AxiosError, CreateAdminUserRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.post('/admin/users', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
    },
  })
}

export function useUpdateAdminUserMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<AdminUser>, AxiosError, UpdateAdminUserRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/users/${id}`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      queryClient.invalidateQueries({ queryKey: ['admin-user', id] })
    },
  })
}

export function useUpdateAdminUserPasswordMutation(id: number) {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, UpdateAdminUserPasswordRequest>({
    mutationFn: async (data) => {
      const response = await axiosInstance.put(`/admin/users/${id}/password`, data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-user', id] })
    },
  })
}

export function useDeleteAdminUserMutation() {
  const queryClient = useQueryClient()

  return useMutation<IResponse<void>, AxiosError, number>({
    mutationFn: async (id) => {
      const response = await axiosInstance.delete(`/admin/users/${id}`)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
    },
  })
}

// 获取所有角色(用于管理员分配角色)
export function useGetAllRolesQuery() {
  return useQuery<IResponse<Role[]>, AxiosError>({
    queryKey: ['all-roles'],
    queryFn: async () => {
      const response = await axiosInstance.get('/admin/roles', { params: { page: 1, page_size: 1000, status: 1 } })
      return { ...response.data, data: response.data.data.list }
    },
  })
}
