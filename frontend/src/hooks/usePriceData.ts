import { useState, useEffect } from "react";
import axios from "axios";
import { PriceData, ApiResponse, PriceState } from "../types/price";

const HOURS_24 = 24 * 60 * 60 * 1000; // 24 hours in milliseconds

export const filterLast24Hours = (data: PriceData[]): PriceData[] => {
  const now = new Date().getTime();
  const twentyFourHoursAgo = now - HOURS_24;

  return data.filter((item) => {
    const timestamp = new Date(item.timestamp).getTime();
    return timestamp >= twentyFourHoursAgo && timestamp <= now;
  });
};

// Global variable to store cryptocurrencies
export let cryptocurrencies: string[] = [];

// Function to fetch cryptocurrencies and update the global variable
export const fetchCryptocurrencies = async () => {
  try {
    const response = await axios.get<ApiResponse<{ symbol: string }[]>>("/api/list/name");

    if (response.data.meta.code === 200) {
      // Extract only the symbol field from API response
      cryptocurrencies = response.data.data.map((crypto) => crypto.symbol);
      console.log("Fetched cryptocurrencies:", cryptocurrencies);
    } else {
      console.error("Failed to fetch cryptocurrencies:", response.data.meta.message);
    }
  } catch (error) {
    console.error("Error fetching cryptocurrencies:", error);
  }
};

// Fetch cryptocurrencies once when the module is loaded
fetchCryptocurrencies();

export const usePriceData = (symbol: string): PriceState => {
  const [state, setState] = useState<PriceState>({
    currentPrice: null,
    priceHistory: [],
    isLoading: false,
    error: null,
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));

        const [currentPriceResponse, historyResponse] = await Promise.all([
          axios.get<ApiResponse<PriceData>>(`/api/price/latest?symbol=${symbol}&currency=USDT`),
          axios.get<ApiResponse<PriceData[]>>(`/api/price/interval?symbol=${symbol}&currency=USDT&interval=24`),
        ]);

        if (currentPriceResponse.data.meta.code !== 200) {
          throw new Error(currentPriceResponse.data.meta.message || "Failed to fetch current price");
        }

        if (historyResponse.data.meta.code !== 200) {
          throw new Error(historyResponse.data.meta.message || "Failed to fetch price history");
        }

        const currentPrice = currentPriceResponse.data.data?.price ?? null;
        const filteredHistory = Array.isArray(historyResponse.data.data)
          ? filterLast24Hours(historyResponse.data.data)
          : [];

        const sortedHistory = filteredHistory.sort(
          (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
        );

        setState({
          currentPrice,
          priceHistory: sortedHistory,
          isLoading: false,
          error: null,
        });
      } catch (error) {
        let errorMessage = "An unexpected error occurred";

        if (axios.isAxiosError(error)) {
          if (error.response?.data?.meta?.message) {
            errorMessage = error.response.data.meta.message;
          } else if (error.response) {
            errorMessage = "Server error occurred";
          } else if (error.request) {
            errorMessage = "No response from server. Please check your connection.";
          } else {
            errorMessage = error.message;
          }
        } else if (error instanceof Error) {
          errorMessage = error.message;
        }

        setState((prev) => ({
          ...prev,
          isLoading: false,
          error: `Error: ${errorMessage}`,
          currentPrice: null,
          priceHistory: [],
        }));

        console.error("Error fetching price data:", error);
      }
    };

    if (symbol) {
      fetchData();
    }
  }, [symbol]);

  return state;
};

export default usePriceData;
