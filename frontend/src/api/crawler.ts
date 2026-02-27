import api from '@/api/request'
import type { ApiResponse } from '@/types'

export interface CrawlRequest {
  url: string
  count: number
}

export interface CrawledArticleSummary {
  id?: number
  title: string
  author: string
  date?: string
  url: string
  saved: boolean
}

export interface CrawlResult {
  total_found: number
  crawled_count: number
  saved_count: number
  articles: CrawledArticleSummary[]
  errors?: string[]
}

export const crawlerApi = {
  crawlArticles: (data: CrawlRequest) => {
    return api.post<ApiResponse<CrawlResult>>('/crawl', data)
  }
}
