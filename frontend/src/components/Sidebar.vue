<template>
  <el-menu
    default-active="1"
    class="sidebar-menu"
    :collapse="isCollapse"
    background-color="#001529"
    text-color="#fff"
    active-text-color="#409EFF"
    router
  >
    <div class="logo">
      <el-icon v-if="!isCollapse"><HomeFilled /></el-icon>
      <span v-if="!isCollapse">PG MANAGER</span>
    </div>

    <el-menu-item index="/dashboard">
      <el-icon><DataBoard /></el-icon>
      <span>Dashboard</span>
    </el-menu-item>

    <el-menu-item index="/properties">
      <el-icon><OfficeBuilding /></el-icon>
      <span>Properties</span>
    </el-menu-item>

    <el-menu-item index="/tenants">
      <el-icon><User /></el-icon>
      <span>Tenants</span>
    </el-menu-item>

    <el-menu-item index="/complaints">
      <el-icon><Warning /></el-icon>
      <span>Complaints</span>
    </el-menu-item>

    <el-menu-item index="/expenditures">
      <el-icon><Money /></el-icon>
      <span>Finance</span>
    </el-menu-item>
  </el-menu>
</template>

<script setup>
import { HomeFilled, DataBoard, OfficeBuilding, User, Warning, Money } from '@element-plus/icons-vue'
defineProps(['isCollapse'])
</script>

<style scoped>
.sidebar-menu {
  height: 100vh;
  border-right: none;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 18px;
  background: #002140;
}
</style>