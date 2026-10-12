<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    /** 标记尺寸（px） */
    size?: number
    /** 是否展示品牌文字 */
    showText?: boolean
    /** 副标题，传空字符串可隐藏 */
    subtitle?: string
    /** 深色底场景（侧边栏）下的浅色文字 */
    inverse?: boolean
  }>(),
  { size: 34, showText: true, subtitle: '学术知识付费平台', inverse: false },
)

const markStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  fontSize: `${Math.round(props.size * 0.45)}px`,
}))
</script>

<template>
  <div class="ym-logo" :class="{ 'ym-logo--inverse': inverse }">
    <span class="ym-logo__mark" :style="markStyle" aria-hidden="true">研</span>
    <span v-if="showText" class="ym-logo__text">
      <b>研脉学术</b>
      <i v-if="subtitle">{{ subtitle }}</i>
    </span>
  </div>
</template>

<style scoped>
.ym-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  line-height: 1.25;
}

.ym-logo__mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: none;
  border-radius: 9px;
  background: var(--ym-primary, #1a56db);
  color: #fff;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.ym-logo__text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.ym-logo__text b {
  font-size: 15px;
  font-weight: 500;
  color: var(--ym-ink, #1a1f2b);
}

.ym-logo__text i {
  font-style: normal;
  font-size: 11px;
  color: var(--ym-ink-4, #98a2b3);
}

.ym-logo--inverse .ym-logo__text b {
  color: #fff;
}

.ym-logo--inverse .ym-logo__text i {
  color: #7c8caa;
}
</style>
