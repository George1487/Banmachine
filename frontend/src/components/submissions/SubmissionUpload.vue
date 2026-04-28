<script setup>
import AppButton from '@/components/button/AppButton.vue'
import Dropzone from '@/components/form/Dropzone.vue'
import SubmissionStatusTimeline from '@/components/submissions/SubmissionStatusTimeline.vue'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  isUploading: {
    type: Boolean,
    default: false,
  },
  progress: {
    type: Number,
    default: 0,
  },
  status: {
    type: String,
    default: '',
  },
  errorMessage: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['file-selected', 'retry'])
</script>

<template>
  <section class="upload card-surface">
    <Dropzone :disabled="disabled || isUploading" @file-selected="emit('file-selected', $event)" />

    <div v-if="isUploading" class="upload__progress">
      <p>Загрузка: {{ progress }}%</p>
      <div class="upload__bar">
        <div class="upload__bar-fill" :style="{ width: `${progress}%` }" />
      </div>
    </div>

    <SubmissionStatusTimeline v-if="status" :status="status" />

    <p v-if="errorMessage" class="upload__error">{{ errorMessage }}</p>

    <AppButton v-if="status === 'failed'" variant="secondary" @click="emit('retry')">
      Загрузить заново
    </AppButton>
  </section>
</template>

<style scoped>
.upload {
  display: grid;
  gap: 12px;
}

.upload__progress {
  display: grid;
  gap: 6px;
}

.upload__bar {
  height: 10px;
  border-radius: var(--radius-pill);
  border: 1px solid var(--line);
  background: var(--surface-muted);
  overflow: hidden;
}

.upload__bar-fill {
  height: 100%;
  background: var(--accent);
  transition: width 160ms ease;
}

.upload__error {
  color: var(--danger);
  font-weight: 600;
}
</style>
