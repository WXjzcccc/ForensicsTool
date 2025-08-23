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
        <t-form-item label="target" class="form-item">
          <t-tooltip :overlay-style="{width:'50vh'}" :content="brutePlace.target">
            <t-input v-model="form.target" :placeholder="brutePlace.target"/>
          </t-tooltip>
        </t-form-item>
        <t-form-item label="region" class="form-item">
          <t-tooltip :overlay-style="{width:'50vh'}" :content="brutePlace.region">
            <t-input v-model="form.region" :placeholder="brutePlace.region"/>
          </t-tooltip>
        </t-form-item>

        </div>
        <div class="form-column">
          <t-form-item label="mac" class="form-item">
            <t-tooltip :overlay-style="{width:'50vh'}" :content="brutePlace.mac">
              <t-input v-model="form.mac" :placeholder="brutePlace.mac"/>
            </t-tooltip>
          </t-form-item>
        <t-form-item label="length"  class="form-item">
          <t-tooltip :overlay-style="{width:'50vh'}" :content="brutePlace.length">
            <t-input v-model="form.length" :placeholder="brutePlace.length"/>
          </t-tooltip>
        </t-form-item>
        </div>
        </div>
      </t-form>
      <div style="height: 1vh"></div>
      <t-button class="button" theme="primary" @click="handleBruteForce"><template #icon><t-icon name="cpu"/></template>爆破</t-button>
      <t-button class="button" theme="danger" @click="handleCancel"><template #icon><t-icon name="stop-circle"/></template>停止</t-button>
      <t-button class="button" theme="primary" @click="handleClear"><template #icon><t-icon name="clear-formatting"/></template>清空输出</t-button>
    </t-card>
    <t-card class="result-card">
      <div class="result-container">
        <div class="result-output" v-html="resultText"/>
      </div>
    </t-card>

  </div>
</template>

<script setup>
import {ref} from 'vue'
import {MessagePlugin} from 'tdesign-vue-next'
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import {usePageDataStore} from "@/store/index.js";
import {watch} from "vue";
import {CancelCrack, CrackAirDrop, CrackWXUin, GetState} from "../../wailsjs/go/cracker/ForensicsCracker.js";
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
const options = ref([
  {label:"AirDrop手机号爆破",value:"1"},
  {label:"微信UIN爆破",value:"2"},
])

watch([form,resultText],()=>{
  store.saveBruteForceData({
    formData:form.value,
    resultData:resultText.value
  })
})

const cracking = ref(false)

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
          MessagePlugin.error("AirDrop必须提供首位各5尾的哈希值",500)
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
}

const handleCancel = () => {
  CancelCrack().then(()=>{
    cracking.value = false
    resultText.value += generateNormalTextOutput(`已手动取消`,"red")
  })
}

function handleResult(result){
  resultText.value += generateSuccessTextOutput("爆破结束",result.result)
  resultText.value += generateNormalTextOutput(`耗时：${result.time}`,"#df9b26")
  resultText.value += generateNormalTextOutput(`错误：${result.error}`,"red")
  cracking.value = false
}

function handleState(){
  cracking.value = true
  resultText.value += generateNormalTextOutput(`开始爆破`,"#ee5310")
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
.page-container {
  display: flex;
  flex-direction: column;
  height: 92vh;
}
.form-column {
  display: flex;
  flex-direction: column;
  width: 55vh;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
}
.result-card {
  flex: 1;
  overflow-y: scroll;
}
.result-output {
  text-align: left;
}
.form-item {
  width: 50vh;
}
.button {
  border-radius: 25px;
  width: 33%
}

</style>