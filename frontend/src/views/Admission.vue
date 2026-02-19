<template>
  <div class="admission-container">
    <div class="page-header">
      <el-button link @click="router.push(`/property/${propertyId}/dashboard`)">
        <el-icon><ArrowLeft /></el-icon> Back to Dashboard
      </el-button>
      <h2>Tenant Admission & Digital KYC</h2>
    </div>

    <el-form :model="form" label-position="top" ref="admissionFormRef" v-loading="loading">
      <el-card class="form-section" shadow="never">
        <template #header><div class="section-title">Step 1: Room & Bed Assignment</div></template>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="Select Available Bed" required>
              <el-select v-model="form.room_id" placeholder="Choose Room" style="width: 100%" @change="handleRoomChange">
                <el-option 
                  v-for="room in availableRooms" 
                  :key="room.ID" 
                  :label="`Room ${room.room_no} - Bed ${(room.occupied || 0) + 1} (Rent: ₹${room.price})`" 
                  :value="room.ID" 
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      <el-card class="form-section" shadow="never">
        <template #header><div class="section-title">Step 2: Personal Information</div></template>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="Full Name" required><el-input v-model="form.name" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="WhatsApp Number" required><el-input v-model="form.phone_number" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="Email ID"><el-input v-model="form.mail_id" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="Father's Name"><el-input v-model="form.father_name" /></el-form-item></el-col>
          <el-col :span="8">
            <el-form-item label="Date of Birth" required>
              <el-date-picker
                v-model="form.dob"
                type="date"
                placeholder="Select Date"
                format="DD-MM-YYYY"
                value-format="DD-MM-YYYY"
                style="width: 100%"
                @change="calculateAge"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Age (Auto-calculated)">
              <el-input-number v-model="form.age" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="Permanent Address"><el-input type="textarea" v-model="form.permanent_address" /></el-form-item>
      </el-card>

      <el-card class="form-section" shadow="never">
        <template #header><div class="section-title">Step 3: Occupation & Logistics</div></template>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="Occupation"><el-input v-model="form.occupation" /></el-form-item></el-col>
          <el-col :span="16"><el-form-item label="Office/College Address"><el-input v-model="form.office_address" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="6"><el-form-item label="Vegetarian?"><el-switch v-model="form.is_vegetarian" active-text="Yes" inactive-text="No" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="Two Wheeler?"><el-switch v-model="form.has_two_wheeler" /></el-form-item></el-col>
          <el-col :span="12" v-if="form.has_two_wheeler"><el-form-item label="Vehicle No"><el-input v-model="form.vehicle_no" /></el-form-item></el-col>
        </el-row>
      </el-card>

      <el-card class="form-section" shadow="never">
        <template #header><div class="section-title">Step 4: ID Proof & Documents</div></template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="ID Type" required>
              <el-select v-model="form.id_proof_type" style="width: 100%">
                <el-option label="Aadhar Card" value="Aadhar"/>
                <el-option label="PAN Card" value="PAN"/>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="ID Number" required><el-input v-model="form.id_proof_no" /></el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Upload ID Proof (JPG, PNG, or PDF)" required>
          <el-upload
            class="kyc-uploader"
            action="#"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleKycUpload"
            accept=".jpg,.jpeg,.png,.pdf" 
          >
            <template v-if="kycPreview">
              <img v-if="isImage" :src="kycPreview" class="kyc-preview-img" />
              <div v-else class="pdf-placeholder">
                <el-icon :size="45" color="#f56c6c"><Document /></el-icon>
                <p>PDF Document Attached</p>
              </div>
            </template>
            <div v-else class="upload-placeholder">
              <el-icon :size="28"><Plus /></el-icon>
              <span>Click to upload ID</span>
            </div>
          </el-upload>
          <p class="upload-hint">Accepted: JPG, PNG, PDF (Max 2MB)</p>
        </el-form-item>

        <el-divider />

        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="Emergency Contact No."><el-input v-model="form.emergency_contact" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Emergency Relation"><el-input v-model="form.emergency_relation" /></el-form-item></el-col>
        </el-row>
      </el-card>

      <div class="sticky-footer">
        <el-button type="primary" size="large" @click="handleOnboard" round>Verify & Generate OTP</el-button>
      </div>
    </el-form>

    <el-dialog v-model="otpDialog" title="Final Verification" width="400px" center>
      <div style="text-align: center">
        <p>Enter the 6-digit OTP from the Go Terminal</p>
        <el-input v-model="otpInput" placeholder="000000" maxlength="6" style="margin: 20px 0; text-align: center; font-size: 24px" />
        <el-button type="success" @click="handleVerifyOTP" style="width: 100%" size="large">Activate Tenant</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Plus, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const route = useRoute()
