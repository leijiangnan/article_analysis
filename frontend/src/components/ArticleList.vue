<template>
  <div class="article-list">
    <div class="header">
      <h2>文章列表</h2>
      <div class="filters">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章标题"
          style="width: 200px; margin-right: 10px"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="selectedAuthor"
          placeholder="选择作者"
          style="width: 150px; margin-right: 10px"
          @change="handleAuthorChange"
          clearable
        >
          <el-option
            v-for="author in authors"
            :key="author"
            :label="author"
            :value="author"
          />
        </el-select>
      </div>
    </div>

    <el-table
      :data="articles"
      v-loading="loading"
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column prop="title" label="标题" min-width="300">
          <template #default="{ row }">
            <el-tooltip 
              :content="getCoreViewpoints(row)" 
              placement="top"
              :disabled="!hasAnalysis(row)"
              :show-after="500"
              :max-width="100"
              effect="dark"
            >
              <span class="article-title">{{ row.title }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
      <el-table-column prop="author" label="作者" width="120" />
      <el-table-column prop="file_size" label="文件大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="analysis_status" label="分析状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="getAnalysisStatusType(row.analysis_status)"
            size="small"
          >
            {{ getAnalysisStatusText(row.analysis_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="upload_time" label="上传时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.upload_time || row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            @click.stop="handleAnalyze(row.id)"
          >
            分析
          </el-button>
          <el-button
            type="info"
            size="small"
            @click.stop="handleView(row.id)"
          >
            查看
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click.stop="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
    
    <div class="total-size" v-if="totalSize > 0">
      总文件大小: {{ formatFileSize(totalSize) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { articleApi } from '@/api/article'
import { analysisApi } from '@/api/analysis'
import type { Article } from '@/types'

const router = useRouter()

const loading = ref(false)
const articles = ref<Article[]>([])
const authors = ref<string[]>([])
const searchQuery = ref('')
const selectedAuthor = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalSize = ref(0)
const coreViewpointsCache = ref<Record<number, string>>({})

const loadArticles = async () => {
  loading.value = true
  try {
    const response = await articleApi.getArticleListWithAnalysis({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchQuery.value || undefined,
      author: selectedAuthor.value || undefined
    })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
    
    // 计算总文件大小
    totalSize.value = articles.value.reduce((sum, article) => sum + (article.file_size || 0), 0)
  } catch (error) {
    ElMessage.error('加载文章列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadAuthors = async () => {
  try {
    const response = await articleApi.getAuthors()
    authors.value = response.data
  } catch (error) {
    console.error('加载作者列表失败', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadArticles()
}

const handleAuthorChange = () => {
  currentPage.value = 1
  loadArticles()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  loadArticles()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadArticles()
}

const handleRowClick = (row: Article) => {
  router.push(`/articles/${row.id}`)
}

const handleView = (id: number) => {
  router.push(`/articles/${id}`)
}

const handleAnalyze = async (id: number) => {
  try {
    await analysisApi.createAnalysis({ article_id: id })
    ElMessage.success('分析任务已创建，请稍候...')
    
    // 重新加载文章列表以更新分析状态
    setTimeout(() => {
      loadArticles()
    }, 1000)
    
  } catch (error) {
    ElMessage.error('创建分析任务失败')
    console.error(error)
  }
}

const handleDelete = async (row: Article) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文章 "${row.title}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await articleApi.deleteArticle(row.id)
    ElMessage.success('文章删除成功')
    loadArticles() // 重新加载文章列表
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除文章失败')
      console.error(error)
    }
  }
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

// 获取文章核心观点
  const getCoreViewpoints = (row: Article) => {
    if (row.analysis_status === 'completed' && row.has_analysis) {
      // 检查缓存中是否有数据
      if (coreViewpointsCache.value[row.id]) {
        return coreViewpointsCache.value[row.id]
      }
      
      // 异步加载核心观点数据
      loadCoreViewpoints(row.id)
      return '正在加载核心观点...'
    } else if (row.analysis_status === 'processing') {
      return '分析中，请稍候...'
    } else if (row.analysis_status === 'failed') {
      return '分析失败'
    } else {
      return '尚未分析'
    }
  }

// 加载核心观点数据
const loadCoreViewpoints = async (articleId: number) => {
  // 避免重复加载
  if (coreViewpointsCache.value[articleId]) {
    return
  }
  
  try {
    const response = await analysisApi.getAnalysisResult(articleId)
    if (response.data && response.data.core_viewpoints) {
      coreViewpointsCache.value[articleId] = response.data.core_viewpoints
    } else {
      coreViewpointsCache.value[articleId] = '暂无核心观点数据'
    }
  } catch (error) {
    console.error(`加载文章 ${articleId} 的核心观点失败:`, error)
    coreViewpointsCache.value[articleId] = '加载核心观点失败'
  }
}

// 检查是否有分析数据
const hasAnalysis = (row: Article) => {
  return row.has_analysis || false
}

// 获取分析状态标签类型
const getAnalysisStatusType = (status: string) => {
  switch (status) {
    case 'pending':
      return 'warning'
    case 'processing':
      return ''
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    default:
      return 'info'
  }
}

// 获取分析状态显示文本
const getAnalysisStatusText = (status: string) => {
  switch (status) {
    case 'pending':
      return '待分析'
    case 'processing':
      return '分析中'
    case 'completed':
      return '已完成'
    case 'failed':
      return '分析失败'
    default:
      return '未分析'
  }
}

onMounted(() => {
  loadArticles()
  loadAuthors()
})
</script>

<style scoped>
.article-list {
  width: 100%;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background-color: #fff;
  padding: 0 0 20px;
}

.filters {
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 0;
  display: flex;
  justify-content: center;
  background-color: #fff;
}

.total-size {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background-color: #f5f7fa;
  padding: 10px 15px;
  border-radius: 6px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  font-size: 14px;
  color: #606266;
  z-index: 1000;
}

.article-title {
  cursor: pointer;
  transition: color 0.3s;
}

.article-title:hover {
  color: #409eff;
}

/* 自定义悬浮框样式 */
:deep(.el-tooltip__popper) {
  max-width: 100px !important;
  line-height: 1.3 !important;
  font-size: 12px !important;
  padding: 6px 8px !important;
  word-break: break-word !important;
  white-space: pre-wrap !important;
}

:deep(.el-tooltip__popper.is-dark) {
  background-color: #303133 !important;
}
</style>