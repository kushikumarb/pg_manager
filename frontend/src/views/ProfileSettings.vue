<template>
  <div class="profile-settings-container">
    <header class="settings-header">
      <h1>Account & Preferences</h1>
      <p>Manage your business branding and system preferences</p>
    </header>

    <el-row :gutter="25">
      <el-col :span="14">
        <el-card shadow="never" class="settings-card">
          <template #header><div class="card-header">Business Branding</div></template>
          <div class="logo-management">
            <el-upload
              class="logo-uploader"
              action="#"
              :auto-upload="false"
              :show-file-list="false"
              :on-change="handleLogoUpdate"
            >
              <img v-if="logoUrl" :src="logoUrl" class="preview-img" />
              <div v-else class="upload-placeholder">
                <el-icon :size="40"><Plus /></el-icon>
                <span>Upload Logo</span>
              </div>
            </el-upload>
            <div class="logo-info">
              <h3>Corporate Logo</h3>
              <p>This logo will appear on your sidebar and generated PDF reports.</p>
              <el-button type="danger" link v-if="logoUrl" @click="removeLogo">Remove Logo</el-button>
            </div>
          </div>
        </el-card>

        <el-card shadow="never" class="settings-card" v-loading="loading">
          <template #header>
            <div class="card-header">
              <span>Owner Profile</span>
              <el-button type="primary" @click="saveProfile" :loading="saving">Save Changes</el-button>
            </div>
          </template>
          <el-form :model="profileForm" label-position="top">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Full Name">
                  <el-input v-model="profileForm.name" placeholder="Enter your full name" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Contact Number">
                  <el-input v-model="profileForm.phone" disabled />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="Email Address">
              <el-input v-model="profileForm.email" placeholder="Enter your business email" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card shadow="never" class="settings-card">
          <template #header><span>System Preferences</span></template>
          
          <div class="pref-item">
            <div class="label-wrap">
              <el-icon><MagicStick /></el-icon>
              <span>Dark Mode Theme</span>
            </div>
            <el-switch 
              v-model="settings.darkMode" 
              @change="toggleTheme" 
              inline-prompt
              active-text="Dark"
              inactive-text="Light"
            />
          </div>

          <el-divider />

          <div class="pref-item">
            <div class="label-wrap">
              <el-icon><ChatDotRound /></el-icon>
              <span>System Language</span>
            </div>
            <el-select v-model="settings.language" style="width: 130px" @change="updateLanguage">
              <el-option label="English" value="en" />
              <el-option label="Kannada" value="kn" />
              <el-option label="Hindi" value="hi" />
            </el-select>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { MagicStick, ChatDotRound, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '../api'

// Import Element Plus dark mode styles
import 'element-plus/theme-chalk/dark/css-vars.css'

const saving = ref(false)
const loading = ref(false)
const logoUrl = ref(localStorage.getItem('business_logo'))

const profileForm = reactive({
  name: '',
  phone: '',
  email: ''
})

const settings = reactive({
  darkMode: localStorage.getItem('theme') === 'dark',
  language: localStorage.getItem('lang') || 'en'
})

// 🛠️ FETCH PROFILE DATA
const fetchProfile = async () => {
  loading.value = true
  try {
    const res = await api.get('/owner/profile')
    // Only update if the request is successful
    profileForm.name = res.data.name
    profileForm.phone = res.data.phone
    profileForm.email = res.data.email
  } catch (e) {
    if (e.response && e.response.status === 401) {
      ElMessage.error('Session expired. Please log in again.')
      // Optional: router.push('/login')
    } else {
      ElMessage.error('Failed to load profile')
    }
  } finally {
    loading.value = false
  }
}

// 🖼️ LOGO MANAGEMENT
const handleLogoUpdate = (file) => {
  const isJPGorPNG = file.raw.type === 'image/jpeg' || file.raw.type === 'image/png';
  if (!isJPGorPNG) return ElMessage.error('Logo must be JPG or PNG!');

  const reader = new FileReader();
  reader.onload = (e) => {
    const base64 = e.target.result;
    logoUrl.value = base64;
    localStorage.setItem('business_logo', base64);
    window.dispatchEvent(new Event('logo-updated'));
    ElMessage.success('Branding logo updated!');
  };
  reader.readAsDataURL(file.raw);
}

const removeLogo = () => {
  logoUrl.value = null;
  localStorage.removeItem('business_logo');
  window.dispatchEvent(new Event('logo-updated'));
}

// 💾 SAVE PROFILE CHANGES
const saveProfile = async () => {
  if (!profileForm.name || !profileForm.email) {
    return ElMessage.warning("Name and Email are required")
  }
  
  saving.value = true
  try {
    await api.put('/owner/profile', {
      name: profileForm.name,
      email: profileForm.email
    })
    localStorage.setItem('owner_name', profileForm.name)
    ElMessage.success('Profile saved successfully')
  } catch (e) {
    ElMessage.error('Failed to update profile')
  } finally {
    saving.value = false
  }
}

// 🌗 THEME TOGGLE
const toggleTheme = (isDark) => {
  const theme = isDark ? 'dark' : 'light'
  localStorage.setItem('theme', theme)
  document.documentElement.className = theme
  
  ElMessage({
    message: `${isDark ? 'Dark' : 'Light'} mode activated`,
    type: 'info',
    duration: 1000
  })
}

const updateLanguage = (val) => {
  localStorage.setItem('lang', val)
  ElMessage.success('Language preferences updated')
}

onMounted(() => {
  fetchProfile()
  // Ensure theme is applied on component load
  if (settings.darkMode) {
    document.documentElement.className = 'dark'
  }
})
</script>

<style scoped>
.profile-settings-container { max-width: 1100px; margin: 0 auto; padding: 20px; }
.settings-header { margin-bottom: 30px; }
.settings-card { border-radius: 12px; margin-bottom: 25px; transition: border-color 0.3s; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: bold; }

.logo-management { display: flex; gap: 25px; align-items: center; }
.logo-uploader {
  width: 120px; height: 120px; border: 2px dashed #dcdfe6;
  border-radius: 12px; cursor: pointer; overflow: hidden;
  display: flex; align-items: center; justify-content: center;
  transition: border-color 0.3s;
}
.logo-uploader:hover { border-color: #409EFF; }
.preview-img { width: 100%; height: 100%; object-fit: contain; }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; color: #909399; font-size: 12px; }

.pref-item { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; }
.label-wrap { display: flex; align-items: center; gap: 12px; font-weight: 500; }
</style>