<template>
  <div class="article-crawler">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>爬取文章</h1>
          <el-button @click="$router.push('/articles')">
            <el-icon><arrow-left /></el-icon>
            返回列表
          </el-button>
        </div>
      </el-header>
      
      <el-main>
        <div class="crawler-section">
          <el-card>
            <template #header>
              <h3>从网页爬取文章</h3>
            </template>
            
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
              <el-form-item label="起始URL" prop="url">
                <el-input 
                  v-model="form.url" 
                  placeholder="请输入要爬取的网页URL，例如: https://example.com/articles"
                  clearable
                />
              </el-form-item>
              
              <el-form-item label="爬取数量" prop="count">
                <el-input-number 
                  v-model="form.count" 
                  :min="1" 
                  :max="100" 
                  :step="1"
                />
                <span class="form-tip">（开发环境限制为2篇）</span>
              </el-form-item>
              
              <el-form-item>
                <el-button 
                  type="primary" 
                  @click="handleCrawl" 
                  :loading="loading"
                  :disabled="!form.url || form.count < 1"
                >
                  <el-icon><download /></el-icon>
                  开始爬取
                </el-button>
                <el-button @click="resetForm">
                  <el-icon><refresh /></el-icon>
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
          
          <el-card v-if="result" class="result-card">
            <template #header>
              <h3>爬取结果</h3>
            </template>
            
            <el-descriptions :column="3" border>
              <el-descriptions-item label="发现链接">{{ result.total_found }}</el-descriptions-item>
              <el-descriptions-item label="成功爬取">{{ result.crawled_count }}</el-descriptions-item>
              <el-descriptions-item label="保存成功">{{ result.saved_count }}</el-descriptions-item>
            </el-descriptions>
            
            <div v-if="result.articles && result.articles.length > 0" class="articles-list">
              <h4>爬取的文章</h4>
              <el-table :data="result.articles" stripe style="width: 100%">
                <el-table-column prop="title" label="标题" min-width="200">
                  <template #default="scope">
                    <span v-if="scope.row.saved" class="saved-title">
                      <el-icon><check /></el-icon>
                      {{ scope.row.title }}
                    </span>
                    <span v-else class="not-saved-title">{{ scope.row.title }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="author" label="作者" width="150" />
                <el-table-column prop="date" label="日期" width="150" />
                <el-table-column prop="saved" label="状态" width="100">
                  <template #default="scope">
                    <el-tag v-if="scope.row.saved" type="success">已保存</el-tag>
                    <el-tag v-else type="info">未保存</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="120">
                  <template #default="scope">
                    <el-button 
                      v-if="scope.row.saved && scope.row.id" 
                      type="primary" 
                      size="small"
                      @click="viewArticle(scope.row.id)"
                    >
                      查看
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            
            <div v-if="result.errors && result.errors.length > 0" class="errors-section">
              <el-alert 
                v-for="(error, index) in result.errors" 
                :key="index"
                :title="error"
                type="warning"
                :closable="false"
                show-icon
                class="error-item"
              />
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ArrowLeft, Download, Refresh, Check } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useRouter } from 'vue-router'
import { crawlerApi } from '@/api/crawler'
import type { CrawlResult } from '@/api/crawler'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const result = ref<CrawlResult | null>(null)

const form = reactive({
  url: '',
  count: 2
})

const rules: FormRules = {
  url: [
    { required: true, message: '请输入URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  count: [
    { required: true, message: '请输入爬取数量', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '数量必须在1-100之间', trigger: 'blur' }
  ]
}

const handleCrawl = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    result.value = null
    
    try {
      const response = await crawlerApi.crawlArticles({
        url: form.url,
        count: form.count
      })
      
      const data = response as any
      if (data.code === 200) {
        result.value = data.data
        ElMessage.success(`爬取完成！成功保存 ${data.data.saved_count} 篇文章`)
      } else {
        ElMessage.error(data.message || '爬取失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '爬取失败，请检查URL是否正确')
    } finally {
      loading.value = false
    }
  })
}

const resetForm = () => {
  formRef.value?.resetFields()
  result.value = null
}

const viewArticle = (id: number) => {
  router.push(`/articles/${id}`)
}
</script>

<style scoped>
.article-crawler {
  min-height: 100vh;
  background-color: #f5f7fa;
  width: 100%;
}

.el-container {
  width: 100%;
}

.el-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
}

.header-content h1 {
  margin: 0;
  color: #303133;
  font-size: 24px;
}

.el-main {
  padding: 40px 20px;
  display: flex;
  justify-content: center;
}

.crawler-section {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.result-card {
  margin-top: 20px;
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}

.articles-list {
  margin-top: 20px;
}

.articles-list h4 {
  margin-bottom: 15px;
  color: #303133;
}

.saved-title {
  color: #67c23a;
  display: flex;
  align-items: center;
  gap: 5px;
}

.not-saved-title {
  color: #909399;
}

.errors-section {
  margin-top: 20px;
}

.error-item {
  margin-bottom: 10px;
}
</style>
