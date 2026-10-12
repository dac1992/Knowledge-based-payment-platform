<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    description?: string
    /** 标题左侧的强调色竖条 */
    accent?: boolean
  }>(),
  { description: '', accent: true },
)
</script>

<template>
  <header class="ym-page-header">
    <div class="ym-page-header__main">
      <h2 class="ym-page-header__title" :class="{ 'is-accent': accent }">{{ title }}</h2>
      <p v-if="description" class="ym-page-header__desc">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="ym-page-header__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.ym-page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 18px;
}

.ym-page-header__title {
  position: relative;
  margin: 0;
  font-size: 19px;
  font-weight: 500;
  color: var(--ym-ink, #1a1f2b);
}

.ym-page-header__title.is-accent {
  padding-left: 12px;
}

.ym-page-header__title.is-accent::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 17px;
  border-radius: 2px;
  background: var(--ym-primary, #1a56db);
}

.ym-page-header__desc {
  margin: 6px 0 0;
  padding-left: 12px;
  font-size: 13px;
  color: var(--ym-ink-3, #6b7688);
}

.ym-page-header__title:not(.is-accent) + .ym-page-header__desc {
  padding-left: 0;
}

.ym-page-header__actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
