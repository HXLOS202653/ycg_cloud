import { QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import React, { useState } from 'react'
import { env } from '@/config/env'
import { createQueryClient } from '@/config/queryClient'

interface QueryProviderProps {
  children: React.ReactNode
}

/**
 * TanStack Query 提供者组件
 */
export const QueryProvider: React.FC<QueryProviderProps> = ({ children }) => {
  // 使用 useState 确保客户端和服务端的 QueryClient 实例一致
  const [queryClient] = useState(() => createQueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      {children}

      {/* 开发工具 - 仅在开发环境显示 */}
      {env.ENABLE_DEVTOOLS && <ReactQueryDevtools initialIsOpen={false} />}
    </QueryClientProvider>
  )
}

// 注意：useQueryClient 已经由 @tanstack/react-query 提供
// 这里不需要重新导出，直接从 @tanstack/react-query 导入即可
