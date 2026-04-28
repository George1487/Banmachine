<script setup>
import { computed } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import Badge from '@/components/ui/Badge.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const roleLabel = computed(() => authStore.user?.role || '')
const roleTone = computed(() => (authStore.user?.role === 'teacher' ? 'warning' : 'success'))
</script>

<template>
  <header class="header card-surface">
    <RouterLink class="header__logo" to="/">BanMachine</RouterLink>

    <div class="header__right">
      <RouterLink class="header__profile" to="/profile">
        <strong>{{ authStore.user?.fullName || 'Гость' }}</strong>
        <Badge v-if="roleLabel" :value="roleLabel" :tone="roleTone" />
      </RouterLink>
      <AppButton variant="ghost" @click="authStore.logout">Выйти</AppButton>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: sticky;
  top: 0;
  z-index: 40;
  height: var(--header-height);
  border-radius: 0;
  border-left: 0;
  border-right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header__logo {
  color: var(--ink);
  text-decoration: none;
  font-size: 1.1rem;
  font-weight: 800;
  letter-spacing: 0.01em;
}

.header__right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header__profile {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: inherit;
  border-radius: var(--radius-control);
  padding: 4px 8px;
  transition: background 120ms;
}

.header__profile:hover {
  background: var(--surface-muted);
}

@media (max-width: 768px) {
  .header {
    padding-inline: 16px;
  }

  .header__profile strong {
    display: none;
  }
}
</style>
