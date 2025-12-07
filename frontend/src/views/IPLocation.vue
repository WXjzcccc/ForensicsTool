<template>
  <div class="page-container">
    <!-- 输入卡片，支持抽屉式收起 -->
    <Card class="form-card" :class="{ 'collapsed': isCardCollapsed }">
      <template #content>
        <!-- 卡片内容区域 -->
        <div class="card-content" :class="{ 'collapsed-content': isCardCollapsed }">
          <!-- 批量IP输入区域 -->
          <div class="mb-3">
            <FloatLabel variant="on">
              <Textarea 
                id="batch-ips"
                v-model="batchIPs" 
                rows="5" 
                class="w-full"
                fluid
              />
              <label for="batch-ips">批量查询 (每行一个IP)</label>
            </FloatLabel>
          </div>

          <!-- 数据库操作和查询按钮行 -->
          <div class="button-row mb-4">
            <Button 
              label="检查数据库" 
              icon="pi pi-database" 
              @click="checkDatabase" 
              :loading="loading.checkDB"
              size="small"
              class="button-item"
            />
            <Button 
              label="更新数据库" 
              icon="pi pi-refresh" 
              @click="updateDatabase" 
              :loading="loading.updateDB"
              severity="secondary"
              size="small"
              class="button-item"
            />
            <Button 
              label="批量查询" 
              icon="pi pi-search" 
              @click="searchBatchIPs"
              :loading="loading.searchBatch"
              :disabled="!batchIPs || !dbLoaded"
              size="small"
              severity="primary"
              class="button-item"
            />
          </div>
        </div>
      </template>
    </Card>

    <!-- 视图切换和折叠按钮组 -->
    <div class="view-controls">
      <!-- 视图切换按钮组 -->
      <div class="view-switch-buttons">
        <Button 
          label="表格视图" 
          icon="pi pi-table" 
          :severity="viewMode === 'table' ? 'primary' : 'secondary'"
          @click="viewMode = 'table'"
          size="small"
        />
        <Button 
          label="地图视图" 
          icon="pi pi-map" 
          :severity="viewMode === 'map' ? 'primary' : 'secondary'"
          @click="viewMode = 'map'"
          size="small"
        />
      </div>
      
      <!-- 折叠按钮 -->
      <div class="collapse-button" @click="toggleCardCollapse">
        <i :class="isCardCollapsed ? 'pi pi-chevron-down' : 'pi pi-chevron-up'"></i>
      </div>
    </div>

    <!-- 结果卡片 -->
    <Card class="result-card" ref="resultCardRef">
      <template #content>
        <div v-if="searchResults.length === 0" class="empty">
          <Empty />
        </div>

        <!-- 表格视图 -->
        <div v-if="viewMode === 'table' && searchResults.length > 0" class="mb-4">
          <DataTable 
            :value="searchResults" 
            responsiveLayout="scroll"
            :loading="loading.search || loading.searchBatch"
            class="p-datatable-sm fixed-header-table"
            scrollable
            :scrollHeight="tableScrollHeight"
            :virtualScrollerOptions="{ itemSize: 46 }"
            :columnResizeMode="'fit'"
            :resizableColumns="true"
            :frozenWidth="tableFrozenWidth"
          >
            <Column field="ip" header="IP地址" :style="{ width: '150px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
            <Column field="country_name" header="国家" :style="{ width: '120px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
            <Column field="region_name" header="省份/地区" :style="{ width: '150px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
            <Column field="city_name" header="城市" :style="{ width: '150px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
            <Column field="district_name" header="区县" :style="{ width: '120px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
            <Column field="isp_domain" header="运营商" :style="{ width: '200px' }">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])" class="cursor-pointer">{{ data[field] }}</span>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- 地图视图 -->
        <div v-if="viewMode === 'map' && searchResults.length > 0" class="map-container">
          <!-- 地图加载动画 -->
          <div v-if="mapLoading" class="map-loading-overlay">
            <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="8" />
            <div class="mt-3">正在加载地图数据...</div>
          </div>
          <div ref="mapChart" class="map-chart"></div>
        </div>
      </template>
    </Card>
    <Toast position="bottom-right"/>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { useToast } from 'primevue/usetoast';
