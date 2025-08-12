import { QueryClient } from '@tanstack/react-query'
import { env } from './env'

/**
 * TanStack Query 配置
 */
export const queryClientConfig = {
  defaultOptions: {
    queries: {
      // 数据默认被认为是过期的时间（5分钟）
      staleTime: 5 * 60 * 1000,

      // 缓存时间（30分钟）
      gcTime: 30 * 60 * 1000,

      // 重试配置
      retry: (failureCount: number, error: any) => {
        // 4xx 错误不重试
        if (error?.status >= 400 && error?.status < 500) {
          return false
        }
        // 最多重试2次
        return failureCount < 2
      },

      // 重试延迟（指数退避）
      retryDelay: (attemptIndex: number) => Math.min(1000 * 2 ** attemptIndex, 30000),

      // 网络重连时重新获取
      refetchOnReconnect: true,

      // 窗口聚焦时重新获取（仅在生产环境）
      refetchOnWindowFocus: !env.ENABLE_DEVTOOLS,

      // 挂载时重新获取
      refetchOnMount: true,
    },

    mutations: {
      // 变更操作的重试配置
      retry: (failureCount: number, error: any) => {
        // 客户端错误不重试
        if (error?.status >= 400 && error?.status < 500) {
          return false
        }
        // 最多重试1次
        return failureCount < 1
      },

      // 变更重试延迟
      retryDelay: (attemptIndex: number) => Math.min(1000 * 2 ** attemptIndex, 5000),
    },
  },
}

/**
 * 创建 Query Client 实例
 */
export const createQueryClient = () => {
  return new QueryClient(queryClientConfig)
}

/**
 * 默认的 Query Client 实例
 */
export const queryClient = createQueryClient()

/**
 * Query Keys 工厂函数
 * 用于生成一致性的查询键
 */
export const queryKeys = {
  // 用户相关
  users: {
    all: ['users'] as const,
    lists: () => [...queryKeys.users.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.users.lists(), { filters }] as const,
    details: () => [...queryKeys.users.all, 'detail'] as const,
    detail: (id: string | number) => [...queryKeys.users.details(), id] as const,
    profile: () => [...queryKeys.users.all, 'profile'] as const,
  },

  // 文件相关
  files: {
    all: ['files'] as const,
    lists: () => [...queryKeys.files.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.files.lists(), { filters }] as const,
    details: () => [...queryKeys.files.all, 'detail'] as const,
    detail: (id: string | number) => [...queryKeys.files.details(), id] as const,
    tree: (folderId?: string | number) => [...queryKeys.files.all, 'tree', folderId] as const,
    search: (keyword: string) => [...queryKeys.files.all, 'search', keyword] as const,
  },

  // 文件夹相关
  folders: {
    all: ['folders'] as const,
    lists: () => [...queryKeys.folders.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.folders.lists(), { filters }] as const,
    details: () => [...queryKeys.folders.all, 'detail'] as const,
    detail: (id: string | number) => [...queryKeys.folders.details(), id] as const,
    tree: (parentId?: string | number) => [...queryKeys.folders.all, 'tree', parentId] as const,
  },

  // 上传任务相关
  uploads: {
    all: ['uploads'] as const,
    lists: () => [...queryKeys.uploads.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.uploads.lists(), { filters }] as const,
    details: () => [...queryKeys.uploads.all, 'detail'] as const,
    detail: (id: string | number) => [...queryKeys.uploads.details(), id] as const,
    progress: (id: string | number) => [...queryKeys.uploads.all, 'progress', id] as const,
  },

  // 系统统计相关
  stats: {
    all: ['stats'] as const,
    overview: () => [...queryKeys.stats.all, 'overview'] as const,
    storage: () => [...queryKeys.stats.all, 'storage'] as const,
    activities: () => [...queryKeys.stats.all, 'activities'] as const,
  },

  // 设置相关
  settings: {
    all: ['settings'] as const,
    user: () => [...queryKeys.settings.all, 'user'] as const,
    system: () => [...queryKeys.settings.all, 'system'] as const,
  },
} as const

/**
 * 查询键类型
 */
export type QueryKeys = typeof queryKeys

/**
 * 无效化查询的工具函数
 */
export const invalidateQueries = {
  // 无效化所有用户相关查询
  users: () => queryClient.invalidateQueries({ queryKey: queryKeys.users.all }),

  // 无效化所有文件相关查询
  files: () => queryClient.invalidateQueries({ queryKey: queryKeys.files.all }),

  // 无效化所有文件夹相关查询
  folders: () => queryClient.invalidateQueries({ queryKey: queryKeys.folders.all }),

  // 无效化所有上传相关查询
  uploads: () => queryClient.invalidateQueries({ queryKey: queryKeys.uploads.all }),

  // 无效化所有统计相关查询
  stats: () => queryClient.invalidateQueries({ queryKey: queryKeys.stats.all }),

  // 无效化所有设置相关查询
  settings: () => queryClient.invalidateQueries({ queryKey: queryKeys.settings.all }),

  // 无效化所有查询
  all: () => queryClient.invalidateQueries(),
}

/**
 * 预取查询的工具函数
 */
export const prefetchQueries = {
  // 预取用户列表
  userList: (filters: Record<string, any> = {}) =>
    queryClient.prefetchQuery({
      queryKey: queryKeys.users.list(filters),
      queryFn: () => {
        // 这里会在后续创建API函数时替换
        return Promise.resolve([])
      },
      staleTime: 5 * 60 * 1000,
    }),

  // 预取文件列表
  fileList: (filters: Record<string, any> = {}) =>
    queryClient.prefetchQuery({
      queryKey: queryKeys.files.list(filters),
      queryFn: () => {
        // 这里会在后续创建API函数时替换
        return Promise.resolve([])
      },
      staleTime: 5 * 60 * 1000,
    }),
}
