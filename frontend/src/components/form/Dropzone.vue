<script setup>
import { ref } from 'vue'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  hint: {
    type: String,
    default: 'Перетащите файл .docx сюда или выберите его вручную',
  },
})

const emit = defineEmits(['file-selected'])

const isOver = ref(false)

function onDrop(event) {
  event.preventDefault()
  isOver.value = false

  if (props.disabled) return

  const file = event.dataTransfer?.files?.[0]
  if (file) emit('file-selected', file)
}

function onPick(event) {
  const file = event.target.files?.[0]
  if (file) emit('file-selected', file)
}
</script>

<template>
  <label
    class="dropzone"
    :class="{ 'dropzone--over': isOver, 'dropzone--disabled': disabled }"
    @dragover.prevent="!disabled && (isOver = true)"
    @dragleave.prevent="isOver = false"
    @drop="onDrop"
  >
    <input
      class="dropzone__input"
      type="file"
      accept=".docx"
      :disabled="disabled"
      @change="onPick"
    />
    <strong>Загрузка отчета</strong>
    <p>{{ hint }}</p>
    <span class="dropzone__action">Выбрать файл</span>
  </label>
</template>

<style scoped>
.dropzone {
  border: 2px dashed var(--line);
  border-radius: var(--radius-card);
  padding: 24px;
  display: grid;
  gap: 8px;
  justify-items: center;
  text-align: center;
  background: color-mix(in srgb, var(--surface), var(--accent-soft) 22%);
  cursor: pointer;
  transition: all 160ms ease;
}

.dropzone--over {
  border-color: var(--accent);
  transform: translateY(-1px);
}

.dropzone--disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.dropzone__input {
  display: none;
}

p {
  color: var(--ink-soft);
}

.dropzone__action {
  min-height: 36px;
  padding: 0 12px;
  border-radius: var(--radius-control);
  border: 1px solid var(--line);
  background: var(--surface);
  display: inline-flex;
  align-items: center;
  font-weight: 700;
}
</style>