import * as echarts from 'echarts';
import { CheckDB, LoadDB, CheckUpdate, UpdateDB, Search, SearchAll } from '../../wailsjs/go/ip/IP';
import { usePageDataStore } from "@/store";
import FloatLabel from 'primevue/floatlabel';
import Empty from '@/components/Empty.vue';
import { useTableHeight } from '@/composables/useTableHeight.js';
import { ClipboardSetText } from "../../wailsjs/runtime/runtime.js";

// 状态变量
const toast = useToast();
const store = usePageDataStore();
const batchIPs = ref(store.ipLocationStore?.formData?.batchIPs || '');
const dbLoaded = ref(store.ipLocationStore?.dbLoaded || false);
const dbUpdateAvailable = ref(store.ipLocationStore?.dbUpdateAvailable || false);
const searchResults = ref(store.ipLocationStore?.searchResults || []);
const viewMode = ref(store.ipLocationStore?.viewMode || 'table'); // table, map
const mapChart = ref(null);
const chartInstance = ref(null);
const resultCardRef = ref(null);
// 从store中获取卡片折叠状态
const isCardCollapsed = ref(store.ipLocationStore?.isCardCollapsed || false);
const mapLoading = ref(false); // 地图加载状态
const mapInitialized = ref(false); // 地图是否已初始化

// 动态计算DataTable的scrollHeight
const tableScrollHeight = ref('60vh');
const tableFrozenWidth = ref('740px'); // 表格冻结宽度，确保所有列宽之和

// 使用公共的表格高度计算函数
const { calculateTableHeight, handleResize } = useTableHeight(resultCardRef, tableScrollHeight);

// 组件挂载后设置监听
onMounted(() => {
  window.addEventListener('resize', handleResize);
  // 初始计算
  nextTick(() => {
    calculateTableHeight();
  });
  // 初始化数据库状态
  checkDatabase();
});

// 组件卸载时移除监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

// 加载状态
const loading = ref({
  checkDB: false,
  loadDB: false,
  updateDB: false,
  search: false,
  searchBatch: false
});

// 显示消息
const showMessage = (severity, summary, detail) => {
  toast.add({ severity, summary, detail, life: 3000 });
};

// 切换卡片收起状态
const toggleCardCollapse = () => {
  isCardCollapsed.value = !isCardCollapsed.value;
  // 折叠状态变化后重新计算表格高度
  nextTick(() => {
    calculateTableHeight();
  });
};

// 保存数据到store
watch([batchIPs, dbLoaded, dbUpdateAvailable, searchResults, viewMode, isCardCollapsed], () => {
  store.saveIPLocationData({
    formData: {
      batchIPs: batchIPs.value
    },
    dbLoaded: dbLoaded.value,
    dbUpdateAvailable: dbUpdateAvailable.value,
    searchResults: searchResults.value,
    viewMode: viewMode.value,
    isCardCollapsed: isCardCollapsed.value
  });
});

// 检查数据库状态
const checkDatabase = () => {
  loading.value.checkDB = true;
  CheckDB()
    .then(() => {
      dbLoaded.value = true;
      showMessage('success', '成功', '数据库已加载');
      
      // 先加载数据库，然后检查是否有更新
      return LoadDB();
    })
    .then(() => {
      return CheckUpdate();
    })
    .then((hasUpdate) => {
      if (hasUpdate) {
        dbUpdateAvailable.value = true;
        showMessage('info', '提示', '数据库有新版本可用');
      }
    })
    .catch((error) => {
      dbLoaded.value = false;
      showMessage('error', '错误', `数据库检查失败: ${error}`);
      // 尝试加载数据库
      loadDatabase();
    })
    .finally(() => {
      loading.value.checkDB = false;
    });
};

