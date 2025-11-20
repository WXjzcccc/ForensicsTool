<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <div class="field full-width">
          <FloatLabel class="field full-width" variant="over">
              <label for="selected">选择任务</label>
            <Select 
              id="selected"
              v-model="form.selected" 
              :options="options" 
              optionLabel="label" 
              optionValue="value"
              placeholder="请选择任务"
              v-tooltip.top="'选择要解密的数据库类型'"
            />
          </FloatLabel>
        </div>
        <div v-for="(fieldGroup, index) in inputFields" :key="index" class="input-row">
          <div v-for="field in fieldGroup" :key="field.name" class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                v-if="field.type === 'input'"
                :id="field.name"
                v-model="form[field.name]" 
                :placeholder="field.name === 'file' ? '请拖入文件或目录' : '解密密码'"
                aria-autocomplete="none"
                @drop.prevent="field.name === 'file' ? null : null"
                v-tooltip.top="field.name === 'file' ? '拖入要解密的数据库文件路径' : '输入解密数据库所需的密码（某些数据库需要）'"
              />
              <label :for="field.name">{{ field.label }}</label>
            </FloatLabel>
          </div>
        </div>
      </form>
      <div class="button-row">
        <Button class="button equal-width" @click="handleDecrypt">
          <i class="pi pi-lock-open"></i>
          解密
        </Button>
        <Button class="button equal-width" severity="secondary" @click="handleClear">
          <i class="pi pi-trash"></i>
          清空输出
        </Button>
      </div>
      </template>
    </Card>
    <Card class="result-card">
      <template #content>
      <div class="empty" v-if="resultText===''">
        <Empty />
      </div>
        <div class="result-output" v-html="resultText"/>
      </template>
    </Card>
    <Toast position="bottom-right"/>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import {
  DecryptAMapDB,
  DecryptDingTalkDB,
  DecryptEnMicroMsg,
  DecryptFTSIndexDB, DecryptNtqqDB, DecryptSQLCipher3DB, DecryptSQLCipher4DB, DecryptSystemDataSQLite, DecryptWCDB
} from "../../wailsjs/go/database/DecryptDatabase.js";
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime.js";
import {usePageDataStore} from "@/store";
import FloatLabel from 'primevue/floatlabel'
import Select from 'primevue/select'
import Empty from '@/components/Empty.vue'

// 任务配置对象，定义每个任务需要的输入字段
const taskConfigs = {
  "1": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "2": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "3": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" }
    ]
  },
  "4": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "5": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "6": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "7": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "8": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  },
  "9": {
    fields: [
      { name: "file", label: "文件/目录", type: "input" },
      { name: "password", label: "密码", type: "input" }
    ]
  }
}

const toast = useToast()
const store = usePageDataStore()
const form = ref(store.databaseDecryptStore?.formData || {
  selected:"",
  password:"",
  file:"",
})
// 计算当前任务配置
const currentTaskConfig = computed(() => {
  return taskConfigs[form.value.selected] || { fields: [] }
})

// 计算输入字段，每两个一组
const inputFields = computed(() => {
  const fields = currentTaskConfig.value.fields
  const groups = []
  for (let i = 0; i < fields.length; i += 2) {
    groups.push(fields.slice(i, i + 2))
  }
  return groups
})

const resultText = ref(store.databaseDecryptStore?.resultData || "")
const options = ref([
  {label:"微信的EnMicroMsg.db",value:"1"},
  {label:"微信的FTS5IndexMicroMsg_encrypt.db",value:"2"},
  {label:"高德的girf_sync.db",value:"3"},
  {label:"钉钉的数据库，密码如HUAWEI P40/armeabi-v7a/P40/qcom/HUAWEIP40",value:"4"},
  {label:"SQLCipher4加密的数据库",value:"5"},
  {label:"SQLCipher3加密的数据库",value:"6"},
  {label:"wcdb加密的数据库",value:"7"},
  {label:"ntqq数据库解密",value:"8"},
  {label:"System.Data.SQLite库加密的数据库",value:"9"},
])

watch([form,resultText],()=>{
  store.saveDatabaseDecryptData({
    formData:form.value,
    resultData:resultText.value
  })
})

const handleClear = () => {
  resultText.value = ""
  toast.add({ severity: 'info', summary: '提示', detail: '已清空输出', life: 3000 })
}

const handleDecrypt = () => {
  const currentConfig = taskConfigs[form.value.selected]
  if (!currentConfig) {
    toast.add({ severity: 'warn', summary: '警告', detail: '请选择一个解密任务', life: 3000 })
    return
  }

  // 动态构建参数对象
  const params = {}
  currentConfig.fields.forEach(field => {
    params[field.name] = form.value[field.name]
  })

  // 从动态参数中获取file和password
  const { file, password } = params

  if (!file) {
    toast.add({ severity: 'warn', summary: '警告', detail: '请输入文件/目录', life: 3000 })
    return
  }

  var func
  
  switch (form.value.selected) {
    case "1":{
      func = DecryptEnMicroMsg;
      break;
    }
    case "2":{
      func = DecryptFTSIndexDB;
      break;
    }
    case "3":{
      func = DecryptAMapDB;
      break;
    }
    case "4":{
      func = DecryptDingTalkDB;
      break;
    }
    case "5":{
      func = DecryptSQLCipher4DB;
      break;
    }
    case "6":{
      func = DecryptSQLCipher3DB;
      break;
    }
    case "7":{
      func = DecryptWCDB;
      break;
    }
    case "8":{
      func = DecryptNtqqDB;
      break;
    }
    case "9":{
      func = DecryptSystemDataSQLite;
      break;
    }
    default:
      toast.add({ severity: 'warn', summary: '警告', detail: '请选择一个解密任务', life: 3000 })
      return
  }
  
  toast.add({ severity: 'info', summary: '提示', detail: '开始解密...', life: 3000 })
  
  if (form.value.selected === "3") {
    func(file).then((result)=>{
      if (result.err !== "") {
        resultText.value += generateNormalTextOutput(result.err,"red")
        toast.add({ severity: 'error', summary: '错误', detail: '解密失败', life: 3000 })
      }else{
        resultText.value += generateSuccessTextOutput("解密成功，解密后的数据库已保存至",result.save_path)
        toast.add({ severity: 'success', summary: '成功', detail: '解密完成', life: 3000 })
      }
    })
  }else{
    func(file,password).then((result)=>{
      if (result.err !== "") {
        resultText.value += generateNormalTextOutput(result.err,"red")
        toast.add({ severity: 'error', summary: '错误', detail: '解密失败', life: 3000 })
      }else{
        resultText.value += generateSuccessTextOutput("解密成功，解密后的数据库已保存至",result.save_path)
        if (form.value.selected === "1") {
          resultText.value += generateSuccessTextOutput("成功提取微信ID：",result.wxid)
        }
        toast.add({ severity: 'success', summary: '成功', detail: '解密完成', life: 3000 })
      }
    })
  }
}

onMounted(() => {
  OnFileDrop((x, y, paths) => {
    if (paths.length > 0) {
      form.value.file = paths[0]
      toast.add({ severity: 'success', summary: '成功', detail: '文件已添加', life: 3000 })
    }
  }, false)
})

watch(resultText, () => {
  const card = document.querySelector('.result-card');
  if (card) {
    card.scrollTop = card.scrollHeight;
  }
});

// 组件卸载时清理文件拖放监听器
onUnmounted(() => {
  OnFileDropOff()
})
</script>

<style scoped>
</style>