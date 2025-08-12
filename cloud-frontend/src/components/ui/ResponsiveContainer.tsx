import React from 'react'
import { useDesignTokens, designTokens } from '@/config/designTokens'

interface ResponsiveContainerProps {
  children: React.ReactNode
  maxWidth?: keyof typeof designTokens.breakpoints
  padding?: boolean
  center?: boolean
  className?: string
  style?: React.CSSProperties
}

/**
 * 响应式容器组件
 */
export const ResponsiveContainer: React.FC<ResponsiveContainerProps> = ({
  children,
  maxWidth = 'xxl',
  padding = true,
  center = true,
  className,
  style = {},
}) => {
  const tokens = useDesignTokens()

  const getMaxWidth = () => {
    const bp = designTokens.breakpoints[maxWidth]
    if ('min' in bp) {
      return bp.min
    }
    return 1200 // 默认最大宽度
  }

  const containerStyle: React.CSSProperties = {
    width: '100%',
    ...(center && { margin: '0 auto' }),
    ...(maxWidth && {
      maxWidth: getMaxWidth(),
    }),
    ...(padding && {
      paddingLeft: tokens.spacing.lg,
      paddingRight: tokens.spacing.lg,
    }),
    ...style,
  }

  return (
    <div className={className} style={containerStyle}>
      {children}
    </div>
  )
}
