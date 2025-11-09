<template>
  <div class="article-detail-container">
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
    </el-card>

    <!-- 并排展示区域 -->
    <div v-if="article" class="content-wrapper">
      <!-- 文章内容区域 -->
      <el-card class="content-card article-content-card">
        <template #header>
          <h3>文章内容</h3>
        </template>
        <div class="content-text">
          {{ article.content }}
        </div>
      </el-card>

      <!-- 分析结果区域 -->
      <el-card v-if="analysis" class="content-card analysis-content-card">
        <template #header>
          <h3>分析结果</h3>
        </template>
        <div class="analysis-content">
          <div class="analysis-section">
            <h4>核心观点</h4>
            <div class="viewpoints-text">
              {{ analysis.core_viewpoints }}
            </div>
          </div>

          <div class="analysis-section">
            <h4>文件结构</h4>
            <div class="structure-text">
              {{ analysis.file_structure }}
            </div>
          </div>

          <div class="analysis-section">
            <h4>作者思路</h4>
            <div class="thoughts-text">
              {{ analysis.author_thoughts }}
            </div>
          </div>

          <div class="analysis-section">
            <h4>相关材料</h4>
            <div class="materials-text">
              {{ analysis.related_materials }}
            </div>
          </div>

          <div class="analysis-section" v-if="analysis.error_message">
            <h4>错误信息</h4>
            <el-alert
              :title="analysis.error_message"
              type="error"
              :closable="false"
            />
          </div>
        </div>
      </el-card>

      <!-- 无分析结果时的占位 -->
      <el-card v-else class="content-card no-analysis-card">
        <template #header>
          <h3>分析结果</h3>
        </template>
        <div class="no-analysis-content">
          <el-empty description="暂无分析结果" />
          <el-button type="primary" @click="handleAnalyze" :loading="analyzing">
            开始分析
          </el-button>
        </div>
      </el-card>
    </div>

    <div v-if="loading" class="loading-container">
      <el-loading text="加载中..." />
    </div>
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



onMounted(() => {
  loadArticle()
  loadAnalysis()
})
</script>

<style scoped>
.article-detail-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0;
  box-sizing: border-box;
}

.article-detail {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
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

/* 并排布局样式 */
.content-wrapper {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 20px;
  margin-top: 20px;
  width: 100%;
  box-sizing: border-box;
}

.content-card {
  height: fit-content;
  max-height: 80vh;
  width: 100%;
  box-sizing: border-box;
}

.article-content-card {
  overflow: hidden;
}

.analysis-content-card {
  overflow: hidden;
}

.content-text {
  white-space: pre-wrap;
  line-height: 1.6;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  max-height: calc(80vh - 120px);
  overflow-y: auto;
}

.analysis-content {
  max-height: calc(80vh - 120px);
  overflow-y: auto;
}

.analysis-section {
  margin-bottom: 20px;
}

.analysis-section h4 {
  margin-bottom: 10px;
  color: #303133;
  font-size: 16px;
}

.viewpoints-text,
.structure-text,
.thoughts-text,
.materials-text {
  line-height: 1.6;
  color: #606266;
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 4px;
  white-space: pre-wrap;
  max-height: 200px;
  overflow-y: auto;
  font-size: 14px;
}

.no-analysis-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.no-analysis-content {
  text-align: center;
  padding: 40px;
}

.no-analysis-content .el-button {
  margin-top: 20px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .content-wrapper {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .content-card {
    max-height: 60vh;
  }
  
  .content-text,
  .analysis-content {
    max-height: calc(60vh - 120px);
  }
}

@media (max-width: 768px) {
  .article-detail {
    padding: 15px;
    max-width: 100%;
  }
  
  .content-wrapper {
    gap: 10px;
  }
  
  .content-card {
    max-height: 50vh;
  }
  
  .content-text,
  .analysis-content {
    max-height: calc(50vh - 100px);
    padding: 10px;
  }
  
  .viewpoints-text,
  .structure-text,
  .thoughts-text,
  .materials-text {
    max-height: 150px;
    font-size: 13px;
  }
}
</style>