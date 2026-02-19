<template>
  <div class="dashboard-content">
    <header class="page-header">
      <div class="title-meta">
        <h1>{{ propertyName }}</h1>
        <p>Property Overview & Financial Performance</p>
      </div>
      <el-button type="primary" size="large" @click="router.push(`/property/${propertyId}/admission`)">
        <el-icon><Plus /></el-icon> New Admission
      </el-button>
    </header>

    <div class="stats-grid" v-loading="loading">
      <el-card class="stat-card room-card" shadow="hover" @click="router.push(`/property/${propertyId}/inventory`)">
        <div class="card-inner">
          <div class="icon-wrap"><el-icon><HomeFilled /></el-icon></div>
          <div class="data">
            <span class="label">Total Rooms</span>
            <h2 class="value">{{ stats.total_rooms || 0 }}</h2>
          </div>
        </div>
        <div class="card-footer">Manage Inventory <el-icon><ArrowRight /></el-icon></div>
      </el-card>

      <el-card class="stat-card tenant-card" shadow="hover" @click="router.push(`/property/${propertyId}/tenants`)">
        <div class="card-inner">
          <div class="icon-wrap"><el-icon><User /></el-icon></div>
          <div class="data">
            <span class="label">Total Tenants</span>
            <h2 class="value">{{ stats.active_tenants || 0 }}</h2>
          </div>
        </div>
        <div class="card-footer">View Residents <el-icon><ArrowRight /></el-icon></div>
      </el-card>

      <el-card class="stat-card payment-card" shadow="hover" @click="router.push(`/property/${propertyId}/payments`)">
        <div class="card-inner">
          <div class="icon-wrap"><el-icon><List /></el-icon></div>
          <div class="data">
            <span class="label">Total Collections</span>
            <h2 class="value text-success">₹{{ (stats.total_revenue || 0).toLocaleString() }}</h2>
          </div>
        </div>
        <div class="card-footer">Audit Payments <el-icon><ArrowRight /></el-icon></div>
      </el-card>

      <el-card class="stat-card expense-card" shadow="hover" @click="router.push(`/property/${propertyId}/expenditure`)">
        <div class="card-inner">
          <div class="icon-wrap"><el-icon><Money /></el-icon></div>
          <div class="data">
            <span class="label">Monthly Expenses</span>
            <h2 class="value text-danger">₹{{ (stats.total_expenditure || 0).toLocaleString() }}</h2>
          </div>
        </div>
        <div class="card-footer">Track Spending <el-icon><ArrowRight /></el-icon></div>
      </el-card>

      <el-card class="stat-card complaint-card" shadow="hover" @click="router.push(`/property/${propertyId}/complaints`)">
        <div class="card-inner">
          <el-badge :value="stats.pending_issues" :hidden="!stats.pending_issues" type="danger">
            <div class="icon-wrap"><el-icon><Warning /></el-icon></div>
          </el-badge>
          <div class="data">
            <span class="label">Pending Complaints</span>
            <h2 class="value" :class="{'text-danger': stats.pending_issues > 0}">
              {{ stats.pending_issues || 0 }}
            </h2>
          </div>
        </div>
        <div class="card-footer">Resolve QR Issues <el-icon><ArrowRight /></el-icon></div>
      </el-card>
    </div>

    <div class="summary-row" v-if="!loading">
       <el-alert
        title="Estimated Net Profit"
        :type="(stats.total_revenue - stats.total_expenditure) >= 0 ? 'success' : 'error'"
        :description="`Current Balance: ₹${((stats.total_revenue || 0) - (stats.total_expenditure || 0)).toLocaleString()}`"
        show-icon
        :closable="false"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { HomeFilled, User, Warning, ArrowRight, Plus, Money, List } from '@element-plus/icons-vue'
import api from '../api'

const route = useRoute()
const router = useRouter()
const stats = ref({})
const loading = ref(false)
const propertyName = ref(localStorage.getItem('selectedPropertyName') || 'Property Dashboard')
const propertyId = ref(route.params.id)

const fetchStats = async () => {
  loading.value = true
  try {
    const res = await api.get(`/dashboard?property_id=${propertyId.value}`)
    stats.value = res.data
  } catch (e) {
    console.error("Failed to load dashboard stats", e)
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, (newId) => {
  propertyId.value = newId
  propertyName.value = localStorage.getItem('selectedPropertyName') || 'Property Dashboard'
  fetchStats()
})

onMounted(fetchStats)
</script>

<style scoped>
.dashboard-content { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(260px, 1fr)); gap: 20px; }
.stat-card { border-radius: 12px; cursor: pointer; border: none; transition: transform 0.3s; height: 100%; position: relative; }
.stat-card:hover { transform: translateY(-5px); box-shadow: 0 8px 24px rgba(0,0,0,0.1); }
.card-inner { display: flex; align-items: center; gap: 20px; padding: 25px; }

.icon-wrap { width: 60px; height: 60px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 30px; }
.room-card .icon-wrap { background: #e6f7ff; color: #1890ff; }
.tenant-card .icon-wrap { background: #f9f0ff; color: #722ed1; }
.payment-card .icon-wrap { background: #f6ffed; color: #52c41a; }
.expense-card .icon-wrap { background: #fff7e6; color: #ffa940; }
.complaint-card .icon-wrap { background: #fff1f0; color: #f5222d; }

.label { font-size: 14px; color: #8c8c8c; text-transform: uppercase; letter-spacing: 0.5px; }
.value { font-size: 28px; margin-top: 5px; color: #262626; font-weight: bold; }
.text-success { color: #52c41a; }
.text-danger { color: #f5222d; }

.card-footer { padding: 12px 25px; background: #fafafa; border-top: 1px solid #f0f0f0; display: flex; justify-content: space-between; align-items: center; font-size: 13px; color: #1890ff; font-weight: 500; }
.summary-row { margin-top: 30px; border-radius: 8px; overflow: hidden; }

/* Custom Badge Alignment */
:deep(.el-badge__content) { top: 10px; right: 10px; }
</style>