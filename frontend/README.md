# BanMachine — Frontend

**Фронтенд-часть платформы для сдачи лабораторных работ и автоматического обнаружения плагиата.**

Студенты загружают отчёты через drag-and-drop. Преподаватели запускают анализ одной кнопкой и получают детальные score-карты по каждому совпадению — с разбивкой по тексту, формулам и изображениям.

> Фронтенд работает с реальным бэкендом через `VITE_API_BASE_URL`.

---

## Стек

| | |
|---|---|
| **Framework** | Vue 3.5 — Composition API, `<script setup>` |
| **Router** | vue-router 5 |
| **State** | Pinia 3 |
| **HTTP** | Axios · token refresh queue · `bindAuthTokenBridge` |
| **Build** | Vite 8 |
| **CSS** | CSS custom properties, без UI-фреймворка |
| **Fonts** | Manrope (UI) · IBM Plex Mono (scores, id, технические значения) |

---

## Быстрый старт

```bash
npm install
npm run dev
```

Открывается на `http://localhost:5173`. Для работы нужен доступный бэкенд.

---

## Сборка и деплой

```bash
npm run build    # → dist/
npm run preview  # превью production-сборки
```

Приложение ходит в реальный бэкенд по `VITE_API_BASE_URL`.

---

## Структура

```
src/
├── services/       # axios-клиенты по домену: auth, labs, submissions, analysis
├── stores/         # Pinia: auth, labs, submissions, analysis, notifications
├── composables/    # usePolling · useFileUpload
├── layouts/        # AuthLayout · MainLayout
├── views/          # страницы: student/* · teacher/* · auth/* · /profile
└── components/     # UI-kit + domain-компоненты (Badge, DataTable, MatchCard, ...)
```

---

## Маршруты

| Роль | Маршруты |
|------|----------|
| **Студент** | `/student/labs` · `/student/labs/:id/submit` · `/student/submissions` |
| **Преподаватель** | `/teacher/labs` · `/teacher/labs/:id/submissions` · `/teacher/labs/:id/analysis` · `/teacher/submissions/:id/matches` |
| **Оба** | `/profile` · `/login` · `/register` |
