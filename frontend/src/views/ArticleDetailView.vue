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
.article-detail-view {
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
  padding: 20px;
  align-self: center;
}
</style>