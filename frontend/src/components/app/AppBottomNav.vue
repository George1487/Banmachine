<script setup>
import { computed } from 'vue'

import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const links = computed(() => {
  if (authStore.isTeacher) {
    return [
      { to: '/teacher/labs', label: 'Лабы' },
      { to: '/teacher/labs/create', label: 'Создать' },
      { to: '/profile', label: 'Профиль' },
    ]
  }

  return [
    { to: '/student/labs', label: 'Лабы' },
    { to: '/student/submissions', label: 'Мои работы' },
    { to: '/profile', label: 'Профиль' },
  ]
})
</script>

<template>
  <nav class="bottom-nav card-surface">
    <RouterLink v-for="link in links" :key="link.to + link.label" :to="link.to">
      {{ link.label }}
    </RouterLink>
  </nav>
</template>

<style scoped>
.bottom-nav {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  bottom: 10px;
  width: min(calc(100% - 16px), 520px);
  padding: 6px;
  display: none;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  z-index: 50;
}

.bottom-nav a {
  min-height: 40px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  text-decoration: none;
  color: var(--ink-soft);
  font-size: 0.82rem;
  font-weight: 700;
}

.bottom-nav a.router-link-active {
  color: var(--accent);
  background: var(--accent-soft);
}

@media (max-width: 768px) {
  .bottom-nav {
    display: grid;
  }
}
</style>
