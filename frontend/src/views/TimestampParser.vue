<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <div class="input-row">
          <div class="field full-width">
            <FloatLabel class="field full-width" variant="over">
              <label for="oldTimezone">原始时区</label>
              <Select 
                id="oldTimezone"
                v-model="form.selectedOld" 
                :options="options" 
                optionLabel="label" 
                optionValue="value"
                placeholder="请选择时区"
                v-tooltip.top="'选择原始时间戳所在的时区'"
              />
              </FloatLabel>
          </div>
          <div class="field full-width">
            <FloatLabel class="field full-width" variant="over">
              <label for="newTimezone">目标时区</label>
              <Select 
                id="newTimezone"
                v-model="form.selectedNew" 
                :options="options" 
                optionLabel="label" 
                optionValue="value"
                placeholder="请选择时区"
                v-tooltip.top="'选择要转换到的目标时区'"
              />
              </FloatLabel>
          </div>
        </div>
        <div class="field full-width">
          <FloatLabel variant="on">
            <InputText 
              id="timestamp"
              v-model="form.ts" 
              aria-autocomplete="none"
              placeholder="请输入时间戳"
              v-tooltip.top="'输入要转换的时间戳，支持Unix时间戳'"
            />
            <label for="timestamp">时间戳</label>
          </FloatLabel>
        </div>
      </form>
        <div class="button-row">
          <Button class="button equal-width" @click="handleTrans">
            <i class="pi pi-history"></i>
            转换
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
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { usePageDataStore } from "@/store"
import { watch } from "vue"
import { ParseTimeStamp } from "../../wailsjs/go/timestamp/TimeStampParser.js"
import FloatLabel from 'primevue/floatlabel'
import Select from 'primevue/select'
import Empty from '@/components/Empty.vue'

const store = usePageDataStore()
const form = ref(store.timestampStore?.formData || {
  selectedOld: "UTC",
  selectedNew: "Asia/Shanghai",
  ts: ""
})
const resultText = ref(store.timestampStore?.resultData || "")
const timezones = Intl.supportedValuesOf('timeZone')
timezones.push('UTC')
const options = ref(timezones.map((item, index) => ({
  label: item,
  value: item
})))

watch([form, resultText], () => {
  store.saveTimestampData({
    formData: form.value,
    resultData: resultText.value
  })
})

const handleTrans = () => {
  let ori = form.value.selectedOld;
  let target = form.value.selectedNew;
  let ts = form.value.ts;
  ParseTimeStamp(ts, ori, target).then((result) => {
    resultText.value += result.replaceAll("\n", "<br>")
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
</style>