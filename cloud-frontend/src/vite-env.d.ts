/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string
  readonly VITE_APP_VERSION: string
  readonly VITE_API_BASE_URL: string
  readonly VITE_API_TIMEOUT: string
  readonly VITE_UPLOAD_MAX_SIZE: string
  readonly VITE_UPLOAD_CHUNK_SIZE: string
  readonly VITE_WS_URL: string
  readonly VITE_ENABLE_MOCK: string
  readonly VITE_ENABLE_PWA: string
  readonly VITE_ENABLE_DEVTOOLS: string
  readonly VITE_SHOW_CONSOLE_LOG: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
