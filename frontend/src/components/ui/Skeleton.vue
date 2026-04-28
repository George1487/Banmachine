<script setup>
defineProps({
  height: {
    type: String,
    default: '16px',
  },
})
</script>

<template>
  <div class="skeleton" :style="{ height }" />
</template>

<style scoped>
.skeleton {
  width: 100%;
  border-radius: 10px;
  background: linear-gradient(
    90deg,
    color-mix(in srgb, var(--line), white 15%) 0%,
    color-mix(in srgb, var(--line), white 45%) 50%,
    color-mix(in srgb, var(--line), white 15%) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s linear infinite;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }

  100% {
    background-position: -200% 0;
  }
}
</style>