// 加载数据库
const loadDatabase = () => {
  loading.value.loadDB = true;
  LoadDB()
    .then(() => {
      dbLoaded.value = true;
      showMessage('success', '成功', '数据库加载成功');
    })
    .catch((error) => {
      showMessage('error', '错误', `数据库加载失败: ${error}`);
    })
    .finally(() => {
      loading.value.loadDB = false;
    });
};

// 更新数据库
const updateDatabase = () => {
  loading.value.updateDB = true;
  // 直接执行更新，然后重新加载
  UpdateDB()
    .then(() => {
      // 更新后需要重新加载
      return LoadDB();
    })
    .then(() => {
      dbUpdateAvailable.value = false;
      dbLoaded.value = true;
      showMessage('success', '成功', '数据库更新成功');
    })
    .catch((error) => {
      showMessage('error', '错误', `数据库更新失败: ${error}`);
    })
    .finally(() => {
      loading.value.updateDB = false;
    });
};

// 搜索批量IP
const searchBatchIPs = () => {
  if (!batchIPs.value.trim()) {
    showMessage('warn', '警告', '请输入IP地址');
    return;
  }
  
  if (!dbLoaded.value) {
    showMessage('warn', '警告', '请先加载数据库');
    return;
  }

  const ipList = batchIPs.value
    .split('\n')
    .map(ip => ip.trim())
    .filter(ip => ip);
    
  if (ipList.length === 0) {
    showMessage('warn', '警告', '请输入有效的IP地址');
    return;
  }

  loading.value.searchBatch = true;
  SearchAll(ipList)
    .then((results) => {
      // 转换结果为数组格式
      const formattedResults = Object.entries(results).map(([ip, data]) => ({
        ip,
        ...data
      }));
      
      searchResults.value = formattedResults;
      showMessage('success', '成功', `成功查询 ${formattedResults.length} 个IP地址`);
      
      // 查询成功后自动收起输入卡片
      if (!isCardCollapsed.value) {
        isCardCollapsed.value = true;
        // 折叠后重新计算表格高度
        nextTick(() => {
          calculateTableHeight();
        });
      }
    })
    .catch((error) => {
      showMessage('error', '错误', `批量查询失败: ${error}`);
    })
    .finally(() => {
      loading.value.searchBatch = false;
    });
};

// 初始化地图
const initMap = () => {
  if (!mapChart.value || mapInitialized.value) return;
  
  // 销毁已存在的图表实例
  if (chartInstance.value) {
    chartInstance.value.dispose();
  }
  
  chartInstance.value = echarts.init(mapChart.value);
  
  // 显示加载动画
  mapLoading.value = true;
  
  // 加载地图数据
  fetch('https://geo.datav.aliyun.com/areas_v3/bound/100000_full_city.json')
    .then((response) => response.json())
    .then((mapData) => {
      echarts.registerMap('china', mapData);
      
      // 设置初始地图选项
      const option = {
        title: {
          text: 'IP归属地分布',
          left: 'center'
        },
        tooltip: {
          trigger: 'item',
          formatter: '{b}: {c} 个IP'
        },
        series: [
          {
            name: 'IP数量',
            type: 'map',
            map: 'china',
            roam: true,
            emphasis: {
              label: {
                show: true
              }
            },
            data: []
          }
        ]
      };
      
      chartInstance.value.setOption(option);
      mapInitialized.value = true;
    })
    .catch((error) => {
      console.error('地图数据加载失败:', error);
      showMessage('error', '错误', '地图数据加载失败');
    })
    .finally(() => {
      mapLoading.value = false;
    });
};

