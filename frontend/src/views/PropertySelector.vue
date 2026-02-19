<template>
  <div class="selector-content">
    <header class="welcome-header">
      <h1>Welcome back!</h1>
      <p>Select a property to manage your tenants and inventory</p>
    </header>

    <div class="property-grid" v-loading="loading">
      <el-card 
        v-for="prop in properties" 
        :key="prop.id || prop.ID" 
        class="property-card" 
        shadow="hover" 
        @click="selectProperty(prop)"
      >
        <el-icon :size="40" color="#409EFF"><OfficeBuilding /></el-icon>
        <h3>{{ prop.name }}</h3>
        <p class="address">{{ prop.address }}</p>
        <div class="manage-link">Manage Building <el-icon><ArrowRight /></el-icon></div>
      </el-card>

      <el-card class="add-card" @click="addDialogVisible = true">
        <el-icon :size="40"><Plus /></el-icon>
        <p>Add New Property</p>
      </el-card>
    </div>

    <el-dialog v-model="addDialogVisible" title="Register New Property" width="450px">
      <el-form :model="propertyForm" label-position="top">
        <el-form-item label="Building Name">
          <el-input v-model="propertyForm.name" placeholder="e.g. BTM Luxury PG" />
        </el-form-item>
        <el-form-item label="Full Address">
          <el-input v-model="propertyForm.address" type="textarea" placeholder="Street, Area, City" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="submitNewProperty" :loading="submitting">Save Property</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { OfficeBuilding, ArrowRight, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const router = useRouter()
const properties = ref([])
const loading = ref(false)
const addDialogVisible = ref(false)
const submitting = ref(false)
const propertyForm = reactive({ name: '', address: '' })

const fetchProperties = async () => {
  loading.value = true
  try {
    const res = await api.get('/properties')
    properties.value = res.data
  } catch (e) {
    ElMessage.error("Failed to load properties")
  } finally { loading.value = false }
}

const submitNewProperty = async () => {
  if (!propertyForm.name || !propertyForm.address) return ElMessage.warning("Fill all fields")
  submitting.value = true
  try {
    await api.post('/properties', propertyForm)
    ElMessage.success('New property added!')
    addDialogVisible.value = false
    propertyForm.name = ''; propertyForm.address = ''
    fetchProperties()
  } catch (e) {
    ElMessage.error("Error saving property")
  } finally { submitting.value = false }
}

const selectProperty = (prop) => {
  localStorage.setItem('selectedPropertyId', prop.id || prop.ID)
  localStorage.setItem('selectedPropertyName', prop.name)
  router.push(`/property/${prop.id || prop.ID}/dashboard`)
}

onMounted(fetchProperties)
</script>

<style scoped>
.selector-content { max-width: 1200px; margin: 0 auto; }
.welcome-header { text-align: center; margin-bottom: 50px; }
.property-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 30px; }
.property-card { text-align: center; padding: 30px; cursor: pointer; border-radius: 16px; transition: all 0.3s ease; }
.property-card:hover { transform: translateY(-10px); box-shadow: 0 12px 32px rgba(0,0,0,0.1); }
.address { color: #909399; font-size: 14px; margin: 10px 0 20px; }
.manage-link { color: #409EFF; font-weight: bold; display: flex; align-items: center; justify-content: center; gap: 5px; }
.add-card { border: 2px dashed #dcdfe6; height: 250px; display: flex; flex-direction: column; align-items: center; justify-content: center; cursor: pointer; color: #909399; border-radius: 16px; }
.add-card:hover { border-color: #409EFF; color: #409EFF; }
</style>