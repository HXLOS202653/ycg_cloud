import { ConfigProvider, App } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import React from 'react'
import { antdTheme } from '@/config/theme'

interface AntdProviderProps {
  children: React.ReactNode
}

/**
 * Ant Design 配置提供者
 */
export const AntdProvider: React.FC<AntdProviderProps> = ({ children }) => {
  return (
    <ConfigProvider
      theme={antdTheme}
      locale={zhCN}
      componentSize="middle"
      // 自动插入空格
      autoInsertSpaceInButton={false}
      // 波纹效果
      wave={{ disabled: false }}
      // 虚拟滚动
      virtual
      // 表单验证消息模板
      form={{
        validateMessages: {
          required: '${label}是必填项',
          types: {
            email: '${label}不是有效的邮箱格式',
            number: '${label}不是有效的数字格式',
            url: '${label}不是有效的URL格式',
          },
          number: {
            range: '${label}必须在${min}和${max}之间',
            min: '${label}不能小于${min}',
            max: '${label}不能大于${max}',
          },
          string: {
            range: '${label}长度必须在${min}和${max}之间',
            min: '${label}长度不能小于${min}',
            max: '${label}长度不能大于${max}',
          },
        },
      }}
      // 输入框配置
      input={{
        autoComplete: 'off',
      }}
      // 选择器配置
      select={{
        showSearch: true,
      }}
      // 其他配置可以通过组件props直接传递
    >
      <App>{children}</App>
    </ConfigProvider>
  )
}
