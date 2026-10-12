<script setup lang="ts">
import { useRouter } from 'vue-router'

import { YmEmptyState } from '@yanmai/ui'

const router = useRouter()
</script>

<template>
  <div class="nf">
    <YmEmptyState title="页面不存在" description="你访问的地址可能已被移动或删除。">
      <el-button type="primary" @click="router.replace('/dashboard')">返回工作台</el-button>
    </YmEmptyState>
  </div>
</template>

<style scoped>
.nf {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: #fff;
}
</style>
