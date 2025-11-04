<template>
  <div v-if="loading" class="text-center py-8">
    <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
  </div>

  <div v-else-if="invoice">
    <v-row class="mb-4">
      <v-col cols="12" class="d-flex justify-space-between align-center">
        <h1 class="text-h4 font-weight-bold text-primary">
          Invoice {{ invoice.invoice_number }}
        </h1>
        <v-chip
          :color="getStatusColor(invoice.status)"
          size="large"
          label
        >
          {{ invoice.status.toUpperCase() }}
        </v-chip>
      </v-col>
    </v-row>

    <v-card elevation="0" class="mb-6 detail-card">
      <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
        <div class="d-flex align-center">
          <v-icon icon="mdi-information" color="primary" class="mr-3"></v-icon>
          <span class="text-h6 font-weight-bold">Invoice Information</span>
        </div>
        <div>
          <v-btn
            variant="outlined"
            size="small"
            prepend-icon="mdi-download"
            color="accent"
            class="mr-2"
            @click="downloadPDF"
          >
            Download PDF
          </v-btn>
          <v-btn
            color="primary"
            size="small"
            prepend-icon="mdi-pencil"
            :to="{ name: 'InvoiceEdit', params: { id: invoice.id } }"
          >
            Edit
          </v-btn>
        </div>
      </v-card-title>

      <v-divider :thickness="1" color="surface-lighter"></v-divider>

      <v-card-text class="pa-6">
        <v-row>
          <v-col cols="12" md="6">
            <div class="info-section mb-6">
              <div class="text-caption text-secondary mb-2">CLIENT NAME</div>
              <div class="text-h6">{{ invoice.client_name }}</div>
            </div>

            <div class="info-section mb-6">
              <div class="text-caption text-secondary mb-2">INVOICE NUMBER</div>
              <div class="text-body-1 font-weight-medium">{{ invoice.invoice_number }}</div>
            </div>

            <div class="info-section">
              <div class="text-caption text-secondary mb-2">STATUS</div>
              <v-chip
                :color="getStatusColor(invoice.status)"
                size="default"
                variant="flat"
              >
                {{ invoice.status.toUpperCase() }}
              </v-chip>
            </div>
          </v-col>

          <v-col cols="12" md="6">
            <div class="info-section mb-6">
              <div class="text-caption text-secondary mb-2">AMOUNT</div>
              <div class="text-h4 text-primary font-weight-bold">
                {{ formatCurrency(invoice.amount) }}
              </div>
            </div>

            <div class="info-section mb-6">
              <div class="text-caption text-secondary mb-2">ISSUE DATE</div>
              <div class="text-body-1 d-flex align-center">
                <v-icon icon="mdi-calendar" size="20" color="secondary" class="mr-2"></v-icon>
                {{ formatDate(invoice.issue_date) }}
              </div>
            </div>

            <div class="info-section">
              <div class="text-caption text-secondary mb-2">DUE DATE</div>
              <div
                class="text-body-1 d-flex align-center"
                :class="{ 'text-error': isOverdue }"
              >
                <v-icon icon="mdi-calendar-clock" size="20" :color="isOverdue ? 'error' : 'secondary'" class="mr-2"></v-icon>
                {{ formatDate(invoice.due_date) }}
                <v-chip
                  v-if="isOverdue && invoice.status !== 'paid'"
                  color="error"
                  size="small"
                  class="ml-3"
                >
                  OVERDUE
                </v-chip>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-card v-if="items.length > 0" elevation="0" class="mb-6 detail-card">
      <v-card-title class="pa-6 pb-4">
        <v-icon icon="mdi-format-list-bulleted" color="primary" class="mr-3"></v-icon>
        <span class="text-h6 font-weight-bold">Invoice Items</span>
      </v-card-title>

      <v-divider :thickness="1" color="surface-lighter"></v-divider>

      <v-card-text class="pa-0">
        <div class="table-wrapper">
          <v-table class="custom-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Description</th>
                <th class="text-center">Quantity</th>
                <th class="text-right">Unit Price</th>
                <th class="text-center">VAT Rate</th>
                <th class="text-right">Net Total</th>
                <th class="text-right">Gross Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in items" :key="item.id">
                <td class="text-secondary">{{ index + 1 }}</td>
                <td>{{ item.description }}</td>
                <td class="text-center">{{ item.quantity }}</td>
                <td class="text-right">{{ formatCurrency(item.unit_price) }}</td>
                <td class="text-center">{{ item.vat_rate }}%</td>
                <td class="text-right">
                  {{ formatCurrency(item.quantity * item.unit_price) }}
                </td>
                <td class="text-right font-weight-bold text-primary">
                  {{ formatCurrency(item.quantity * item.unit_price * (1 + item.vat_rate / 100)) }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="total-row">
                <td colspan="6" class="text-right font-weight-bold text-h6 pa-4">Total:</td>
                <td class="text-right font-weight-bold text-h5 text-primary pa-4">
                  {{ formatCurrency(calculateTotal()) }}
                </td>
              </tr>
            </tfoot>
          </v-table>
        </div>
      </v-card-text>
    </v-card>

    <v-card v-else elevation="0" class="mb-6 empty-card">
      <v-card-text class="text-center py-12">
        <v-icon icon="mdi-package-variant" size="80" color="surface-lighter"></v-icon>
        <div class="text-h6 mt-6 mb-2">No items added yet</div>
        <p class="text-body-2 text-secondary mb-6">Add items to this invoice to get started</p>
        <v-btn
          color="secondary"
          class="mt-2"
          prepend-icon="mdi-plus"
          size="large"
          :to="{ name: 'InvoiceEdit', params: { id: invoice.id } }"
        >
          Add Items
        </v-btn>
      </v-card-text>
    </v-card>

    <div class="d-flex justify-space-between">
      <v-btn
        variant="outlined"
        prepend-icon="mdi-arrow-left"
        @click="$router.push({ name: 'InvoiceList' })"
      >
        Back to List
      </v-btn>

      <v-btn
        color="error"
        variant="outlined"
        prepend-icon="mdi-delete"
        @click="confirmDelete"
      >
        Delete Invoice
      </v-btn>
    </div>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog.show" max-width="400">
      <v-card>
        <v-card-title class="text-h6">Confirm Deletion</v-card-title>
        <v-card-text>
          Are you sure you want to delete invoice
          <strong>{{ invoice.invoice_number }}</strong>?
          This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="deleteDialog.show = false">Cancel</v-btn>
          <v-btn
            color="error"
            :loading="deleteDialog.loading"
            @click="deleteInvoice"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { invoiceService } from '@/services/api'

const router = useRouter()
const route = useRoute()
const showNotification = inject('showNotification')

const loading = ref(true)
const invoice = ref(null)
const items = ref([])

const deleteDialog = reactive({
  show: false,
  loading: false
})

const isOverdue = computed(() => {
  if (!invoice.value) return false
  return new Date(invoice.value.due_date) < new Date() &&
    new Date(invoice.value.due_date).toDateString() !== new Date().toDateString()
})

const getStatusColor = (status) => {
  const colors = {
    draft: 'grey',
    sent: 'warning',
    paid: 'success'
  }
  return colors[status] || 'grey'
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const calculateTotal = () => {
  return items.value.reduce((sum, item) => {
    return sum + (item.quantity * item.unit_price * (1 + item.vat_rate / 100))
  }, 0)
}

const loadInvoice = async () => {
  loading.value = true
  try {
    const [invoiceRes, itemsRes] = await Promise.all([
      invoiceService.getInvoice(route.params.id),
      invoiceService.getInvoiceItems(route.params.id)
    ])

    invoice.value = invoiceRes.data.data
    items.value = itemsRes.data.data || []
  } catch (error) {
    showNotification(error.message, 'error')
    router.push({ name: 'InvoiceList' })
  } finally {
    loading.value = false
  }
}

const downloadPDF = async () => {
  try {
    const response = await invoiceService.getInvoicePDF(invoice.value.id)
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `invoice-${invoice.value.invoice_number}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
    showNotification('PDF downloaded successfully', 'success')
  } catch (error) {
    showNotification(error.message, 'error')
  }
}

const confirmDelete = () => {
  deleteDialog.show = true
}

const deleteInvoice = async () => {
  deleteDialog.loading = true
  try {
    await invoiceService.deleteInvoice(invoice.value.id)
    showNotification('Invoice deleted successfully', 'success')
    router.push({ name: 'InvoiceList' })
  } catch (error) {
    showNotification(error.message, 'error')
  } finally {
    deleteDialog.loading = false
  }
}

onMounted(() => {
  loadInvoice()
})
</script>