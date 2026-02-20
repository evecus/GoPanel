import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth.js'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: () => import('../views/Login.vue'), meta: { public: true } },
    {
      path: '/', component: () => import('../views/Layout.vue'),
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', component: () => import('../views/Dashboard.vue') },
        { path: 'processes', component: () => import('../views/Processes.vue') },
        { path: 'docker',    component: () => import('../views/Docker.vue') },
        { path: 'services',  component: () => import('../views/Services.vue') },
        { path: 'settings',  component: () => import('../views/Settings.vue') },
      ]
    }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) return '/login'
})

export default router
