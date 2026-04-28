export const analysisJobs = [
  {
    jobId: 'job_1',
    labId: 'lab_3',
    status: 'Finished',
    createdBy: 'user_1',
    createdAt: '2026-03-21T10:00:00Z',
    startedAt: '2026-03-21T10:00:05Z',
    finishedAt: '2026-03-21T10:03:00Z',
    errorMessage: null,
  },
  {
    jobId: 'job_2',
    labId: 'lab_1',
    status: 'Finished',
    createdBy: 'user_1',
    createdAt: '2026-03-27T14:00:00Z',
    startedAt: '2026-03-27T14:00:03Z',
    finishedAt: '2026-03-27T14:02:30Z',
    errorMessage: null,
  },
]

export const analysisResults = {
  lab_1: {
    lastAnalysisJob: {
      id: 'job_2',
      status: 'Finished',
      createdAt: '2026-03-27T14:00:00Z',
      finishedAt: '2026-03-27T14:02:30Z',
    },
    stats: {
      totalSubmissions: 5,
      actualSubmissions: 5,
      parsedSubmissions: 3,
      highRiskCount: 1,
      mediumRiskCount: 1,
      lowRiskCount: 1,
      maxFinalScore: 0.84,
    },
    items: [
      {
        submissionId: 'sub_1',
        student: { studentId: 'user_10', fullName: 'Иванов Иван Сергеевич' },
        topMatchSubmissionId: 'sub_2',
        topMatchScore: 0.84,
        finalScoreRiskLevel: 'High',
      },
      {
        submissionId: 'sub_2',
        student: { studentId: 'user_11', fullName: 'Смирнова Мария Александровна' },
        topMatchSubmissionId: 'sub_1',
        topMatchScore: 0.84,
        finalScoreRiskLevel: 'High',
      },
      {
        submissionId: 'sub_3',
        student: { studentId: 'user_12', fullName: 'Козлов Дмитрий Андреевич' },
        topMatchSubmissionId: 'sub_1',
        topMatchScore: 0.45,
        finalScoreRiskLevel: 'Medium',
      },
    ],
  },
  lab_3: {
    lastAnalysisJob: {
      id: 'job_1',
      status: 'Finished',
      createdAt: '2026-03-21T10:00:00Z',
      finishedAt: '2026-03-21T10:03:00Z',
    },
    stats: {
      totalSubmissions: 3,
      actualSubmissions: 3,
      parsedSubmissions: 3,
      highRiskCount: 0,
      mediumRiskCount: 1,
      lowRiskCount: 2,
      maxFinalScore: 0.52,
    },
    items: [
      {
        submissionId: 'sub_8',
        student: { studentId: 'user_10', fullName: 'Иванов Иван Сергеевич' },
        topMatchSubmissionId: 'sub_9',
        topMatchScore: 0.52,
        finalScoreRiskLevel: 'Medium',
      },
      {
        submissionId: 'sub_9',
        student: { studentId: 'user_11', fullName: 'Смирнова Мария Александровна' },
        topMatchSubmissionId: 'sub_8',
        topMatchScore: 0.52,
        finalScoreRiskLevel: 'Medium',
      },
      {
        submissionId: 'sub_10',
        student: { studentId: 'user_12', fullName: 'Козлов Дмитрий Андреевич' },
        topMatchSubmissionId: 'sub_8',
        topMatchScore: 0.28,
        finalScoreRiskLevel: 'Low',
      },
    ],
  },
}

export const matches = {
  sub_1: {
    submissionId: 'sub_1',
    analysisJobId: 'job_2',
    matches: [
      {
        otherSubmissionId: 'sub_2',
        student: { studentId: 'user_11', fullName: 'Смирнова Мария Александровна' },
        textScore: 0.91,
        calculationScore: 0.76,
        imagesScore: 0.52,
        finalScore: 0.84,
        riskLevel: 'High',
      },
      {
        otherSubmissionId: 'sub_3',
        student: { studentId: 'user_12', fullName: 'Козлов Дмитрий Андреевич' },
        textScore: 0.38,
        calculationScore: 0.55,
        imagesScore: 0.21,
        finalScore: 0.40,
        riskLevel: 'Medium',
      },
    ],
  },
  sub_2: {
    submissionId: 'sub_2',
    analysisJobId: 'job_2',
    matches: [
      {
        otherSubmissionId: 'sub_1',
        student: { studentId: 'user_10', fullName: 'Иванов Иван Сергеевич' },
        textScore: 0.91,
        calculationScore: 0.76,
        imagesScore: 0.52,
        finalScore: 0.84,
        riskLevel: 'High',
      },
      {
        otherSubmissionId: 'sub_3',
        student: { studentId: 'user_12', fullName: 'Козлов Дмитрий Андреевич' },
        textScore: 0.30,
        calculationScore: 0.42,
        imagesScore: 0.15,
        finalScore: 0.32,
        riskLevel: 'Low',
      },
    ],
  },
  sub_3: {
    submissionId: 'sub_3',
    analysisJobId: 'job_2',
    matches: [
      {
        otherSubmissionId: 'sub_1',
        student: { studentId: 'user_10', fullName: 'Иванов Иван Сергеевич' },
        textScore: 0.38,
        calculationScore: 0.55,
        imagesScore: 0.21,
        finalScore: 0.40,
        riskLevel: 'Medium',
      },
      {
        otherSubmissionId: 'sub_2',
        student: { studentId: 'user_11', fullName: 'Смирнова Мария Александровна' },
        textScore: 0.30,
        calculationScore: 0.42,
        imagesScore: 0.15,
        finalScore: 0.32,
        riskLevel: 'Low',
      },
    ],
  },
  sub_8: {
    submissionId: 'sub_8',
    analysisJobId: 'job_1',
    matches: [
      {
        otherSubmissionId: 'sub_9',
        student: { studentId: 'user_11', fullName: 'Смирнова Мария Александровна' },
        textScore: 0.55,
        calculationScore: 0.60,
        imagesScore: 0.30,
        finalScore: 0.52,
        riskLevel: 'Medium',
      },
      {
        otherSubmissionId: 'sub_10',
        student: { studentId: 'user_12', fullName: 'Козлов Дмитрий Андреевич' },
        textScore: 0.25,
        calculationScore: 0.35,
        imagesScore: 0.18,
        finalScore: 0.28,
        riskLevel: 'Low',
      },
    ],
  },
}
