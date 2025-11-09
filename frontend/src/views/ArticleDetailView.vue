<template>
  <div class="article-detail-view">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>文章详情</h1>
          <nav class="nav-menu">
            <el-button @click="$router.push('/articles')">
              <el-icon><arrow-left /></el-icon>
              返回列表
            </el-button>
            <el-button @click="$router.push('/')">
              <el-icon><house /></el-icon>
              首页
            </el-button>
          </nav>
        </div>
      </el-header>
      
      <el-main>
        <ArticleDetail />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ArrowLeft, House } from '@element-plus/icons-vue'
import ArticleDetail from '@/components/ArticleDetail.vue'
</script>

<style scoped>
/* 全局宽度设置 */
:global(html, body) {
  width: 100%;
  max-width: none;
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
.article-detail-view {
  min-height: 100vh;
  background-color: #f5f7fa;
  width: 100%;
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.el-container {
  width: 100%;
  max-width: 100%;
}

.el-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0;
  width: 100%;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  max-width: 1400px;
  margin: 0 auto;
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
  padding: 20px;
  width: 100%;
  max-width: none;
  display: block;
  margin: 0;
  box-sizing: border-box;
}

/* 确保内容区域正常扩展 */
:deep(.article-detail) {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

/* 覆盖Element Plus的默认布局 */
:deep(.el-container) {
  width: 100% !important;
  max-width: 100% !important;
}

:deep(.el-main) {
  width: 100% !important;
  max-width: none !important;
  padding: 20px !important;
  margin: 0 !important;
  box-sizing: border-box !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 15px;
  }
  
  .header-content h1 {
    font-size: 20px;
  }
  
  .el-main {
    padding: 15px;
  }
}
</style>