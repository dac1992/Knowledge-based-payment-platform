<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    description?: string
    /** 视觉尺寸 */
    size?: 'sm' | 'md'
  }>(),
  { title: '暂无数据', description: '', size: 'md' },
)
</script>

<template>
  <div class="ym-empty" :class="`is-${size}`">
    <div class="ym-empty__icon" aria-hidden="true">
      <svg viewBox="0 0 48 48" width="44" height="44">
        <rect
          x="7"
          y="10"
          width="34"
          height="24"
          rx="3"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        />
        <path
          d="M7 28h10l3 4h8l3-4h10"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linejoin="round"
        />
        <path
          d="M17 18h14M17 23h9"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
        />
      </svg>
    </div>
    <p class="ym-empty__title">{{ title }}</p>
    <p v-if="description" class="ym-empty__desc">{{ description }}</p>
    <div v-if="$slots.default" class="ym-empty__extra">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.ym-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 40px 20px;
  text-align: center;
  color: var(--ym-ink-4, #98a2b3);
}

.ym-empty.is-sm {
  padding: 22px 16px;
}

.ym-empty__icon {
  color: var(--ym-ink-4, #98a2b3);
  opacity: 0.75;
}

.ym-empty__title {
  margin: 4px 0 0;
  font-size: 14px;
  color: var(--ym-ink-2, #3d4759);
}

.ym-empty__desc {
  margin: 0;
  font-size: 12.5px;
  color: var(--ym-ink-4, #98a2b3);
  max-width: 420px;
  line-height: 1.7;
}

.ym-empty__extra {
  margin-top: 10px;
}
</style>
