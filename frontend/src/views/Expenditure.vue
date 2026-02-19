<template>
  <div class="expenditure-container">
    <div class="page-header">
      <div class="title-section">
        <el-button link @click="router.push(`/property/${propertyId}/dashboard`)" class="back-link">
          <el-icon><ArrowLeft /></el-icon> Back to Dashboard
        </el-button>
        <h2>Expenditure Tracking: {{ propertyName }}</h2>
      </div>
      <el-button type="warning" @click="dialogVisible = true">
        <el-icon><Plus /></el-icon> Add Expense
      </el-button>
    </div>

    <el-card shadow="never" class="table-card">
      <el-table :data="expenses" v-loading="loading" style="width: 100%" border stripe>
        <el-table-column prop="date" label="Date" width="150">
          <template #default="scope">
            {{ new Date(scope.row.date).toLocaleDateString() }}
          </template>
        </el-table-column>
        <el-table-column prop="category" label="Category" width="150">
          <template #default="scope">
            <el-tag>{{ scope.row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="Description" min-width="250" />
        <el-table-column prop="amount" label="Amount (₹)" width="150">
          <template #default="scope">
            <strong>₹{{ scope.row.amount }}</strong>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="Record New Expense" width="450px">
      <el-form :model="form" label-position="top">
        <el-form-item label="Category" required>
          <el-select v-model="form.category" placeholder="Select Category" style="width: 100%">
            <el-option label="Electricity" value="Electricity" />
            <el-option label="Water" value="Water" />
            <el-option label="Maintenance" value="Maintenance" />
            <el-option label="Cleaning" value="Cleaning" />
            <el-option label="Repairs" value="Repairs" />
            <el-option label="Other" value="Other" />
          </el-select>
        </el-form-item>
        <el-form-item label="Amount (₹)" required>
          <el-input-number v-model="form.amount" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="form.description" type="textarea" placeholder="e.g. Fixed water tap in Room 102" />
        </el-form-item>
        <el-form-item label="Date">
          <el-date-picker v-model="form.date" type="date" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitExpense" :loading="submitting">Save Expense</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const route = useRoute()
const router = useRouter()
const propertyId = route.params.id
const propertyName = localStorage.getItem('selectedPropertyName')

const expenses = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)

const form = reactive({
  property_id: parseInt(propertyId),
  amount: 0,
  category: '',
  description: '',
  date: new Date()
})

const fetchExpenses = async () => {
  loading.value = true
  try {
    const res = await api.get(`/expenditures?property_id=${propertyId}`)
    expenses.value = res.data || []
  } finally { loading.value = false }
}

const submitExpense = async () => {
  if (!form.category || !form.amount) return ElMessage.warning("Fill required fields")
  submitting.value = true
  try {
    await api.post('/expenditures', form)
    ElMessage.success("Expense recorded")
    dialogVisible.value = false
    fetchExpenses()
  } finally { submitting.value = false }
}

onMounted(fetchExpenses)
</script>

<style scoped>
.expenditure-container { padding: 10px; }
.page-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 25px; }
.back-link { font-weight: bold; margin-bottom: 8px; }
.table-card { border-radius: 12px; }
</style>