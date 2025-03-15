export interface PriceData {
  timestamp: string
  symbol: string
  currency: string
  price: number
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