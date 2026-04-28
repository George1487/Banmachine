<script setup>
import { useAuthStore } from '@/stores/auth'
import AppBottomNav from './AppBottomNav.vue'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'

const authStore = useAuthStore()
</script>

<template>
  <div class="shell">
    <AppHeader />

    <div class="shell__body container" :class="{ 'shell__body--with-sidebar': authStore.isTeacher }">
      <AppSidebar v-if="authStore.isTeacher" class="shell__sidebar" />

      <main class="shell__content">
        <slot />
      </main>
    </div>

    <AppBottomNav />
  </div>
</template>

<style scoped>
.shell {
  min-height: 100dvh;
}

.shell__body {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  padding: 16px 0 28px;
}

.shell__body--with-sidebar {
  grid-template-columns: minmax(250px, 280px) minmax(0, 1fr);
}

.shell__content {
  min-width: 0;
  display: grid;
  gap: 16px;
  align-content: start;
}

@media (max-width: 1024px) {
  .shell__body {
    grid-template-columns: 1fr;
    padding-bottom: 88px;
  }

  .shell__sidebar {
    display: none;
  }
}
</style>