// 更新地图数据
const updateMapData = () => {
  if (!chartInstance.value || !searchResults.value.length) return;
  
  // 统计各城市的IP数量
  const cityStats = {};
  searchResults.value.forEach(item => {
    if (item.city_name) {
      cityStats[item.city_name] = (cityStats[item.city_name] || 0) + 1;
    }
  });
  
  // 转换为ECharts数据格式
  const mapData = Object.entries(cityStats).map(([name, value]) => ({
    name,
    value
  }));
  
  // 更新地图
  chartInstance.value.setOption({
    series: [
      {
        data: mapData
      }
    ],
    visualMap: {
      max: Math.max(...mapData.map(item => item.value), 10)
    }
  });
};

// 监听搜索结果变化，更新地图
watch(searchResults, () => {
  if (viewMode.value === 'map') {
    nextTick(() => {
      updateMapData();
    });
  }
});

// 处理单元格点击事件
const handleCellClick = (value) => {
  // 复制单元格内容
  ClipboardSetText(value).then((ok) => {
    if (ok) {
      toast.add({ severity: 'success', summary: '成功', detail: '已复制到剪贴板', life: 2000 });
    }
  });
};

// 监听视图模式变化
watch(viewMode, (newMode, oldMode) => {
  // 如果从表格视图切换到地图视图，初始化地图
  if (newMode === 'map' && oldMode === 'table') {
    nextTick(() => {
      if (!mapInitialized.value) {
        initMap();
      } else if (chartInstance.value) {
        updateMapData();
      }
    });
  }
  // 如果从地图视图切换到表格视图，销毁地图实例以释放资源
  else if (oldMode === 'map' && newMode === 'table') {
    nextTick(() => {
      if (chartInstance.value) {
        chartInstance.value.dispose();
        chartInstance.value = null;
        mapInitialized.value = false;
      }
      // 切换到表格视图后重新计算表格高度
      calculateTableHeight();
    });
  }
});
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 500px;
  position: relative;
}

.map-chart {
  width: 100%;
  height: 100%;
}

.map-loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.8);
  z-index: 10;
}

.p-dark .map-loading-overlay {
  background-color: rgba(0, 0, 0, 0.7);
}

/* 视图控制区域样式 */
.view-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: var(--surface-card);
  border-top: 1px solid var(--surface-border);
  border-bottom: 1px solid var(--surface-border);
  padding: 0.5rem 1rem;
}

/* 视图切换按钮组样式 */
.view-switch-buttons {
  display: flex;
  gap: 0.5rem;
}

/* 折叠按钮样式 */
.collapse-button {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 30px;
  height: 30px;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.collapse-button:hover {
  background-color: var(--surface-hover);
}

.collapse-button i {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
}

/* 卡片收起状态样式 */
.form-card {
  transition: all 0.3s ease;
  overflow: hidden;
  margin-bottom: 0;
}

.form-card.collapsed {
  height: 0;
  opacity: 0;
  margin: 0;
  padding: 0;
  border: none;
}

.card-content {
  transition: all 0.3s ease;
  max-height: 500px;
  opacity: 1;
}

.card-content.collapsed-content {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
  margin: 0;
  padding: 0;
}

/* 按钮行样式 */
.button-row {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.button-item {
  flex: 1;
  min-width: 120px;
}

/* 固定表头宽度样式 */
.fixed-header-table .p-datatable-thead > tr > th {
  position: sticky !important;
  top: 0 !important;
  z-index: 10 !important;
  background-color: var(--surface-card) !important;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1) !important;
  min-width: inherit !important;
  max-width: inherit !important;
  width: inherit !important;
}

.fixed-header-table .p-datatable-tbody > tr > td {
  min-width: inherit !important;
  max-width: inherit !important;
  width: inherit !important;
}

/* 确保表格在滚动时列宽保持一致 */
.fixed-header-table .p-datatable-scrollable-header {
  overflow: hidden !important;
}

.fixed-header-table .p-datatable-scrollable-body {
  overflow: auto !important;
}

/* 防止内容溢出 */
.fixed-header-table .p-datatable-scrollable-body td {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .button-row {
    flex-direction: column;
  }
  
  .button-item {
    width: 100%;
  }
}
</style>