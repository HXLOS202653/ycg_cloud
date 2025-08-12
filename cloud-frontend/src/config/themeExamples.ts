import type { ThemeConfig } from 'antd'
import { designTokens } from './designTokens'

/**
 * 主题配置示例 - 展示如何轻松自定义布局和样式
 */

// 紧凑型布局主题
export const compactLayoutTheme = {
  ...designTokens,
  layout: {
    ...designTokens.layout,
    responsive: {
      // 紧凑布局 - 右侧占更多空间
      mainContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 12,
        xl: 14,
        xxl: 16,
      },
      sideContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 12,
        xl: 10,
        xxl: 8,
      },
      statCard: {
        xs: 12,
        sm: 6,
        md: 6,
        lg: 6,
        xl: 6,
        xxl: 4, // 更紧凑的统计卡片
      },
    },
    container: {
      maxWidth: '1400px', // 限制最大宽度
      padding: { xs: 12, sm: 16, md: 20, lg: 20, xl: 20, xxl: 20 },
      margin: '0 auto',
    },
  },
}

// 宽屏优化主题
export const wideScreenTheme = {
  ...designTokens,
  layout: {
    ...designTokens.layout,
    responsive: {
      // 宽屏优化 - 左侧占据更多空间
      mainContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 16,
        xl: 18,
        xxl: 20,
      },
      sideContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 8,
        xl: 6,
        xxl: 4,
      },
      statCard: {
        xs: 24,
        sm: 12,
        md: 8,
        lg: 6,
        xl: 4,
        xxl: 4, // 宽屏下更多列
      },
    },
    container: {
      maxWidth: 'none', // 无限制宽度
      padding: { xs: 16, sm: 24, md: 32, lg: 32, xl: 32, xxl: 40 },
      margin: '0',
    },
  },
}

// 移动端优化主题
export const mobileOptimizedTheme = {
  ...designTokens,
  layout: {
    ...designTokens.layout,
    responsive: {
      // 移动端优化 - 全部单列布局
      mainContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 24,
        xl: 16,
        xxl: 18,
      },
      sideContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 24,
        xl: 8,
        xxl: 6,
      },
      statCard: {
        xs: 24,
        sm: 12,
        md: 12,
        lg: 6,
        xl: 6,
        xxl: 6, // 移动端大卡片
      },
    },
    container: {
      maxWidth: 'none',
      padding: { xs: 8, sm: 12, md: 16, lg: 20, xl: 24, xxl: 24 },
      margin: '0',
    },
  },
}

// Ant Design 主题配置示例
export const customAntdTheme: ThemeConfig = {
  token: {
    // 可以通过环境变量或用户设置动态配置
    colorPrimary: import.meta.env.VITE_THEME_PRIMARY_COLOR || '#1677ff',
    borderRadius: 8,

    // 布局相关的令牌
    controlHeight: 40,
    fontSize: 14,

    // 间距令牌
    padding: 16,
    paddingLG: 24,
    paddingXS: 8,
  },

  components: {
    Layout: {
      // 可以通过主题动态设置Layout样式
      headerBg: '#ffffff',
      bodyBg: '#f5f5f5',
      siderBg: '#ffffff',
    },

    Card: {
      // 卡片样式可主题化
      paddingLG: 24,
      borderRadiusLG: 12,
    },
  },
}

/**
 * 主题应用函数示例
 */
export const applyLayoutTheme = (themeName: 'compact' | 'wide' | 'mobile' | 'default') => {
  let themeConfig = designTokens

  switch (themeName) {
    case 'compact':
      themeConfig = compactLayoutTheme
      break
    case 'wide':
      themeConfig = wideScreenTheme
      break
    case 'mobile':
      themeConfig = mobileOptimizedTheme
      break
    default:
      themeConfig = designTokens
  }

  // 应用CSS变量（用于运行时主题切换）
  const root = document.documentElement
  root.style.setProperty(
    '--layout-max-width',
    typeof themeConfig.layout.container.maxWidth === 'number'
      ? `${themeConfig.layout.container.maxWidth}px`
      : themeConfig.layout.container.maxWidth
  )

  // 可以设置更多CSS变量...
  Object.entries(themeConfig.layout.container.padding).forEach(([breakpoint, padding]) => {
    root.style.setProperty(`--layout-padding-${breakpoint}`, `${padding}px`)
  })

  return themeConfig
}

/**
 * 用户自定义主题接口
 */
export interface UserThemeConfig {
  layout?: {
    responsive?: Partial<typeof designTokens.layout.responsive>
    container?: Partial<typeof designTokens.layout.container>
  }
  colors?: {
    primary?: string
    success?: string
    warning?: string
    error?: string
  }
  spacing?: Partial<typeof designTokens.grid.gutters>
}

/**
 * 合并用户自定义主题
 */
export const mergeUserTheme = (userConfig: UserThemeConfig) => {
  return {
    ...designTokens,
    layout: {
      ...designTokens.layout,
      ...(userConfig.layout && {
        responsive: {
          ...designTokens.layout.responsive,
          ...userConfig.layout.responsive,
        },
        container: {
          ...designTokens.layout.container,
          ...userConfig.layout.container,
        },
      }),
    },
  }
}
