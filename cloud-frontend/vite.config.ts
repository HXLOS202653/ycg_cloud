import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  
  // 开发服务器配置
  server: {
    port: 5173,
    host: true, // 允许外部访问
    open: true, // 自动打开浏览器
    cors: true, // 启用CORS
    // API代理配置
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端API地址
        changeOrigin: true,
        secure: false,
      },
    },
  },
  
  // 路径别名配置
  resolve: {
    alias: {
      '@': '/src',
      '@components': '/src/components',
      '@pages': '/src/pages',
      '@utils': '/src/utils',
      '@types': '/src/types',
      '@assets': '/src/assets',
      '@hooks': '/src/hooks',
      '@services': '/src/services',
      '@store': '/src/store',
    },
  },
  
  // 构建配置
  build: {
    target: 'esnext', // 构建目标
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: true, // 生成sourcemap
    minify: 'esbuild', // 使用esbuild压缩
    rollupOptions: {
      output: {
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: '[ext]/[name]-[hash].[ext]',
        manualChunks: {
          // 第三方库分包
          'react-vendor': ['react', 'react-dom'],
          'antd-vendor': ['antd', '@ant-design/icons', '@ant-design/colors'],
          'query-vendor': ['@tanstack/react-query', '@tanstack/react-query-devtools'],
        },
      },
    },
    // 构建优化
    chunkSizeWarningLimit: 1000,
  },
  
  // 预览服务器配置
  preview: {
    port: 4173,
    host: true,
  },
  
  // 环境变量前缀
  envPrefix: 'VITE_',
  
  // CSS配置
  css: {
    devSourcemap: true, // 开发环境CSS sourcemap
    preprocessorOptions: {
      scss: {
        additionalData: '@import "@/styles/variables.scss";',
      },
    },
  },
  
  // 依赖预构建优化
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'antd',
      '@ant-design/icons',
      '@ant-design/colors',
      '@tanstack/react-query',
      '@tanstack/react-query-devtools',
    ],
  },
})
