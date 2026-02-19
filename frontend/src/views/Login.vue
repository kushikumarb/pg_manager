<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="login-header">
        <el-icon :size="40" color="#409EFF"><HomeFilled /></el-icon>
        <h2>PG Manager</h2>
        <p>Owner Administration Portal</p>
      </div>

      <el-form :model="loginForm" :rules="rules" ref="loginRef" label-position="top">
        <el-form-item label="Email Address" prop="email">
          <el-input 
            v-model="loginForm.email" 
            placeholder="Enter registered email"
            prefix-icon="Message"
          />
        </el-form-item>

        <el-form-item label="Password" prop="password">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="Enter password"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-button 
          type="primary" 
          class="login-button" 
          @click="handleLogin" 
          :loading="loading"
        >
          Login to Dashboard
        </el-button>
      </el-form>
      
      <div class="login-footer">
        <span>Don't have an account? </span>
        <el-link type="primary">Register Property</el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { HomeFilled, Phone, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const loginRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  email: '',
  password: ''
})

const rules = {
  phone: [{ required: true, message: 'Please input phone number', trigger: 'blur' }],
  password: [{ required: true, message: 'Please input password', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginRef.value) return
  
  await loginRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      const success = await authStore.login(loginForm)
      loading.value = false
      
      if (success) {
        ElMessage.success('Welcome back, Owner!')
        router.push('/property-selector')
      } else {
        ElMessage.error('Invalid credentials. Please try again.')
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
  border-radius: 12px;
  padding: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  margin: 10px 0 5px;
  color: #303133;
}

.login-header p {
  color: #909399;
  font-size: 14px;
}

.login-button {
  width: 100%;
  margin-top: 10px;
  height: 45px;
  font-size: 16px;
}

.login-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
}
</style>