import {
  StatGroup,
  Stat,
  StatLabel,
  StatHelpText,
  Text,
  useColorModeValue,
} from '@chakra-ui/react'

interface CurrentPriceProps {
  price: number
  symbol: string
}

const CurrentPrice: React.FC<CurrentPriceProps> = ({ price, symbol }) => {
  const statBg = useColorModeValue('blue.50', 'blue.900')

  return (
    <StatGroup>
      <Stat
        bg={statBg}
        p={4}
        borderRadius="lg"
        boxShadow="sm"
        textAlign="center"
        w="100%"
      >
        <StatLabel fontSize="lg" color="gray.600">
          {symbol}/USDT Price
        </StatLabel>
        <Text 
          fontSize="3xl" 
          fontWeight="bold" 
          bgGradient="linear(to-r, blue.400, purple.500)" 
          bgClip="text"
          my={2}
        >
          ${price.toLocaleString()}
        </Text>
        <StatHelpText fontSize="sm">
          Last updated: {new Date().toLocaleTimeString()}
        </StatHelpText>
      </Stat>
    </StatGroup>
  )
}

export default CurrentPrice 