<template>
  <div class="public-form-container">
    <el-card class="form-card" shadow="always">
      <div class="form-header">
        <div class="logo-icon">
          <el-icon :size="40" color="#F56C6C"><Warning /></el-icon>
        </div>
        <h2>Property Help Desk</h2>
        <p>Verified tenants can report issues below.</p>
      </div>

      <el-form :model="form" label-position="top" @submit.prevent="submitComplaint">
        <el-row :gutter="10">
          <el-col :span="12">
            <el-form-item label="Room No." required>
              <el-input v-model="form.room_no" placeholder="e.g. 102" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Your Name" required>
              <el-input v-model="form.tenant_name" placeholder="Name" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Registered Phone Number" required>
          <el-input 
            v-model="form.phone_number" 
            placeholder="Enter your 10-digit number"
            type="tel"
            maxlength="10"
          >
            <template #prefix>+91</template>
          </el-input>
          <small class="helper-text">Must match your admission record</small>
        </el-form-item>

        <el-form-item label="Category" required>
          <el-select v-model="form.category" placeholder="What is the issue?" style="width: 100%">
            <el-option label="Plumbing (Tap, Leakage)" value="Plumbing" />
            <el-option label="Electrical (Fan, Light)" value="Electrical" />
            <el-option label="Internet / WiFi" value="Internet" />
            <el-option label="Cleaning" value="Cleaning" />
            <el-option label="Food / Mess" value="Food" />
            <el-option label="Other" value="Other" />
          </el-select>
        </el-form-item>

        <el-form-item label="Problem Details">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="3" 
            placeholder="Describe the issue in detail..." 
          />
        </el-form-item>

        <el-button 
          type="danger" 
          class="submit-btn" 
          :loading="loading" 
          @click="submitComplaint"
        >
          Verify & Submit
        </el-button>
      </el-form>

      <div class="footer-note">
        Powered by SensiLens PG Manager
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Warning } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const route = useRoute()
const loading = ref(false)

const form = reactive({
  property_id: 0,
  room_no: '',
  tenant_name: '',
  phone_number: '', // Added field
  category: '',
  description: ''
})

onMounted(() => {
  if (route.params.id) {
    form.property_id = parseInt(route.params.id)
  }
})

const submitComplaint = async () => {
  // Validate basic fields
  if (!form.room_no || !form.tenant_name || !form.phone_number || !form.category) {
    return ElMessage.warning("Please fill in all required fields")
  }

  // Validate phone length
  if (form.phone_number.length !== 10) {
    return ElMessage.error("Please enter a valid 10-digit phone number")
  }

  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  
  loading.value = true
  try {
    const response = await axios.post(`${baseURL}/api/public/complaint`, form)
    
    ElMessage.success("Complaint submitted! Our team has verified your record.")
    
    // Clear form
    form.room_no = ''; form.tenant_name = ''; form.phone_number = ''; 
    form.category = ''; form.description = ''
  } catch (err) {
    // Check if backend rejected due to unverified phone
    if (err.response && err.response.status === 403) {
      ElMessage.error("Verification Failed: Phone number not found in this PG's records.")
    } else {
      ElMessage.error("Failed to submit. Please try again.")
    }
    console.error("Submission Error:", err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.public-form-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f2f5;
  padding: 15px;
}
.form-card {
  width: 100%;
  max-width: 400px;
  border-radius: 15px;
}
.form-header {
  text-align: center;
  margin-bottom: 20px;
}
.helper-text {
  color: #909399;
  font-size: 11px;
}
.submit-btn {
  width: 100%;
  padding: 12px;
  font-weight: bold;
  font-size: 16px;
  margin-top: 10px;
}
.footer-note {
  text-align: center;
  margin-top: 20px;
  font-size: 12px;
  color: #c0c4cc;
}
</style>