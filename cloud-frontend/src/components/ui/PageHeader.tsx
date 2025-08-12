import { ArrowLeftOutlined } from '@ant-design/icons'
import { Breadcrumb, Typography, Space, Button, Divider } from 'antd'
import React from 'react'
import { useDesignTokens } from '@/config/designTokens'

const { Title } = Typography

interface BreadcrumbItem {
  title: string
  href?: string
}

interface PageHeaderProps {
  title: string
  subTitle?: string
  breadcrumb?: BreadcrumbItem[]
  extra?: React.ReactNode
  onBack?: () => void
  children?: React.ReactNode
}

/**
 * 页面头部组件
 */
export const PageHeader: React.FC<PageHeaderProps> = ({
  title,
  subTitle,
  breadcrumb,
  extra,
  onBack,
  children,
}) => {
  const tokens = useDesignTokens()

  return (
    <div style={{ marginBottom: tokens.spacing.lg }}>
      {/* 面包屑导航 */}
      {breadcrumb && breadcrumb.length > 0 && (
        <Breadcrumb
          style={{ marginBottom: tokens.spacing.md }}
          items={breadcrumb.map((item) => ({
            title: item.href ? (
              <a href={item.href} style={{ color: tokens.colors.text.secondary }}>
                {item.title}
              </a>
            ) : (
              <span style={{ color: tokens.colors.text.primary }}>{item.title}</span>
            ),
          }))}
        />
      )}

      {/* 标题区域 */}
      <div
        style={{
          display: 'flex',
          alignItems: 'flex-start',
          justifyContent: 'space-between',
          marginBottom: children ? tokens.spacing.md : 0,
          flexWrap: 'wrap',
          gap: tokens.spacing.sm,
        }}
      >
        <Space align="center" size={tokens.spacing.sm}>
          {onBack && <Button type="text" icon={<ArrowLeftOutlined />} onClick={onBack} />}
          <div>
            <Title
              level={2}
              style={{
                margin: 0,
                fontSize: tokens.sizes.font.xl,
                color: tokens.colors.text.primary,
              }}
            >
              {title}
            </Title>
            {subTitle && (
              <div
                style={{
                  color: tokens.colors.text.secondary,
                  fontSize: tokens.sizes.font.sm,
                  marginTop: tokens.spacing.xs,
                  lineHeight: 1.5,
                }}
              >
                {subTitle}
              </div>
            )}
          </div>
        </Space>

        {extra && (
          <Space size={tokens.spacing.sm} wrap>
            {extra}
          </Space>
        )}
      </div>

      {/* 内容区域 */}
      {children && (
        <>
          <Divider style={{ margin: `${tokens.spacing.md}px 0` }} />
          {children}
        </>
      )}
    </div>
  )
}
