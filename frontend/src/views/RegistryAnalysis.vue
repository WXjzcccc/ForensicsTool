<template>
  <div class="page-container">
    <t-card class="form-card">
      <t-form layout="inline" label-width="calc(2em + 7vh)"  label-align="left">
        <t-form-item label="file" style="width: 100vh">
          <t-input v-model="form.file" placeholder="请拖入文件或目录，目录包含SYSTEM、SAM、SOFRWARE和用户注册表文件" @drop.prevent="handleDrop"
                   @dragover.prevent/>
        </t-form-item>
      </t-form>
      <div style="height: 1vh"></div>
      <t-button class="button" theme="primary" @click="handleExtract"><template #icon><t-icon name="search"/></template>解析</t-button>
      <t-button class="button" theme="primary" @click="handleClear"><template #icon><t-icon name="clear-formatting"/></template>清空输出</t-button>
    </t-card>

    <t-card class="result-card" :loading="loading">
      <t-tabs v-model="activeTab" @change="handleTabChange">
        <t-tab-panel v-for="tab in tabs" :key="tab.value" :value="tab.value" :label="tab.label">
          <t-table
              :data="tableData[tab.value]"
              :columns="columns[tab.value]"
              row-key="id"
              stripe
              hover
              size="medium"
              fixed-rows="[0,1]"
              @cell-click="handleCellClick"
              lazy-load
              bordered
              max-height="68vh"
              :scroll="{ type: 'virtual' }"
              dragSort='col'
              @drag-sort="onDragSort"
          >
            <template #empty>
              <div class="empty">暂无数据</div>
            </template>
          </t-table>
        </t-tab-panel>
      </t-tabs>
    </t-card>
  </div>
</template>

<script setup>
import {ref} from 'vue'
import {MessagePlugin} from 'tdesign-vue-next'
import {usePageDataStore} from "@/store/index.js";
import {watch} from "vue";
import {ClipboardSetText, OnFileDrop} from "../../wailsjs/runtime/runtime.js";
import {
  ExtractDbeaver, 
  ExtractFinalShell, ExtractHawk2, ExtractMetaMask,
  ExtractMobaXterm,
  ExtractNavicat, ExtractXShell
} from "../../wailsjs/go/extractor/InfoExtractor.js";
import {AnalyzeWinReg} from "../../wailsjs/go/winreg/Reg.js";
const store = usePageDataStore()
const form = ref(store.registryStore?.formData || {
  file: '',
})
const tableData = ref(store.registryStore?.tableData || {
})
const tabs = ref(store.registryStore?.tabsData || [
])
const activeTab = ref(store.registryStore?.tabData || '')
const columns = ref(store.registryStore?.columnsData || {
})
watch([form,tableData,tabs,activeTab],()=>{
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
  AnalyzeWinReg(file).then((result) => {
      if (result.err !== "") {
        MessagePlugin.error(result.err)
        loading.value = false
      }else {
        handleResult(result)
      }
  })
}

function handleResult(result){
  columns.value = {}
  tabs.value = []
  tableData.value = {}
  let idx = 0
  for (let tab in result.data) {
    if (idx === 0){
      activeTab.value = tab
      idx++
    }
    tabs.value.push({
      value: tab, label: tab
    })
    let inf = result.data[tab][0]
    tableData.value[tab] = {}
    columns.value[tab] = []
    tableData.value[tab] = result.data[tab]
    for (let key in inf) {
      columns.value[tab].push({
        colKey: key, title: key,ellipsis: true, width:"20vh"
      })
    }
  }
  loading.value = false
}

const handleClear = () => {
  tableData.value = {}
  tabs.value = []
  activeTab.value = ""
  columns.value = {}
}

const handleTabChange = (value) => {
  activeTab.value = value
}

const handleCellClick = (context) => {
  // 复制单元格内容
  ClipboardSetText(context.row[context.col.colKey]).then((ok)=>{
    if (ok){
      MessagePlugin.success('已复制到剪贴板',500)
    }
  })
}
const handleDrop = (event) => {
  OnFileDrop((x,y,paths)=>{
    if (paths.length > 0) {
      form.value.file = paths[0]
    }
  },false)
}

const onDragSort = ({ currentIndex, targetIndex, current, target, data, newData, e, sort }) => {
  console.log('交换行', currentIndex, targetIndex, current, target, data, newData, e, sort);
  if (sort === 'col') {
    columns.value[activeTab.value] = newData;
  }
};
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  height: 92vh;
}
.form-column {
  display: flex;
  flex-direction: column;
  width: 50%;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
}
.result-container {
  display: flex;
  flex-direction: column;
}
.form-item {
  width: 50vh;
}
.result-card {
  flex: 1;
  overflow-y: hidden;
}
.button {
  border-radius: 25px;
  width: 50%
}

</style>