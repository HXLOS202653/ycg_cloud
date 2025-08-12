import { theme } from 'antd'

// 设计令牌系统 - 避免硬编码
export const useDesignTokens = () => {
  const { token } = theme.useToken()

  return {
    // 颜色令牌
    colors: {
      primary: token.colorPrimary,
      success: token.colorSuccess,
      warning: token.colorWarning,
      error: token.colorError,
      info: token.colorInfo,

      text: {
        primary: token.colorText,
        secondary: token.colorTextSecondary,
        tertiary: token.colorTextTertiary,
        quaternary: token.colorTextQuaternary,
      },

      background: {
        base: token.colorBgBase,
        container: token.colorBgContainer,
        elevated: token.colorBgElevated,
        layout: token.colorBgLayout,
        spotlight: token.colorBgSpotlight,
        mask: token.colorBgMask,
      },

      border: {
        primary: token.colorBorder,
        secondary: token.colorBorderSecondary,
      },
    },

    // 尺寸令牌
    sizes: {
      unit: token.sizeUnit,
      step: token.sizeStep,

      control: {
        xs: token.controlHeightXS,
        sm: token.controlHeightSM,
        md: token.controlHeight,
        lg: token.controlHeightLG,
      },

      font: {
        xs: 12,
        sm: token.fontSizeSM,
        md: token.fontSize,
        lg: token.fontSizeLG,
        xl: token.fontSizeXL,
      },
    },

    // 间距令牌
    spacing: {
      xs: token.sizeXS,
      sm: token.sizeSM,
      md: token.size,
      lg: token.sizeLG,
      xl: token.sizeXL,
      xxl: token.sizeXXL,
    },

    // 圆角令牌
    radius: {
      xs: token.borderRadiusXS,
      sm: token.borderRadiusSM,
      md: token.borderRadius,
      lg: token.borderRadiusLG,
    },

    // 阴影令牌
    shadows: {
      primary: token.boxShadow,
      secondary: token.boxShadowSecondary,
      tertiary: token.boxShadowTertiary,
    },

    // 动画令牌
    motion: {
      fast: token.motionDurationFast,
      mid: token.motionDurationMid,
      slow: token.motionDurationSlow,
    },

    // 层级令牌
    zIndex: {
      base: 1,
      dropdown: 1050,
      sticky: 1020,
      fixed: 1030,
      modal: 1050,
      popover: 1060,
      tooltip: 1070,
      notification: 1080,
    },

    // 断点令牌
    breakpoints: {
      xs: 480,
      sm: 576,
      md: 768,
      lg: 992,
      xl: 1200,
      xxl: 1600,
    },
  }
}

// 静态设计令牌（不依赖Hook）
export const designTokens = {
  // 布局相关
  layout: {
    headerHeight: 64,
    siderWidth: 256,
    siderCollapsedWidth: 80,
    footerHeight: 64,
    contentPadding: 24,

    // 响应式布局配置
    responsive: {
      // 左侧主要内容区域
      mainContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 14,
        xl: 16,
        xxl: 18,
      },
      // 右侧辅助内容区域
      sideContent: {
        xs: 24,
        sm: 24,
        md: 24,
        lg: 10,
        xl: 8,
        xxl: 6,
      },
      // 统计卡片
      statCard: {
        xs: 24,
        sm: 12,
        md: 6,
        lg: 6,
        xl: 6,
        xxl: 6,
      },
    },

    // 容器配置
    container: {
      maxWidth: 'none', // 可设置为具体数值如 1200
      padding: {
        xs: 16,
        sm: 20,
        md: 24,
        lg: 24,
        xl: 24,
        xxl: 24,
      },
      margin: '0 auto',
    },
  },

  // 网格系统
  grid: {
    columns: 24,
    gutters: {
      xs: 8,
      sm: 16,
      md: 24,
      lg: 32,
      xl: 40,
      xxl: 48,
    },
  },

  // 响应式断点
  breakpoints: {
    xs: { max: 575 },
    sm: { min: 576, max: 767 },
    md: { min: 768, max: 991 },
    lg: { min: 992, max: 1199 },
    xl: { min: 1200, max: 1599 },
    xxl: { min: 1600 },
  },

  // 组件尺寸
  components: {
    card: {
      padding: {
        xs: 12,
        sm: 16,
        md: 24,
        lg: 32,
      },
      margin: {
        xs: 8,
        sm: 12,
        md: 16,
        lg: 24,
      },
    },

    button: {
      height: {
        xs: 24,
        sm: 32,
        md: 40,
        lg: 48,
      },
    },

    input: {
      height: {
        xs: 24,
        sm: 32,
        md: 40,
        lg: 48,
      },
    },
  },
}

// 响应式工具函数
export const createResponsiveStyle = (
  styles: Partial<Record<keyof typeof designTokens.breakpoints, React.CSSProperties>>
): React.CSSProperties => {
  const mediaQueries: React.CSSProperties = {}

  Object.entries(styles).forEach(([breakpoint, style]) => {
    const bp = designTokens.breakpoints[breakpoint as keyof typeof designTokens.breakpoints]

    if ('max' in bp && 'min' in bp) {
      // @ts-ignore
      mediaQueries[`@media (min-width: ${bp.min}px) and (max-width: ${bp.max}px)`] = style
    } else if ('min' in bp) {
      // @ts-ignore
      mediaQueries[`@media (min-width: ${bp.min}px)`] = style
    } else if ('max' in bp) {
      // @ts-ignore
      mediaQueries[`@media (max-width: ${bp.max}px)`] = style
    }
  })

  return mediaQueries
}
