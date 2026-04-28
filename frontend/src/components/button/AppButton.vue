<script setup>
import { computed } from 'vue'

const props = defineProps({
  as: {
    type: [String, Object],
    default: 'button',
  },
  variant: {
    type: String,
    default: 'primary',
  },
  size: {
    type: String,
    default: 'md',
  },
  type: {
    type: String,
    default: 'button',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  block: {
    type: Boolean,
    default: false,
  },
})

const classes = computed(() => [
  'app-button',
  `app-button--${props.variant}`,
  `app-button--${props.size}`,
  props.block && 'app-button--block',
])
</script>

<template>
  <component
    :is="as"
    class="focus-ring"
    :class="classes"
    :type="as === 'button' ? type : undefined"
    :disabled="disabled"
  >
    <slot />
  </component>
</template>

<style scoped>
.app-button {
  min-height: 44px;
  padding: 0 16px;
  border-radius: var(--radius-control);
  border: 1px solid transparent;
  font-size: 0.95rem;
  font-weight: 700;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 160ms ease;
}

.app-button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.app-button--primary {
  background: var(--accent);
  color: #fff;
}

.app-button--primary:hover:not(:disabled) {
  background: var(--accent-hover);
}

.app-button--secondary {
  background: var(--surface);
  border-color: var(--line);
  color: var(--ink);
}

.app-button--secondary:hover:not(:disabled) {
  background: var(--surface-muted);
}

.app-button--ghost {
  background: transparent;
  color: var(--ink);
}

.app-button--ghost:hover:not(:disabled) {
  background: #ffffff66;
}

.app-button--danger {
  background: var(--danger);
  color: #fff;
}

.app-button--danger:hover:not(:disabled) {
  filter: brightness(0.95);
}

.app-button--sm {
  min-height: 36px;
  padding-inline: 12px;
  font-size: 0.85rem;
}

.app-button--block {
  width: 100%;
}
</style>
