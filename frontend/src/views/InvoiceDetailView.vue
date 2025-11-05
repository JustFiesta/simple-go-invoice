<template>
  <div>
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
      <p class="text-body-1 mt-4">Loading invoice...</p>
    </div>

    <!-- Invoice Content -->
    <div v-else-if="invoice">
      <!-- Page Header -->
      <v-row class="mb-6">
        <v-col cols="12" class="d-flex justify-space-between align-center">
          <div>
            <h1 class="text-h4 font-weight-bold text-primary mb-2">
              Invoice {{ invoice.invoice_number }}
            </h1>
            <p class="text-body-1 text-secondary">
              Created {{ formatDate(invoice.created_at) }}
            </p>
          </div>
          <InvoiceStatusBadge :status="invoice.status" size="large" />
        </v-col>
      </v-row>

      <v-row>
        <!-- Main Content -->
        <v-col cols="12" lg="8">
          <!-- Invoice Information Card -->
          <v-card elevation="0" class="detail-card mb-6">
            <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <v-icon icon="mdi-information" color="primary" class="mr-3"></v-icon>
                <span class="text-h6 font-weight-bold">Invoice Information</span>
              </div>
              <div class="d-flex gap-2">
                <v-btn
                  variant="outlined"
                  size="small"
                  prepend-icon="mdi-download"
                  color="accent"
                  @click="downloadPDF"
                  :loading="downloadingPdf"
                >
                  PDF
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
                <v-col cols="12" sm="6" md="4">
                  <div class="info-field mb-4">
                    <div class="text-caption text-secondary mb-1">INVOICE NUMBER</div>
                    <div class="text-h6">{{ invoice.invoice_number }}</div>
                  </div>
                </v-col>

                <v-col cols="12" sm="6" md="4">
                  <div class="info-field mb-4">
                    <div class="text-caption text-secondary mb-1">CLIENT NAME</div>
                    <div class="text-h6">{{ invoice.client_name }}</div>
                  </div>
                </v-col>

                <v-col cols="12" sm="6" md="4">
                  <div class="info-field mb-4">
                    <div class="text-caption text-secondary mb-1">STATUS</div>
                    <InvoiceStatusBadge :status="invoice.status" />
                  </div>
                </v-col>

                <v-col cols="12" sm="6" md="4">
                  <div class="info-field mb-4">
                    <div class="text-caption text-secondary mb-1">ISSUE DATE</div>
                    <div class="text-body-1 d-flex align-center">
                      <v-icon icon="mdi-calendar" size="18" color="secondary" class="mr-2"></v-icon>
                      {{ formatDate(invoice.issue_date) }}
                    </div>
                  </div>
                </v-col>

                <v-col cols="12" sm="6" md="4">
                  <div class="info-field mb-4">
                    <div class="text-caption text-secondary mb-1">DUE DATE</div>
                    <div
                      class="text-body-1 d-flex align-center"
                      :class="{ 'text-error': isOverdue }"
                    >
                      <v-icon
                        icon="mdi-calendar-clock"
                        size="18"
                        :color="isOverdue ? 'error' : 'secondary'"
                        class="mr-2"
                      ></v-icon>
                      {{ formatDate(invoice.due_date) }}
                      <v-chip
                        v-if="isOverdue && invoice.status !== 'paid'"
                        color="error"
                        size="x-small"
                        class="ml-2"
                      >
                        OVERDUE
                      </v-chip>
                    </div>
                  </div>
                </v-col>

                <v-col cols="12" sm="6" md="4">
                  <div class="info-field">
                    <div class="text-caption text-secondary mb-1">TOTAL AMOUNT</div>
                    <div class="text-h5 text-primary font-weight-bold">
                      {{ formatCurrency(invoice.amount) }}
                    </div>
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>

          <!-- Invoice Items -->
          <v-card elevation="0" class="detail-card mb-6">
            <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <v-icon icon="mdi-format-list-bulleted" color="primary" class="mr-3"></v-icon>
                <span class="text-h6 font-weight-bold">Invoice Items</span>
                <v-chip v-if="items.length > 0" size="small" class="ml-3" color="secondary">
                  {{ items.length }}
                </v-chip>
              </div>
              <v-btn
                v-if="items.length > 0"
                size="small"
                variant="outlined"
                color="secondary"
                prepend-icon="mdi-plus"
                :to="{ name: 'InvoiceEdit', params: { id: invoice.id } }"
              >
                Manage Items
              </v-btn>
            </v-card-title>

            <v-divider :thickness="1" color="surface-lighter"></v-divider>

            <v-card-text v-if="items.length > 0" class="pa-0">
              <div class="table-wrapper">
                <v-table class="custom-table">
                  <thead>
                    <tr>
                      <th style="width: 50px">#</th>
                      <th>Description</th>
                      <th class="text-center" style="width: 100px">Qty</th>
                      <th class="text-right" style="width: 130px">Unit Price</th>
                      <th class="text-center" style="width: 100px">VAT %</th>
                      <th class="text-right" style="width: 130px">Net Total</th>
                      <th class="text-right" style="width: 150px">Gross Total</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item, index) in items" :key="item.id">
                      <td class="text-secondary">{{ index + 1 }}</td>
                      <td>{{ item.description }}</td>
                      <td class="text-center">{{ item.quantity }}</td>
                      <td class="text-right">{{ formatCurrency(item.unit_price) }}</td>
                      <td class="text-center">{{ item.vat_rate }}%</td>
                      <td class="text-right font-weight-medium">
                        {{ formatCurrency(item.quantity * item.unit_price) }}
                      </td>
                      <td class="text-right font-weight-bold text-primary">
                        {{ formatCurrency(calculateItemGross(item)) }}
                      </td>
                    </tr>
                  </tbody>
                  <tfoot>
                    <tr class="total-row">
                      <td colspan="6" class="text-right font-weight-bold text-h6 pa-4">
                        Total:
                      </td>
                      <td class="text-right font-weight-bold text-h5 text-primary pa-4">
                        {{ formatCurrency(calculateTotal()) }}
                      </td>
                    </tr>
                  </tfoot>
                </v-table>
              </div>
            </v-card-text>

            <!-- Empty State -->
            <v-card-text v-else class="text-center py-12">
              <v-icon icon="mdi-package-variant" size="80" color="surface-lighter"></v-icon>
              <h4 class="text-h6 mt-6 mb-2">No items added yet</h4>
              <p class="text-body-2 text-secondary mb-6">
                Add items to this invoice to complete it
              </p>
              <v-btn
                color="secondary"
                prepend-icon="mdi-plus"
                size="large"
                :to="{ name: 'InvoiceEdit', params: { id: invoice.id } }"
              >
                Add Items
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Sidebar -->
        <v-col cols="12" lg="4">
          <InvoiceSummary :items="items" class="mb-6" />

          <!-- Actions Card -->
          <v-card elevation="0" class="actions-card">
            <v-card-title class="pa-4 text-h6 font-weight-bold">
              Actions
            </v-card-title>
            <v-divider :thickness="1" color="surface-lighter"></v-divider>
            <v-card-text class="pa-4">
              <v-btn
                block
                color="primary"
                prepend-icon="mdi-pencil"
                class="mb-3"
                :to="{ name: 'InvoiceEdit', params: { id: invoice.id } }"
              >
                Edit Invoice
              </v-btn>

              <v-btn
                block
                color="accent"
                variant="outlined"
                prepend-icon="mdi-download"
                class="mb-3"
                :loading="downloadingPdf"
                @click="downloadPDF"
              >
                Download PDF
              </v-btn>

              <v-divider class="my-3"></v-divider>

              <v-btn
                block
                color="error"
                variant="outlined"
                prepend-icon="mdi-delete"
                @click="confirmDelete"
              >
                Delete Invoice
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Bottom Actions -->
      <v-row class="mt-6">
        <v-col cols="12">
          <v-btn
            variant="outlined"
            size="large"
            prepend-icon="mdi-arrow-left"
            @click="$router.push({ name: 'InvoiceList' })"
          >
            Back to List
          </v-btn>
        </v-col>
      </v-row>
    </div>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog.show" max-width="450">
      <v-card class="dialog-card">
        <v-card-title class="text-h6 pa-6">Confirm Deletion</v-card-title>
        <v-card-text class="px-6 pb-2">
          Are you sure you want to delete invoice
          <strong class="text-primary">{{ invoice?.invoice_number }}</strong>?
          This action cannot be undone.
        </v-card-text>
        <v-card-actions class="pa-6 pt-4">
          <v-spacer></v-spacer>
          <v-btn variant="outlined" @click="deleteDialog.show = false">Cancel</v-btn>
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
import InvoiceStatusBadge from '@/components/InvoiceStatusBadge.vue'
import InvoiceSummary from '@/components/InvoiceSummary.vue'

