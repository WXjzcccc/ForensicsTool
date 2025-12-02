<script setup>
import AppTopbar from "./components/AppTopbar.vue";
import AppFooter from "./components/AppFooter.vue";
import AppSidebar from "./components/AppSidebar.vue";
import { onMounted, onUnmounted } from 'vue';

// 全局拖放事件处理，阻止浏览器默认行为
const handleGlobalDragOver = (e) => {
  // 阻止默认行为，但不阻止事件冒泡
  e.preventDefault();
};

const handleGlobalDrop = (e) => {
  // 阻止默认行为，但不阻止事件冒泡
  e.preventDefault();
};

// 组件挂载时添加全局事件监听
onMounted(() => {
  document.addEventListener('dragover', handleGlobalDragOver);
  document.addEventListener('drop', handleGlobalDrop);
});

// 组件卸载时移除全局事件监听
onUnmounted(() => {
  document.removeEventListener('dragover', handleGlobalDragOver);
  document.removeEventListener('drop', handleGlobalDrop);
});
</script>

<template>
    <div class="app-layout">
        <AppTopbar />
        <div class="app-body">
            <AppSidebar />
            <div class="app-content">
                <transition name="page-transition" mode="out-in">
                    <router-view />
                </transition>
            </div>
        </div>
        <!-- <AppFooter /> -->
    </div>
</template>

<style>
/* 页面过渡效果样式 */
.page-transition-enter-active,
.page-transition-leave-active {
    transition: all 400ms ease;
}

.page-transition-enter-from {
    transform: scale(0.9);
}

.page-transition-leave-to {
    transform: scale(1.1);
}

.page-transition-enter-active,
.page-transition-leave-active {
    transition: transform 400ms ease, opacity 400ms ease;
}

.page-transition-enter-from,
.page-transition-leave-to {
    opacity: 0;
}

</style>
