import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'InvoiceList',
    component: () => import('@/views/InvoiceList.vue')
  },
  {
    path: '/invoices/new',
    name: 'InvoiceCreate',
    component: () => import('@/views/InvoiceForm.vue')
  },
  {
    path: '/invoices/:id',
    name: 'InvoiceDetail',
    component: () => import('@/views/InvoiceDetail.vue')
  },
  {
    path: '/invoices/:id/edit',
    name: 'InvoiceEdit',
    component: () => import('@/views/InvoiceForm.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router