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
            <div class="hero-section">
              <div class="hero-image">
                <img src="@/assets/analysis-illustration.svg" alt="文档分析插图" class="analysis-illustration"/>
              </div>
              <div class="hero-description">
                <p class="description-text">
                  智能文档分析平台，帮助您快速处理和分析各类文章文件。
                  支持多种格式上传，AI智能分析，一键生成详细报告。
                </p>
                <div class="hero-stats">
                  <div class="stat-item">
                    <div class="stat-number">50+</div>
                    <div class="stat-label">支持格式</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-number">AI</div>
                    <div class="stat-label">智能分析</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-number">秒级</div>
                    <div class="stat-label">响应速度</div>
                  </div>
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
import { UploadFilled, Document } from '@element-plus/icons-vue'
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.el-header {
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  padding: 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
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
  padding: 0;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-section {
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  padding: 20px 40px;
}

.welcome-card {
  text-align: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border-radius: 20px;
  overflow: hidden;
}

.welcome-card h2 {
  margin: 0;
  color: #2c3e50;
  font-size: 32px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-section {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 80px;
  padding: 80px 60px;
  min-height: 600px;
  max-width: 100%;
}

.hero-image {
  flex: 1.5;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 400px;
}

.analysis-illustration {
  max-width: 100%;
  height: auto;
  max-height: 500px;
  width: auto;
  max-width: 600px;
  filter: drop-shadow(0 16px 32px rgba(0, 0, 0, 0.2));
  transition: transform 0.3s ease, filter 0.3s ease;
}

.analysis-illustration:hover {
  transform: scale(1.05);
  filter: drop-shadow(0 12px 24px rgba(0, 0, 0, 0.15));
}

.hero-description {
  flex: 2;
  text-align: left;
  min-width: 400px;
}

.description-text {
  font-size: 24px;
  line-height: 1.7;
  color: #4a5568;
  margin-bottom: 60px;
  text-align: center;
  font-weight: 400;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.hero-stats {
  display: flex;
  justify-content: center;
  gap: 50px;
  margin-top: 20px;
}

.stat-item {
  text-align: center;
  padding: 30px 25px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: white;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  transition: transform 0.3s ease;
  min-width: 120px;
}

.stat-item:hover {
  transform: translateY(-5px);
}

.stat-number {
  font-size: 42px;
  font-weight: bold;
  margin-bottom: 12px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.stat-label {
  font-size: 16px;
  opacity: 0.9;
  font-weight: 500;
  letter-spacing: 0.5px;
}

@media (max-width: 768px) {
  .hero-section {
    flex-direction: column;
    gap: 30px;
  }
  
  .hero-stats {
    gap: 20px;
  }
  
  .stat-item {
    padding: 15px;
  }
  
  .stat-number {
    font-size: 24px;
  }
}
</style>
