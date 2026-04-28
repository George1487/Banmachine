<script setup>
import AppButton from '@/components/button/AppButton.vue'
import Badge from '@/components/ui/Badge.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<template>
  <section class="view">
    <PageHeader title="Профиль" />

    <div class="profile card-surface">
      <div class="profile__row">
        <span class="profile__label">Имя</span>
        <span class="profile__value">{{ authStore.user?.fullName || '—' }}</span>
      </div>
      <div class="profile__row">
        <span class="profile__label">Email</span>
        <span class="profile__value">{{ authStore.user?.email || '—' }}</span>
      </div>
      <div v-if="authStore.user?.groupName" class="profile__row">
        <span class="profile__label">Группа</span>
        <span class="profile__value">{{ authStore.user.groupName }}</span>
      </div>
      <div class="profile__row">
        <span class="profile__label">Роль</span>
        <Badge :value="authStore.user?.role" />
      </div>
    </div>

    <AppButton variant="danger" @click="authStore.logout">Выйти из аккаунта</AppButton>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.profile {
  display: grid;
  gap: 0;
}

.profile__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-bottom: 1px solid var(--line);
}

.profile__row:last-child {
  border-bottom: none;
}

.profile__label {
  color: var(--ink-soft);
  font-size: 0.875rem;
}

.profile__value {
  font-weight: 600;
}
</style>
