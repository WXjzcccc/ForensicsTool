<template>
  <div class="page-container">
    <t-card class="form-card">
      <t-form layout="inline" label-width="calc(2em + 4vw)"  label-align="left">
        <t-form-item label="选择任务" style="width: 100vw">
          <t-select v-model="form.selected" placeholder="请选择任务">
            <t-option v-for="item in options" :key="item.value" :value="item.value" :label="item.label"></t-option>
          </t-select>
        </t-form-item>
        <div class="form-grid">
        <div class="form-column">
        <t-form-item label="微信uin" class="form-item">
          <t-tooltip :overlay-style="{width:'35vw'}" :content="calPlace.uin">
            <t-input v-model="form.uin" :placeholder="calPlace.uin"/>
          </t-tooltip>
        </t-form-item>

        <t-form-item label="imei" class="form-item">
          <t-tooltip :overlay-style="{width:'35vw'}" :content="calPlace.imei">
            <t-input v-model="form.imei" :placeholder="calPlace.imei"/>
          </t-tooltip>
        </t-form-item>

        <t-form-item label="wxid" class="form-item">
          <t-tooltip :overlay-style="{width:'35vw'}" :content="calPlace.wxid">
            <t-input v-model="form.wxid" :placeholder="calPlace.wxid"/>
          </t-tooltip>
        </t-form-item>
        </div>
        <div class="form-column">
        <t-form-item label="token"  class="form-item">
          <t-tooltip :overlay-style="{width:'35vw'}" :content="calPlace.token">
            <t-input v-model="form.token" :placeholder="calPlace.token"/>
          </t-tooltip>
        </t-form-item>

        <t-form-item label="uid" class="form-item">
          <t-tooltip :overlay-style="{width:'35vw'}" :content="calPlace.uid">
            <t-input v-model="form.uid" :placeholder="calPlace.uid"/>
          </t-tooltip>
        </t-form-item>
        </div>
        </div>
      </t-form>
      <div style="height: 1vh"></div>
      <t-space>
        <t-button class="button" theme="primary" @click="handleCalculate"><template #icon><t-icon name="calculator"/></template>计算</t-button>
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
import {CalMostone, CalTiktok, CalWechat, CalWechatIndex, CalWildFire} from "../../wailsjs/go/passwdCalc/PasswdCalc.js";
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import {usePageDataStore} from "@/store/index.js";
import {watch} from "vue";
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
<!-- 如果只在style.css中定义，编译后不生效，暂时不知道原因 -->
<style scoped>
.empty{
  margin-top:20vh
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