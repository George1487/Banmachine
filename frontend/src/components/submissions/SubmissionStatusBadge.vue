<script setup>
import Badge from '@/components/ui/Badge.vue'

defineProps({
  status: {
    type: String,
    default: 'uploaded',
  },
})
</script>

<template>
  <Badge :value="status" />
</template>
