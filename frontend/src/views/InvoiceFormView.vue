<template>
  <div>
    <!-- Page Header -->
    <v-row class="mb-6">
      <v-col cols="12" class="d-flex justify-space-between align-center">
        <div>
          <h1 class="text-h4 font-weight-bold text-primary mb-2">
            {{ isEditMode ? 'Edit Invoice' : 'Create New Invoice' }}
          </h1>
          <p class="text-body-1 text-secondary">
            {{ isEditMode ? 'Update invoice details and manage items' : 'Fill in the details to create a new invoice' }}
          </p>
        </div>
        <v-chip
          v-if="isEditMode && invoice"
          :color="getStatusColor(invoice.status)"
          size="large"
          label
        >
          {{ invoice.status.toUpperCase() }}
        </v-chip>
      </v-col>
    </v-row>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
      <p class="text-body-1 mt-4">Loading invoice...</p>
    </div>

    <!-- Main Form -->
    <div v-else>
      <InvoiceForm
        :key="invoiceId"
        ref="invoiceFormRef"
        :invoice="invoice"
        :items="items"
        :is-edit-mode="isEditMode"
        :invoice-id="invoiceId"
        @add-item="handleAddItem"
        @update-item="handleUpdateItem"
        @delete-item="handleDeleteItem"
      />

      <!-- Action Buttons -->
      <v-row class="mt-6">
        <v-col cols="12" class="d-flex justify-space-between">
          <v-btn
            variant="outlined"
            size="large"
            prepend-icon="mdi-arrow-left"
            @click="handleCancel"
          >
            Cancel
          </v-btn>

          <div class="d-flex gap-2">
            <v-btn
              v-if="isEditMode && invoice"
              variant="outlined"
              size="large"
              color="accent"
              prepend-icon="mdi-eye"
              :to="{ name: 'InvoiceDetail', params: { id: invoiceId } }"
            >
              View Details
            </v-btn>

            <v-btn
              color="primary"
              size="large"
              prepend-icon="mdi-content-save"
              :loading="saving"
              @click="handleSubmit"
            >
              {{ isEditMode ? 'Update Invoice' : 'Create Invoice' }}
            </v-btn>
          </div>
        </v-col>
      </v-row>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { invoiceService } from '@/services/api'
import InvoiceForm from '@/components/InvoiceForm.vue'

const router = useRouter()
const route = useRoute()
const showNotification = inject('showNotification')

const invoiceFormRef = ref(null)
const loading = ref(false)
const saving = ref(false)
const invoice = ref(null)
const items = ref([])

const isEditMode = computed(() => route.name === 'InvoiceEdit')
const invoiceId = computed(() => route.params.id)

const getStatusColor = (status) => {
  const colors = { draft: 'grey', sent: 'warning', paid: 'success' }
  return colors[status] || 'grey'
}

const loadInvoice = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const [invoiceRes, itemsRes] = await Promise.all([
      invoiceService.getInvoice(invoiceId.value),
      invoiceService.getInvoiceItems(invoiceId.value)
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

const handleSubmit = async () => {
  console.log('Starting form submission...')
  
  const isValid = await invoiceFormRef.value.validate()
  console.log('Form validation result:', isValid)
  
  if (!isValid) {
    showNotification('Please fill in all required fields correctly', 'error')
    return
  }

  saving.value = true
  try {
    console.log('Getting data from form...')
    const data = invoiceFormRef.value.getData()
    console.log('Data to be submitted:', JSON.stringify(data, null, 2))

    if (isEditMode.value) {
      console.log('Updating invoice...')
      await invoiceService.updateInvoice(invoiceId.value, data)
      showNotification('Invoice updated successfully', 'success')
      router.push({ name: 'InvoiceDetail', params: { id: invoiceId.value } })
    } else {
      console.log('Creating invoice...')
      const response = await invoiceService.createInvoice(data)
      console.log('Backend response:', response.data)
      showNotification('Invoice created successfully', 'success')
      const newId = response.data.data.id
      router.push({ name: 'InvoiceEdit', params: { id: newId } })
    }
  } catch (error) {
    console.error('Full error details:', error)
    console.error('Error response:', error.response?.data)
    
    if (error.response?.data?.error) {
      const apiError = error.response.data.error
      showNotification(`${apiError.message}: ${apiError.details}`, 'error')
    } else if (error.message) {
      showNotification(error.message, 'error')
    } else {
      showNotification('Failed to save invoice', 'error')
    }
  } finally {
    saving.value = false
  }
}

const handleAddItem = async (itemData) => {
  try {
    await invoiceService.addInvoiceItem(invoiceId.value, itemData)
    showNotification('Item added successfully', 'success')
    await loadInvoice()
  } catch (error) {
    showNotification(error.message, 'error')
    throw error
  }
}

const handleUpdateItem = async (itemId, itemData) => {
  try {
    await invoiceService.updateInvoiceItem(invoiceId.value, itemId, itemData)
    showNotification('Item updated successfully', 'success')
    await loadInvoice()
  } catch (error) {
    showNotification(error.message, 'error')
    throw error
  }
}

const handleDeleteItem = async (itemId) => {
  try {
    await invoiceService.deleteInvoiceItem(invoiceId.value, itemId)
    showNotification('Item deleted successfully', 'success')
    await loadInvoice()
  } catch (error) {
    showNotification(error.message, 'error')
    throw error
  }
}

const handleCancel = () => {
  if (isEditMode.value) {
    router.push({ name: 'InvoiceDetail', params: { id: invoiceId.value } })
  } else {
    router.push({ name: 'InvoiceList' })
  }
}

onMounted(() => {
  loadInvoice()
})
</script>

<style scoped>
.gap-2 {
  gap: 8px;
}
</style>