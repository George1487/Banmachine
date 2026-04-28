import { onBeforeUnmount, ref } from 'vue'

export function usePolling(fetchFn, options = {}) {
  const interval = options.interval || 3000
  const stopWhen = options.stopWhen || (() => false)

  const data = ref(null)
  const error = ref(null)
  const isPolling = ref(false)

  let timerId = null
  let inFlight = false

  async function tick() {
    if (inFlight) return

    inFlight = true

    try {
      const result = await fetchFn()
      data.value = result
      error.value = null

      if (stopWhen(result)) {
        stop()
      }
    } catch (err) {
      error.value = err
    } finally {
      inFlight = false
    }
  }

  async function start() {
    if (isPolling.value) return

    isPolling.value = true
    await tick()

    if (!isPolling.value) return

    timerId = window.setInterval(tick, interval)
  }

  function stop() {
    isPolling.value = false

    if (timerId) {
      window.clearInterval(timerId)
      timerId = null
    }
  }

  onBeforeUnmount(stop)

  return {
    data,
    error,
    isPolling,
    start,
    stop,
  }
}
