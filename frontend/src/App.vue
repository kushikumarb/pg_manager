<template>
  <router-view />
</template>

<script setup>
// Global styles or logic can go here
</script>

<style>
/* Global professional font and resets */
body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  background-color: #f0f2f5;
}

/* Custom scrollbar for a "Premium" feel */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 10px;
}
</style>