import api from '@/api/request'
import type { Article, ArticleListResponse, ApiResponse } from '@/types'

export interface UploadArticleRequest {
  file: File
}

export interface CreateArticleRequest {
  title: string
  author: string
  content: string
}

export interface ArticleListParams {
  page?: number
  page_size?: number
  keyword?: string
  author?: string
}

export const articleApi = {
  // 上传文章
  uploadArticle: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post<ApiResponse<Article>>('/articles/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 创建文章（直接输入文本）
  createArticle: (data: CreateArticleRequest) => {
    return api.post<ApiResponse<Article>>('/articles/create', data)
  },

  // 获取文章列表
  getArticleList: (params: ArticleListParams = {}) => {
    return api.get<ApiResponse<ArticleListResponse>>('/articles', { params })
  },

  // 获取文章列表及分析状态
  getArticleListWithAnalysis: (params: ArticleListParams = {}) => {
    return api.get<ApiResponse<ArticleListResponse>>('/articles/with-analysis', { params })
  },

  // 获取文章详情
  getArticleDetail: (id: number) => {
    return api.get<ApiResponse<Article>>(`/articles/${id}`)
  },

  // 获取作者列表
  getAuthors: () => {
    return api.get<ApiResponse<string[]>>('/articles/authors')
  },

  // 删除文章
  deleteArticle: (id: number) => {
    return api.delete<ApiResponse<null>>(`/articles/${id}`)
  }
}