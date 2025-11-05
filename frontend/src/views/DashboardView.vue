<template>
  <div>
    <!-- Page Header -->
    <v-row class="mb-6">
      <v-col cols="12">
        <h1 class="text-h4 font-weight-bold text-primary mb-2">Dashboard</h1>
        <p class="text-body-1 text-secondary">Overview of your invoicing activity</p>
      </v-col>
    </v-row>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <!-- Dashboard Content -->
    <div v-else>
      <!-- Summary Cards -->
      <v-row class="mb-6">
        <v-col cols="12" sm="6" lg="3">
          <v-card elevation="2" class="stat-card">
            <v-card-text class="pa-6">
              <div class="d-flex justify-space-between align-center mb-3">
                <v-icon icon="mdi-file-document-multiple" size="40" color="primary"></v-icon>
                <v-chip color="primary" size="small">Total</v-chip>
              </div>
              <div class="text-h3 font-weight-bold text-primary mb-1">
                {{ stats.total }}
              </div>
              <div class="text-body-2 text-secondary">Total Invoices</div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" sm="6" lg="3">
          <v-card elevation="2" class="stat-card">
            <v-card-text class="pa-6">
              <div class="d-flex justify-space-between align-center mb-3">
                <v-icon icon="mdi-file-document-edit" size="40" color="grey"></v-icon>
                <v-chip color="grey" size="small">Draft</v-chip>
              </div>
              <div class="text-h3 font-weight-bold mb-1">{{ stats.draft }}</div>
              <div class="text-body-2 text-secondary">Draft Invoices</div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" sm="6" lg="3">
          <v-card elevation="2" class="stat-card">
            <v-card-text class="pa-6">
              <div class="d-flex justify-space-between align-center mb-3">
                <v-icon icon="mdi-send" size="40" color="warning"></v-icon>
                <v-chip color="warning" size="small">Sent</v-chip>
              </div>
              <div class="text-h3 font-weight-bold mb-1">{{ stats.sent }}</div>
              <div class="text-body-2 text-secondary">Sent Invoices</div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" sm="6" lg="3">
          <v-card elevation="2" class="stat-card">
            <v-card-text class="pa-6">
              <div class="d-flex justify-space-between align-center mb-3">
                <v-icon icon="mdi-check-circle" size="40" color="success"></v-icon>
                <v-chip color="success" size="small">Paid</v-chip>
              </div>
              <div class="text-h3 font-weight-bold mb-1">{{ stats.paid }}</div>
              <div class="text-body-2 text-secondary">Paid Invoices</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Financial Summary -->
      <v-row class="mb-6">
        <v-col cols="12" md="4">
          <v-card elevation="2" class="financial-card">
            <v-card-text class="pa-6">
              <div class="text-caption text-secondary mb-2">TOTAL REVENUE</div>
              <div class="text-h4 font-weight-bold text-primary mb-2">
                {{ formatCurrency(financials.total) }}
              </div>
              <div class="text-caption">All invoices combined</div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" md="4">
          <v-card elevation="2" class="financial-card">
            <v-card-text class="pa-6">
              <div class="text-caption text-secondary mb-2">PENDING PAYMENT</div>
              <div class="text-h4 font-weight-bold text-warning mb-2">
                {{ formatCurrency(financials.pending) }}
              </div>
              <div class="text-caption">Sent but not paid</div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" md="4">
          <v-card elevation="2" class="financial-card">
            <v-card-text class="pa-6">
              <div class="text-caption text-secondary mb-2">PAID AMOUNT</div>
              <div class="text-h4 font-weight-bold text-success mb-2">
                {{ formatCurrency(financials.paid) }}
              </div>
              <div class="text-caption">Successfully collected</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Recent Invoices -->
      <v-row>
        <v-col cols="12">
          <v-card elevation="0" class="recent-card">
            <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <v-icon icon="mdi-clock-outline" color="primary" class="mr-3"></v-icon>
                <span class="text-h6 font-weight-bold">Recent Invoices</span>
              </div>
              <v-btn
                variant="outlined"
                color="secondary"
                size="small"
                :to="{ name: 'InvoiceList' }"
              >
                View All
              </v-btn>
            </v-card-title>

            <v-divider :thickness="1" color="surface-lighter"></v-divider>

            <v-card-text class="pa-0">
              <v-list v-if="recentInvoices.length > 0" class="bg-transparent">
                <v-list-item
                  v-for="invoice in recentInvoices"
                  :key="invoice.id"
                  :to="{ name: 'InvoiceDetail', params: { id: invoice.id } }"
                  class="invoice-list-item"
                >
                  <template #prepend>
                    <v-icon icon="mdi-file-document" color="primary"></v-icon>
                  </template>

                  <v-list-item-title class="font-weight-medium">
                    {{ invoice.invoice_number }}
                  </v-list-item-title>

                  <v-list-item-subtitle>
                    {{ invoice.client_name }}
                  </v-list-item-subtitle>

                  <template #append>
                    <div class="d-flex align-center">
                      <div class="text-right mr-4">
                        <div class="font-weight-bold text-primary">
                          {{ formatCurrency(invoice.amount) }}
                        </div>
                        <div class="text-caption text-secondary">
                          {{ formatDate(invoice.issue_date) }}
                        </div>
                      </div>
                      <InvoiceStatusBadge :status="invoice.status" size="small" />
                    </div>
                  </template>
                </v-list-item>
              </v-list>

              <div v-else class="text-center py-12">
                <v-icon icon="mdi-file-document-outline" size="64" color="surface-lighter"></v-icon>
                <h4 class="text-h6 mt-4 mb-2">No invoices yet</h4>
                <p class="text-body-2 text-secondary mb-4">Create your first invoice to get started</p>
                <v-btn
                  color="primary"
                  prepend-icon="mdi-plus"
                  :to="{ name: 'InvoiceCreate' }"
                >
                  Create Invoice
                </v-btn>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, inject } from 'vue'
