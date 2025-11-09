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
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="author" label="作者" width="120" />
      <el-table-column prop="file_size" label="文件大小" width="100">
        <template #default="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>
      </el-table-column>
      <el-table-column prop="upload_time" label="上传时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.upload_time || row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
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

const loadArticles = async () => {
  loading.value = true
  try {
    const response = await articleApi.getArticleList({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchQuery.value || undefined,
      author: selectedAuthor.value || undefined
    })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
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
    ElMessage.success('分析任务已创建')
    router.push(`/analysis/${id}`)
  } catch (error) {
    ElMessage.error('创建分析任务失败')
    console.error(error)
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

onMounted(() => {
  loadArticles()
  loadAuthors()
})
</script>

<style scoped>
.article-list {
  padding: 20px;
  width: 100%;
  max-width: 1400px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>