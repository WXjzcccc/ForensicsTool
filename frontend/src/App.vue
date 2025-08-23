<template>
  <t-layout>
    <t-layout class="full-height-layout">
      <SideBar />
      <MainContent class="content-layout"/>
      <t-sticky-tool
          :style="{ position: '', overflow: 'hidden' }"
          placement="left-bottom"
          shape="round"
          :offset="[-30,30]"
          @click="handleFloatButtonClick"
      >
        <t-sticky-item label="Tips">
          <template #icon><t-icon name="sticky-note" /></template>
        </t-sticky-item>
      </t-sticky-tool>
      <t-dialog v-model:visible="visible" header="Tips" class="dialog"
                :cancelBtn="null"
                :confirmBtn="null"
                :preventScrollThrough="false"
      >
        <t-list size="small" :split="true" class="tip-container">
          <t-list-item>1.小红书，密码为xhsdev，选择sqlcipher3直接解密</t-list-item>
          <t-list-item>2.微信的imei，现在可以通过files/KeyInfo.bin获取了，需要解密文件，算法为RC4，密钥为_wEcHAT_</t-list-item>
          <t-list-item>3.MosGram(泡泡)，密码为cust_id的md5值，在sp目录的account_config.xml文件中，使用SQLCipher4参数解</t-list-item>
          <t-list-item>4.悟空IM系列的聊天数据库，数据库名为wk_用户ID.db，解密密码即为用户ID，使用SQLCipher4参数解密</t-list-item>
        </t-list>
      </t-dialog>
    </t-layout>
  </t-layout>
</template>
<script setup>
import MainContent from "@/components/MainContent.vue";
import SideBar from "@/components/SideBar.vue";
import {ref} from "vue";
const visible = ref(false)
const handleFloatButtonClick = ()=>{
  visible.value = true
}
</script>
<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100vh;
}

.full-height-layout {
  height: 100vh;
}

.content-layout {
  flex: 1;
  min-height: 0; /* 关键属性，解决flex容器滚动问题 */
}

/* 响应侧边栏折叠状态 */
.t-layout--collapsed .floating-button-container {
  left: 64px; /* 侧边栏折叠后的宽度 */
}
.dialog{
  text-align: left;
}

.tip-container{
  overflow-y: scroll;
  height: 20vh;
}
</style>
