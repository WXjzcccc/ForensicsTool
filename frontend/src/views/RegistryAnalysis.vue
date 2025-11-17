<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <div class="field full-width">
          <FloatLabel variant="on">
            <InputText 
              id="file"
              v-model="form.file" 
              aria-autocomplete="none"
              placeholder="请拖入目录，目录包含SYSTEM、SAM、SOFRWARE和用户注册表文件" 
              @drop.prevent="handleDrop"
              @dragover.prevent
              v-tooltip.top="'请拖入目录，目录包含SYSTEM、SAM、SOFRWARE和用户注册表文件'"
            />
            <label for="file">文件夹</label>
          </FloatLabel>
        </div>
      </form>
      <div class="button-row">
        <Button class="button equal-width" @click="handleExtract">
          <i class="pi pi-search"></i>
          解析
        </Button>
        <Button class="button equal-width" severity="secondary" @click="handleClear">
          <i class="pi pi-trash"></i>
          清空输出
        </Button>
      </div>
      </template>
    </Card>
    <Card class="result-card" ref="resultCardRef">
      <template #content>
      <div v-if="Object.keys(tableData).length === 0" class="empty">
        <Empty />
      </div>
      <Tabs v-if="Object.keys(tableData).length !== 0" :value="getActiveTabIndex" @tab-change="handleTabChange">
        <TabList>
          <Tab v-for="(tab,index) in tabs" :key="tab.label" :value="index">{{ tab.label }}</Tab>
        </TabList>
        <TabPanels>
        <TabPanel v-for="(tab,index) in tabs" :key="tab.value" :value="index">
          <DataTable
            :value="tableData[tab.value]"
            :columns="columns[tab.value]"
            dataKey="id"
            stripedRows
            removableSort
            :resizableColumns="true"
            columnResizeMode="fit"
            scrollable
            :scrollHeight="tableScrollHeight"
            @column-reorder="onColumnReorder"
            :reorderableColumns="true"
          >
            <Column v-for="col in columns[tab.value]" :key="col.colKey" :field="col.colKey" :header="col.title" :sortable="true">
              <template #body="{ data, field }">
                <span @click="handleCellClick(data[field])">{{ data[field] }}</span>
              </template>
            </Column>
            <template #empty>
              <Empty />
            </template>
          </DataTable>
        </TabPanel>
        </TabPanels>
      </Tabs>
      </template>
    </Card>
    <Toast position="bottom-right"/>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useToast } from 'primevue/usetoast'
import { usePageDataStore } from "@/store"
import { watch } from "vue"
import { ClipboardSetText, OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime.js"
import {
  ExtractDbeaver, 
  ExtractFinalShell, ExtractHawk2, ExtractMetaMask,
  ExtractMobaXterm,
  ExtractNavicat, ExtractXShell
} from "../../wailsjs/go/extractor/InfoExtractor.js"
import { AnalyzeWinReg } from "../../wailsjs/go/winreg/Reg.js"
import FloatLabel from 'primevue/floatlabel'
import Tabs from 'primevue/tabs'
import Empty from '@/components/Empty.vue'
import { useTableHeight } from '@/composables/useTableHeight.js'

const toast = useToast()
const store = usePageDataStore()
const form = ref(store.registryStore?.formData || {
  file: '',
})
const tableData = ref(store.registryStore?.tableData || {})
const tabs = ref(store.registryStore?.tabsData || [])
const activeTab = ref(store.registryStore?.tabData || '')
const columns = ref(store.registryStore?.columnsData || {})

// 动态计算DataTable的scrollHeight
const tableScrollHeight = ref('54vh')
const resultCardRef = ref(null)

// 使用公共的表格高度计算函数
const { calculateTableHeight, handleResize } = useTableHeight(resultCardRef, tableScrollHeight)

// 计算属性，用于获取当前活动标签的索引
const getActiveTabIndex = computed(() => {
  const index = tabs.value.findIndex(tab => tab.value === activeTab.value)
  return index >= 0 ? index : 0
})

// 组件挂载后设置监听
onMounted(() => {
  window.addEventListener('resize', handleResize)
  // 初始计算
  nextTick(() => {
    calculateTableHeight()
  })
  
})

// 组件卸载时移除监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 清理文件拖放监听器
  OnFileDropOff()
})

watch([form, tableData, tabs, activeTab], () => {
  store.saveRegistryData({
    formData: form.value,
    tableData: tableData.value,
    tabsData: tabs.value,
    tabData: activeTab.value,
    columnsData: columns.value
  })
})

const loading = ref(false)

const handleExtract = () => {
  loading.value = true
  let file = form.value.file
  if (file === "" || file === undefined || file === null) {
    toast.add({ severity: 'error', summary: '错误', detail: '参数错误！', life: 3000 })
    loading.value = false
    return
  }
  AnalyzeWinReg(file).then((result) => {
  
    try {
      if (result.err !== "") {
        toast.add({ severity: 'error', summary: '错误', detail: result.err, life: 3000 })
      } else {
        handleResult(result)
      }
    } finally {
      loading.value = false
    }
  })
}

const handleResult = (result) => {
  columns.value = {}
  tabs.value = []
  tableData.value = {}
  let idx = 0
  for (let tab in result.data) {
    if (idx === 0) {
      activeTab.value = tab
      idx++
    }
    tabs.value.push({
      value: tab, label: tab
    })
    if (result.data[tab] == null) {
        toast.add({ severity: 'warn', summary: '提示', detail: tab + '无数据', life: 3000 })
      continue
    }
    let inf = result.data[tab][0]
    tableData.value[tab] = {}
    columns.value[tab] = []
    tableData.value[tab] = result.data[tab]
    for (let key in inf) {
      columns.value[tab].push({
        colKey: key, title: key, ellipsis: true, width: "15vw"
      })
    }
  }
  
  // 数据加载完成后重新计算表格高度
  nextTick(() => {
    calculateTableHeight()
  })
}

const handleClear = () => {
  tableData.value = {}
  tabs.value = []
  activeTab.value = ""
  columns.value = {}
}

const handleTabChange = (event) => {
  activeTab.value = tabs.value[event.index].value
}

//TODO 实现复制，现在有问题
const handleCellClick = (value) => {
  // 复制单元格内容
  ClipboardSetText(value).then((ok) => {
    if (ok) {
      toast.add({ severity: 'success', summary: '成功', detail: '已复制到剪贴板', life: 2000 })
    }
  })
}

const handleDrop = (event) => {
  OnFileDrop((x, y, paths) => {
    if (paths.length > 0) {
      form.value.file = paths[0]
      toast.add({ severity: 'success', summary: '成功', detail: '文件已添加', life: 3000 })
    }
  }, false)
}

const onColumnReorder = (event) => {
  columns.value[activeTab.value] = event.columns
}
</script>

<style scoped>
</style>