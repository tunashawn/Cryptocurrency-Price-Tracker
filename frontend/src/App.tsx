import { Box, Flex, Heading, useColorModeValue } from '@chakra-ui/react'
import CryptoPriceTracker from './components/CryptoPriceTracker'
import ErrorBoundary from './components/ErrorBoundary'

function App() {
  const bgGradient = useColorModeValue(
    'linear(to-br, blue.50, purple.50)',
    'linear(to-br, gray.900, purple.900)'
  )

  return (
    <Flex 
      minH="100vh" 
      w="100vw"
      bgGradient={bgGradient}
      alignItems="center"
      justifyContent="center"
      p={{ base: 4, md: 8 }}
    >
      <Box
        bg={useColorModeValue('white', 'gray.800')}
        p={{ base: 6, md: 8 }}
        borderRadius="xl"
        boxShadow="2xl"
        maxW="600px"
        w="100%"
      >
        <Heading 
          textAlign="center" 
          mb={8}
          bgGradient="linear(to-r, blue.400, purple.500)"
          bgClip="text"
          fontSize={{ base: "2xl", md: "3xl" }}
          letterSpacing="tight"
        >
          Cryptocurrency Price Tracker
        </Heading>
        <ErrorBoundary>
          <CryptoPriceTracker />
        </ErrorBoundary>
      </Box>
    </Flex>
  )
}

export default App
