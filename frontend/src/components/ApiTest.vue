<template>
  <div class="api-test">
    <h3>API测试工具</h3>
    <el-button @click="testDirectApi" type="primary">直接测试API</el-button>
    <el-button @click="testWithProxy" type="success">测试代理API</el-button>
    <div v-if="result" class="result">
      <h4>测试结果:</h4>
      <pre>{{ JSON.stringify(result, null, 2) }}</pre>
    </div>
    <div v-if="error" class="error">
      <h4>错误信息:</h4>
      <pre>{{ error }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const result = ref<any>(null)
const error = ref<string>('')

const testDirectApi = async () => {
  try {
    error.value = ''
    result.value = null
    
    console.log('开始直接API测试...')
    const response = await axios.post('http://localhost:8080/api/articles/create', {
      title: '直接API测试',
      author: '测试用户',
      content: '这是直接API测试的内容'
    })
    
    console.log('直接API响应:', response)
    result.value = response
  } catch (err: any) {
    console.log('直接API测试失败:', err)
    error.value = `错误: ${err.message}\n响应: ${JSON.stringify(err.response?.data, null, 2)}`
  }
}

const testWithProxy = async () => {
  try {
    error.value = ''
    result.value = null
    
    console.log('开始代理API测试...')
    const response = await axios.post('/api/articles/create', {
      title: '代理API测试',
      author: '测试用户',
      content: '这是代理API测试的内容'
    })
    
    console.log('代理API响应:', response)
    result.value = response
  } catch (err: any) {
    console.log('代理API测试失败:', err)
    error.value = `错误: ${err.message}\n响应: ${JSON.stringify(err.response?.data, null, 2)}`
  }
}
</script>

<style scoped>
.api-test {
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin: 20px 0;
}

.result, .error {
  margin-top: 20px;
  padding: 15px;
  border-radius: 4px;
}

.result {
  background-color: #f0f9ff;
  border: 1px solid #1890ff;
}

.error {
  background-color: #fff2f0;
  border: 1px solid #ff4d4f;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>