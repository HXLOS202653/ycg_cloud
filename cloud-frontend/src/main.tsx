import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AntdProvider } from '@/components'
import { QueryProvider } from '@/components/QueryProvider'
import './index.css'
import './styles/antd-override.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryProvider>
      <AntdProvider>
        <App />
      </AntdProvider>
    </QueryProvider>
  </StrictMode>
)
