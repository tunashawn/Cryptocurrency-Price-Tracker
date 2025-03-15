export interface PriceData {
  timestamp: string
  symbol: string
  currency: string
  latest_price: number
}

export interface ApiResponse<T> {
  meta: {
    code: number
    message: string
  }
  data: T
}

export interface PriceState {
  currentPrice: number | null
  priceHistory: PriceData[]
  isLoading: boolean
  error: string | null
}