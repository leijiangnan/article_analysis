import api from '@/api/request'
import type { ArticleAnalysis, ApiResponse } from '@/types'

export interface CreateAnalysisRequest {
  article_id: string
}

export const analysisApi = {
  // 创建分析任务
  createAnalysis: (data: CreateAnalysisRequest) => {
    return api.post<ApiResponse<{ task_id: string }>>(`/articles/${data.article_id}/analyze`, {})
  },

  // 获取分析结果
  getAnalysisResult: (articleId: string) => {
    return api.get<ApiResponse<ArticleAnalysis>>(`/articles/${articleId}/analysis`)
  },

  // 获取分析状态
  getAnalysisStatus: (taskId: string) => {
    return api.get<ApiResponse<{ status: string; result?: ArticleAnalysis }>>(`/analysis/status/${taskId}`)
  }
}