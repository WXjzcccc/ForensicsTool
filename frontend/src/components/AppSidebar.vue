<template>
    <div class="sidebar" :class="{ 'sidebar-collapsed': isCollapsed }" style="--wails-draggable: drag;">
        <Menu :model="menuItems" class="sidebar-menu">
            <template #start>
                <div class="sidebar-header">
                    <Button @click="toggleSidebar" class="toggle-button"
                        :icon="isCollapsed ? 'pi pi-bars' : 'pi pi-bars'" text />
                </div>
            </template>

            <template #item="{ item, props }">
                <a v-ripple class="menu-link" v-bind="props.action" @click="navigateTo(item.route)"
                    :class="{ 'active': $route.name === item.name }">
                    <span :class="item.icon" />
                    <span v-if="!isCollapsed">{{ item.label }}</span>
                </a>
            </template>

            <template #end>
                <div class="sidebar-footer" style="--wails-draggable: no-drag;">
                    <Button v-if="!isCollapsed" @click="drawerVisible = true" text> {{ version }} </Button>
                    <Button v-else @click="drawerVisible = true" icon="pi pi-info-circle" text />
                </div>
            </template>
        </Menu>
    </div>
    <Drawer v-model:visible="drawerVisible" header="小Tips" position="right" class="sidebar-tips">
        <Message size="small">
            1.小红书，密码为xhsdev或mmkv文件com.xingin.xhs_preferences中msg_db_password_updated的值，选择sqlcipher3直接解密
        </Message>
        <Message size="small">2.微信的imei，现在可以通过files/KeyInfo.bin获取了，需要解密文件，算法为RC4，密钥为_wEcHAT_</Message>
        <Message size="small">3.MosGram(泡泡)，密码为cust_id的md5值，在sp目录的account_config.xml文件中，使用SQLCipher4参数解</Message>
        <Message size="small">4.悟空IM系列的聊天数据库，数据库名为wk_用户ID.db，解密密码即为用户ID，使用SQLCipher4参数解密</Message>
        <Message size="small">5.抖音的aweme_database_数据库，密码为aweme_database_passphrase，使用wcdb参数解密</Message>
    </Drawer>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useLayout } from '../composables/useLayout';
import { GetVersion } from "../../wailsjs/go/main/App"

const { isDarkMode, toggleDarkMode } = useLayout();
const router = useRouter();
const route = useRoute();
const isCollapsed = ref(false);
const windowWidth = ref(window.innerWidth);
const version = ref('');
const drawerVisible = ref(false);

// 菜单项数据 - 移除了分组，所有菜单项放在同一层级
const menuItems = ref([
    {
        label: '密钥计算',
        icon: 'pi pi-key',
        command: () => navigateTo('/KeyCalculation'),
        name: 'KeyCalculation'
    },
    {
        label: '数据库解密',
        icon: 'pi pi-lock-open',
        command: () => navigateTo('/DatabaseDecrypt'),
        name: 'DatabaseDecrypt'
    },
    {
        label: '数据提取',
        icon: 'pi pi-database',
        command: () => navigateTo('/DataExtraction'),
        name: 'DataExtraction'
    },
    {
        label: '暴力破解',
        icon: 'pi pi-lock',
        command: () => navigateTo('/BruteForce'),
        name: 'BruteForce'
    },
    {
        label: '注册表分析',
        icon: 'pi pi-microsoft',
        command: () => navigateTo('/RegistryAnalysis'),
        name: 'RegistryAnalysis'
    },
    {
        label: '时间戳转换',
        icon: 'pi pi-clock',
        command: () => navigateTo('/TimestampParser'),
        name: 'TimestampParser'
    },
    {
        label: '文件读取',
        icon: 'pi pi-file-arrow-up',
        command: () => navigateTo('/FileReader'),
        name: 'FileReader'
    },
    {
        label: 'IP归属地查询',
        icon: 'pi pi-globe',
        command: () => navigateTo('/IPLocation'),
        name: 'IPLocation'
    },
    {
        label: '关于',
        icon: 'pi pi-info-circle',
        command: () => navigateTo('/About'),
        name: 'About'
    }
]);

// 切换侧边栏状态
const toggleSidebar = () => {
    isCollapsed.value = !isCollapsed.value;
};

// 导航到指定路由
const navigateTo = (routePath) => {
    if (routePath) {
        router.push(routePath);
    }
};

