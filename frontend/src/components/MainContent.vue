<template>
  <t-content class="main-content-container">
    <router-view />
  </t-content>
</template>

<script setup>
</script>

<style scoped>
t-content {
  background: var(--td-bg-color-container);
  height: 100vh;
  overflow: auto;
}
.main-content-container {
  flex: 1;
  overflow-y: auto;
  padding: 3vh 1vw;
  background: var(--td-bg-color-container);
}
</style>