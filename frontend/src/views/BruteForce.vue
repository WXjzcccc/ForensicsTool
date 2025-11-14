<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <!-- 下拉选择框独占一行 -->
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
        
        <!-- 文本输入框一行两个 -->
        <div class="input-row">
          <div class="field half-width">
              <FloatLabel variant="on">
                  <InputText v-tooltip.top=brutePlace.target
                    id="target"
                    v-model="form.target"
                    :placeholder="brutePlace.target"
                  />
                  <label for="on_label">爆破目标</label>
              </FloatLabel>
          </div>
          <div class="field half-width">
            <FloatLabel variant="on">
            <label for="region">区号</label>
              <InputText 
                v-tooltip.top=brutePlace.region
                id="region"
                v-model="form.region" 
                :placeholder="brutePlace.region"
              />
              </FloatLabel>
          </div>
        </div>
        
        <div class="input-row">
          <div class="field half-width">
            <FloatLabel variant="on">
            <label for="mac">号段</label>
              <InputText 
                v-tooltip.top=brutePlace.mac
                id="mac"
                v-model="form.mac" 
                :placeholder="brutePlace.mac"
              />
              </FloatLabel>
          </div>
          <div class="field half-width">
            <FloatLabel variant="on">
            <label for="length">长度</label>
              <InputText 
               v-tooltip.top=brutePlace.length
                id="length"
                v-model="form.length" 
                :placeholder="brutePlace.length"
              />
              </FloatLabel>
          </div>
        </div>
      </form>
      
      <!-- 按钮等分在同一行 -->
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
  </div>
</template>

<script setup>
import {ref, watch, onMounted} from 'vue'
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
const form = ref(store.bruteForceStore?.formData || {
  selected:"",
  target:"",
  mac:"",
  region:"",
  length:"",
})
const resultText = ref(store.bruteForceStore?.resultData || "")
const cracking = ref(store.getBruteForceCrackingState())
const options = ref([
  {label:"AirDrop手机号爆破",value:"1"},
  {label:"微信UIN爆破",value:"2"},
])

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
  var region = form.value.region;
  var target = form.value.target;
  var mac = form.value.mac;
  var length = form.value.length;
  switch (form.value.selected) {
    case "1":{
      let head = ""
      let tail = ""
      let macs = []
      if (target !== "" && target.indexOf(",") !== -1){
        let tmp = target.split(",")
        if(tmp.length !== 2){
          toast.add({ severity: 'error', summary: '错误', detail: 'AirDrop必须提供首位各5尾的哈希值', life: 5000 })
          return
        }
        head = tmp[0]
        tail = tmp[1]
      }

      if (mac !== ""){
        if(mac.indexOf(",") !== -1){
          macs = mac.split(",")
        }else{
          macs.push(mac)
        }
      }
      handleState()
      CrackAirDrop(head,tail,region,macs,parseInt(length)).then((result)=>{
        handleResult(result)
      })
      break;
    }
    case "2":{
      let para = []
      if (target !== ""){
        if (target.indexOf(",") !== -1){
          para.push(target)
        }else{
          para = target.split(",")
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