const router = useRouter()
const propertyId = route.params.id
const propertyName = localStorage.getItem('selectedPropertyName')

const loading = ref(false)
const otpDialog = ref(false)
const otpInput = ref('')
const currentTenantID = ref(null)
const rooms = ref([])
const availableRooms = ref([])
const kycPreview = ref(null)
const isImage = ref(true)

const form = reactive({
  property_id: parseInt(propertyId),
  room_id: null,
  name: '',
  phone_number: '',
  status: 'pending',
  father_name: '', dob: '', age: 0, permanent_address: '',
  emergency_contact: '', emergency_relation: '', mail_id: '',
  is_vegetarian: false, has_two_wheeler: false, vehicle_no: '',
  education: '', occupation: '', office_address: '',
  id_proof_type: 'Aadhar', id_proof_no: '', id_proof_image: '',
  monthly_rent: 0, deposit: 0, maintenance_charges: 0
})

// Logic to auto-calculate age based on DOB selection
const calculateAge = (val) => {
  if (!val) return;
  const parts = val.split('-');
  const birthDate = new Date(parts[2], parts[1] - 1, parts[0]);
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const m = today.getMonth() - birthDate.getMonth();
  if (m < 0 || (m === 0 && today.getDate() < birthDate.getDate())) {
    age--;
  }
  form.age = age;
}

const handleKycUpload = (file) => {
  const fileType = file.raw.type
  const isAllowed = ['image/jpeg', 'image/png', 'application/pdf'].includes(fileType)
  const isLt2M = file.raw.size / 1024 / 1024 < 2

  if (!isAllowed) {
    ElMessage.error('Only JPG, PNG, or PDF files are allowed!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('File size cannot exceed 2MB!')
    return false
  }

  isImage.value = fileType.startsWith('image/')
  const reader = new FileReader()
  reader.onload = (e) => {
    kycPreview.value = e.target.result
    form.id_proof_image = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

const fetchRooms = async () => {
  try {
    const res = await api.get(`/properties/${propertyId}/rooms`)
    rooms.value = res.data
    availableRooms.value = res.data.filter(r => (r.occupied || 0) < r.capacity)
  } catch (e) {
    ElMessage.error("Failed to load inventory")
  }
}

const handleRoomChange = (val) => {
  const room = rooms.value.find(r => r.ID === val)
  if (room) {
    form.monthly_rent = room.price
    form.deposit = room.deposit || (room.price * 2)
  }
}

const handleOnboard = async () => {
  if (!form.room_id || !form.name || !form.phone_number) {
    return ElMessage.warning("Please fill mandatory fields (Room, Name, Phone)")
  }
  loading.value = true
  try {
    const res = await api.post('/tenants/onboard', form)
    currentTenantID.value = res.data.tenant_id
    otpDialog.value = true
    ElMessage.success("OTP Printed in Terminal!")
  } catch (e) {
    ElMessage.error(e.response?.data?.error || "Error")
  } finally { loading.value = false }
}

const handleVerifyOTP = async () => {
  try {
    await api.post('/tenants/verify', { tenant_id: currentTenantID.value, otp: otpInput.value })
    ElMessage.success("Admission Successful!")
    router.push(`/property/${propertyId}/tenants`)
  } catch (e) { ElMessage.error("Invalid OTP") }
}

onMounted(fetchRooms)
</script>

<style scoped>
.admission-container { padding: 20px; max-width: 900px; margin: 0 auto; background: #f9fafc; }
.form-section { margin-bottom: 25px; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.03); }
.section-title { font-weight: 700; color: #303133; font-size: 16px; }
.kyc-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  width: 280px;
  height: 180px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  background: #fafafa;
  overflow: hidden;
}
.kyc-uploader:hover { border-color: #409eff; }
.kyc-preview-img { width: 100%; height: 100%; object-fit: contain; }
.pdf-placeholder { display: flex; flex-direction: column; align-items: center; gap: 10px; color: #606266; }
.upload-placeholder { color: #8c939d; text-align: center; display: flex; flex-direction: column; gap: 8px; }
.upload-hint { font-size: 12px; color: #909399; margin-top: 8px; }
.sticky-footer { position: sticky; bottom: 20px; text-align: right; background: white; padding: 20px; border-radius: 12px; box-shadow: 0 -4px 15px rgba(0,0,0,0.08); z-index: 10; }
</style>