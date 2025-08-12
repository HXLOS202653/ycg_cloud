import {
  CloudOutlined,
  PlusOutlined,
  UploadOutlined,
  DownloadOutlined,
  UserOutlined,
  SettingOutlined,
} from '@ant-design/icons'
import {
  Layout,
  Button,
  Card,
  Space,
  Typography,
  Divider,
  Tag,
  Alert,
  Statistic,
  Row,
  Col,
  Spin,
} from 'antd'
import { useState } from 'react'
import { PageHeader, Loading, EmptyState } from '@/components'
import { useDesignTokens, designTokens } from '@/config/designTokens'
import { env } from '@/config/env'
import { useStats } from '@/hooks/useApi'

const { Header, Content, Footer } = Layout
const { Title, Paragraph } = Typography

function App() {
  const [loading, setLoading] = useState(false)
  const [count, setCount] = useState(0)
  const tokens = useDesignTokens()

  // 使用 TanStack Query 获取数据
  const { data: statsData, isLoading: statsLoading } = useStats.overview({
    // 在没有真实API时禁用查询，避免网络错误
    enabled: false,
  })

  // 暂时注释掉未使用的查询，避免编译错误
  // const {
  //   data: filesData,
  //   isLoading: filesLoading,
  // } = useFiles.list(
  //   { page: 1, pageSize: 10 },
  //   { enabled: false }
  // )

  // const {
  //   data: foldersData,
  //   isLoading: foldersLoading,
  // } = useFolders.list(
  //   { page: 1, pageSize: 10 },
  //   { enabled: false }
  // )

  const handleTest = () => {
    setLoading(true)
    setTimeout(() => {
      setCount(count + 1)
      setLoading(false)
    }, 1000)
  }

  // 使用模拟数据或真实数据
  const stats = statsData || {
    totalFiles: 1128,
    storageUsed: 93,
    storageTotal: 100,
    todayUploads: 25,
    todayDownloads: 156,
  }

  // 响应式样式
  const headerStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    background: tokens.colors.background.container,
    boxShadow: tokens.shadows.secondary,
    padding: `0 ${tokens.spacing.lg}px`,
    height: designTokens.layout.headerHeight,
  }

  const iconStyle: React.CSSProperties = {
    fontSize: tokens.sizes.font.xl,
    color: tokens.colors.primary,
  }

  const titleStyle: React.CSSProperties = {
    margin: 0,
    color: tokens.colors.primary,
    fontSize: tokens.sizes.font.xl,
  }

  return (
    <Layout
      style={{
        minHeight: '100vh',
        width: '100%',
        maxWidth: 'none',
      }}
    >
      <Header style={headerStyle}>
        <Space align="center" size={tokens.spacing.md}>
          <CloudOutlined style={iconStyle} />
          <Title level={3} style={titleStyle}>
            {env.APP_TITLE}
          </Title>
          <Tag color="blue">v{env.APP_VERSION}</Tag>
        </Space>

        <Space size={tokens.spacing.sm}>
          <Button icon={<UserOutlined />}>用户中心</Button>
          <Button icon={<SettingOutlined />}>设置</Button>
        </Space>
      </Header>

      <Content
        style={{
          padding: 0,
          background: tokens.colors.background.layout,
          minHeight: `calc(100vh - ${designTokens.layout.headerHeight + designTokens.layout.footerHeight}px)`,
        }}
      >
        <div
          style={{
            width: '100%',
            padding: `${tokens.spacing.lg}px`,
            maxWidth: 'none',
          }}
        >
          <PageHeader
            title="欢迎使用网络云盘系统"
            subTitle="基于 React 19 + TypeScript 5.8 + Ant Design 5.27 构建"
            breadcrumb={[{ title: '首页' }, { title: '仪表板' }]}
            extra={
              <Space size={tokens.spacing.sm}>
                <Button type="primary" icon={<UploadOutlined />}>
                  上传文件
                </Button>
                <Button icon={<PlusOutlined />}>新建文件夹</Button>
              </Space>
            }
          />

          <Row
            gutter={[
              {
                xs: designTokens.grid.gutters.xs,
                sm: designTokens.grid.gutters.sm,
                md: designTokens.grid.gutters.md,
                lg: designTokens.grid.gutters.lg,
              },
              {
                xs: designTokens.grid.gutters.xs,
                sm: designTokens.grid.gutters.sm,
                md: designTokens.grid.gutters.md,
                lg: designTokens.grid.gutters.lg,
              },
            ]}
            style={{ marginBottom: tokens.spacing.lg }}
          >
            <Col {...designTokens.layout.responsive.statCard}>
              <Card>
                <Spin spinning={statsLoading} size="small">
                  <Statistic
                    title="总文件数"
                    value={stats.totalFiles}
                    prefix={<CloudOutlined />}
                    valueStyle={{ color: tokens.colors.success }}
                  />
                </Spin>
              </Card>
            </Col>
            <Col {...designTokens.layout.responsive.statCard}>
              <Card>
                <Spin spinning={statsLoading} size="small">
                  <Statistic
                    title="存储空间"
                    value={stats.storageUsed}
                    suffix={`/ ${stats.storageTotal} GB`}
                    valueStyle={{ color: tokens.colors.error }}
                  />
                </Spin>
              </Card>
            </Col>
            <Col {...designTokens.layout.responsive.statCard}>
              <Card>
                <Spin spinning={statsLoading} size="small">
                  <Statistic
                    title="今日上传"
                    value={stats.todayUploads}
                    prefix={<UploadOutlined />}
                    valueStyle={{ color: tokens.colors.primary }}
                  />
                </Spin>
              </Card>
            </Col>
            <Col {...designTokens.layout.responsive.statCard}>
              <Card>
                <Spin spinning={statsLoading} size="small">
                  <Statistic
                    title="今日下载"
                    value={stats.todayDownloads}
                    prefix={<DownloadOutlined />}
                    valueStyle={{ color: tokens.colors.info }}
                  />
                </Spin>
              </Card>
            </Col>
          </Row>

          <Row
            gutter={[
              {
                xs: designTokens.grid.gutters.xs,
                sm: designTokens.grid.gutters.sm,
                md: designTokens.grid.gutters.md,
                lg: designTokens.grid.gutters.lg,
              },
              {
                xs: designTokens.grid.gutters.xs,
                sm: designTokens.grid.gutters.sm,
                md: designTokens.grid.gutters.md,
                lg: designTokens.grid.gutters.lg,
              },
            ]}
          >
            <Col {...designTokens.layout.responsive.mainContent}>
              <Card
                title="系统状态"
                style={{
                  minHeight: '400px',
                  height: 'auto',
                }}
              >
                <Alert
                  message="系统运行正常"
                  description="所有服务正常运行，API响应时间良好。"
                  type="success"
                  showIcon
                  style={{ marginBottom: tokens.spacing.md }}
                />

                <Divider>功能测试</Divider>

                <Space direction="vertical" style={{ width: '100%' }} size={tokens.spacing.md}>
                  <Paragraph style={{ color: tokens.colors.text.primary }}>
                    这是一个基于最新技术栈构建的网络云盘系统演示界面：
                  </Paragraph>
                  <ul
                    style={{
                      color: tokens.colors.text.secondary,
                      lineHeight: 1.8,
                      paddingLeft: tokens.spacing.lg,
                    }}
                  >
                    <li>✅ React 19.1.1 - 最新版本React框架</li>
                    <li>✅ TypeScript 5.8.3 - 强类型支持</li>
                    <li>✅ Ant Design 5.27.0 - 企业级UI组件库</li>
                    <li>✅ Vite 7.1.2 - 极速构建工具</li>
                    <li>✅ 完整的开发环境配置</li>
                  </ul>

                  <Space size={tokens.spacing.sm}>
                    <Button type="primary" onClick={handleTest} loading={loading}>
                      测试按钮 (点击次数: {count})
                    </Button>
                    <Button>取消</Button>
                  </Space>
                </Space>
              </Card>
            </Col>

            <Col {...designTokens.layout.responsive.sideContent}>
              <Card
                title="快速操作"
                style={{
                  minHeight: '400px',
                  height: 'auto',
                }}
              >
                {count === 0 ? (
                  <EmptyState
                    title="暂无操作记录"
                    description="点击左侧测试按钮开始体验"
                    showCreateButton
                    createButtonText="开始体验"
                    onCreateClick={handleTest}
                  />
                ) : (
                  <Space direction="vertical" style={{ width: '100%' }} size={tokens.spacing.md}>
                    <Alert message={`已测试 ${count} 次`} type="info" showIcon />
                    <Button block icon={<UploadOutlined />}>
                      上传文件
                    </Button>
                    <Button block icon={<PlusOutlined />}>
                      新建文件夹
                    </Button>
                    <Button block icon={<DownloadOutlined />}>
                      批量下载
                    </Button>
                  </Space>
                )}
              </Card>
            </Col>
          </Row>
        </div>

        {loading && (
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              background: tokens.colors.background.mask,
              zIndex: tokens.zIndex.modal,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Loading tip="正在处理请求..." />
          </div>
        )}
      </Content>

      <Footer
        style={{
          textAlign: 'center',
          background: tokens.colors.background.container,
          color: tokens.colors.text.secondary,
          height: designTokens.layout.footerHeight,
          lineHeight: `${designTokens.layout.footerHeight}px`,
          borderTop: `1px solid ${tokens.colors.border.secondary}`,
        }}
      >
        网络云盘系统 ©2025 Created with ❤️ by Development Team
      </Footer>
    </Layout>
  )
}

export default App
