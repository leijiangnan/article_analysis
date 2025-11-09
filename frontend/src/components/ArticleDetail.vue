<template>
  <div class="article-detail">
    <el-card v-if="article" class="article-card">
      <template #header>
        <div class="card-header">
          <h2>{{ article.title }}</h2>
          <div class="header-actions">
            <el-button
              type="primary"
              @click="handleAnalyze"
              :loading="analyzing"
            >
              分析文章
            </el-button>
            <el-button @click="goBack">
              返回列表
            </el-button>
          </div>
        </div>
      </template>

      <div class="article-info">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="info-item">
              <label>作者：</label>
              <span>{{ article.author }}</span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="info-item">
              <label>文件大小：</label>
              <span>{{ formatFileSize(article.file_size) }}</span>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="info-item">
              <label>上传时间：</label>
              <span>{{ formatDate(article.created_at) }}</span>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="article-content">
        <h3>文章内容</h3>
        <div class="content-text">
          {{ article.content }}
        </div>
      </div>
    </el-card>

    <!-- 分析结果 -->
    <el-card v-if="analysis" class="analysis-card">
      <template #header>
        <h3>分析结果</h3>
      </template>

      <div class="analysis-content">
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="analysis-section">
              <h4>基本信息</h4>
              <div class="info-grid">
                <div class="info-item">
                  <label>分类：</label>
                  <el-tag>{{ analysis.category }}</el-tag>
                </div>
                <div class="info-item">
                  <label>情感倾向：</label>
                  <el-tag :type="getSentimentType(analysis.sentiment)">
                    {{ analysis.sentiment }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <label>字数：</label>
                  <span>{{ analysis.word_count }}</span>
                </div>
                <div class="info-item">
                  <label>预计阅读时间：</label>
                  <span>{{ analysis.reading_time }} 分钟</span>
                </div>
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="analysis-section">
              <h4>文章摘要</h4>
              <div class="summary-text">
                {{ analysis.summary }}
              </div>
            </div>
          </el-col>
        </el-row>

        <div class="analysis-section">
          <h4>关键要点</h4>
          <el-tag
            v-for="(point, index) in analysis.key_points"
            :key="index"
            type="info"
            class="key-point-tag"
          >
            {{ point }}
          </el-tag>
        </div>
      </div>
    </el-card>

    <div v-if="loading" class="loading-container">
      <el-loading text="加载中..." />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { articleApi } from '@/api/article'
import { analysisApi } from '@/api/analysis'
import type { Article, ArticleAnalysis } from '@/types'

const router = useRouter()
const route = useRoute()

const loading = ref(false)
const analyzing = ref(false)
const article = ref<Article | null>(null)
const analysis = ref<ArticleAnalysis | null>(null)

const articleId = Number(route.params.id)

const loadArticle = async () => {
  loading.value = true
  try {
    const response = await articleApi.getArticleDetail(articleId)
    article.value = response.data
  } catch (error) {
    ElMessage.error('加载文章失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadAnalysis = async () => {
  try {
    const response = await analysisApi.getAnalysisResult(articleId)
    analysis.value = response.data
  } catch (error) {
    // 分析结果不存在是正常的，不显示错误
    console.log('分析结果不存在')
  }
}

const handleAnalyze = async () => {
  analyzing.value = true
  try {
    await analysisApi.createAnalysis({ article_id: articleId })
    ElMessage.success('分析任务已创建，请稍候...')
    
    // 轮询检查分析状态
    const checkStatus = async () => {
      try {
        const response = await analysisApi.getAnalysisResult(articleId)
        analysis.value = response.data
        ElMessage.success('分析完成！')
      } catch (error) {
        // 如果分析还未完成，继续轮询
        setTimeout(checkStatus, 2000)
      }
    }
    
    setTimeout(checkStatus, 3000)
  } catch (error) {
    ElMessage.error('创建分析任务失败')
    console.error(error)
  } finally {
    analyzing.value = false
  }
}

const goBack = () => {
  router.push('/articles')
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const getSentimentType = (sentiment: string) => {
  const typeMap: Record<string, string> = {
    'positive': 'success',
    'negative': 'danger',
    'neutral': 'info'
  }
  return typeMap[sentiment] || 'info'
}

onMounted(() => {
  loadArticle()
  loadAnalysis()
})
</script>

<style scoped>
.article-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.article-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.article-info {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.info-item {
  margin-bottom: 10px;
}

.info-item label {
  font-weight: bold;
  margin-right: 8px;
  color: #606266;
}

.article-content {
  margin-top: 20px;
}

.content-text {
  white-space: pre-wrap;
  line-height: 1.6;
  max-height: 400px;
  overflow-y: auto;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.analysis-card {
  margin-top: 20px;
}

.analysis-section {
  margin-bottom: 20px;
}

.analysis-section h4 {
  margin-bottom: 15px;
  color: #303133;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.summary-text {
  line-height: 1.6;
  color: #606266;
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
}

.key-point-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}
</style>