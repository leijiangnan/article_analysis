<template>
  <div class="article-upload">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>录入文章</h1>
          <el-button @click="$router.push('/articles')">
            <el-icon><arrow-left /></el-icon>
            返回列表
          </el-button>
        </div>
      </el-header>
      
      <el-main>
        <div class="upload-section">
          <el-card>
            <template #header>
              <el-tabs v-model="activeTab" class="upload-tabs">
                <el-tab-pane label="文件上传" name="file">
                  <h3>选择文件上传</h3>
                </el-tab-pane>
                <el-tab-pane label="文本输入" name="text">
                  <h3>直接输入文章内容</h3>
                </el-tab-pane>
              </el-tabs>
            </template>
            
            <!-- 文件上传面板 -->
            <div v-if="activeTab === 'file'" class="tab-content">
              <FileUpload @upload-success="handleUploadSuccess" @upload-error="handleUploadError" />
            </div>
            
            <!-- 文本输入面板 -->
            <div v-if="activeTab === 'text'" class="tab-content">
              <TextInput @create-success="handleCreateSuccess" @create-error="handleCreateError" />
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
import FileUpload from '@/components/FileUpload.vue'
import TextInput from '@/components/TextInput.vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeTab = ref('file')

const handleUploadSuccess = (data: any) => {
  ElMessage.success('文件上传成功！')
  // 上传成功后跳转到文章详情页
  setTimeout(() => {
    router.push(`/articles/${data.id}`)
  }, 1000)
}

const handleUploadError = (error: any) => {
  console.error('上传失败:', error)
}

const handleCreateSuccess = (data: any) => {
  ElMessage.success('文章创建成功！')
  // 创建成功后跳转到文章详情页
  setTimeout(() => {
    router.push(`/articles/${data.id}`)
  }, 1000)
}

const handleCreateError = (error: any) => {
  console.error('创建失败:', error)
  console.error('创建失败详情:', {
    message: error.message,
    response: error.response,
    data: error.response?.data
  })
}
</script>

<style scoped>
.article-upload {
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

.upload-section {
  max-width: 1600px;
  margin: 0 auto;
  width: 100%;
}

.upload-tabs {
  margin: -20px -20px 20px -20px;
}

.upload-tabs .el-tabs__header {
  margin-bottom: 0;
}

.tab-content {
  min-height: 400px;
}
</style>