export function formatDateTime(value) {
  if (!value) return '—'

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '—'

  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatShortDate(value) {
  if (!value) return '—'

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '—'

  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

export function formatPercent(value) {
  if (value === null || value === undefined || Number.isNaN(Number(value))) {
    return '—'
  }

  return `${Math.round(Number(value) * 100)}%`
}

export function formatScore(value) {
  if (value === null || value === undefined || Number.isNaN(Number(value))) {
    return '—'
  }

  return Number(value).toFixed(2)
}

export function mapRiskLabel(risk) {
  if (risk === 'high') return 'Высокий риск'
  if (risk === 'medium') return 'Средний риск'
  if (risk === 'low') return 'Низкий риск'
  return '—'
}

export function mapSubmissionStatus(status) {
  const labels = {
    uploaded: 'Загружено',
    parsing: 'Обрабатывается',
    parsed: 'Готово',
    failed: 'Ошибка',
  }

  return labels[status] || status || '—'
}

export function mapJobStatus(status) {
  const labels = {
    pending: 'В очереди',
    processing: 'Выполняется',
    done: 'Завершено',
    failed: 'Ошибка',
  }

  return labels[status] || status || '—'
}
