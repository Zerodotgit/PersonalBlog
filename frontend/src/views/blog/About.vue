<template>
  <div class="about">
    <h1>关于</h1>
    <p>
      你好，这是我的个人博客。我在这里记录开发过程中遇到的问题、解决方案以及一些技术思考，
      内容涵盖 Go、前端、数据库与工程化等方向。
    </p>
    <p>本站基于 Go + Gin 后端与 Vue 3 + Vite 前端搭建，文章使用 Markdown 编写。</p>
    <p>欢迎通过下方方式交流（示例占位，可按需替换）：</p>
    <ul>
      <li>GitHub：your-github-id</li>
      <li>邮箱：your@email.com</li>
    </ul>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.about {
  max-width: var(--content-max);
  margin: 0 auto;
  line-height: 1.8;
}
.about h1 {
  font-size: 28px;
  margin-bottom: 16px;
}
</style>
