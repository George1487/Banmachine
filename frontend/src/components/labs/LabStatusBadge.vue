<script setup>
import Badge from '@/components/ui/Badge.vue'

defineProps({
  status: {
    type: String,
    default: 'active',
  },
})
</script>

<template>
  <Badge :value="status" />
</template>
