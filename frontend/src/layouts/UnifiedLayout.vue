<template>
  <div class="layout-container">
    <aside :class="['sidebar', { 'is-collapsed': isCollapsed }]">
      <div class="sidebar-header">
        <div class="logo-display-box" v-if="!isCollapsed">
          <img v-if="logoUrl" :src="logoUrl" class="business-logo-sidebar" />
          <div v-else class="logo-placeholder-sidebar">
            <el-icon :size="24"><Picture /></el-icon>
          </div>
        </div>
        
        <div class="menu-toggle" @click="isCollapsed = !isCollapsed">
          <el-icon :size="24">
            <Fold v-if="!isCollapsed" />
            <Expand v-else />
          </el-icon>
        </div>
      </div>

      <el-menu 
        class="menu" 
        :default-active="route.path" 
        :collapse="isCollapsed" 
        background-color="#001529" 
        text-color="#fff"
      >
        <el-menu-item index="/property-selector" @click="router.push('/property-selector')">
          <el-icon><MenuIcon /></el-icon>
          <template #title>My Properties</template>
        </el-menu-item>

        <el-menu-item index="/profile-settings" @click="router.push('/profile-settings')">
          <el-icon><User /></el-icon>
          <template #title>Profile Settings</template>
        </el-menu-item>

        <template v-if="propertyId">
          <div class="menu-divider" v-if="!isCollapsed">Property Management</div>
          <el-menu-item 
            :index="`/property/${propertyId}/dashboard`" 
            @click="router.push(`/property/${propertyId}/dashboard`)"
          >
            <el-icon><OfficeBuilding /></el-icon>
            <template #title>Building Dashboard</template>
          </el-menu-item>
          <el-menu-item 
            :index="`/property/${propertyId}/archives`" 
            @click="router.push(`/property/${propertyId}/archives`)"
          >
            <el-icon><Collection /></el-icon>
            <template #title>Backup Records</template>
          </el-menu-item>
        </template>
        
        <div class="sidebar-spacer"></div>
        
        <el-menu-item @click="handleLogout" class="logout-item">
          <el-icon><SwitchButton /></el-icon>
          <template #title>Logout</template>
        </el-menu-item>
      </el-menu>
    </aside>
    
    <main class="content">
      <router-view :key="route.fullPath" />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  Menu as MenuIcon, SwitchButton, Fold, Expand, Picture, 
  User, Collection, OfficeBuilding 
} from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const isCollapsed = ref(false)
const propertyId = ref(route.params.id)
const logoUrl = ref(null)

const updateLogo = () => {
  logoUrl.value = localStorage.getItem('business_logo')
}

// Listen for custom events to update logo in real-time
onMounted(() => {
  updateLogo()
  window.addEventListener('logo-updated', updateLogo)
})

watch(() => route.params.id, (newId) => {
  propertyId.value = newId
})

const handleLogout = () => {
  ElMessageBox.confirm('Confirm Logout?', 'Warning', { type: 'warning' }).then(() => {
    localStorage.clear()
    router.push('/login')
  })
}
</script>

<style scoped>
.layout-container { display: flex; height: 100vh; width: 100vw; overflow: hidden; }
.sidebar { width: 260px; background-color: #001529; transition: width 0.3s ease; display: flex; flex-direction: column; flex-shrink: 0; }
.sidebar.is-collapsed { width: 80px; }
.sidebar-header { padding: 20px; display: flex; flex-direction: column; align-items: center; gap: 15px; }

.logo-display-box { width: 100%; display: flex; justify-content: center; margin-bottom: 10px; }
.business-logo-sidebar { width: 80px; height: 80px; object-fit: contain; }
.logo-placeholder-sidebar { width: 60px; height: 60px; border-radius: 50%; background: #1890ff20; display: flex; align-items: center; justify-content: center; color: #1890ff; }

.menu-toggle { cursor: pointer; color: white; align-self: flex-end; }
.menu { border: none; background: transparent; flex: 1; display: flex; flex-direction: column; }
.menu-divider { color: #595959; font-size: 12px; padding: 15px 20px 5px; text-transform: uppercase; letter-spacing: 1px; }
.sidebar-spacer { flex: 1; }
.content { flex: 1; background-color: #f5f7fa; overflow-y: auto; padding: 30px; }
.logout-item { color: #f56c6c !important; }
</style>