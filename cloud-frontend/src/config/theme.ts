import type { ThemeConfig } from 'antd'
import { theme } from 'antd'

// 自定义主题配置
export const antdTheme: ThemeConfig = {
  // 使用算法主题
  algorithm: theme.defaultAlgorithm,

  // 设计令牌
  token: {
    // 主色
    colorPrimary: '#1677ff',

    // 成功色
    colorSuccess: '#52c41a',

    // 警告色
    colorWarning: '#faad14',

    // 错误色
    colorError: '#ff4d4f',

    // 信息色
    colorInfo: '#1677ff',

    // 文本色
    colorTextBase: '#000000',
    colorTextSecondary: '#666666',
    colorTextTertiary: '#999999',
    colorTextQuaternary: '#cccccc',

    // 背景色
    colorBgBase: '#ffffff',
    colorBgContainer: '#ffffff',
    colorBgElevated: '#ffffff',
    colorBgLayout: '#f5f5f5',
    colorBgSpotlight: '#ffffff',
    colorBgMask: 'rgba(0, 0, 0, 0.45)',

    // 边框色
    colorBorder: '#d9d9d9',
    colorBorderSecondary: '#f0f0f0',

    // 字体
    fontFamily: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji'`,
    fontSize: 14,
    fontSizeLG: 16,
    fontSizeSM: 12,
    fontSizeXL: 20,

    // 尺寸
    sizeUnit: 4,
    sizeStep: 4,
    sizePopupArrow: 16,
    controlHeight: 32,
    controlHeightLG: 40,
    controlHeightSM: 24,
    controlHeightXS: 16,

    // 圆角
    borderRadius: 6,
    borderRadiusLG: 8,
    borderRadiusSM: 4,
    borderRadiusXS: 2,

    // 阴影
    boxShadow: `
      0 6px 16px 0 rgba(0, 0, 0, 0.08),
      0 3px 6px -4px rgba(0, 0, 0, 0.12),
      0 9px 28px 8px rgba(0, 0, 0, 0.05)
    `,
    boxShadowSecondary: `
      0 6px 16px 0 rgba(0, 0, 0, 0.08),
      0 3px 6px -4px rgba(0, 0, 0, 0.12),
      0 9px 28px 8px rgba(0, 0, 0, 0.05)
    `,

    // 动画
    motionDurationFast: '0.1s',
    motionDurationMid: '0.2s',
    motionDurationSlow: '0.3s',
  },

  // 组件令牌
  components: {
    // Button 组件
    Button: {
      colorPrimary: '#1677ff',
      algorithm: true,
    },

    // Menu 组件
    Menu: {
      itemBg: 'transparent',
      subMenuItemBg: 'transparent',
      itemSelectedBg: '#e6f4ff',
      itemSelectedColor: '#1677ff',
    },

    // Layout 组件
    Layout: {
      headerBg: '#ffffff',
      siderBg: '#ffffff',
      bodyBg: '#f5f5f5',
    },

    // Table 组件
    Table: {
      headerBg: '#fafafa',
      rowHoverBg: '#f5f5f5',
    },

    // Form 组件
    Form: {
      labelColor: '#000000',
      labelRequiredMarkColor: '#ff4d4f',
    },

    // Input 组件
    Input: {
      hoverBorderColor: '#4096ff',
      activeBorderColor: '#1677ff',
    },

    // Select 组件
    Select: {
      optionSelectedBg: '#e6f4ff',
    },

    // Card 组件
    Card: {
      headerBg: 'transparent',
      boxShadow:
        '0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
    },

    // Modal 组件
    Modal: {
      headerBg: '#ffffff',
      contentBg: '#ffffff',
      footerBg: 'transparent',
    },
  },
}

// 暗色主题配置
export const darkTheme: ThemeConfig = {
  algorithm: theme.darkAlgorithm,
  token: {
    colorPrimary: '#1677ff',
    colorBgBase: '#141414',
    colorBgContainer: '#1f1f1f',
    colorBgElevated: '#262626',
    colorBgLayout: '#000000',
    colorTextBase: '#ffffff',
    colorBorder: '#424242',
    colorBorderSecondary: '#303030',
  },
  components: {
    Layout: {
      headerBg: '#1f1f1f',
      siderBg: '#1f1f1f',
      bodyBg: '#141414',
    },
    Menu: {
      itemBg: 'transparent',
      subMenuItemBg: 'transparent',
      itemSelectedBg: '#1677ff',
    },
  },
}

// 主题切换工具
export const getThemeConfig = (isDark: boolean): ThemeConfig => {
  return isDark ? darkTheme : antdTheme
}
