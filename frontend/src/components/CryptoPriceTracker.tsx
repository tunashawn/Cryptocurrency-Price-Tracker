import { ChangeEvent, useState } from "react";
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
  List,
  ListItem,
  Box,
  useColorMode,
} from "@chakra-ui/react";
import { SearchIcon } from "@chakra-ui/icons";
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler } from "chart.js";
import CurrentPrice from "./CurrentPrice/CurrentPrice";
import PriceChart from "./PriceChart/PriceChart";
import { usePriceData } from "../hooks/usePriceData";
import { cryptocurrencies } from "../hooks/usePriceData";

// Register Chart.js components
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler);

const CryptoPriceTracker = () => {
  const [inputSymbol, setInputSymbol] = useState<string>("BTC");
  const [symbol, setSymbol] = useState<string>("BTC");
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [showDropdown, setShowDropdown] = useState(false);
  const { currentPrice, priceHistory, isLoading, error } = usePriceData(symbol);

  const { colorMode } = useColorMode();
  const bgColor = colorMode === "light" ? "white" : "gray.800";
  const borderColor = colorMode === "light" ? "gray.300" : "gray.600";
  const hoverBg = colorMode === "light" ? "gray.100" : "gray.700";

  const handleSearch = () => {
    if (inputSymbol.trim()) {
      setSymbol(inputSymbol.trim().toUpperCase());
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
      setShowDropdown(false);
    }
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.toUpperCase();
    setInputSymbol(value);
  
    if (value) {
      const filteredSuggestions = cryptocurrencies.filter((symbol) => 
        symbol.startsWith(value)
      );
      setSuggestions(filteredSuggestions);
      setShowDropdown(filteredSuggestions.length > 0); // Show dropdown only if there are suggestions
    } else {
      setSuggestions([]);
      setShowDropdown(false);
    }
  };
  

  const handleSuggestionClick = (symbol: string) => {
    setInputSymbol(symbol);
    setSymbol(symbol);
    setSuggestions([]);
    setShowDropdown(false);
  };
    

  return (
    <VStack spacing={6} align="stretch">
      {/* Search Box Container */}
      <Box position="relative" width="full">
        <InputGroup size="lg">
          <InputLeftElement pointerEvents="none">
            <SearchIcon color="gray.400" />
          </InputLeftElement>
          <Input
            placeholder="Enter cryptocurrency symbol (e.g., BTC)"
            value={inputSymbol}
            onChange={handleChange}
            onKeyDown={handleKeyPress}
            borderWidth={2}
            borderRightRadius="0"
            _focus={{
              borderColor: "blue.400",
              boxShadow: "0 0 0 1px blue.400",
            }}
            onFocus={() => setShowDropdown(true)}
            onBlur={() => setTimeout(() => setShowDropdown(false), 200)}
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

        {/* Dropdown Suggestions */}
        {showDropdown && suggestions.length > 0 && (
          <Box
            position="absolute"
            top="100%" // Ensures dropdown appears right below input
            left={0}
            width="100%"
            bg={bgColor}
            border="1px solid"
            borderColor={borderColor}
            borderRadius="md"
            zIndex={10}
            boxShadow="lg"
            mt={1}
          >
            <List maxHeight="200px" overflowY="auto">
              {suggestions.map((suggestion) => (
                <ListItem
                  key={suggestion}
                  onMouseDown={() => handleSuggestionClick(suggestion)}
                  px={4}
                  py={2}
                  cursor="pointer"
                  _hover={{ bg: hoverBg }}
                  borderBottom="1px solid"
                  borderColor={borderColor}
                >
                  {suggestion}
                </ListItem>
              ))}
            </List>
          </Box>
        )}
      </Box>

      {/* Error Message */}
      {error && (
        <Alert status="error" borderRadius="md">
          <AlertIcon />
          {error}
        </Alert>
      )}

      {/* Loading Spinner */}
      {isLoading ? (
        <Center py={6}>
          <Spinner size="xl" color="blue.500" thickness="4px" />
        </Center>
      ) : (
        <>
          {/* Current Price */}
          {currentPrice ? (
            <CurrentPrice price={currentPrice} symbol={symbol} />
          ) : (
            !isLoading &&
            !error && (
              <Alert status="warning" borderRadius="md">
                <AlertIcon />
                No price data found for {symbol}
              </Alert>
            )
          )}

          {/* Price Chart */}
          {priceHistory.length > 0 ? (
            <PriceChart priceHistory={priceHistory} symbol={symbol} />
          ) : (
            !isLoading &&
            !error && (
              <Alert status="warning" borderRadius="md">
                <AlertIcon />
                No historical data found for {symbol}
              </Alert>
            )
          )}
        </>
      )}
    </VStack>
  );
};

export default CryptoPriceTracker;