<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    required: true,
  },
})
</script>

<template>
  <article class="tile card-surface">
    <p class="tile__label">{{ label }}</p>
    <p class="tile__value mono">{{ value }}</p>
  </article>
</template>

<style scoped>
.tile {
  padding: 14px 16px;
  display: grid;
  gap: 8px;
}

.tile__label {
  color: var(--ink-soft);
  font-size: 0.85rem;
}

.tile__value {
  color: var(--ink);
  font-size: 1.45rem;
  font-weight: 700;
}
</style>
