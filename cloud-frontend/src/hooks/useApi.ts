import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryOptions, UseMutationOptions } from '@tanstack/react-query'
import { message } from 'antd'
import { queryKeys, invalidateQueries } from '@/config/queryClient'
import { apiClient } from '@/services/api'
import type { ApiError, PaginatedResponse, PaginationParams } from '@/services/api'

/**
 * 用户数据类型
 */
export interface User {
  id: string
  username: string
  email: string
  avatar?: string
  role: string
  status: 'active' | 'inactive'
  createdAt: string
  updatedAt: string
}

/**
 * 文件数据类型
 */
export interface FileItem {
  id: string
  name: string
  size: number
  type: string
  mimeType: string
  path: string
  folderId?: string
  createdAt: string
  updatedAt: string
  downloadUrl?: string
  previewUrl?: string
}

/**
 * 文件夹数据类型
 */
export interface Folder {
  id: string
  name: string
  path: string
  parentId?: string
  fileCount: number
  size: number
  createdAt: string
  updatedAt: string
}

/**
 * 系统统计数据类型
 */
export interface SystemStats {
  totalFiles: number
  totalSize: number
  todayUploads: number
  todayDownloads: number
  storageUsed: number
  storageTotal: number
}

/**
 * 用户相关 API Hooks
 */
export const useUsers = {
  /**
   * 获取用户列表
   */
  list: (
    params: PaginationParams & { keyword?: string } = {},
    options?: Omit<UseQueryOptions<PaginatedResponse<User>, ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.users.list(params),
      queryFn: ({ signal }) => apiClient.get<PaginatedResponse<User>>('/users', params, signal),
      ...options,
    })
  },

  /**
   * 获取用户详情
   */
  detail: (id: string, options?: Omit<UseQueryOptions<User, ApiError>, 'queryKey' | 'queryFn'>) => {
    return useQuery({
      queryKey: queryKeys.users.detail(id),
      queryFn: ({ signal }) => apiClient.get<User>(`/users/${id}`, undefined, signal),
      enabled: !!id,
      ...options,
    })
  },

  /**
   * 获取当前用户信息
   */
  profile: (options?: Omit<UseQueryOptions<User, ApiError>, 'queryKey' | 'queryFn'>) => {
    return useQuery({
      queryKey: queryKeys.users.profile(),
      queryFn: ({ signal }) => apiClient.get<User>('/users/profile', undefined, signal),
      ...options,
    })
  },

  /**
   * 更新用户信息
   */
  update: (
    options?: Omit<
      UseMutationOptions<User, ApiError, { id: string; data: Partial<User> }>,
      'mutationFn'
    >
  ) => {
    const queryClient = useQueryClient()

    return useMutation({
      mutationFn: ({ id, data }) => apiClient.put<User>(`/users/${id}`, data),
      onSuccess: (data, { id }) => {
        // 更新缓存
        queryClient.setQueryData(queryKeys.users.detail(id), data)
        queryClient.setQueryData(queryKeys.users.profile(), data)

        // 无效化列表查询
        invalidateQueries.users()

        message.success('用户信息更新成功')
      },
      onError: (error) => {
        message.error(error.message || '更新失败')
      },
      ...options,
    })
  },
}

/**
 * 文件相关 API Hooks
 */
