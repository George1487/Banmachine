<template>
  <section class="auth-layout container">
    <div class="auth-layout__left card-surface">
      <RouterView />
    </div>

    <aside class="auth-layout__right card-surface">
      <h2>BanMachine</h2>
      <p>
        Платформа для приёма лабораторных работ и анализа совпадений в рамках одной лабораторной.
      </p>
    </aside>
  </section>
</template>

<style scoped>
.auth-layout {
  min-height: 100dvh;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  align-items: center;
  padding: 24px 0;
}

.auth-layout__left,
.auth-layout__right {
  padding: 24px;
}

.auth-layout__right {
  display: grid;
  gap: 12px;
}

.auth-layout__right p {
  color: var(--ink-soft);
}

@media (max-width: 900px) {
  .auth-layout {
    grid-template-columns: 1fr;
    align-items: start;
  }

  .auth-layout__right {
    display: none;
  }
}
</style>
