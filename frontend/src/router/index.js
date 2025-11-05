import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: {
      title: 'Dashboard'
    }
  },
  {
    path: '/invoices',
    name: 'InvoiceList',
    component: () => import('@/views/InvoiceListView.vue'),
    meta: {
      title: 'All Invoices'
    }
  },
  {
    path: '/invoices/new',
    name: 'InvoiceCreate',
    component: () => import('@/views/InvoiceFormView.vue'),
    meta: {
      title: 'Create Invoice'
    }
  },
  {
    path: '/invoices/:id',
    name: 'InvoiceDetail',
    component: () => import('@/views/InvoiceDetailView.vue'),
    meta: {
      title: 'Invoice Details'
    }
  },
  {
    path: '/invoices/:id/edit',
    name: 'InvoiceEdit',
    component: () => import('@/views/InvoiceFormView.vue'),
    meta: {
      title: 'Edit Invoice'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title 
    ? `${to.meta.title} | Invoice Manager` 
    : 'Invoice Manager'
  next()
})

export default router