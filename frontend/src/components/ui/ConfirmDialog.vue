<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: 'Подтвердите действие',
  },
  message: {
    type: String,
    default: 'Вы уверены?',
  },
})

const emit = defineEmits(['update:modelValue', 'confirm'])

function close() {
  emit('update:modelValue', false)
}

function confirm() {
  emit('confirm')
  close()
}
</script>

<template>
  <teleport to="body">
    <transition name="quick-fade">
      <div v-if="modelValue" class="dialog-overlay" @click.self="close">
        <section class="dialog card-surface" role="dialog" aria-modal="true">
          <h3>{{ title }}</h3>
          <p>{{ message }}</p>
          <div class="dialog__actions">
            <button class="dialog__button focus-ring" type="button" @click="close">Отмена</button>
            <button class="dialog__button dialog__button--danger focus-ring" type="button" @click="confirm">
              Подтвердить
            </button>
          </div>
        </section>
      </div>
    </transition>
  </teleport>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.4);
  display: grid;
  place-items: center;
  padding: 16px;
  z-index: 60;
}

.dialog {
  width: min(520px, 100%);
  padding: 20px;
  box-shadow: var(--shadow-modal);
  display: grid;
  gap: 12px;
}

p {
  color: var(--ink-soft);
}

.dialog__actions {
  display: flex;
  justify-content: end;
  gap: 8px;
}

.dialog__button {
  min-height: 40px;
  border-radius: var(--radius-control);
  border: 1px solid var(--line);
  padding: 0 14px;
  background: var(--surface);
  cursor: pointer;
}

.dialog__button--danger {
  border-color: var(--danger);
  background: var(--danger);
  color: #fff;
}
</style>
