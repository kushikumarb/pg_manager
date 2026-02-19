<template>
  <div class="rooms-container">
    <div class="page-header">
      <div class="title-section">
        <el-button link @click="router.push(`/property/${propertyId}/dashboard`)" class="back-link">
          <el-icon><ArrowLeft /></el-icon> Back to Dashboard
        </el-button>
        <h2>Room Management: {{ propertyName }}</h2>
      </div>
      <el-button type="primary" @click="roomFormVisible = true">
        <el-icon><Plus /></el-icon> Add New Room
      </el-button>
    </div>

    <el-card shadow="never" class="table-card">
      <el-table :data="rooms" style="width: 100%" v-loading="loading" border stripe>
        <el-table-column prop="room_no" label="Room #" width="100" />
        <el-table-column prop="capacity" label="Total Beds" width="110" />
        <el-table-column prop="price" label="Rent (₹)" width="120" />
        
        <el-table-column prop="deposit" label="Security Deposit (₹)" width="160">
           <template #default="scope">
             <span style="color: #E6A23C; font-weight: bold;">₹{{ scope.row.deposit }}</span>
           </template>
        </el-table-column>
        
        <el-table-column label="Occupancy Status" min-width="200">
          <template #default="scope">
            <div class="bed-visual-container">
              <div class="bed-dots">
                <div 
                  v-for="i in scope.row.capacity" 
                  :key="i"
                  class="bed-dot"
                  :class="i <= (scope.row.occupied || 0) ? 'occupied' : 'available'"
                >
                  <el-tooltip :content="i <= (scope.row.occupied || 0) ? 'Occupied' : 'Available'" placement="top">
                    <div class="dot-inner"></div>
                  </el-tooltip>
                </div>
              </div>
              <span class="occupancy-label">({{ scope.row.occupied || 0 }}/{{ scope.row.capacity }})</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="Actions" width="120" fixed="right">
          <template #default="scope">
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(scope.row)"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="roomFormVisible" title="Add New Room" width="420px">
      <el-form :model="roomForm" label-position="top">
        <el-form-item label="Room Number" required>
          <el-input v-model="roomForm.room_no" placeholder="e.g. 101" />
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Total Beds" required>
              <el-input-number v-model="roomForm.capacity" :min="1" :max="10" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Rent per Bed (₹)" required>
              <el-input v-model="roomForm.price" type="number" placeholder="5000" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Fixed Security Deposit (₹)" required>
          <el-input v-model="roomForm.deposit" type="number" placeholder="10000" prefix-icon="Money" />
          <p class="helper-text">This amount will be auto-applied during new admissions.</p>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="roomFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitRoom" :loading="submitting">Save Room</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Plus, ArrowLeft, Money } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const route = useRoute()
const router = useRouter()
const propertyId = route.params.id
const propertyName = localStorage.getItem('selectedPropertyName')

const rooms = ref([])
const loading = ref(false)
const roomFormVisible = ref(false)
const submitting = ref(false)

const roomForm = reactive({ 
  room_no: '', 
  capacity: 1, 
  price: '', 
  deposit: '' 
})

const fetchRooms = async () => {
  loading.value = true
  try {
    const res = await api.get(`/properties/${propertyId}/rooms`)
    rooms.value = res.data
  } catch (e) {
    ElMessage.error('Failed to load rooms')
  } finally {
    loading.value = false
  }
}

const submitRoom = async () => {
  if (!roomForm.room_no || !roomForm.price || !roomForm.deposit) {
    return ElMessage.warning("Please fill all required fields")
  }
  
  submitting.value = true
  try {
    const payload = { 
      ...roomForm, 
      property_id: parseInt(propertyId), 
      price: parseFloat(roomForm.price),
      deposit: parseFloat(roomForm.deposit)
    }
    await api.post('/rooms', payload)
    ElMessage.success('Room Added Successfully')
    
    Object.assign(roomForm, { room_no: '', capacity: 1, price: '', deposit: '' })
    roomFormVisible.value = false
    fetchRooms()
  } catch (e) {
    ElMessage.error('Error adding room')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row) => {
  // 1. Frontend Check: Check if room is empty
  if (row.occupied > 0) {
    return ElMessageBox.alert(
      `Room ${row.room_no} cannot be deleted because it currently has ${row.occupied} active tenant(s). Please offboard them first.`,
      'Cannot Delete Room',
      { type: 'error', confirmButtonText: 'OK' }
    )
  }

  // 2. Confirmation before calling API
  try {
    await ElMessageBox.confirm(
      `Are you sure you want to delete Room ${row.room_no}? This action cannot be undone.`,
      'Confirm Deletion',
      {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    )

    loading.value = true
    // Note: Assuming your backend route is DELETE /api/rooms/:id
    await api.delete(`/rooms/${row.ID}`)
    ElMessage.success('Room deleted successfully')
    fetchRooms()
  } catch (e) {
    if (e !== 'cancel') {
      const errorMsg = e.response?.data?.error || 'Failed to delete room'
      ElMessage.error(errorMsg)
    }
  } finally {
    loading.value = false
  }
}

onMounted(fetchRooms)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 25px; }
.back-link { margin-bottom: 8px; font-weight: bold; }
.bed-visual-container { display: flex; align-items: center; gap: 12px; }
.bed-dots { display: flex; gap: 6px; }
.bed-dot { width: 14px; height: 14px; border-radius: 50%; border: 1px solid rgba(0,0,0,0.1); }
.available { background-color: #67c23a; }
.occupied { background-color: #f56c6c; }
.occupancy-label { font-size: 13px; font-weight: 600; color: #606266; }
.helper-text { font-size: 11px; color: #909399; margin-top: 5px; }
.table-card { border-radius: 12px; }
</style>