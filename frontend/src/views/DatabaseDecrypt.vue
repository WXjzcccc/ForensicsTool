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
            <t-form-item label="file" class="form-item">
              <t-input v-model="form.file" placeholder="请拖入文件或目录" @drop.prevent="handleDrop"
                       @dragover.prevent/>
            </t-form-item>
          </div>
          <div class="form-column">
            <t-form-item label="password"  class="form-item">
              <t-input v-model="form.password" placeholder="解密密码"/>
            </t-form-item>
          </div>

        </div>
      </t-form>
      <div style="height: 1vh"></div>
      <t-button class="button" theme="primary" @click="handleDecrypt"><template #icon><t-icon name="lock-off"/></template>解密</t-button>
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
import {
  DecryptAMapDB,
  DecryptDingTalkDB,
  DecryptEnMicroMsg,
  DecryptFTSIndexDB, DecryptNtqqDB, DecryptSQLCipher3DB, DecryptSQLCipher4DB, DecryptSystemDataSQLite, DecryptWCDB
} from "../../wailsjs/go/database/DecryptDatabase.js";
import {generateNormalTextOutput, generateSuccessTextOutput} from "@/utils.js";
import {OnFileDrop} from "../../wailsjs/runtime/runtime.js";
import {usePageDataStore} from "@/store/index.js";
import {watch} from "vue";
const store = usePageDataStore()
const form = ref(store.databaseDecryptStore?.formData || {
  selected:"",
  password:"",
  file:"",
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
}

const handleDecrypt = () => {
  var file = form.value.file;
  var password = form.value.password;
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
  }
  if (form.value.selected === "3") {
    func(file).then((result)=>{
      if (result.err !== "") {
        resultText.value += generateNormalTextOutput(result.err,"red")
      }else{
        resultText.value += generateSuccessTextOutput("解密成功，解密后的数据库已保存至",result.save_path)
      }
    })
  }else{
    func(file,password).then((result)=>{
      if (result.err !== "") {
        resultText.value += generateNormalTextOutput(result.err,"red")
      }else{
        resultText.value += generateSuccessTextOutput("解密成功，解密后的数据库已保存至",result.save_path)
        if (form.value.selected === "1") {
          resultText.value += generateSuccessTextOutput("成功提取微信ID：",result.wxid)
        }
      }
    })
  }
}

const handleDrop = (event) => {
  OnFileDrop((x,y,paths)=>{
    if (paths.length > 0) {
      form.value.file = paths[0]
    }
  },false)
}
watch(resultText, () => {
  const card = document.querySelector('.result-card');
  if (card) {
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
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
}
.result-container {
  display: flex;
  flex-direction: column;
}
.result-card {
  flex: 1;
  overflow-y: scroll;
}
.result-output {
  text-align: left;
}
.button {
  border-radius: 25px;
  width: 50%
}
.form-item {
  width: 50vh;
}
</style>