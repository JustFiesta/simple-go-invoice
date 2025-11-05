<template>
  <div>
    <!-- Page Header -->
    <v-row class="mb-6">
      <v-col cols="12" class="d-flex justify-space-between align-center">
        <div>
          <h1 class="text-h4 font-weight-bold text-primary mb-2">All Invoices</h1>
          <p class="text-body-1 text-secondary">
            Manage and track all your invoices
          </p>
        </div>
        <v-btn
          color="primary"
          size="large"
          prepend-icon="mdi-plus"
          :to="{ name: 'InvoiceCreate' }"
        >
          Create Invoice
        </v-btn>
      </v-col>
    </v-row>

    <!-- Filters Card -->
    <v-card class="mb-6 filter-card" elevation="0">
      <v-card-text class="pa-6">
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="filters.search"
              label="Search invoices"
              placeholder="Invoice number or client name"
              prepend-inner-icon="mdi-magnify"
              clearable
              bg-color="surface-light"
              variant="outlined"
              @update:model-value="debouncedSearch"
            ></v-text-field>
          </v-col>

          <v-col cols="6" md="2">
            <v-select
              v-model="filters.status"
              label="Status"
              :items="statusOptions"
              clearable
              bg-color="surface-light"
              variant="outlined"
              prepend-inner-icon="mdi-tag"
              @update:model-value="fetchInvoices"
            ></v-select>
          </v-col>

          <v-col cols="6" md="2">
            <v-select
              v-model="filters.sort"
              label="Sort By"
              :items="sortOptions"
              bg-color="surface-light"
              variant="outlined"
              prepend-inner-icon="mdi-sort"
              @update:model-value="fetchInvoices"
            ></v-select>
          </v-col>

          <v-col cols="6" md="2">
            <v-select
              v-model="filters.order"
              label="Order"
              :items="orderOptions"
              bg-color="surface-light"
              variant="outlined"
              prepend-inner-icon="mdi-sort-variant"
              @update:model-value="fetchInvoices"
            ></v-select>
          </v-col>

          <v-col cols="6" md="2">
            <v-select
              v-model="pagination.limit"
              label="Per Page"
              :items="[6, 9, 12, 24]"
              bg-color="surface-light"
              variant="outlined"
              prepend-inner-icon="mdi-view-grid"
              @update:model-value="handleLimitChange"
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Loading Progress -->
    <v-progress-linear
      v-if="loading"
      indeterminate
      color="primary"
      class="mb-4"
    ></v-progress-linear>

    <!-- Results Info -->
    <div v-if="!loading && invoices.length > 0" class="mb-4 d-flex justify-space-between align-center">
      <div class="text-body-2 text-secondary">
        Showing {{ ((pagination.page - 1) * pagination.limit) + 1 }} - 
        {{ Math.min(pagination.page * pagination.limit, pagination.total) }} 
        of {{ pagination.total }} invoices
      </div>
      <div class="text-body-2 text-secondary">
        Page {{ pagination.page }} of {{ pagination.total_pages }}
      </div>
    </div>

    <!-- Invoice Grid -->
    <v-row v-if="!loading && invoices.length > 0">
      <v-col
        v-for="invoice in invoices"
        :key="invoice.id"
        cols="12"
        sm="6"
        lg="4"
      >
        <InvoiceCard
          :invoice="invoice"
          @view="viewInvoice"
          @edit="editInvoice"
          @delete="confirmDelete"
          @download="downloadPDF"
        />
      </v-col>
    </v-row>

    <!-- Empty State -->
    <v-card
      v-if="!loading && invoices.length === 0"
      class="pa-12 text-center empty-state"
      elevation="0"
    >
      <v-icon
        :icon="filters.search || filters.status ? 'mdi-magnify-close' : 'mdi-file-document-outline'"
        size="80"
        color="surface-lighter"
      ></v-icon>
      <h3 class="text-h5 mt-6 mb-3">
        {{ filters.search || filters.status ? 'No matching invoices' : 'No invoices yet' }}
      </h3>
      <p class="text-body-1 text-secondary mb-6">
        {{ filters.search || filters.status 
          ? 'Try adjusting your filters or search terms' 
          : 'Create your first invoice to get started' 
        }}
      </p>
      <v-btn
        v-if="!filters.search && !filters.status"
        :to="{ name: 'InvoiceCreate' }"
        color="primary"
        size="large"
        prepend-icon="mdi-plus"
      >
        Create Invoice
      </v-btn>
      <v-btn
        v-else
        variant="outlined"
        size="large"
        prepend-icon="mdi-refresh"
        @click="clearFilters"
      >
        Clear Filters
      </v-btn>
    </v-card>

    <!-- Pagination -->
    <v-row v-if="pagination.total_pages > 1" class="mt-6">
      <v-col cols="12" class="d-flex justify-center">
        <v-pagination
          v-model="pagination.page"
          :length="pagination.total_pages"
          :total-visible="7"
          color="primary"
          @update:model-value="fetchInvoices"
        ></v-pagination>
      </v-col>
    </v-row>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog.show" max-width="450">
      <v-card class="dialog-card">
        <v-card-title class="text-h6 pa-6 pb-4">
          <v-icon icon="mdi-alert-circle" color="error" class="mr-2"></v-icon>
          Confirm Deletion
        </v-card-title>
        <v-card-text class="px-6 pb-2">
          Are you sure you want to delete invoice
          <strong class="text-primary">{{ deleteDialog.invoice?.invoice_number }}</strong>?
          <br><br>
          <span class="text-error">This action cannot be undone.</span>
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
import { ref, reactive, onMounted, inject } from 'vue'
import { useRouter } from 'vue-router'
import { invoiceService } from '@/services/api'
import InvoiceCard from '@/components/InvoiceCard.vue'