import { invoiceService } from '@/services/api'
import InvoiceStatusBadge from '@/components/InvoiceStatusBadge.vue'

const showNotification = inject('showNotification')

const loading = ref(true)
const recentInvoices = ref([])

const stats = reactive({
  total: 0,
  draft: 0,
  sent: 0,
  paid: 0
})

const financials = reactive({
  total: 0,
  pending: 0,
  paid: 0
})

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('pl-PL', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const loadDashboardData = async () => {
  loading.value = true
  try {
    // Fetch invoices with pagination info
    const response = await invoiceService.getInvoices({ 
      page: 1, 
      limit: 100 
    })

    const allInvoices = response.data.data || []
    const meta = response.data.meta

    // Calculate stats from fetched invoices
    stats.total = meta?.total || allInvoices.length
    stats.draft = allInvoices.filter(inv => inv.status === 'draft').length
    stats.sent = allInvoices.filter(inv => inv.status === 'sent').length
    stats.paid = allInvoices.filter(inv => inv.status === 'paid').length

    // Calculate financials
    financials.total = allInvoices.reduce((sum, inv) => sum + (inv.amount || 0), 0)
    financials.pending = allInvoices
      .filter(inv => inv.status === 'sent')
      .reduce((sum, inv) => sum + (inv.amount || 0), 0)
    financials.paid = allInvoices
      .filter(inv => inv.status === 'paid')
      .reduce((sum, inv) => sum + (inv.amount || 0), 0)

    // Get recent invoices (first 5)
    recentInvoices.value = allInvoices.slice(0, 5)
  } catch (error) {
    showNotification(error.message || 'Failed to load dashboard data', 'error')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.stat-card,
.financial-card,
.recent-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover,
.financial-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(255, 123, 62, 0.2);
}

.invoice-list-item {
  border-bottom: 1px solid #353535;
  transition: background-color 0.2s ease;
}

.invoice-list-item:hover {
  background-color: rgba(75, 240, 240, 0.05);
}

.invoice-list-item:last-child {
  border-bottom: none;
}
</style>