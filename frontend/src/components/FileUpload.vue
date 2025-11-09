<template>
  <div class="upload-container">
    <el-upload
      class="upload-demo"
      drag
      :action="uploadAction"
      :before-upload="beforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      :headers="uploadHeaders"
      accept=".txt,.pdf,.doc,.docx"
      :limit="1"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        拖拽文件到此处或 <em>点击上传</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          支持 txt/pdf/doc/docx 格式，文件大小不超过 10MB
        </div>
      </template>
    </el-upload>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { UploadProps } from 'element-plus'

const emit = defineEmits<{
  uploadSuccess: [data: any]
  uploadError: [error: any]
}>()

const uploadAction = ref('http://localhost:8080/api/v1/articles/upload')
const uploadHeaders = ref({
  'Accept': 'application/json'
})

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  const allowedTypes = ['text/plain', 'application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document']
  
  if (!allowedTypes.includes(file.type) && !file.name.match(/\.(txt|pdf|doc|docx)$/)) {
    ElMessage.error('请上传正确的文件格式!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB!')
    return false
  }
  return true
}

const handleSuccess: UploadProps['onSuccess'] = (response) => {
  emit('uploadSuccess', response.data)
}

const handleError: UploadProps['onError'] = (error) => {
  ElMessage.error('文件上传失败: ' + error.message)
  emit('uploadError', error)
}
</script>

<style scoped>
.upload-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.upload-demo {
  width: 100%;
}
</style>