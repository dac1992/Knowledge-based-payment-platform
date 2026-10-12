<script setup lang="ts">
withDefaults(
  defineProps<{
    label: string
    value: string | number
    hint?: string
    /** 语义色调：neutral 中性 / up 正向（红）/ down 负向（绿）/ warn 预警 */
    tone?: 'neutral' | 'up' | 'down' | 'warn'
    /** 趋势文案，例如 "+12.4%" */
    trend?: string
    /** 加载中骨架 */
    loading?: boolean
  }>(),
  { hint: '', tone: 'neutral', trend: '', loading: false },
)
</script>

<template>
  <div class="ym-stat-card" :class="`is-${tone}`">
    <div class="ym-stat-card__label">{{ label }}</div>
    <div class="ym-stat-card__value">
      <template v-if="loading">
        <span class="ym-stat-card__skeleton" />
      </template>
      <template v-else>{{ value }}</template>
      <span v-if="trend && !loading" class="ym-stat-card__trend">{{ trend }}</span>
    </div>
    <div v-if="hint" class="ym-stat-card__hint">{{ hint }}</div>
  </div>
</template>

<style scoped>
.ym-stat-card {
  background: var(--ym-card, #fff);
  border: 1px solid var(--ym-line, #e4e8ef);
  border-radius: var(--ym-radius, 10px);
  padding: 16px 18px;
  box-shadow: var(--ym-shadow-sm, 0 1px 2px rgba(26, 31, 43, 0.06));
}

.ym-stat-card__label {
  font-size: 12.5px;
  color: var(--ym-ink-3, #6b7688);
}

.ym-stat-card__value {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-top: 8px;
  font-size: 26px;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  color: var(--ym-ink, #1a1f2b);
}

.ym-stat-card__trend {
  font-size: 12.5px;
  font-weight: 400;
}

.ym-stat-card__hint {
  margin-top: 6px;
  font-size: 12px;
  color: var(--ym-ink-4, #98a2b3);
}

.ym-stat-card__skeleton {
  display: inline-block;
  width: 76px;
  height: 24px;
  border-radius: 6px;
  background: var(--ym-line-light, #eef1f6);
}

.ym-stat-card.is-up .ym-stat-card__value {
  color: var(--ym-up, #dc2626);
}

.ym-stat-card.is-down .ym-stat-card__value {
  color: var(--ym-down, #16a34a);
}

.ym-stat-card.is-warn .ym-stat-card__value {
  color: var(--ym-warn, #d97706);
}

.ym-stat-card.is-up .ym-stat-card__trend {
  color: var(--ym-up, #dc2626);
}

.ym-stat-card.is-down .ym-stat-card__trend {
  color: var(--ym-down, #16a34a);
}
</style>
