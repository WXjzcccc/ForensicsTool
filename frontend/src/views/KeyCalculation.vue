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
              v-tooltip.top="'选择要计算的密钥类型'"
            />
          </FloatLabel>
        </div>
        
        <!-- 文本输入框一行两个 -->
        <div class="input-row">
          <div class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                id="uin"
                v-model="form.uin" 
                :placeholder="calPlace.uin"
                v-tooltip.top="calPlace.uin"
              />
              <label for="uin">微信uin</label>
            </FloatLabel>
          </div>
          <div class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                id="imei"
                v-model="form.imei" 
                :placeholder="calPlace.imei"
                v-tooltip.top="calPlace.imei"
              />
              <label for="imei">imei</label>
            </FloatLabel>
          </div>
        </div>
        
        <div class="input-row">
          <div class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                id="wxid"
                v-model="form.wxid" 
                :placeholder="calPlace.wxid"
                v-tooltip.top="calPlace.wxid"
              />
              <label for="wxid">wxid</label>
            </FloatLabel>
          </div>
          <div class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                id="token"
                v-model="form.token" 
                :placeholder="calPlace.token"
                v-tooltip.top="calPlace.token"
              />
              <label for="token">token</label>
            </FloatLabel>
          </div>
        </div>
        
        <div class="input-row">
          <div class="field half-width">
            <FloatLabel variant="on">
              <InputText 
                id="uid"
                v-model="form.uid" 
                :placeholder="calPlace.uid"
                v-tooltip.top="calPlace.uid"
              />
              <label for="uid">uid</label>
            </FloatLabel>
          </div>
        </div>
      </form>
      
      <!-- 按钮等分在同一行 -->
      <div class="button-row">
        <Button class="button equal-width" @click="handleCalculate">
          <i class="pi pi-calculator"></i>
          计算
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
import {ref} from 'vue'
import {CalMostone, CalTiktok, CalWechat, CalWechatIndex, CalWildFire} from "../../wailsjs/go/passwdCalc/PasswdCalc.js";
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import {usePageDataStore} from "@/store";
import {watch} from "vue";
import FloatLabel from 'primevue/floatlabel'
import Select from 'primevue/select'
import Empty from '@/components/Empty.vue'
const calPlace = {
  uin:"微信用户的uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值",
  imei:"微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF，可以为空",
  wxid:"数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_preferences.xml中提取login_weixin_username的值",
  token:"野火IM系应用的用户token，shared_prefs/config.xml的token的值",
  uid:"默往（通常在shared_prefs/im.xml中的userId的值）、抖音（数据库文件名中的id）计算密钥需要的内容、QQ（msf_mmkv_file中QQ号对应的uid）"
}
const store = usePageDataStore()
const form = ref(store.keyCalculationStore?.formData || {
  selected:"",
  uin:"",
  imei:"",
  wxid:"",
  token:"",
  uid:"",
})
const resultText = ref(store.keyCalculationStore?.resultData || "")
const options = ref([
  {label:"微信的EnMicroMsg.db",value:"1"},
  {label:"微信的FTS5IndexMicroMsg_encrypt.db",value:"2"},
  {label:"野火IM系应用的data",value:"3"},
  {label:"默往APP的msg.db",value:"4"},
  {label:"抖音的聊天数据库",value:"5"},
])

watch([form,resultText],()=>{
  store.saveKeyCalculationData({
    formData:form.value,
    resultData:resultText.value
  })
})

const handleCalculate = () => {
  var uin = form.value.uin;
  var imei = form.value.imei;
  var wxid = form.value.wxid;
  var token = form.value.token;
  var uid = form.value.uid;
  switch (form.value.selected) {
    case "1":{
      CalWechat(uin,imei).then((result)=>{
        resultText.value += generateSuccessTextOutput("成功获取到安卓微信EnMicroMsg.db数据库密钥",result)
      })
      break;
    }
    case "2":{
      CalWechatIndex(uin,wxid,imei).then((result)=>{
        if (result !== ""){
          resultText.value += generateSuccessTextOutput("成功获取到安卓微信FTS5IndexMicroMsg_encrypt.db数据库密钥",result)
        }else{
          resultText.value += generateNormalTextOutput("计算失败，请确认参数","red")
        }
      })
      break;
    }
    case "3":{
      CalWildFire(token).then((result)=>{
        resultText.value += generateSuccessTextOutput("成功获取到野火IM数据库密钥",result[0])
        resultText.value += generateNormalTextOutput(result[1],"green")
      })
      break;
    }
    case "4":{
      CalMostone(uid).then((result)=>{
        resultText.value += generateSuccessTextOutput("成功获取到默往msg.db数据库密钥",result[0])
        resultText.value += generateNormalTextOutput(result[1],"green")
      })
      break;
    }
    case "5":{
      CalTiktok(uid).then((result)=>{
        resultText.value += generateSuccessTextOutput("成功获取到抖音聊天数据库密钥",result[0])
        resultText.value += generateNormalTextOutput(result[1],"green")
      })
      break;
    }
  }
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