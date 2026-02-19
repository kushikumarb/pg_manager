<template>
  <div class="tenants-container">
    <div class="page-header">
      <div class="title-section">
        <el-button link @click="router.push(`/property/${propertyId}/dashboard`)" class="back-link">
          <el-icon><ArrowLeft /></el-icon> Back to Dashboard
        </el-button>
        <h2>Resident List: {{ propertyName }}</h2>
      </div>
      <div class="action-buttons">
        <el-button type="primary" @click="router.push(`/property/${propertyId}/admission`)">
          <el-icon><Plus /></el-icon> New Admission
        </el-button>
      </div>
    </div>

    <el-card shadow="never" class="table-card">
      <div v-if="loading" v-loading="loading" style="height: 300px"></div>
      <el-table v-else :data="tenants" style="width: 100%" border stripe>
        <el-table-column prop="user_id" label="ID" width="80" />
        <el-table-column prop="name" label="Full Name" min-width="180" />
        <el-table-column prop="room_no" label="Room #" width="100" />
        <el-table-column prop="phone_number" label="Contact" width="160" />
        <el-table-column prop="status" label="Status" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'warning'">
              {{ scope.row.status ? scope.row.status.toUpperCase() : 'PENDING' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Actions" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openTenantDrawer(scope.row.user_id)">Manage</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-drawer v-model="drawerVisible" :title="`Profile: ${selectedTenant.name || ''}`" size="500px">
      <div v-if="drawerLoading" v-loading="drawerLoading" style="height: 200px"></div>
      <div v-else class="drawer-inner">
        <div class="finance-banner">
          <div class="stat">
            <label>Balance Due</label>
            <h3 :class="{ 'text-danger': selectedTenant.balance > 0 }">₹{{ selectedTenant.balance || 0 }}</h3>
          </div>
          <div class="stat">
            <label>Monthly Rent</label>
            <h3>₹{{ selectedTenant.monthly_rent }}</h3>
          </div>
        </div>

        <el-divider content-position="left">Personal Info</el-divider>
        <div class="profile-info">
          <p><strong>Father:</strong> {{ selectedTenant.father_name || 'N/A' }}</p>
          <p><strong>Contact:</strong> {{ selectedTenant.phone_number }}</p>
          <p><strong>Address:</strong> {{ selectedTenant.permanent_address }}</p>
        </div>

        <el-divider content-position="left">KYC Documents</el-divider>
        <div class="kyc-container" @click="previewImage(selectedTenant.id_proof_image)">
          <template v-if="selectedTenant.id_proof_image">
            <img v-if="selectedTenant.id_proof_image.startsWith('data:image')" :src="selectedTenant.id_proof_image" class="kyc-display" />
            <div v-else class="pdf-link-box">
              <el-icon :size="40" color="#f56c6c"><Document /></el-icon>
              <span>View ID (PDF)</span>
            </div>
          </template>
          <el-empty v-else description="No Document" :image-size="40" />
        </div>

        <div class="drawer-actions">
          <el-button 
            type="success" 
            size="large" 
            style="width: 100%; margin-bottom: 12px;" 
            :disabled="selectedTenant.balance <= 0"
            @click="recordPayment"
          >
            {{ selectedTenant.balance <= 0 ? 'No Due Balance' : 'Receive Rent Payment' }}
          </el-button>

          <el-button 
            type="danger" 
            plain
            size="large" 
            style="width: 100%" 
            @click="handleOffboard"
            :loading="offboardLoading"
            :disabled="selectedTenant.status === 'inactive'"
          >
            {{ selectedTenant.status === 'inactive' ? 'Already Checked Out' : 'Offboard / Checkout Tenant' }}
          </el-button>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="paymentDialog" title="Record Payment" width="400px" center>
      <el-form :model="paymentForm" label-position="top">
        <el-form-item :label="`Amount Received (Current Balance: ₹${selectedTenant.balance})`">
          <el-input-number 
            v-model="paymentForm.amount" 
            :min="0"
            style="width: 100%"
            controls-position="right"
          />
        </el-form-item>
        <el-form-item label="Payment Method">
          <el-select v-model="paymentForm.method" style="width: 100%">
            <el-option label="Cash" value="Cash" />
            <el-option label="UPI (PhonePe/GPay)" value="UPI" />
            <el-option label="Bank Transfer" value="Bank" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="paymentDialog = false">Cancel</el-button>
        <el-button type="success" @click="submitPayment" :loading="paymentLoading">Confirm</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Plus, Document } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const route = useRoute(); const router = useRouter()
const propertyId = route.params.id; const propertyName = localStorage.getItem('selectedPropertyName') || 'Property'
const tenants = ref([]); const loading = ref(false)
const drawerVisible = ref(false); const drawerLoading = ref(false); const selectedTenant = ref({})
const paymentDialog = ref(false); const paymentLoading = ref(false)
const offboardLoading = ref(false)
const paymentForm = reactive({ amount: 0, method: 'UPI' })

const fetchTenants = async () => {
  loading.value = true
  try { const res = await api.get(`/tenants?property_id=${propertyId}`); tenants.value = res.data || [] }
  finally { loading.value = false }
}

const openTenantDrawer = async (userId) => {
  drawerVisible.value = true; drawerLoading.value = true
  try { const res = await api.get(`/tenants/${userId}`); selectedTenant.value = res.data }
  finally { drawerLoading.value = false }
}

const previewImage = (data) => { if (data) window.open().document.write(data.includes('pdf') ? `<iframe width='100%' height='100%' src='${data}'></iframe>` : `<img src="${data}" />`) }

const recordPayment = () => {
  paymentForm.amount = selectedTenant.value.balance 
  paymentDialog.value = true
}

const submitPayment = async () => {
  const enteredAmount = Number(paymentForm.amount);
  const currentBalance = Number(selectedTenant.value.balance);

  if (enteredAmount > currentBalance) {
    ElMessage.error(`Cannot accept ₹${enteredAmount}. Max due is ₹${currentBalance}.`);
    paymentForm.amount = currentBalance;
    return;
  }

  if (enteredAmount <= 0) return ElMessage.warning("Please enter a valid amount.");

  paymentLoading.value = true
  try {
    await api.post(`/tenants/${selectedTenant.value.user_id}/pay`, {
      amount: enteredAmount,
      method: paymentForm.method
    })
    ElMessage.success(enteredAmount < currentBalance ? "Partial payment recorded!" : "Full payment recorded!");
    paymentDialog.value = false
    await openTenantDrawer(selectedTenant.value.user_id); await fetchTenants()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || "Error recording payment")
  } finally { paymentLoading.value = false }
}

