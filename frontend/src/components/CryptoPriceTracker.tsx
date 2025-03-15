import { ChangeEvent, useState } from 'react'
import {
  VStack,
  Input,
  InputGroup,
  InputLeftElement,
  Button,
  Alert,
  AlertIcon,
  Center,
  Spinner,
} from '@chakra-ui/react'
import { SearchIcon } from '@chakra-ui/icons'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler } from 'chart.js'
import CurrentPrice from './CurrentPrice/CurrentPrice'
import PriceChart from './PriceChart/PriceChart'
import usePriceData from '../hooks/usePriceData'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const CryptoPriceTracker = () => {
  const [inputSymbol, setInputSymbol] = useState<string>('BTC')
  const [symbol, setSymbol] = useState<string>('BTC')
  const { currentPrice, priceHistory, isLoading, error } = usePriceData(symbol)

  const handleSearch = () => {
    if (inputSymbol.trim()) {
      setSymbol(inputSymbol.trim().toUpperCase())
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  return (
    <VStack spacing={6} align="stretch">
      <InputGroup size="lg">
        <InputLeftElement pointerEvents="none">
          <SearchIcon color="gray.400" />
        </InputLeftElement>
        <Input
          placeholder="Enter cryptocurrency symbol (e.g., BTC)"
          value={inputSymbol}
          onChange={(e: ChangeEvent<HTMLInputElement>) => setInputSymbol(e.target.value.toUpperCase())}
          onKeyPress={handleKeyPress}
          borderWidth={2}
          borderRightRadius="0"
          _focus={{
            borderColor: "blue.400",
            boxShadow: "0 0 0 1px blue.400"
          }}
        />
        <Button
          onClick={handleSearch}
          colorScheme="blue"
          size="lg"
          borderLeftRadius="0"
          px={8}
          isLoading={isLoading}
        >
          Search
        </Button>
      </InputGroup>

      {error && (
        <Alert status="error" borderRadius="md">
          <AlertIcon />
          {error}
        </Alert>
      )}
      
      {isLoading ? (
        <Center py={6}>
          <Spinner size="xl" color="blue.500" thickness="4px" />
        </Center>
      ) : (
        <>
          {currentPrice ? (
            <CurrentPrice price={currentPrice} symbol={symbol} />
          ) : !isLoading && !error && (
            <Alert status="warning" borderRadius="md">
              <AlertIcon />
              No price data found for {symbol}
            </Alert>
          )}

          {priceHistory.length > 0 ? (
            <PriceChart priceHistory={priceHistory} symbol={symbol} />
          ) : !isLoading && !error && (
            <Alert status="warning" borderRadius="md">
              <AlertIcon />
              No historical data found for {symbol}
            </Alert>
          )}
        </>
      )}
    </VStack>
  )
}

export default CryptoPriceTracker 