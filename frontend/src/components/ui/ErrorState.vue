<script setup>
const props = defineProps({
  title: {
    type: String,
    default: 'Не удалось загрузить данные',
  },
  message: {
    type: String,
    default: 'Попробуйте обновить страницу или повторить действие.',
  },
})

const emit = defineEmits(['retry'])
</script>

<template>
  <section class="error-state card-surface">
    <h3>{{ props.title }}</h3>
    <p>{{ props.message }}</p>
    <button class="error-state__retry focus-ring" type="button" @click="emit('retry')">
      Повторить
    </button>
  </section>
</template>

<style scoped>
.error-state {
  padding: 24px;
  text-align: center;
  display: grid;
  gap: 10px;
}

p {
  color: var(--ink-soft);
}

.error-state__retry {
  justify-self: center;
  border: 1px solid var(--line);
  background: var(--surface);
  border-radius: var(--radius-control);
  min-height: 40px;
  padding: 0 14px;
  cursor: pointer;
}
</style>
