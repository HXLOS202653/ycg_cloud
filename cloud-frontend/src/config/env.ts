/**
 * 环境配置管理
 */

// 环境变量类型
interface EnvConfig {
  APP_TITLE: string
  APP_VERSION: string
  API_BASE_URL: string
  API_TIMEOUT: number
  UPLOAD_MAX_SIZE: number
  UPLOAD_CHUNK_SIZE: number
  WS_URL: string
  ENABLE_MOCK: boolean
  ENABLE_PWA: boolean
  ENABLE_DEVTOOLS: boolean
  SHOW_CONSOLE_LOG: boolean
}

// 默认配置
const defaultConfig: EnvConfig = {
  APP_TITLE: '网络云盘系统',
  APP_VERSION: '1.0.0',
  API_BASE_URL: 'http://localhost:8080/api',
  API_TIMEOUT: 10000,
  UPLOAD_MAX_SIZE: 100 * 1024 * 1024, // 100MB
  UPLOAD_CHUNK_SIZE: 2 * 1024 * 1024, // 2MB
  WS_URL: 'ws://localhost:8080/ws',
  ENABLE_MOCK: false,
  ENABLE_PWA: true,
  ENABLE_DEVTOOLS: import.meta.env.DEV,
  SHOW_CONSOLE_LOG: import.meta.env.DEV,
}

// 获取环境变量值
function getEnvValue<T>(key: string, defaultValue: T, parser?: (value: string) => T): T {
  const envKey = `VITE_${key}`
  const value = import.meta.env[envKey]

  if (value === undefined) {
    return defaultValue
  }

  if (parser) {
    try {
      return parser(value)
    } catch {
      return defaultValue
    }
  }

  return value as T
}

// 解析布尔值
const parseBoolean = (value: string): boolean => {
  return value.toLowerCase() === 'true'
}

// 解析数字
const parseNumber = (value: string): number => {
  const num = parseInt(value, 10)
  if (isNaN(num)) {
    throw new Error(`Invalid number: ${value}`)
  }
  return num
}

// 环境配置
export const env: EnvConfig = {
  APP_TITLE: getEnvValue('APP_TITLE', defaultConfig.APP_TITLE),
  APP_VERSION: getEnvValue('APP_VERSION', defaultConfig.APP_VERSION),
  API_BASE_URL: getEnvValue('API_BASE_URL', defaultConfig.API_BASE_URL),
  API_TIMEOUT: getEnvValue('API_TIMEOUT', defaultConfig.API_TIMEOUT, parseNumber),
  UPLOAD_MAX_SIZE: getEnvValue('UPLOAD_MAX_SIZE', defaultConfig.UPLOAD_MAX_SIZE, parseNumber),
  UPLOAD_CHUNK_SIZE: getEnvValue('UPLOAD_CHUNK_SIZE', defaultConfig.UPLOAD_CHUNK_SIZE, parseNumber),
  WS_URL: getEnvValue('WS_URL', defaultConfig.WS_URL),
  ENABLE_MOCK: getEnvValue('ENABLE_MOCK', defaultConfig.ENABLE_MOCK, parseBoolean),
  ENABLE_PWA: getEnvValue('ENABLE_PWA', defaultConfig.ENABLE_PWA, parseBoolean),
  ENABLE_DEVTOOLS: getEnvValue('ENABLE_DEVTOOLS', defaultConfig.ENABLE_DEVTOOLS, parseBoolean),
  SHOW_CONSOLE_LOG: getEnvValue('SHOW_CONSOLE_LOG', defaultConfig.SHOW_CONSOLE_LOG, parseBoolean),
}

// 开发环境检查
export const isDev = import.meta.env.DEV
export const isProd = import.meta.env.PROD
export const isTest = import.meta.env.MODE === 'test'

// 环境信息
export const envInfo = {
  mode: import.meta.env.MODE,
  dev: isDev,
  prod: isProd,
  test: isTest,
}

// 控制台输出环境信息
if (env.SHOW_CONSOLE_LOG) {
  console.log('🚀 环境配置:', {
    ...envInfo,
    config: env,
  })
}
