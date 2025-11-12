<template>
  <div class="text-input-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文本输入</span>
          <el-button @click="handleReset" type="info" size="small">重置</el-button>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="文章标题" prop="title">
          <el-input 
            v-model="form.title" 
            placeholder="请输入文章标题（可选，不输入将自动提取）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="作者" prop="author">
          <el-input 
            v-model="form.author" 
            placeholder="请输入作者（可选，不输入将自动提取）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="文章内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="20"
            placeholder="请直接输入或粘贴文章内容到此处..."
            maxlength="10485760"
            show-word-limit
            :autosize="{ minRows: 15, maxRows: 30 }"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            <el-icon><document-add /></el-icon>
            创建文章
          </el-button>
          <el-button @click="handleReset">
            <el-icon><refresh /></el-icon>
            重置
          </el-button>
          <el-button @click="testApi" type="warning" size="small">
            测试API
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- 字数统计 -->
      <div class="word-count" v-if="form.content">
        <el-text type="info">
          字数：{{ form.content.length }} | 
          字符数：{{ form.content.length }} |
          预估阅读时间：{{ Math.ceil(form.content.length / 500) }} 分钟
        </el-text>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { DocumentAdd, Refresh } from '@element-plus/icons-vue'
import { articleApi } from '@/api/article'

const emit = defineEmits<{
  createSuccess: [data: any]
  createError: [error: any]
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  title: '',
  author: '',
  content: ''
})

const rules: FormRules = {
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  console.log('开始处理表单提交...')
  if (!formRef.value) {
    console.log('表单引用不存在')
    return
  }
  
  await formRef.value.validate(async (valid: boolean) => {
    console.log('表单验证结果:', valid)
    if (!valid) {
      console.log('表单验证失败')
      return
    }
    
    console.log('表单验证通过，开始提交数据...')
    console.log('表单数据:', {
      title: form.title,
      author: form.author,
      contentLength: form.content.length
    })
    
    // 内容非空检查
    if (!form.content.trim()) {
      ElMessage.error('请输入有效的文章内容')
      return
    }
    
    loading.value = true
    
    try {
      const requestData = {
        title: form.title.trim(),
        author: form.author.trim(),
        content: form.content.trim()
      }
      console.log('请求数据:', requestData)
      
      const response = await articleApi.createArticle(requestData)
      console.log('API响应数据:', response)
      
      // 检查响应数据结构 - 响应拦截器已经返回了response.data
      if (response) {
        console.log('响应数据:', response)
        
        // 直接处理响应数据
        console.log('文章创建成功，响应数据:', response)
        // 不再显示成功消息，由父组件处理
        emit('createSuccess', response)
        // 重置表单
        handleReset()
      } else {
        console.log('响应数据格式错误:', response)
        throw new Error('服务器响应格式错误')
      }
    } catch (error: any) {
      console.error('创建文章失败:', error)
      console.error('错误详情:', {
        message: error.message,
        response: error.response,
        data: error.response?.data,
        status: error.response?.status,
        stack: error.stack
      })
      const errorMessage = error.response?.data?.message || error.message || '创建文章失败'
      ElMessage.error(errorMessage)
      emit('createError', error)
    } finally {
      loading.value = false
    }
  })
}

const handleReset = () => {
  form.title = ''
  form.author = ''
  form.content = ''
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 测试API连接
const testApi = async () => {
  try {
    console.log('开始测试API连接...')
    const response = await articleApi.createArticle({
      title: 'API测试文章',
      author: 'API测试',
      content: '这是一个API测试文章内容。'
    })
    console.log('API测试成功:', response)
    ElMessage.success('API连接测试成功！')
  } catch (error: any) {
    console.error('API测试失败:', error)
    console.error('API测试失败详情:', {
      message: error.message,
      response: error.response,
      data: error.response?.data,
      status: error.response?.status
    })
    ElMessage.error('API连接测试失败: ' + (error.message || '未知错误'))
  }
}
</script>

<style scoped>
.text-input-container {
  padding: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.word-count {
  margin-top: 20px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
  text-align: center;
}

:deep(.el-textarea__inner) {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  line-height: 1.6;
  resize: vertical;
}
</style>