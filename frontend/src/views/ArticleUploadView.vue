<template>
  <div class="article-upload">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>上传文章</h1>
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
              <h3>选择文件上传</h3>
            </template>
            <FileUpload @upload-success="handleUploadSuccess" @upload-error="handleUploadError" />
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ArrowLeft } from '@element-plus/icons-vue'
import FileUpload from '@/components/FileUpload.vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const router = useRouter()

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
  max-width: 800px;
  margin: 0 auto;
}
</style>