const router = useRouter()
const route = useRoute()
const showNotification = inject('showNotification')

const loading = ref(true)
const downloadingPdf = ref(false)
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

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const calculateItemGross = (item) => {
  const net = (item.quantity || 0) * (item.unit_price || 0)
  const vat = net * ((item.vat_rate || 0) / 100)
  return net + vat
}

const calculateTotal = () => {
  return items.value.reduce((sum, item) => sum + calculateItemGross(item), 0)
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
  downloadingPdf.value = true
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
  } finally {
    downloadingPdf.value = false
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

<style scoped>
.detail-card,
.actions-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
}

.info-field {
  padding: 8px 0;
}

.table-wrapper {
  overflow-x: auto;
}

.custom-table {
  background: transparent;
}

.custom-table thead th {
  background-color: rgba(75, 240, 240, 0.1);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  padding: 16px 12px;
  border-bottom: 2px solid #404040;
}

.custom-table tbody td {
  padding: 12px;
  border-bottom: 1px solid #353535;
}

.custom-table tfoot .total-row {
  background-color: rgba(255, 123, 62, 0.1);
  border-top: 2px solid #404040;
}

.dialog-card {
  background: #2A2A2A;
  border: 1px solid #404040;
}

.gap-2 {
  gap: 8px;
}
</style>