// 更新窗口宽度
const updateWindowWidth = () => {
    windowWidth.value = window.innerWidth;
    // 根据窗口宽度自动调整侧边栏状态
    if (windowWidth.value < 768) {
        isCollapsed.value = true;
    }
    // 展开状态使用CSS中的相对尺寸（15vw）
};

// 监听窗口大小变化
onMounted(() => {
    updateWindowWidth();
    window.addEventListener('resize', updateWindowWidth);
    nextTick(() => {
        if (GetVersion == undefined) {
            version.value = '1.0.0'
        } else {
            GetVersion().then(res => {
                version.value = res
            })
        }
    })
});

onUnmounted(() => {
    window.removeEventListener('resize', updateWindowWidth);
});
</script>

<style scoped>
.sidebar {
    width: 15vw;
    min-width: 60px;
    max-width: 300px;
    background-color: var(--p-surface-100);
    border-right: 1px solid var(--p-surface-200);
    display: flex;
    flex-direction: column;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    height: 100vh;
    position: relative;
    will-change: width;
}

.p-dark .sidebar {
    background-color: var(--p-surface-900);
    border-right-color: var(--p-surface-700);
}

.sidebar-collapsed {
    width: 5vw;
}

.sidebar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    border-bottom: 1px solid var(--p-surface-200);
    height: 8vh;
    min-height: 50px;
    max-height: 70px;
}

.p-dark .sidebar-header {
    border-bottom-color: var(--p-surface-700);
}

.toggle-button {
    color: var(--p-surface-600);
}

.p-dark .toggle-button {
    color: var(--p-surface-400);
}

.sidebar-menu {
    width: 100%;
    height: 100%;
    border: none;
    background: transparent;
}

.sidebar-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    border-top: 1px solid var(--p-surface-200);
}

.p-dark .sidebar-footer {
    border-top-color: var(--p-surface-700);
}

.version-text {
    font-size: 0.75rem;
    color: var(--p-primary-600);
}

.p-dark .version-text {
    color: var(--p-primary-400);
}

.menu-link {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    color: var(--p-surface-700);
    text-decoration: none;
    transition: background-color 0.2s, color 0.2s;
    border-radius: 0;
    cursor: pointer;
    width: 100%;
    box-sizing: border-box;
    font-size: 0.8rem;
}

.p-dark .menu-link {
    color: var(--p-surface-300);
}

.menu-link:hover {
    background-color: var(--p-surface-200);
    color: var(--p-surface-900);
}

.p-dark .menu-link:hover {
    background-color: var(--p-surface-800);
    color: var(--p-surface-0);
}

.menu-link.active {
    background-color: var(--p-primary-100);
    color: var(--p-primary-700);
    border-right: 3px solid var(--p-primary-500);
}

.p-dark .menu-link.active {
    background-color: color-mix(in srgb, var(--p-primary-400), transparent 80%);
    color: var(--p-primary-300);
}

.p-message {
    margin-top: .5rem;
}

/* 收起状态下的样式调整 */
.sidebar-collapsed .sidebar-header {
    padding: 1rem 0;
}

.sidebar-collapsed .menu-link {
    justify-content: center;
    padding: 0.75rem 0;
}

/* 确保收起状态下菜单项宽度不超过侧边栏宽度 */
.sidebar-collapsed :deep(.p-menu) {
    width: 100%;
    min-width: 100%;
    max-width: 100%;
}

.sidebar-collapsed :deep(.p-menuitem-link) {
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    min-width: 100%;
    max-width: 100%;
    box-sizing: border-box;
}

.sidebar-collapsed :deep(.p-menuitem-content) {
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    min-width: 100%;
    max-width: 100%;
    box-sizing: border-box;
}

.sidebar-collapsed :deep(.p-menuitem) {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    min-width: 100%;
    max-width: 100%;
    box-sizing: border-box;
}

.sidebar-collapsed :deep(.p-menuitem-icon) {
    margin: 0;
}

.sidebar-collapsed :deep(.p-menuitem-text) {
    display: none;
}

/* 调整PrimeVue Menu组件的样式 */
:deep(.p-menu) {
    border: none;
    background: transparent;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    height: 100%;
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.p-menu-list) {
    flex: 1;
    overflow-y: auto;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box;
}

:deep(.p-menuitem) {
    margin: 0;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box;
}

:deep(.p-menuitem-content) {
    padding: 0;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box;
}

:deep(.p-submenu-header) {
    padding: 0;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box;
}

:deep(.p-menuitem-link) {
    padding: 0;
    width: 100% !important;
    min-width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>