const router = useRouter()
const showNotification = inject('showNotification')

const loading = ref(false)
const invoices = ref([])

const filters = reactive({
  search: '',
  status: null,
  sort: 'created_at',
  order: 'desc'
})

const pagination = reactive({
  page: 1,
  limit: 9,
  total: 0,
  total_pages: 0
})

const deleteDialog = reactive({
  show: false,
  invoice: null,
  loading: false
})

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Sent', value: 'sent' },
  { title: 'Paid', value: 'paid' }
]

const sortOptions = [
  { title: 'Created Date', value: 'created_at' },
  { title: 'Updated Date', value: 'updated_at' },
  { title: 'Invoice Number', value: 'invoice_number' },
  { title: 'Client Name', value: 'client_name' },
  { title: 'Amount', value: 'amount' },
  { title: 'Status', value: 'status' },
  { title: 'Issue Date', value: 'issue_date' },
  { title: 'Due Date', value: 'due_date' }
]

const orderOptions = [
  { title: 'Descending', value: 'desc' },
  { title: 'Ascending', value: 'asc' }
]

let searchTimeout = null
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.page = 1
    fetchInvoices()
  }, 500)
}

const fetchInvoices = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      sort: filters.sort,
      order: filters.order.toUpperCase()
    }

    if (filters.search) params.search = filters.search
    if (filters.status) params.status = filters.status

    const response = await invoiceService.getInvoices(params)
    invoices.value = response.data.data || []
    
    if (response.data.meta) {
      pagination.total = response.data.meta.total
      pagination.total_pages = response.data.meta.total_pages
    }
  } catch (error) {
    showNotification(error.message, 'error')
  } finally {
    loading.value = false
  }
}

const handleLimitChange = () => {
  pagination.page = 1
  fetchInvoices()
}

const clearFilters = () => {
  filters.search = ''
  filters.status = null
  pagination.page = 1
  fetchInvoices()
}

const viewInvoice = (id) => {
  router.push({ name: 'InvoiceDetail', params: { id } })
}

const editInvoice = (id) => {
  router.push({ name: 'InvoiceEdit', params: { id } })
}

const confirmDelete = (invoice) => {
  deleteDialog.invoice = invoice
  deleteDialog.show = true
}

const deleteInvoice = async () => {
  deleteDialog.loading = true
  try {
    await invoiceService.deleteInvoice(deleteDialog.invoice.id)
    showNotification('Invoice deleted successfully', 'success')
    deleteDialog.show = false
    
    if (invoices.value.length === 1 && pagination.page > 1) {
      pagination.page--
    }
    
    fetchInvoices()
  } catch (error) {
    showNotification(error.message, 'error')
  } finally {
    deleteDialog.loading = false
  }
}

const downloadPDF = async (invoice) => {
  try {
    const response = await invoiceService.getInvoicePDF(invoice.id)
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `invoice-${invoice.invoice_number}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
    showNotification('PDF downloaded successfully', 'success')
  } catch (error) {
    showNotification(error.message, 'error')
  }
}

onMounted(() => {
  fetchInvoices()
})
</script>

<style scoped>
.filter-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
}

.empty-state {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
}

.dialog-card {
  background: #2A2A2A;
  border: 1px solid #404040;
}
</style>