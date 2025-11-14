<script setup>
import { ref } from 'vue';
import { useLayout } from "../composables/useLayout";
import AppConfig from "./AppConfig.vue";

const { isDarkMode, toggleDarkMode } = useLayout();
const isMaximized = ref(false);
const isWailsEnvironment = ref(false);

// 检查是否在Wails环境中
if (typeof window !== 'undefined' && window.runtime) {
  isWailsEnvironment.value = true;
  
  // 导入Wails运行时函数
  import("../../wailsjs/runtime/runtime").then(runtime => {
    const { WindowMinimise, WindowToggleMaximise, Quit } = runtime;
    
    // 最小化窗口
    window.minimizeWindow = () => {
      WindowMinimise();
    };

    window.quitApplication = () => {
      Quit();
    };
    
    // 最大化/还原窗口 - 使用WindowToggleMaximise函数
    window.toggleMaximizeWindow = () => {
      WindowToggleMaximise();
      // 切换后更新状态
      isMaximized.value = !isMaximized.value;
    };
  });
}

// 非Wails环境下的空函数
const minimizeWindow = () => {
  if (isWailsEnvironment.value && window.minimizeWindow) {
    window.minimizeWindow();
  } else {
    console.log("最小化功能仅在Wails环境中可用");
  }
};

const toggleMaximizeWindow = () => {
  if (isWailsEnvironment.value && window.toggleMaximizeWindow) {
    window.toggleMaximizeWindow();
  } else {
    console.log("最大化功能仅在Wails环境中可用");
  }
};

// 关闭窗口（暂时空着）
const closeWindow = () => {
  window.quitApplication();
};
</script>

<template>
    <div class="topbar" style="--wails-draggable: drag;">
        <div class="topbar-container">
            <div class="topbar-brand">
                <span class="topbar-brand-text" style="--wails-draggable: no-drag;">
                    <span class="topbar-title">ForensicsTool</span>
                    <span class="topbar-subtitle">by WXjzc</span>
                </span>
            </div>
            <div class="topbar-actions">
                <Button type="button" class="topbar-theme-button" @click="toggleDarkMode" text rounded>
                    <i :class="['pi ', 'pi ', { 'pi-moon': isDarkMode, 'pi-sun': !isDarkMode }]" />
                </Button>
                <div class="relative">
                    <Button
                        v-styleclass="{
                            selector: '@next',
                            enterFromClass: 'hidden',
                            enterActiveClass: 'animate-scalein',
                            leaveToClass: 'hidden',
                            leaveActiveClass: 'animate-fadeout',
                            hideOnOutsideClick: true,
                        }"
                        icon="pi pi-palette"
                        text
                        rounded
                        aria-label="Settings"
                    />
                    <AppConfig />

                </div>
                <Button type="button" class="topbar-theme-button" @click="minimizeWindow" icon="pi pi-minus" text rounded />
                    <Button type="button" class="topbar-theme-button" @click="toggleMaximizeWindow" text rounded>
                        <i :class="['pi', isMaximized ? 'pi-window-minimize' : 'pi-window-maximize']" />
                    </Button>
                    <Button type="button" class="topbar-theme-button topbar-close-button" icon="pi pi-times" @click="closeWindow" text rounded >
                        <i class="pi pi-times" style="color: red;"/>
                    </Button>
            </div>
        </div>
    </div>
</template>

