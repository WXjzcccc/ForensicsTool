<template>
  <div class="page-container">
    <t-card class="form-card">
      <t-form layout="inline" label-width="calc(2em + 7vh)"  label-align="left">
        <t-form-item label="选择任务" style="width: 100vh">
          <t-select v-model="form.selected" placeholder="请选择任务">
            <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label"></t-option>
          </t-select>
        </t-form-item>
        <div class="form-grid">
          <div class="form-column">
            <t-form-item class="form-item" label="file">
              <t-input v-model="form.file" placeholder="请拖入文件或目录" @drop.prevent="handleDrop"
                       @dragover.prevent/>
            </t-form-item>
          </div>
          <div class="form-column">
            <t-form-item class="form-item" label="password">
              <t-input v-model="form.password" placeholder="解密密码"/>
            </t-form-item>
          </div>
        </div>
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
              max-height="60vh"
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
const store = usePageDataStore()
const form = ref(store.dataExtractionStore?.formData || {
  selected: '',
  file: '',
  password: '',
})
const tableData = ref(store.dataExtractionStore?.tableData || {
})
const tabs = ref(store.dataExtractionStore?.tabsData || [
])
const activeTab = ref(store.dataExtractionStore?.tabData || '')
const columns = ref(store.dataExtractionStore?.columnsData || {
})
watch([form,tableData,tabs,activeTab],()=>{
  store.saveDataExtractionData({
    formData: form.value,
    tableData: tableData.value,
    tabsData: tabs.value,
    tabData: activeTab.value,
    columnsData: columns.value
  })
})

const options = ref([
  {label:"Navicat连接信息提取，指定文件为用户注册表文件NTUSER.DAT",value:"1"},
  {label:"MobaXterm连接信息解密，指定文件为用户注册表文件NTUSER.DAT并给出主密码",value:"2"},
  {label:"Dbeaver连接信息解密，指定文件为data-sources.json的父目录",value:"3"},
  {label:"FinalShell连接信息解密，指定文件为conn文件夹",value:"4"},
  {label:"XShell、XFtp连接信息解密，指定文件为session文件夹并给出密码（用户名+sid）",value:"5"},
  {label:"Hawk2.xml数据解密，给出密码（crypto.KEY_128(256).xml中的base64值）",value:"6"},
  {label:"MetaMask解析，指定文件为persis-root",value:"7"},
])

const loading = ref(false)

const handleExtract = () => {
  loading.value = true
  let file = form.value.file
  let password = form.value.password
  let func = null
  switch (form.value.selected) {
    case "1":
      func = ExtractNavicat
      break;
    case "2":
      func = ExtractMobaXterm
      break;
    case "3":
      func = ExtractDbeaver
      break;
    case "4":
      func = ExtractFinalShell
      break;
    case "5":
      func = ExtractXShell
      break;
    case "6":
      func = ExtractHawk2
      break;
    case "7":
      func = ExtractMetaMask
      break;
  }
  if (form.value.selected === "2" ||
      form.value.selected === "5" ||
  form.value.selected === "6"){
    func(file, password).then((result) => {
      if (result.err !== "") {
        MessagePlugin.error(result.err)
      }else {
        handleResult(result)
      }
    })
  }else {
    func(file).then((result) => {
      if (result.err !== "") {
        MessagePlugin.error(result.err)
      }else {
        handleResult(result)
      }
    })
  }
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