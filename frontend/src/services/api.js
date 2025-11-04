import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  config => config,
  error => Promise.reject(error)
)

// Response interceptor
api.interceptors.response.use(
  response => response,
  error => {
    const message = error.response?.data?.error?.message || error.message || 'An error occurred'
    return Promise.reject(new Error(message))
  }
)

export const invoiceService = {
  // Get all invoices with filters
  getInvoices(params = {}) {
    return api.get('/invoices', { params })
  },

  // Get single invoice
  getInvoice(id) {
    return api.get(`/invoices/${id}`)
  },

  // Create invoice
  createInvoice(data) {
    return api.post('/invoices', data)
  },

  // Update invoice
  updateInvoice(id, data) {
    return api.put(`/invoices/${id}`, data)
  },

  // Delete invoice
  deleteInvoice(id) {
    return api.delete(`/invoices/${id}`)
  },

  // Get invoice PDF
  getInvoicePDF(id) {
    return api.get(`/invoices/${id}/pdf`, {
      responseType: 'blob'
    })
  },

  // Invoice items
  getInvoiceItems(invoiceId) {
    return api.get(`/invoices/${invoiceId}/items`)
  },

  addInvoiceItem(invoiceId, data) {
    return api.post(`/invoices/${invoiceId}/items`, data)
  },

  updateInvoiceItem(invoiceId, itemId, data) {
    return api.put(`/invoices/${invoiceId}/items/${itemId}`, data)
  },

  deleteInvoiceItem(invoiceId, itemId) {
    return api.delete(`/invoices/${invoiceId}/items/${itemId}`)
  }
}

export const healthService = {
  check() {
    return api.get('/health')
  }
}

export default api