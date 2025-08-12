import { PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import { Empty, Button, Space } from 'antd'
import React from 'react'
import { useDesignTokens } from '@/config/designTokens'

interface EmptyStateProps {
  title?: string
  description?: string
  image?: React.ReactNode
  actions?: React.ReactNode
  showCreateButton?: boolean
  showRefreshButton?: boolean
  onCreateClick?: () => void
  onRefreshClick?: () => void
  createButtonText?: string
  refreshButtonText?: string
}

/**
 * 空状态组件
 */
export const EmptyState: React.FC<EmptyStateProps> = ({
  title = '暂无数据',
  description,
  image,
  actions,
  showCreateButton = false,
  showRefreshButton = false,
  onCreateClick,
  onRefreshClick,
  createButtonText = '新建',
  refreshButtonText = '刷新',
}) => {
  const tokens = useDesignTokens()

  const defaultActions = (
    <Space size={tokens.spacing.sm}>
      {showRefreshButton && (
        <Button icon={<ReloadOutlined />} onClick={onRefreshClick}>
          {refreshButtonText}
        </Button>
      )}
      {showCreateButton && (
        <Button type="primary" icon={<PlusOutlined />} onClick={onCreateClick}>
          {createButtonText}
        </Button>
      )}
    </Space>
  )

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '300px',
        padding: `${tokens.spacing.xl}px ${tokens.spacing.lg}px`,
        textAlign: 'center',
      }}
    >
      <Empty
        image={image || Empty.PRESENTED_IMAGE_SIMPLE}
        description={
          <div>
            <div
              style={{
                fontSize: tokens.sizes.font.md,
                color: tokens.colors.text.secondary,
                marginBottom: tokens.spacing.xs,
                fontWeight: 500,
              }}
            >
              {title}
            </div>
            {description && (
              <div
                style={{
                  fontSize: tokens.sizes.font.sm,
                  color: tokens.colors.text.tertiary,
                  lineHeight: 1.5,
                }}
              >
                {description}
              </div>
            )}
          </div>
        }
      >
        {actions || ((showCreateButton || showRefreshButton) && defaultActions)}
      </Empty>
    </div>
  )
}
