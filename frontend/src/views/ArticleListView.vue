<template>
  <div class="article-list-view">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1>文章列表</h1>
          <nav class="nav-menu">
            <el-button type="primary" @click="$router.push('/articles/upload')">
              <el-icon><upload-filled /></el-icon>
              上传文章
            </el-button>
            <el-button @click="$router.push('/')">
              <el-icon><house /></el-icon>
              首页
            </el-button>
          </nav>
        </div>
      </el-header>
      
      <el-main>
        <div class="content-section">
          <ArticleList />
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { UploadFilled, House } from '@element-plus/icons-vue'
import ArticleList from '@/components/ArticleList.vue'
</script>

<style scoped>
.article-list-view {
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

.nav-menu {
  display: flex;
  gap: 10px;
}

.el-main {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
}

.content-section {
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  padding: 20px 40px;
  box-sizing: border-box;
}
</style>