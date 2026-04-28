<script setup>
import { storeToRefs } from 'pinia'

import { useNotificationsStore } from '@/stores/notifications'
import Toast from '@/components/ui/Toast.vue'

const notificationsStore = useNotificationsStore()
const { items } = storeToRefs(notificationsStore)
</script>

<template>
  <section class="toast-container" aria-live="polite">
    <transition-group name="fade-up" tag="div" class="toast-container__list">
      <Toast
        v-for="item in items"
        :key="item.id"
        :type="item.type"
        :message="item.message"
        @close="notificationsStore.dismiss(item.id)"
      />
    </transition-group>
  </section>
</template>

<style scoped>
.toast-container {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: 80;
}

.toast-container__list {
  display: grid;
  gap: 8px;
}
</style>
