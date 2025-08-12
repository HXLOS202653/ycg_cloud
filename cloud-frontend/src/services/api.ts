import { env } from '@/config/env'

/**
 * API 响应类型
 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  success: boolean
  timestamp: number
}

/**
 * API 错误类型
 */
export interface ApiError {
  code: number
  message: string
  details?: any
}

/**
 * 分页响应类型
 */
export interface PaginatedResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

/**
 * 分页请求参数
 */
export interface PaginationParams {
  page?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

/**
 * API 客户端类
 */
class ApiClient {
  private baseURL: string
  private timeout: number
  private defaultHeaders: Record<string, string>

  constructor() {
    this.baseURL = env.API_BASE_URL
    this.timeout = env.API_TIMEOUT
    this.defaultHeaders = {
      'Content-Type': 'application/json',
    }
  }

  /**
   * 设置认证令牌
   */
  setAuthToken(token: string | null) {
    if (token) {
      this.defaultHeaders['Authorization'] = `Bearer ${token}`
    } else {
      delete this.defaultHeaders['Authorization']
    }
  }

  /**
   * 构建请求URL
   */
  private buildURL(endpoint: string, params?: Record<string, any>): string {
    const url = new URL(endpoint, this.baseURL)

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value))
        }
      })
    }

    return url.toString()
  }

  /**
   * 处理响应
   */
  private async handleResponse<T>(response: Response): Promise<T> {
    const contentType = response.headers.get('content-type')
    const isJson = contentType?.includes('application/json')

    let data: any
    try {
      data = isJson ? await response.json() : await response.text()
    } catch (error) {
      throw new Error('Failed to parse response')
    }

    if (!response.ok) {
      const apiError: ApiError = {
        code: response.status,
        message: data?.message || response.statusText || 'Unknown error',
        details: data?.details || data,
      }
      throw apiError
    }

    // 如果是标准API响应格式，返回data字段
    if (data && typeof data === 'object' && 'data' in data) {
      return data.data
    }

    return data
  }

  /**
   * 通用请求方法
   */
  private async request<T>(
    method: string,
    endpoint: string,
    options: {
      params?: Record<string, any>
      data?: any
      headers?: Record<string, string>
      signal?: AbortSignal
    } = {}
  ): Promise<T> {
    const { params, data, headers = {}, signal } = options

    const url = this.buildURL(endpoint, params)

    const requestInit: RequestInit = {
      method,
      headers: {
        ...this.defaultHeaders,
        ...headers,
      },
      signal,
    }

    // 添加请求体
    if (data) {
      if (data instanceof FormData) {
        // FormData 时不设置 Content-Type，让浏览器自动设置
        const headers = requestInit.headers as Record<string, string>
        delete headers['Content-Type']
        requestInit.body = data
      } else {
        requestInit.body = JSON.stringify(data)
      }
    }

    // 设置超时
    let timeoutId: number | undefined
    if (!signal?.aborted) {
      timeoutId = setTimeout(() => {
        // 超时处理：直接忽略，让fetch自然超时
      }, this.timeout) as unknown as number
    }

    try {
      const response = await fetch(url, requestInit)
      if (timeoutId) clearTimeout(timeoutId)
      return await this.handleResponse<T>(response)
    } catch (error) {
      if (timeoutId) clearTimeout(timeoutId)

      // 处理网络错误
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Network error: Please check your internet connection')
      }

      // 处理超时
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error('Request timeout')
      }

      throw error
    }
  }

  /**
   * GET 请求
   */
  async get<T>(endpoint: string, params?: Record<string, any>, signal?: AbortSignal): Promise<T> {
    return this.request<T>('GET', endpoint, { params, signal })
  }

  /**
   * POST 请求
   */
  async post<T>(endpoint: string, data?: any, signal?: AbortSignal): Promise<T> {
    return this.request<T>('POST', endpoint, { data, signal })
  }

  /**
   * PUT 请求
   */
  async put<T>(endpoint: string, data?: any, signal?: AbortSignal): Promise<T> {
    return this.request<T>('PUT', endpoint, { data, signal })
  }

  /**
   * PATCH 请求
   */
  async patch<T>(endpoint: string, data?: any, signal?: AbortSignal): Promise<T> {
    return this.request<T>('PATCH', endpoint, { data, signal })
  }

  /**
   * DELETE 请求
   */
  async delete<T>(endpoint: string, signal?: AbortSignal): Promise<T> {
    return this.request<T>('DELETE', endpoint, { signal })
  }

  /**
   * 上传文件
   */
  async upload<T>(
    endpoint: string,
    formData: FormData,
    onProgress?: (progress: number) => void,
    signal?: AbortSignal
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()

      // 设置超时
      xhr.timeout = this.timeout

      // 上传进度
      if (onProgress) {
        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            const progress = Math.round((event.loaded / event.total) * 100)
            onProgress(progress)
          }
        }
      }

      // 完成处理
      xhr.onload = () => {
        try {
          const response = JSON.parse(xhr.responseText)
          if (xhr.status >= 200 && xhr.status < 300) {
            resolve(response.data || response)
          } else {
            reject({
              code: xhr.status,
              message: response.message || 'Upload failed',
              details: response,
            })
          }
        } catch (error) {
          reject({
            code: xhr.status,
            message: 'Failed to parse upload response',
            details: error,
          })
        }
      }

      // 错误处理
      xhr.onerror = () => {
        reject({
          code: 0,
          message: 'Network error during upload',
        })
      }

      // 超时处理
      xhr.ontimeout = () => {
        reject({
          code: 0,
          message: 'Upload timeout',
        })
      }

      // 取消处理
      if (signal) {
        signal.addEventListener('abort', () => {
          xhr.abort()
          reject({
            code: 0,
            message: 'Upload cancelled',
          })
        })
      }

      // 设置请求头
      Object.entries(this.defaultHeaders).forEach(([key, value]) => {
        if (key !== 'Content-Type') {
          // FormData 时不设置 Content-Type
          xhr.setRequestHeader(key, value)
        }
      })

      // 发送请求
      const url = this.buildURL(endpoint)
      xhr.open('POST', url, true)
      xhr.send(formData)
    })
  }
}

/**
 * API 客户端实例
 */
export const apiClient = new ApiClient()

/**
 * 创建带有取消令牌的请求函数
 */
export const createRequestWithCancel = () => {
  const controller = new AbortController()

  return {
    signal: controller.signal,
    cancel: () => controller.abort(),
  }
}
