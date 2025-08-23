<template>
  <div class="page-container">
    <t-card class="form-card">
      <t-form layout="inline" label-width="calc(2em + 4vw)"  label-align="left">
        <div class="form-grid">
          <div class="form-column">
            <t-form-item label="原始时区" style="width: 35vw">
              <t-select v-model="form.selectedOld" placeholder="请选择时区">
                <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label"></t-option>
              </t-select>
            </t-form-item>
          </div>
          <div class="form-column">
            <t-form-item label="目标时区" style="width: 35vw">
              <t-select v-model="form.selectedNew" placeholder="请选择时区">
                <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label"></t-option>
              </t-select>
            </t-form-item>
          </div>
        </div>
        <t-form-item label="时间戳" style="width: 35vw">
          <t-input v-model="form.ts" placeholder="请输入时间戳"></t-input>
        </t-form-item>
      </t-form>
      <div style="height: 1vh"></div>
      <t-space>
        <t-button class="button" theme="primary" @click="handleTrans"><template #icon><t-icon name="history"/></template>转换</t-button>
        <t-button class="button" theme="primary" @click="handleClear"><template #icon><t-icon name="clear-formatting"/></template>清空输出</t-button>
      </t-space>
    </t-card>
    <div style="height: 1vh"></div>
    <t-card class="result-card">
      <t-empty v-if="resultText===''" class="empty"/>
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
import {ParseTimeStamp} from "../../wailsjs/go/timestamp/TimeStampParser.js";
const store = usePageDataStore()
const form = ref(store.timestampData?.formData || {
  selectedOld:"UTC",
  selectedNew:"Asia/Shanghai",
  ts:""
})
const resultText = ref(store.timestampData?.resultData || "")
const timezones = Intl.supportedValuesOf('timeZone')
timezones.push('UTC')
const options = ref(timezones.map((item, index) => ({
  label: item,
  value: index + 1
})))

watch([form,resultText],()=>{
  store.saveTimestampData({
    formData:form.value,
    resultData:resultText.value
  })
})


const handleTrans = () => {
  let ori = form.value.selectedOld;
  let target = form.value.selectedNew;
  let ts = form.value.ts;
  ParseTimeStamp(ts,ori,target).then((result)=>{
    resultText.value += result.replaceAll("\n","<br>")
  })
}

const handleClear = () => {
  resultText.value = ""
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
.empty{
  margin-top:25vh
}
.form-card {
  border-color: blue;
  border-width: 3px;
}
.result-card {
  flex: 1;
  overflow-y: hidden;
  border-color: blue;
  border-width: 3px;
}
</style>