<template>
  <div class="complaints-container">
    <header class="page-header">
      <div class="header-left">
        <el-button 
          circle 
          @click="router.push(`/property/${propertyId}/dashboard`)" 
          class="back-btn"
        >
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        
        <div class="title-meta">
          <h1>Tenant Complaints</h1>
          <p>Manage and resolve issues reported via QR codes</p>
        </div>
      </div>
      
      <div class="stats-mini">
        <el-tag type="danger" effect="dark">Pending: {{ pendingCount }}</el-tag>
      </div>
    </header>

    <el-card class="table-card">
      <el-table :data="complaints" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="room_no" label="Room" width="100">
          <template #default="scope">
            <span class="room-tag">#{{ scope.row.room_no }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="tenant_name" label="Tenant" width="180" />
        
        <el-table-column prop="category" label="Category" width="150">
          <template #default="scope">
            <el-tag :type="getCategoryType(scope.row.category)">{{ scope.row.category }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="description" label="Issue Description" min-width="250" show-overflow-tooltip />
        
        <el-table-column label="Status" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'Pending' ? 'danger' : 'success'" border>
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="Date" width="150">
          <template #default="scope">
            {{ new Date(scope.row.created_at).toLocaleDateString() }}
          </template>
        </el-table-column>

        <el-table-column label="Action" width="150" fixed="right">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status === 'Pending'"
              type="success" 
              size="small" 
              plain
              @click="resolveIssue(scope.row.id)"
            >
              Resolve
            </el-button>
            <span v-else class="resolved-text">Resolved ✅</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router' // Import useRouter
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue' // Import Icon
import api from '../api'

const route = useRoute()
const router = useRouter() // Initialize router
const complaints = ref([])
const loading = ref(false)
const propertyId = route.params.id

// ... (Rest of your script logic remains same) ...

const pendingCount = computed(() => 
  complaints.value.filter(c => c.status === 'Pending').length
)

const fetchComplaints = async () => {
  loading.value = true
  try {
    const res = await api.get(`/complaints?property_id=${propertyId}`)
    complaints.value = res.data
  } catch (e) {
    ElMessage.error("Failed to load complaints")
  } finally {
    loading.value = false
  }
}

const resolveIssue = (id) => {
  ElMessageBox.confirm('Has this issue been fixed?', 'Confirm Resolution', {
    confirmButtonText: 'Yes, Resolved',
    cancelButtonText: 'Cancel',
    type: 'success',
  }).then(async () => {
    try {
      await api.put(`/complaints/${id}/resolve`)
      ElMessage.success("Status updated")
      fetchComplaints()
    } catch (e) {
      ElMessage.error("Failed to update status")
    }
  })
}

const getCategoryType = (cat) => {
  const map = {
    'Plumbing': 'warning',
    'Electrical': 'danger',
    'Internet': 'info',
    'Cleaning': 'success'
  }
  return map[cat] || ''
}

onMounted(fetchComplaints)
</script>

<style scoped>
.complaints-container { max-width: 1200px; margin: 0 auto; padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }

/* New Header Layout Styles */
.header-left { display: flex; align-items: center; gap: 20px; }
.back-btn { font-size: 18px; border: none; background: #fff; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
.back-btn:hover { background: #f5f7fa; color: #409EFF; }

.room-tag { font-weight: bold; color: #409EFF; }
.resolved-text { color: #67C23A; font-size: 13px; font-weight: 500; }
.table-card { border-radius: 12px; }
</style>