import { useState, useEffect } from 'react'
import { designTokens } from '@/config/designTokens'

// 响应式断点类型
export type Breakpoint = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'xxl'

interface ResponsiveInfo {
  xs: boolean
  sm: boolean
  md: boolean
  lg: boolean
  xl: boolean
  xxl: boolean
  current: Breakpoint
  isMobile: boolean
  isTablet: boolean
  isDesktop: boolean
}

/**
 * 响应式Hook - 监听屏幕尺寸变化
 */
export const useResponsive = (): ResponsiveInfo => {
  const [screenSize, setScreenSize] = useState<ResponsiveInfo>({
    xs: false,
    sm: false,
    md: false,
    lg: false,
    xl: false,
    xxl: false,
    current: 'md',
    isMobile: false,
    isTablet: false,
    isDesktop: true,
  })

  useEffect(() => {
    const updateScreenSize = () => {
      const width = window.innerWidth

      // 计算当前断点
      let current: Breakpoint = 'xs'
      if (width >= designTokens.breakpoints.xxl.min) current = 'xxl'
      else if (width >= designTokens.breakpoints.xl.min) current = 'xl'
      else if (width >= designTokens.breakpoints.lg.min) current = 'lg'
      else if (width >= designTokens.breakpoints.md.min) current = 'md'
      else if (width >= designTokens.breakpoints.sm.min) current = 'sm'

      // 计算各断点状态
      const newScreenSize: ResponsiveInfo = {
        xs: width <= designTokens.breakpoints.xs.max,
        sm: width >= designTokens.breakpoints.sm.min && width <= designTokens.breakpoints.sm.max,
        md: width >= designTokens.breakpoints.md.min && width <= designTokens.breakpoints.md.max,
        lg: width >= designTokens.breakpoints.lg.min && width <= designTokens.breakpoints.lg.max,
        xl: width >= designTokens.breakpoints.xl.min && width <= designTokens.breakpoints.xl.max,
        xxl: width >= designTokens.breakpoints.xxl.min,
        current,
        isMobile: width <= designTokens.breakpoints.xs.max,
        isTablet:
          width >= designTokens.breakpoints.sm.min && width <= designTokens.breakpoints.md.max,
        isDesktop: width >= designTokens.breakpoints.lg.min,
      }

      setScreenSize(newScreenSize)
    }

    // 初始化
    updateScreenSize()

    // 监听窗口大小变化
    window.addEventListener('resize', updateScreenSize)

    // 清理事件监听器
    return () => {
      window.removeEventListener('resize', updateScreenSize)
    }
  }, [])

  return screenSize
}

/**
 * 获取响应式值的Hook
 */
export const useResponsiveValue = <T>(values: {
  xs?: T
  sm?: T
  md?: T
  lg?: T
  xl?: T
  xxl?: T
}): T | undefined => {
  const responsive = useResponsive()

  // 按优先级返回对应断点的值
  if (responsive.xxl && values.xxl !== undefined) return values.xxl
  if (responsive.xl && values.xl !== undefined) return values.xl
  if (responsive.lg && values.lg !== undefined) return values.lg
  if (responsive.md && values.md !== undefined) return values.md
  if (responsive.sm && values.sm !== undefined) return values.sm
  if (responsive.xs && values.xs !== undefined) return values.xs

  // 如果当前断点没有值，则向下查找
  const breakpoints: Breakpoint[] = ['xxl', 'xl', 'lg', 'md', 'sm', 'xs']
  const currentIndex = breakpoints.indexOf(responsive.current)

  for (let i = currentIndex; i < breakpoints.length; i++) {
    const bp = breakpoints[i]
    if (values[bp] !== undefined) return values[bp]
  }

  return undefined
}

/**
 * 获取响应式网格列数
 */
export const useResponsiveColumns = (
  columns: {
    xs?: number
    sm?: number
    md?: number
    lg?: number
    xl?: number
    xxl?: number
  } = {
    xs: 1,
    sm: 2,
    md: 3,
    lg: 4,
    xl: 5,
    xxl: 6,
  }
): number => {
  const value = useResponsiveValue(columns)
  return value || 1
}
