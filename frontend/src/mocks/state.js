import { users as usersData } from './data/users.js'
import { labs as labsData } from './data/labs.js'
import { submissions as submissionsData } from './data/submissions.js'
import { analysisJobs as jobsData, analysisResults as resultsData, matches as matchesData } from './data/analysis.js'

function clone(obj) {
  return JSON.parse(JSON.stringify(obj))
}

export const state = {
  users: clone(usersData),
  labs: clone(labsData),
  submissions: clone(submissionsData),
  analysisJobs: clone(jobsData),
  analysisResults: clone(resultsData),
  matches: clone(matchesData),

  // accessToken/refreshToken → userId
  sessions: {},

  // Счётчики для генерации id
  nextUserId: 20,
  nextSubId: 20,
  nextJobId: 10,
  nextLabId: 10,

  // Счётчик polling-запросов для имитации прогресса job: jobId → count
  jobPollCount: {},
}
