<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <div class="field full-width">
          <FloatLabel class="field full-width" variant="over">
              <label for="task">选择任务</label>
          <Select 
            id="task"
            v-model="form.selected" 
            :options="options" 
            optionLabel="label" 
            optionValue="value"
            placeholder="请选择任务"
            v-tooltip.top="'选择要执行的数据提取任务类型'"
          />
          </FloatLabel>
        </div>
        
        <div v-for="(fieldGroup, index) in inputFields" :key="index" class="input-row">
          <div v-for="field in fieldGroup" :key="field.name" class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                v-if="field.type === 'input'"
                v-tooltip.top="brutePlace[field.name]"
                aria-autocomplete="none"
                :id="field.name"
                v-model="form[field.name]" 
                :placeholder="brutePlace[field.name]"
              />
              <label :for="field.name">{{ field.label }}</label>
            </FloatLabel>
          </div>
        </div>
      </form>
      <div class="button-row">
        <Button class="button equal-width" @click="handleBruteForce" :disabled="cracking">
          <i class="pi pi-cog"></i>
          爆破
        </Button>
        <Button class="button equal-width" severity="danger" @click="handleCancel" :disabled="!cracking">
          <i class="pi pi-stop-circle"></i>
          停止
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
      <div v-if="resultText === ''" class="empty">
        <Empty />
      </div>
        <div class="result-output" v-html="resultText"/>
      </template>
    </Card>
    <Toast position="bottom-right"/>
    <ProgressBar v-if="cracking" mode="indeterminate" style="height: 1vh" />
  </div>
</template>

<script setup>
import {ref, watch, onMounted, computed} from 'vue'
import {useToast} from 'primevue/usetoast'
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import {usePageDataStore} from "@/store";
import {CancelCrack, CrackAirDrop, CrackWXUin, GetState} from "../../wailsjs/go/cracker/ForensicsCracker.js";
import { Select } from 'primevue';
import Empty from '@/components/Empty.vue';

const toast = useToast()
const store = usePageDataStore()
const brutePlace = {
  target:"目标列表，以,分隔",
  region:"AirDrop要爆破的区号，如86、85、1",
  mac:"AirDrop要爆破的号段，以,进行分隔，如139,138",
  length:"AirDrop要爆破的手机号长度（除去区号和号段）",
}
// 任务配置对象，定义每个任务需要的输入字段
const taskConfigs = {
  "1": {
    fields: [
      { name: "target", label: "爆破目标", type: "input" },
      { name: "region", label: "区号", type: "input" },
      { name: "mac", label: "号段", type: "input" },
      { name: "length", label: "长度", type: "input" }
    ]
  },
  "2": {
    fields: [
      { name: "target", label: "爆破目标", type: "input" }
    ]
  }
}

const form = ref(store.bruteForceStore?.formData || {
  selected:"",
  target:"",
  mac:"",
  region:"",
  length:"",
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
const resultText = ref(store.bruteForceStore?.resultData || "")
const options = ref([
  {label:"AirDrop手机号爆破",value:"1"},
  {label:"微信UIN爆破",value:"2"},
])
const cracking = ref(store.getBruteForceCrackingState())

// 组件挂载时，如果爆破状态为true，则重新启动状态检查
onMounted(() => {
  if (cracking.value) {
    // 重新启动状态检查
    const timer = setInterval(() => {
      if (cracking.value) {
        GetState().then((result)=>{
          resultText.value += generateNormalTextOutput(result.replace("\n","<br>"),"#1937e5")
        })
      } else {
        clearInterval(timer);
      }
    }, 5000);
  }
})

watch([form,resultText],()=>{
  store.saveBruteForceData({
    formData:form.value,
    resultData:resultText.value,
    cracking: cracking.value
  })
})

// 监听cracking状态变化，更新store
watch(cracking, (newValue) => {
  store.updateBruteForceCrackingState(newValue)
})

const handleBruteForce = () => {
  const currentConfig = currentTaskConfig.value
  const params = {}
  
  // 根据当前任务配置构建参数
  currentConfig.fields.forEach(field => {
    params[field.name] = form.value[field.name]
  })
  
  switch (form.value.selected) {
    case "1":{
      let head = ""
      let tail = ""
      let macs = []
      if (params.target !== "" && params.target.indexOf(",") !== -1){
        let tmp = params.target.split(",")
        if(tmp.length !== 2){
          toast.add({ severity: 'error', summary: '错误', detail: 'AirDrop必须提供首位各5尾的哈希值', life: 5000 })
          return
        }
        head = tmp[0]
        tail = tmp[1]
      }

      if (params.mac !== ""){
        if(params.mac.indexOf(",") !== -1){
          macs = params.mac.split(",")
        }else{
          macs.push(params.mac)
        }
      }
      handleState()
      CrackAirDrop(head,tail,params.region,macs,parseInt(params.length)).then((result)=>{
        handleResult(result)
      })
      break;
    }
    case "2":{
      let para = []
      if (params.target !== ""){
        if (params.target.indexOf(",") !== -1){
          para.push(params.target)
        }else{
          para = params.target.split(",")
        }
      }
      handleState()
      CrackWXUin(para).then((result)=>{
        handleResult(result)
      })
      break;
    }
  }
}

const handleClear = () => {
  resultText.value = ""
  toast.add({ severity: 'info', summary: '信息', detail: '已清空输出', life: 3000 })
}

const handleCancel = () => {
  CancelCrack().then(()=>{
    cracking.value = false
    resultText.value += generateNormalTextOutput(`已手动取消`,"red")
    toast.add({ severity: 'warn', summary: '警告', detail: '已手动取消', life: 3000 })
  })
}

function handleResult(result){
  resultText.value += generateSuccessTextOutput("爆破结束",result.result)
  resultText.value += generateNormalTextOutput(`耗时：${result.time}`,"#df9b26")
  resultText.value += generateNormalTextOutput(`错误：${result.error}`,"red")
  cracking.value = false
  toast.add({ severity: 'success', summary: '成功', detail: '爆破完成', life: 3000 })
}

function handleState(){
  cracking.value = true
  resultText.value += generateNormalTextOutput(`开始爆破`,"#ee5310")
  toast.add({ severity: 'info', summary: '信息', detail: '开始爆破', life: 3000 })
  const timer = setInterval(() => {
    if (cracking.value) {
      GetState().then((result)=>{
        resultText.value += generateNormalTextOutput(result.replace("\n","<br>"),"#1937e5")
      })
    } else {
      clearInterval(timer);
    }
  }, 5000);
}

watch(resultText, () => {
  const card = document.querySelector('.result-card');
  if (card) {
    void card.offsetHeight;
    card.scrollTop = card.scrollHeight;
  }
});

</script>

<style scoped>
</style>