export const useFiles = {
  /**
   * 获取文件列表
   */
  list: (
    params: PaginationParams & { folderId?: string; keyword?: string } = {},
    options?: Omit<UseQueryOptions<PaginatedResponse<FileItem>, ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.files.list(params),
      queryFn: ({ signal }) => apiClient.get<PaginatedResponse<FileItem>>('/files', params, signal),
      ...options,
    })
  },

  /**
   * 获取文件详情
   */
  detail: (
    id: string,
    options?: Omit<UseQueryOptions<FileItem, ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.files.detail(id),
      queryFn: ({ signal }) => apiClient.get<FileItem>(`/files/${id}`, undefined, signal),
      enabled: !!id,
      ...options,
    })
  },

  /**
   * 搜索文件
   */
  search: (
    keyword: string,
    options?: Omit<UseQueryOptions<FileItem[], ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.files.search(keyword),
      queryFn: ({ signal }) => apiClient.get<FileItem[]>('/files/search', { keyword }, signal),
      enabled: !!keyword && keyword.length > 1,
      ...options,
    })
  },

  /**
   * 上传文件
   */
  upload: (
    options?: Omit<
      UseMutationOptions<
        FileItem,
        ApiError,
        { file: File; folderId?: string; onProgress?: (progress: number) => void }
      >,
      'mutationFn'
    >
  ) => {
    return useMutation({
      mutationFn: ({ file, folderId, onProgress }) => {
        const formData = new FormData()
        formData.append('file', file)
        if (folderId) {
          formData.append('folderId', folderId)
        }

        return apiClient.upload<FileItem>('/files/upload', formData, onProgress)
      },
      onSuccess: () => {
        // 无效化文件列表
        invalidateQueries.files()
        invalidateQueries.stats()

        message.success('文件上传成功')
      },
      onError: (error) => {
        message.error(error.message || '上传失败')
      },
      ...options,
    })
  },

  /**
   * 删除文件
   */
  delete: (options?: Omit<UseMutationOptions<void, ApiError, string>, 'mutationFn'>) => {
    const queryClient = useQueryClient()

    return useMutation({
      mutationFn: async (id: string) => {
        await apiClient.delete(`/files/${id}`)
      },
      onSuccess: (_, id) => {
        // 移除缓存
        queryClient.removeQueries({ queryKey: queryKeys.files.detail(id) })

        // 无效化列表查询
        invalidateQueries.files()
        invalidateQueries.stats()

        message.success('文件删除成功')
      },
      onError: (error) => {
        message.error(error.message || '删除失败')
      },
      ...options,
    })
  },
}

/**
 * 文件夹相关 API Hooks
 */
export const useFolders = {
  /**
   * 获取文件夹列表
   */
  list: (
    params: PaginationParams & { parentId?: string } = {},
    options?: Omit<UseQueryOptions<PaginatedResponse<Folder>, ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.folders.list(params),
      queryFn: ({ signal }) => apiClient.get<PaginatedResponse<Folder>>('/folders', params, signal),
      ...options,
    })
  },

  /**
   * 获取文件夹树
   */
  tree: (
    parentId?: string,
    options?: Omit<UseQueryOptions<Folder[], ApiError>, 'queryKey' | 'queryFn'>
  ) => {
    return useQuery({
      queryKey: queryKeys.folders.tree(parentId),
      queryFn: ({ signal }) => apiClient.get<Folder[]>('/folders/tree', { parentId }, signal),
      ...options,
    })
  },

  /**
   * 创建文件夹
   */
  create: (
    options?: Omit<
      UseMutationOptions<Folder, ApiError, { name: string; parentId?: string }>,
      'mutationFn'
    >
  ) => {
    return useMutation({
      mutationFn: (data) => apiClient.post<Folder>('/folders', data),
      onSuccess: () => {
        // 无效化文件夹查询
        invalidateQueries.folders()

        message.success('文件夹创建成功')
      },
      onError: (error) => {
        message.error(error.message || '创建失败')
      },
      ...options,
    })
  },
}

/**
 * 系统统计相关 API Hooks
 */
export const useStats = {
  /**
   * 获取系统概览统计
   */
  overview: (options?: Omit<UseQueryOptions<SystemStats, ApiError>, 'queryKey' | 'queryFn'>) => {
    return useQuery({
      queryKey: queryKeys.stats.overview(),
      queryFn: ({ signal }) => apiClient.get<SystemStats>('/stats/overview', undefined, signal),
      refetchInterval: 30000, // 30秒自动刷新
      ...options,
    })
  },

  /**
   * 获取存储统计
   */
  storage: (
    options?: Omit<
      UseQueryOptions<{ used: number; total: number; percentage: number }, ApiError>,
      'queryKey' | 'queryFn'
    >
  ) => {
    return useQuery({
      queryKey: queryKeys.stats.storage(),
      queryFn: ({ signal }) => apiClient.get('/stats/storage', undefined, signal),
      ...options,
    })
  },
}

/**
 * 通用错误处理Hook
 */
export const useErrorHandler = () => {
  return {
    handleError: (error: ApiError | Error) => {
      console.error('API Error:', error)

      if ('code' in error) {
        // API错误
        switch (error.code) {
          case 401:
            message.error('请先登录')
            // 可以添加跳转到登录页的逻辑
            break
          case 403:
            message.error('权限不足')
            break
          case 404:
            message.error('资源不存在')
            break
          case 500:
            message.error('服务器错误，请稍后重试')
            break
          default:
            message.error(error.message || '操作失败')
        }
      } else {
        // 网络或其他错误
        message.error(error.message || '网络错误')
      }
    },
  }
}
