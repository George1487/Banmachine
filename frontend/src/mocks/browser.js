import { setupWorker } from 'msw/browser'
import { authHandlers } from './handlers/auth.js'
import { labsHandlers } from './handlers/labs.js'
import { submissionsHandlers } from './handlers/submissions.js'
import { analysisHandlers } from './handlers/analysis.js'

export const worker = setupWorker(
  ...authHandlers,
  ...labsHandlers,
  ...submissionsHandlers,
  ...analysisHandlers,
)
