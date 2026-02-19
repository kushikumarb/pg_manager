<template>
  <div class="history-container">
    <header class="page-header">
      <div class="title-section">
        <el-button link @click="router.back()" class="back-link">
          <el-icon><ArrowLeft /></el-icon> Back
        </el-button>
        <h2>Financial Audit Log</h2>
      </div>
      
      <div class="controls">
        <el-input
          v-model="searchQuery"
          placeholder="Search by Tenant Name..."
          class="search-bar"
          clearable
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" @click="fetchHistory" :loading="loading">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </div>
    </header>

    <el-card shadow="never" class="table-card">
      <el-table 
        :data="filteredHistory" 
        style="width: 100%" 
        v-loading="loading"
        border
        stripe
      >
        <el-table-column label="Date & Time" width="200">
          <template #default="scope">
            <div class="date-cell">
              <el-icon><Calendar /></el-icon>
              <span>{{ formatDate(scope.row.date) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="tenant_name" label="Tenant Name" min-width="180">
          <template #default="scope">
            <span class="tenant-name">{{ scope.row.tenant_name || 'N/A' }}</span>
            <small class="id-tag">(ID: {{ scope.row.tenant_id }})</small>
          </template>
        </el-table-column>

        <el-table-column label="Amount Paid" width="150">
          <template #default="scope">
            <span class="amount-positive">₹{{ scope.row.amount.toLocaleString() }}</span>
          </template>
        </el-table-column>

        <el-table-column label="Method" width="140">
          <template #default="scope">
            <el-tag :type="getMethodTag(scope.row.method)" effect="dark">
              {{ scope.row.method }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="payment_type" label="Payment Category" width="160">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.payment_type }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Search, Refresh, Calendar } from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const history = ref([])
const loading = ref(false)
const searchQuery = ref('')

const fetchHistory = async () => {
  loading.value = true
  try {
    const res = await api.get('/payments/history')
    history.value = res.data
  } catch (e) {
    console.error("Fetch error:", e)
  } finally {
    loading.value = false
  }
}

const filteredHistory = computed(() => {
  if (!searchQuery.value) return history.value
  return history.value.filter(item => 
    item.tenant_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('en-IN', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

const getMethodTag = (method) => {
  if (method === 'Razorpay') return 'primary'
  if (method === 'Cash') return 'success'
  return 'warning'
}

onMounted(fetchHistory)
</script>

<style scoped>
.history-container { padding: 20px; max-width: 1200px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
.title-section h2 { margin: 0; font-size: 24px; color: #303133; }
.back-link { margin-bottom: 5px; padding: 0; }
.controls { display: flex; gap: 12px; }
.search-bar { width: 300px; }
.table-card { border-radius: 8px; }
.tenant-name { font-weight: 600; color: #409EFF; display: block; }
.id-tag { color: #909399; font-size: 11px; }
.amount-positive { color: #67C23A; font-weight: bold; font-size: 16px; }
.date-cell { display: flex; align-items: center; gap: 8px; color: #606266; }
</style>