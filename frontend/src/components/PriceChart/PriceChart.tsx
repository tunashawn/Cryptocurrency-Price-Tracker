import { Box, useColorModeValue } from '@chakra-ui/react'
import { Line } from 'react-chartjs-2'
import { ChartData, ScriptableContext } from 'chart.js'
import { PriceData } from '../../types/price'
import { createChartOptions, formatTooltipTitle, formatTooltipLabel } from './config'

interface PriceChartProps {
  priceHistory: PriceData[]
  symbol: string
}

const PriceChart: React.FC<PriceChartProps> = ({ priceHistory, symbol }) => {
  const borderColor = useColorModeValue('gray.200', 'gray.700')
  const chartBg = useColorModeValue('white', 'gray.800')

  const chartData: ChartData<'line'> = {
    labels: priceHistory.map(data => {
      const date = new Date(data.timestamp)
      return date.toLocaleTimeString([], { 
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      })
    }),
    datasets: [
      {
        label: `${symbol}/USDT Price`,
        data: priceHistory.map(data => data.latest_price),
        fill: true,
        backgroundColor: (context: ScriptableContext<'line'>) => {
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 200);
          gradient.addColorStop(0, 'rgba(52, 152, 219, 0.3)');
          gradient.addColorStop(1, 'rgba(52, 152, 219, 0.0)');
          return gradient;
        },
        borderColor: 'rgba(52, 152, 219, 1)',
        borderWidth: 2,
        pointRadius: 0,
        pointHoverRadius: 6,
        pointBackgroundColor: 'rgba(52, 152, 219, 1)',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: 'rgba(52, 152, 219, 1)',
        pointHoverBorderWidth: 2,
        tension: 0.3,
      },
    ],
  }

  const chartOptions = {
    ...createChartOptions(useColorModeValue),
    plugins: {
      ...createChartOptions(useColorModeValue).plugins,
      tooltip: {
        ...createChartOptions(useColorModeValue).plugins?.tooltip,
        callbacks: {
          title: (tooltipItems: any) => formatTooltipTitle(priceHistory[tooltipItems[0].dataIndex].timestamp),
          label: (context: any) => {
            const price = Number(context.parsed.y)
            const prevPrice = context.dataIndex > 0 
              ? Number(context.dataset.data[context.dataIndex - 1]) 
              : null
            return formatTooltipLabel(price, prevPrice)
          }
        }
      }
    }
  }

  return (
    <Box
      height="400px"
      p={6}
      borderRadius="xl"
      borderWidth={1}
      borderColor={borderColor}
      bg={chartBg}
      boxShadow="sm"
      _hover={{
        boxShadow: "md",
        transition: "box-shadow 0.2s"
      }}
    >
      <Line data={chartData} options={chartOptions} />
    </Box>
  )
}

export default PriceChart 