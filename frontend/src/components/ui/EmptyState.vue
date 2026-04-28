<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <section class="empty-state card-surface">
    <h3>{{ title }}</h3>
    <p v-if="description">{{ description }}</p>
    <slot name="action" />
  </section>
</template>

<style scoped>
.empty-state {
  padding: 28px;
  text-align: center;
  display: grid;
  gap: 10px;
}

p {
  color: var(--ink-soft);
}
</style>