const handleOffboard = async () => {
  if (selectedTenant.value.balance > 0) {
    try {
      await ElMessageBox.confirm(
        `This tenant still has a pending balance of ₹${selectedTenant.value.balance}. Are you sure you want to checkout without clearing the dues?`,
        'Pending Dues Detected',
        { confirmButtonText: 'Checkout Anyway', cancelButtonText: 'Cancel', type: 'error' }
      )
    } catch { return }
  } else {
    try {
      await ElMessageBox.confirm(
        `Are you sure you want to checkout ${selectedTenant.value.name}? This will free up their bed in Room ${selectedTenant.value.room_no}.`,
        'Confirm Checkout',
        { confirmButtonText: 'Confirm', cancelButtonText: 'Cancel', type: 'warning' }
      )
    } catch { return }
  }

  offboardLoading.value = true
  try {
    await api.post(`/tenants/${selectedTenant.value.user_id}/offboard`)
    ElMessage.success("Tenant checked out successfully!")
    drawerVisible.value = false
    fetchTenants()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || "Error during offboarding")
  } finally { offboardLoading.value = false }
}

onMounted(fetchTenants)
</script>

<style scoped>
.tenants-container { padding: 10px; }
.page-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 25px; }
.action-buttons { display: flex; gap: 12px; }
.back-link { font-weight: bold; color: #409EFF; }
.table-card { border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
.finance-banner { display: flex; gap: 20px; background: #f8f9fb; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
.text-danger { color: #f56c6c; font-weight: bold; }
.kyc-container { border: 1px dashed #dcdfe6; border-radius: 8px; padding: 10px; text-align: center; cursor: pointer; margin-bottom: 20px;}
.kyc-display { max-width: 100%; border-radius: 4px; }
.drawer-actions { margin-top: 30px; display: flex; flex-direction: column; }
.pdf-link-box { display: flex; flex-direction: column; align-items: center; gap: 8px; color: #606266; }
</style>