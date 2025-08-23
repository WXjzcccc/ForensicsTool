<template>
  <div class="about-container">
    <t-card class="about-card">
      <t-space direction="vertical">
        <h4>关于软件</h4>
        <t-list>
          <t-list-item>个人用于学习电子取证、提升效率的小软件。</t-list-item>
          <t-list-item><span><t-icon name="logo-github"></t-icon>开源地址: <t-link theme="primary" underline href="https://github.com/WXjzcccc/ForensicsTool" target="_blank">https://github.com/WXjzcccc/ForensicsTool</t-link></span></t-list-item>
        </t-list>
      </t-space>
      <t-divider />
      <t-space direction="vertical">
        <h4>贡献者</h4>
        <t-list>
          <t-list-item><span>[1]WXjzc,<t-icon name="logo-github"></t-icon><t-link theme="primary" underline href="https://github.com/WXjzcccc" target="_blank">https://github.com/WXjzcccc</t-link></span></t-list-item>
          <t-list-item><span>[2]b3nguang,<t-icon name="logo-github"></t-icon><t-link theme="primary" underline href="https://github.com/b3nguang" target="_blank">https://github.com/b3nguang</t-link></span></t-list-item>
        </t-list>
      </t-space>
      <t-divider />
    </t-card>
  </div>
</template>

<script setup>

</script>

<style scoped>
.about-container {
  max-width: 85vw;
  margin: 0 auto;
}
.about-card {
  text-align: left;
  border-color: blue;
  border-width: 3px;
  height: 92vh;
  margin-left: 0.2vw;
  overflow-y: scroll;
}
</style>