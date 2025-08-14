<template>
  <t-layout>
    <t-layout class="full-height-layout">
      <SideBar />
      <MainContent class="content-layout"/>
      <div class="floating-button-container">
        <t-button
            shape="circle"
            theme="primary"
            class="floating-button"
            @click="handleFloatButtonClick"
        >
          <template #icon>
            <t-icon name="sticky-note" />
          </template>
        </t-button>
      </div>
    </t-layout>
    <t-dialog v-model:visible="visible" header="Tips" class="dialog"
      :cancelBtn="null"
      :confirmBtn="null"
    >
      <div class="tip-container">
        <p>1.小红书，密码为xhsdev，选择sqlcipher3直接解密</p>
        <p>2.微信的imei，现在可以通过files/KeyInfo.bin获取了，需要解密文件，算法为RC4，密钥为_wEcHAT_</p>
        <p>3.MosGram(泡泡)，密码为cust_id的md5值，在sp目录的account_config.xml文件中，使用SQLCipher4参数解密</p>
        <p>4.悟空IM系列的聊天数据库，数据库名为wk_用户ID.db，解密密码即为用户ID，使用SQLCipher4参数解密</p>
      </div>
    </t-dialog>
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
.floating-button-container {
  position: fixed;
  left: 5vh; /* 与侧边栏宽度一致 */
  bottom: 5vh;
  z-index: 100;
  transition: left 0.2s; /* 侧边栏折叠时平滑过渡 */
}

/* 悬浮按钮样式 */
.floating-button {
  width: 8vh;
  height: 8vh;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-radius: 50%;
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
