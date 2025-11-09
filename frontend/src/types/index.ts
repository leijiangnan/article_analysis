export interface Article {
  id: number
  title: string
  author: string
  content: string
  file_path: string
  file_size: number
  created_at: string
  updated_at: string
}

export interface ArticleAnalysis {
  id: number
  article_id: number
  summary: string
  key_points: string[]
  sentiment: string
  category: string
  word_count: number
  reading_time: number
  created_at: string
  updated_at: string
}

export interface PaginationRequest {
  page?: number
  page_size?: number
}

export interface PaginationResponse {
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface ArticleListResponse extends PaginationResponse {
  list: Article[]
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}