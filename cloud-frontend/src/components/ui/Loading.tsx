import { LoadingOutlined } from '@ant-design/icons'
import { Spin } from 'antd'
import React from 'react'
import { useDesignTokens } from '@/config/designTokens'

interface LoadingProps {
  size?: 'small' | 'default' | 'large'
  tip?: string
  spinning?: boolean
  children?: React.ReactNode
}

/**
 * 加载组件
 */
export const Loading: React.FC<LoadingProps> = ({
  size = 'default',
  tip = '加载中...',
  spinning = true,
  children,
}) => {
  const tokens = useDesignTokens()

  const getIconSize = () => {
    switch (size) {
      case 'large':
        return tokens.sizes.font.xl
      case 'small':
        return tokens.sizes.font.sm
      default:
        return tokens.sizes.font.lg
    }
  }

  const antIcon = <LoadingOutlined style={{ fontSize: getIconSize() }} spin />

  if (children) {
    return (
      <Spin spinning={spinning} tip={tip} indicator={antIcon} size={size}>
        {children}
      </Spin>
    )
  }

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '200px',
        flexDirection: 'column',
        gap: tokens.spacing.md,
        textAlign: 'center',
      }}
    >
      <Spin indicator={antIcon} size={size} />
      {tip && (
        <div
          style={{
            color: tokens.colors.text.secondary,
            fontSize: tokens.sizes.font.sm,
            lineHeight: 1.5,
          }}
        >
          {tip}
        </div>
      )}
    </div>
  )
}
