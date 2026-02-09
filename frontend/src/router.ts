import { createRouter, createWebHistory } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
  }
}

const CameraView = () => import('./views/CameraView.vue')
const TimelapsesView = () => import('./views/TimelapsesView.vue')
const NotFound = () => import('./views/NotFound.vue')

const routes = [
  {
    path: '/',
    redirect: '/camera'
  },
  {
    path: '/camera',
    name: 'Camera',
    component: CameraView,
    meta: { title: 'Camera' }
  },
  {
    path: '/timelapses/:filename?',
    name: 'Timelapses',
    component: TimelapsesView,
    meta: { title: 'Timelapses' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: { title: 'Not Found' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - Printer` : 'Printer'
  next()
})

export default router
