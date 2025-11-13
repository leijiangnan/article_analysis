export interface Article {
  id: string
  title: string
  author: string
  content: string
  file_path: string
  file_size: number
  created_at: string
  updated_at: string
  analysis_status?: string
  has_analysis?: boolean
}

export interface ArticleAnalysis {
  id: string
  article_id: string
  core_viewpoints: string
  file_structure: string
  author_thoughts: string
  related_materials: string
  analysis_status: string
  analysis_time?: string
  error_message: string
  created_at: string
  updated_at: string
  article?: Article
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