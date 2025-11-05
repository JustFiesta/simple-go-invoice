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
  config => {
    console.log(`[API Request] ${config.method.toUpperCase()} ${config.url}`, config.params)
    return config
  },
  error => {
    console.error('[API Request Error]', error)
    return Promise.reject(error)
  }
)

// Response interceptor with better error handling
api.interceptors.response.use(
  response => {
    console.log(`[API Response] ${response.config.url}`, response.status)
    return response
  },
  error => {
    console.error('[API Error]', error)
    
    let message = 'An error occurred'
    
    // Handle connection errors
    if (error.code === 'ECONNREFUSED' || error.code === 'ERR_NETWORK') {
      message = 'Cannot connect to backend server. Please ensure the backend is running on http://localhost:8080'
    } 
    // Handle timeout
    else if (error.code === 'ECONNABORTED') {
      message = 'Request timeout. Please try again.'
    }
    // Handle HTTP errors
    else if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      if (data?.error?.message) {
        message = data.error.message
      } else {
        switch (status) {
          case 400:
            message = 'Bad request. Please check your input.'
            break
          case 401:
            message = 'Unauthorized. Please log in.'
            break
          case 403:
            message = 'Forbidden. You do not have permission.'
            break
          case 404:
            message = 'Resource not found.'
            break
          case 409:
            message = data?.error?.message || 'Conflict. Resource already exists.'
            break
          case 422:
            message = data?.error?.message || 'Validation failed.'
            break
          case 500:
            message = 'Server error. Please try again later.'
            break
          default:
            message = `Server error: ${status}`
        }
      }
    }
    // Handle other errors
    else if (error.message) {
      message = error.message
    }
    
    return Promise.reject(new Error(message))
  }
)

export const invoiceService = {
  // Get all invoices with filters
  getInvoices(params = {}) {
    // Ensure page and limit have defaults
    const queryParams = {
      page: params.page || 1,
      limit: params.limit || 10,
      ...params
    }
    return api.get('/invoices', { params: queryParams })
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