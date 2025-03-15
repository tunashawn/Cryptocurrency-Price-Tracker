import { ChartOptions } from 'chart.js'

export const createChartOptions = (
  useColorModeValue: (light: string, dark: string) => string
): ChartOptions<'line'> => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false,
    },
    title: {
      display: true,
      text: 'Last 24 Hours Price History',
      color: useColorModeValue('#2D3748', '#CBD5E0'),
      font: {
        size: 16,
        weight: 'bold',
        family: "'Inter', sans-serif",
      },
      padding: {
        top: 20,
        bottom: 20,
      },
    },
    tooltip: {
      mode: 'index',
      intersect: false,
      backgroundColor: useColorModeValue('rgba(255, 255, 255, 0.9)', 'rgba(45, 55, 72, 0.9)'),
      titleColor: useColorModeValue('#1A202C', '#FFFFFF'),
      bodyColor: useColorModeValue('#4A5568', '#E2E8F0'),
      titleFont: {
        size: 14,
        weight: 'bold',
        family: "'Inter', sans-serif",
      },
      bodyFont: {
        size: 13,
        family: "'Inter', sans-serif",
      },
      padding: 12,
      cornerRadius: 8,
      displayColors: false,
      borderColor: useColorModeValue('rgba(0, 0, 0, 0.1)', 'rgba(255, 255, 255, 0.1)'),
      borderWidth: 1,
    },
  },
  scales: {
    x: {
      grid: {
        display: false,
      },
      border: {
        display: false,
      },
      ticks: {
        font: {
          size: 11,
          family: "'Inter', sans-serif",
        },
        color: useColorModeValue('#4A5568', '#A0AEC0'),
        maxRotation: 45,
        minRotation: 45,
        autoSkip: true,
        maxTicksLimit: 8,
      },
    },
    y: {
      position: 'right',
      grid: {
        color: useColorModeValue('rgba(0, 0, 0, 0.06)', 'rgba(255, 255, 255, 0.06)'),
        display: false,
      },
      border: {
        display: false,
      },
      ticks: {
        font: {
          size: 11,
          family: "'Inter', sans-serif",
        },
        color: useColorModeValue('#4A5568', '#A0AEC0'),
        padding: 8,
        callback: (value) => `$${value.toLocaleString()}`,
      },
    },
  },
  interaction: {
    mode: 'nearest',
    axis: 'x',
    intersect: false,
  },
  hover: {
    mode: 'nearest',
    axis: 'x',
    intersect: false,
  },
  animation: {
    duration: 750,
    easing: 'easeInOutQuart',
  },
})

export const formatTooltipTitle = (timestamp: string): string => {
  const date = new Date(timestamp)
  return date.toLocaleString([], {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

export const formatTooltipLabel = (
  price: number,
  prevPrice: number | null
): string[] => {
  const priceStr = `Price: $${price.toLocaleString()}`
  
  if (prevPrice === null) {
    return [priceStr]
  }

  const change = price - prevPrice
  const changePercent = ((change / prevPrice) * 100).toFixed(2)
  const sign = change >= 0 ? '+' : ''
  
  return [
    priceStr,
    `Change: ${sign}$${change.toLocaleString()} (${sign}${changePercent}%)`
  ]
} 