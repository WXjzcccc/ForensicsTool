<template>
  <div class="page-container">
    <Card class="form-card">
      <template #content>
      <form class="form-layout">
        <div v-for="(fieldGroup, index) in inputFields" :key="index" class="input-row">
          <div v-for="field in fieldGroup" :key="field.name" class="field half-width">
            <FloatLabel v-if="field.type === 'select'" class="field full-width" variant="over">
              <label :for="field.name">{{ field.label }}</label>
              <Select 
                :id="field.name"
                v-model="form[field.name]" 
                :options="options" 
                optionLabel="label" 
                optionValue="value"
                :placeholder="field.name === 'selectedOld' ? '请选择时区' : '请选择时区'"
                v-tooltip.top="field.name === 'selectedOld' ? '选择原始时间戳所在的时区' : '选择要转换到的目标时区'"
              />
            </FloatLabel>
            <FloatLabel v-else-if="field.type === 'input'" variant="on">
              <InputText 
                :id="field.name"
                v-model="form[field.name]" 
                aria-autocomplete="none"
                placeholder="请输入时间戳"
                v-tooltip.top="'输入要转换的时间戳，支持Unix时间戳'"
              />
              <label :for="field.name">{{ field.label }}</label>
            </FloatLabel>
          </div>
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
import { ref, computed } from 'vue'
import { usePageDataStore } from "@/store"
import { watch } from "vue"
import { ParseTimeStamp } from "../../wailsjs/go/timestamp/TimeStampParser.js"
import FloatLabel from 'primevue/floatlabel'
import Select from 'primevue/select'
import Empty from '@/components/Empty.vue'

// 任务配置对象，定义每个任务需要的输入字段
const taskConfigs = {
  "1": {
    fields: [
      { name: "selectedOld", label: "原始时区", type: "select" },
      { name: "selectedNew", label: "目标时区", type: "select" },
      { name: "ts", label: "时间戳", type: "input" }
    ]
  }
}

const store = usePageDataStore()
const form = ref(store.timestampStore?.formData || {
  selectedOld: "UTC",
  selectedNew: "Asia/Shanghai",
  ts: ""
})
// 计算当前任务配置
const currentTaskConfig = computed(() => {
  return taskConfigs["1"] || { fields: [] }
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
  const currentConfig = taskConfigs["1"]
  if (!currentConfig) {
    return
  }

  // 动态构建参数对象
  const params = {}
  currentConfig.fields.forEach(field => {
    params[field.name] = form.value[field.name]
  })

  // 从动态参数中获取selectedOld、selectedNew和ts
  const { selectedOld, selectedNew, ts } = params

  ParseTimeStamp(ts, selectedOld, selectedNew).then((result) => {
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