<template>
  <section class="not-found card-surface">
    <h1>404</h1>
    <p>Страница не найдена.</p>
    <RouterLink to="/">Перейти на главную</RouterLink>
  </section>
</template>

<style scoped>
.not-found {
  margin: 60px auto;
  max-width: 520px;
  padding: 28px;
  text-align: center;
  display: grid;
  gap: 12px;
}

p {
  color: var(--ink-soft);
}
</style>
