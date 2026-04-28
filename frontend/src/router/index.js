import { createRouter, createWebHistory } from 'vue-router'

import AuthLayout from '@/layouts/AuthLayout.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import StudentLabSubmitView from '@/views/student/StudentLabSubmitView.vue'
import StudentLabsView from '@/views/student/StudentLabsView.vue'
import StudentSubmissionsView from '@/views/student/StudentSubmissionsView.vue'
import TeacherDashboardView from '@/views/teacher/TeacherDashboardView.vue'
import TeacherLabAnalysisView from '@/views/teacher/TeacherLabAnalysisView.vue'
import TeacherLabCreateView from '@/views/teacher/TeacherLabCreateView.vue'
import TeacherLabEditView from '@/views/teacher/TeacherLabEditView.vue'
import TeacherLabSubmissionsView from '@/views/teacher/TeacherLabSubmissionsView.vue'
import TeacherSubmissionMatchesView from '@/views/teacher/TeacherSubmissionMatchesView.vue'
import ProfileView from '@/views/ProfileView.vue'
import { useAuthStore } from '@/stores/auth'

let sessionRestored = false

function getRoleHome(role) {
  return role === 'teacher' ? '/teacher/labs' : '/student/labs'
}

const routes = [
  {
    path: '/',
    component: AuthLayout,
    children: [
      {
        path: 'login',
        name: 'login',
        component: LoginView,
        meta: { guest: true },
      },
      {
        path: 'register',
        name: 'register',
        component: RegisterView,
        meta: { guest: true },
      },
    ],
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/student/labs',
      },
      {
        path: 'student/labs',
        name: 'student-labs',
        component: StudentLabsView,
        meta: { requiresAuth: true, role: 'student' },
      },
      {
        path: 'student/labs/:labId/submit',
        name: 'student-lab-submit',
        component: StudentLabSubmitView,
        meta: { requiresAuth: true, role: 'student' },
      },
      {
        path: 'student/submissions',
        name: 'student-submissions',
        component: StudentSubmissionsView,
        meta: { requiresAuth: true, role: 'student' },
      },
      {
        path: 'teacher/labs',
        name: 'teacher-dashboard',
        component: TeacherDashboardView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'teacher/labs/create',
        name: 'teacher-lab-create',
        component: TeacherLabCreateView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'teacher/labs/:labId/edit',
        name: 'teacher-lab-edit',
        component: TeacherLabEditView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'teacher/labs/:labId/submissions',
        name: 'teacher-lab-submissions',
        component: TeacherLabSubmissionsView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'teacher/labs/:labId/analysis',
        name: 'teacher-lab-analysis',
        component: TeacherLabAnalysisView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'teacher/submissions/:submissionId/matches',
        name: 'teacher-submission-matches',
        component: TeacherSubmissionMatchesView,
        meta: { requiresAuth: true, role: 'teacher' },
      },
      {
        path: 'profile',
        name: 'profile',
        component: ProfileView,
        meta: { requiresAuth: true },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (!sessionRestored) {
    sessionRestored = true

    if (!authStore.user && authStore.accessToken) {
      await authStore.restoreSession()
    }

    if (authStore.isAuthenticated && !authStore.user) {
      try {
        await authStore.fetchMe()
      } catch {
        authStore.clearSession()
      }
    }
  }

  if (to.path === '/') {
    return authStore.isAuthenticated ? getRoleHome(authStore.user?.role) : '/login'
  }

  if (to.path === '/student') {
    return '/student/labs'
  }

  if (to.path === '/teacher') {
    return '/teacher/labs'
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    return getRoleHome(authStore.user?.role)
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      path: '/login',
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.requiresAuth && !authStore.user) {
    try {
      await authStore.fetchMe()
    } catch {
      authStore.clearSession()
      return {
        path: '/login',
        query: { redirect: to.fullPath },
      }
    }
  }

  if (to.meta.role && authStore.user?.role !== to.meta.role) {
    return getRoleHome(authStore.user?.role)
  }

  return true
})

export default router
