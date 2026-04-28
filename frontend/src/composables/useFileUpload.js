import { ref } from 'vue'

export function useFileUpload(uploadFn) {
  const isUploading = ref(false)
  const progress = ref(0)
  const error = ref(null)

  async function upload(...args) {
    isUploading.value = true
    progress.value = 0
    error.value = null

    try {
      const result = await uploadFn(...args, (value) => {
        progress.value = value
      })

      progress.value = 100
      return result
    } catch (err) {
      error.value = err
      throw err
    } finally {
      isUploading.value = false
    }
  }

  return {
    upload,
    isUploading,
    progress,
    error,
  }
}
