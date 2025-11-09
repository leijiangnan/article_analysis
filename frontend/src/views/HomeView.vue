<template>
  <div class="home">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>文章分析系统</h1>
          <nav class="nav-menu">
            <el-button type="primary" @click="$router.push('/articles/upload')">
              <el-icon><upload-filled /></el-icon>
              上传文章
            </el-button>
            <el-button @click="$router.push('/articles')">
              <el-icon><document /></el-icon>
              文章列表
            </el-button>
          </nav>
        </div>
      </el-header>
      
      <el-main>
        <div class="welcome-section">
          <el-card class="welcome-card">
            <template #header>
              <h2>欢迎使用文章分析系统</h2>
            </template>
            <div class="feature-list">
              <div class="feature-item">
                <el-icon><upload-filled /></el-icon>
                <div class="feature-content">
                  <h3>文件上传</h3>
                  <p>支持上传 txt、pdf、doc、docx 格式的文章文件</p>
                </div>
              </div>
              <div class="feature-item">
                <el-icon><search /></el-icon>
                <div class="feature-content">
                  <h3>智能分析</h3>
                  <p>基于 AI 技术对文章进行深度分析，提取关键信息</p>
                </div>
              </div>
              <div class="feature-item">
                <el-icon><trend-charts /></el-icon>
                <div class="feature-content">
                  <h3>分析结果</h3>
                  <p>提供文章摘要、关键要点、情感分析等详细报告</p>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { UploadFilled, Document, Search, TrendCharts } from '@element-plus/icons-vue'
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f7fa;
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

.nav-menu {
  display: flex;
  gap: 10px;
}

.el-main {
  padding: 40px 20px;
}

.welcome-section {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-card {
  text-align: center;
}

.welcome-card h2 {
  margin: 0;
  color: #303133;
  font-size: 28px;
}

.feature-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  margin-top: 40px;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  transition: transform 0.2s;
}

.feature-item:hover {
  transform: translateY(-2px);
}

.feature-item .el-icon {
  font-size: 32px;
  color: #409eff;
  margin-top: 5px;
}

.feature-content h3 {
  margin: 0 0 10px 0;
  color: #303133;
  font-size: 18px;
}

.feature-content p {
  margin: 0;
  color: #606266;
  line-height: 1.5;
}
</style>
