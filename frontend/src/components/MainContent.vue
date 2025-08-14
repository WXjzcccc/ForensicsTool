<template>
  <t-content class="main-content-container">
    <router-view />
  </t-content>
</template>

<script setup>
</script>

<style scoped>
t-content {
  padding: 20px;
  background: var(--td-bg-color-container);
  height: calc(100vh - 80px);
  overflow: auto;
}
.main-content-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: var(--td-bg-color-container);
}
</style>