<template>
  <aside class="sidebar card-surface">
    <p class="sidebar__title">Навигация</p>
    <nav class="sidebar__nav">
      <RouterLink to="/teacher/labs">Лабораторные</RouterLink>
      <RouterLink to="/teacher/labs/create">Создать лабораторную</RouterLink>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  position: sticky;
  top: calc(var(--header-height) + 16px);
  padding: 14px;
  min-width: 250px;
  max-height: calc(100dvh - var(--header-height) - 32px);
}

.sidebar__title {
  color: var(--ink-soft);
  font-size: 0.82rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 8px;
}

.sidebar__nav {
  display: grid;
  gap: 4px;
}

.sidebar__nav a {
  text-decoration: none;
  color: var(--ink);
  border-radius: 10px;
  min-height: 40px;
  display: flex;
  align-items: center;
  padding: 0 10px;
}

.sidebar__nav a.router-link-active {
  color: var(--accent);
  background: var(--accent-soft);
  font-weight: 700;
}
</style>
