<script setup>
const props = defineProps({
  type: {
    type: String,
    default: 'success',
  },
  message: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['close'])
</script>

<template>
  <div class="toast" :class="`toast--${props.type}`">
    <p>{{ message }}</p>
    <button class="toast__close" type="button" @click="emit('close')">×</button>
  </div>
</template>

<style scoped>
.toast {
  min-width: 260px;
  max-width: 420px;
  border: 1px solid var(--line);
  border-radius: var(--radius-control);
  background: var(--surface);
  box-shadow: var(--shadow-card);
  padding: 12px 14px;
  display: flex;
  gap: 12px;
  align-items: start;
  justify-content: space-between;
}

.toast p {
  margin: 0;
  color: var(--ink);
  line-height: 1.4;
}

.toast__close {
  border: 0;
  background: transparent;
  font-size: 18px;
  line-height: 1;
  color: var(--ink-soft);
  cursor: pointer;
}

.toast--success {
  border-left: 4px solid var(--success);
}

.toast--warning {
  border-left: 4px solid var(--warning);
}

.toast--error {
  border-left: 4px solid var(--danger);
}
</style>
