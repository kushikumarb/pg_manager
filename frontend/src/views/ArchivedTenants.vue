<template>
  <div class="archives-container">
    <div class="page-header">
      <el-button link @click="router.push(`/property/${propertyId}/dashboard`)" class="back-link">
        <el-icon><ArrowLeft /></el-icon> Back to Dashboard
      </el-button>
      <h2>Past Resident Backup (Off-boarded)</h2>
    </div>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="archives" style="width: 100%" border stripe>
        <el-table-column prop="name" label="Full Name" min-width="150" />
        <el-table-column prop="phone_number" label="Phone" width="130" />
        <el-table-column label="Stay Duration" min-width="200">
          <template #default="scope">
            {{ formatDate(scope.row.admission_date) }} to {{ formatDate(scope.row.checkout_date) }}
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="100">
          <template #default="scope">
            <el-button type="primary" link @click="showDetails(scope.row)">View Info</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="Archived Tenant Profile" width="500px">
      <div v-if="selected" class="profile-details">
        <p><strong>Father:</strong> {{ selected.father_name }}</p>
        <p><strong>ID No:</strong> {{ selected.id_proof_no }}</p>
        <p><strong>Address:</strong> {{ selected.permanent_address }}</p>
        <el-divider>KYC Document</el-divider>
        <img :src="selected.id_proof_image" class="kyc-img" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import api from '../api'

const route = useRoute(); const router = useRouter()
const propertyId = route.params.id
const archives = ref([]); const loading = ref(false)
const dialogVisible = ref(false); const selected = ref(null)

const fetchArchives = async () => {
  loading.value = true
  try {
    const res = await api.get(`/tenants/archives?property_id=${propertyId}`)
    archives.value = res.data
  } finally { loading.value = false }
}

const showDetails = (row) => { selected.value = row; dialogVisible.value = true }
const formatDate = (d) => new Date(d).toLocaleDateString('en-IN')

onMounted(fetchArchives)
</script>

<style scoped>
.archives-container { padding: 20px; }
.kyc-img { width: 100%; border-radius: 8px; }
</style>