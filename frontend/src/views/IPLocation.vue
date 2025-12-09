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
              <label for="batch-ips">批量查询 (每行一个IP，数据来源纯真IP数据库)</label>
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
import { useLayout } from '@/composables/useLayout';

// 状态变量
const toast = useToast();
const store = usePageDataStore();
const { isDarkMode, primary } = useLayout(); // 获取主题状态
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
        showMessage('info', '提示', '数据库有新版本可用，请点击更新按钮下载');
      } else {
        showMessage('info', '提示', '数据库已是最新版本');
      }
    })
    .catch((error) => {
      dbLoaded.value = false;
      showMessage('error', '错误', `数据库检查失败: ${error}`);
      // 检查数据库失败时，直接检查更新
      return CheckUpdate();
    })
    .then((hasUpdate) => {
      // 如果是catch分支过来的，这里会检查是否有更新
      if (hasUpdate !== undefined) {
        if (hasUpdate) {
          dbUpdateAvailable.value = true;
          showMessage('info', '提示', '数据库有新版本可用，请点击更新按钮下载');
        } else {
          showMessage('info', '提示', '数据库检查失败，但无可用更新');
        }
      }
    })
    .catch((error) => {
      showMessage('error', '错误', `检查更新失败: ${error}`);
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
      // 注册地图数据
      echarts.registerMap('china', mapData);
      
      // 根据主题设置颜色
      const isDark = isDarkMode.value;
      const primaryColor = primary.value;
      
      // 定义主题颜色
      const getThemeColors = () => {
        if (isDark) {
          // 深色主题颜色
          return {
            bgColors: ['#1e3a8a', '#1e40af', '#2563eb', '#3b82f6', '#60a5fa'],
            textColor: '#e5e7eb',
            titleColor: '#f3f4f6',
            graphicColor: '#9ca3af'
          };
        } else {
          // 浅色主题颜色
          return {
            bgColors: ['#e0f2fe', '#bae6fd', '#7dd3fc', '#38bdf8', '#0ea5e9'],
            textColor: '#374151',
            titleColor: '#111827',
            graphicColor: '#6b7280'
          };
        }
      };
      
      const themeColors = getThemeColors();
      
      // 设置初始地图选项
      const option = {
        backgroundColor: 'transparent',
        title: {
          text: 'IP归属地分布',
          left: 'center',
          textStyle: {
            color: themeColors.titleColor
          }
        },
        tooltip: {
          trigger: 'item',
          formatter: function(params) {
            // 如果值为0或undefined，不显示
            if (!params.value || params.value === 0) {
              return `${params.name}: 暂无数据`;
            }
            return `${params.name}: ${params.value} 个IP`;
          },
          backgroundColor: isDark ? 'rgba(31, 41, 55, 0.9)' : 'rgba(255, 255, 255, 0.9)',
          borderColor: isDark ? 'rgba(75, 85, 99, 0.5)' : 'rgba(229, 231, 235, 0.5)',
          textStyle: {
            color: themeColors.textColor
          }
        },
        visualMap: {
          min: 0,
          max: 10,
          left: 'left',
          top: 'bottom',
          text: ['高', '低'],
          calculable: true,
          textStyle: {
            color: themeColors.textColor
          },
          inRange: {
            color: themeColors.bgColors
          },
          handleStyle: {
            color: isDark ? '#4b5563' : '#d1d5db'
          }
        },
        graphic: {
          type: 'text',
          right: 10,
          bottom: 10,
          style: {
            text: '地图数据来源：阿里DataV平台',
            fontSize: 12,
            fill: themeColors.graphicColor
          }
        },
        series: [
          {
            name: 'IP数量',
            type: 'map',
            map: 'china',
            roam: true,
            emphasis: {
              label: {
                show: true,
                color: themeColors.textColor
              },
              itemStyle: {
                areaColor: isDark ? '#374151' : '#f3f4f6'
              }
            },
            select: {
              label: {
                color: themeColors.textColor
              },
              itemStyle: {
                areaColor: isDark ? '#4b5563' : '#e5e7eb'
              }
            },
            itemStyle: {
              areaColor: isDark ? '#1f2937' : '#f9fafb',
              borderColor: isDark ? '#374151' : '#d1d5db'
            },
            data: [],
            // 存储地图区域数据，用于后续匹配
            mapData: mapData.features.map(feature => ({
              name: feature.properties.name,
              adcode: feature.properties.adcode,
              level: feature.properties.level
            }))
          }
        ]
      };
      
      chartInstance.value.setOption(option);
      mapInitialized.value = true;
      
      // 如果已有搜索结果，更新地图数据
      if (searchResults.value.length > 0) {
        updateMapData();
      }
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
  
  // 统计各城市/区县的IP数量
  const cityStats = {};
  // 中国的直辖市列表
  const municipalities = ['北京市', '天津市', '上海市', '重庆市'];
  // 直辖市简称列表
  const municipalityShortNames = ['北京', '天津', '上海', '重庆'];
  
  searchResults.value.forEach(item => {
    if (item.city_name) {
      // 处理直辖市情况
      let targetCity = null;
      
      // 检查是否是直辖市（全称或简称）
      const isMunicipality = municipalities.includes(item.region_name) || 
                             municipalityShortNames.includes(item.region_name);
      
      if (isMunicipality) {
        // 如果是直辖市，优先使用city_name（区县）
        // 因为地图数据中直辖市区域只显示区县
        targetCity = item.district_name;
      } else if (item.city_name && item.city_name !== item.region_name) {
        // 如果city_name和region_name不同，可能是直辖市的区县
        // 检查region_name是否是直辖市
        const regionIsMunicipality = municipalities.includes(item.region_name) || 
                                    municipalityShortNames.includes(item.region_name);
        
        if (regionIsMunicipality) {
          // 如果region_name是直辖市，则使用city_name（区县）
          targetCity = item.district_name;
        } else {
          // 否则使用city_name
          targetCity = item.city_name;
        }
      } else {
        // 其他情况使用city_name
        targetCity = item.city_name;
      }
      
      // 统计到目标城市/区县
      if (targetCity) {
        cityStats[targetCity] = (cityStats[targetCity] || 0) + 1;
      }
    }
  });
  
  // 获取地图中的所有区域名称
  const mapOption = chartInstance.value.getOption();
  const mapRegions = mapOption.series[0].mapData || [];
  const regionNames = mapRegions.map(region => region.name);
  
  // 创建城市名称映射表，处理可能的名称差异
  const cityMapping = {};
  // 直辖市名称映射表 - 这里不需要将区县映射到直辖市
  // 因为地图数据中直辖市区域只显示区县
  
  // 尝试匹配城市名称
  Object.keys(cityStats).forEach(cityName => {
    // 直接匹配
    if (regionNames.includes(cityName)) {
      cityMapping[cityName] = cityName;
    } else {
      // 尝试模糊匹配
      const matchedRegion = regionNames.find(region => 
        region.includes(cityName) || cityName.includes(region)
      );
      if (matchedRegion) {
        cityMapping[cityName] = matchedRegion;
      } else {
        // 如果没有匹配，使用原始名称
        cityMapping[cityName] = cityName;
      }
    }
  });
  
  // 转换为ECharts数据格式，使用映射后的名称
  const mapData = Object.entries(cityStats).map(([name, value]) => ({
    name: cityMapping[name] || name,
    value: value || 0  // 确保值不为undefined或NaN
  }));
  
  // 过滤掉值为0的数据项，避免在地图上显示"暂无数据"
  const filteredMapData = mapData.filter(item => item.value > 0);
  
  console.log('城市统计数据:', cityStats);
  console.log('地图区域名称:', regionNames);
  console.log('城市名称映射:', cityMapping);
  console.log('地图数据:', filteredMapData);
  
  // 更新地图
  chartInstance.value.setOption({
    series: [
      {
        data: filteredMapData
      }
    ],
    visualMap: {
      max: Math.max(...filteredMapData.map(item => item.value), 10)
    }
  });
  
  // 保存地图区域数据到全局变量，以便后续使用
  if (!window.mapRegions) {
    window.mapRegions = mapRegions;
  }
};

// 监听搜索结果变化，更新地图
watch(searchResults, () => {
  if (viewMode.value === 'map') {
    nextTick(() => {
      updateMapData();
    });
  }
});

// 监听主题变化，重新初始化地图
watch([isDarkMode, primary], () => {
  if (viewMode.value === 'map' && mapInitialized.value) {
    // 保存当前地图数据
    const currentOption = chartInstance.value.getOption();
    const currentData = currentOption.series[0].data;
    
    // 重新初始化地图
    mapInitialized.value = false;
    initMap();
    
    // 等待地图初始化完成后恢复数据
    nextTick(() => {
      setTimeout(() => {
        if (chartInstance.value && currentData.length > 0) {
          chartInstance.value.setOption({
            series: [{
              data: currentData
            }]
          });
        }
      }, 500);
    });
  }
});

// 监听窗口大小变化，在地图视图下重新计算高度
watch(() => tableScrollHeight.value, () => {
  if (viewMode.value === 'map' && chartInstance.value) {
    // 延迟执行，确保DOM已更新
    setTimeout(() => {
      chartInstance.value.resize();
    }, 100);
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
      // 切换到地图视图后重新计算地图高度
      calculateTableHeight();
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
/* 地图容器样式 */
.map-container {
  width: 100%;
  height: v-bind(tableScrollHeight);
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

/* 卡片过渡效果 */
.form-card {
  transition: all 0.3s ease;
  overflow: hidden;
  margin-bottom: 0;
}

.card-content {
  transition: all 0.3s ease;
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