<template>
  <div class="page-container">
    <Card class="result-card">
      <template #title>关于</template>
      <template #content>
      
      <div class="space-vertical">
        <h4>关于软件</h4>
        <ul class="list">
          <li>个人用于学习电子取证、提升效率的小软件。</li>
          <li>
            <span><i class="pi pi-github"></i>开源地址: 
              <a href="https://github.com/WXjzcccc/ForensicsTool" target="_blank" class="link">https://github.com/WXjzcccc/ForensicsTool</a>
            </span>
          </li>
        </ul>
      </div>
      <Divider />
      <div class="space-vertical">
        <h4>贡献者</h4>
        <ul class="list">
          <li>
            <span>[1]WXjzc,
              <i class="pi pi-github"></i>
              <a href="https://github.com/WXjzcccc" target="_blank" class="link">https://github.com/WXjzcccc</a>
            </span>
          </li>
          <li>
            <span>[2]b3nguang,
              <i class="pi pi-github"></i>
              <a href="https://github.com/b3nguang" target="_blank" class="link">https://github.com/b3nguang</a>
            </span>
          </li>
        </ul>
      </div>
      <Divider />
      </template>
    </Card>
  </div>
</template>

<script setup>
</script>

<style scoped>
.space-vertical {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.list li {
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--surface-border);
}

.list li:last-child {
  border-bottom: none;
}

.link {
  color: var(--primary-color);
  text-decoration: underline;
}

.link:hover {
  color: var(--primary-color-text);
}

.pi-github {
  margin-right: 0.5rem;
}
</style>