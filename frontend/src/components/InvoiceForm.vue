<template>
  <v-form ref="formRef" v-model="formValid">
    <v-row>
      <!-- Left Column: Header + Items -->
      <v-col cols="12" lg="8">
        <!-- Invoice Header -->
        <InvoiceHeaderFields
          v-model:invoice-number="invoiceData.invoice_number"
          v-model:client-name="invoiceData.client_name"
          v-model:status="invoiceData.status"
          v-model:issue-date="invoiceData.issue_date"
          v-model:due-date="invoiceData.due_date"
          class="mb-6"
        />

        <!-- Invoice Items Manager -->
        <InvoiceItemsManager
          v-if="isEditMode"
          :items="items"
          :invoice-id="invoiceId"
          @add="handleAddItem"
          @update="handleUpdateItem"
          @delete="handleDeleteItem"
          class="mb-6"
        />

        <v-alert
          v-else
          type="info"
          variant="tonal"
          class="mb-6"
          icon="mdi-information"
        >
          Save the invoice first to add items
        </v-alert>
      </v-col>

      <!-- Right Column: Summary (Sticky) -->
      <v-col cols="12" lg="4">
        <InvoiceSummary
          ref="summaryRef"
          :items="items"
        />
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import InvoiceHeaderFields from './InvoiceHeaderFields.vue'
import InvoiceItemsManager from './InvoiceItemsManager.vue'
import InvoiceSummary from './InvoiceSummary.vue'

const props = defineProps({
  invoice: {
    type: Object,
    default: null
  },
  items: {
    type: Array,
    default: () => []
  },
  isEditMode: {
    type: Boolean,
    default: false
  },
  invoiceId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['add-item', 'update-item', 'delete-item'])

const formRef = ref(null)
const summaryRef = ref(null)
const formValid = ref(false)

const invoiceData = reactive({
  invoice_number: '',
  client_name: '',
  status: 'draft',
  issue_date: new Date().toISOString().split('T')[0],
  due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
  amount: 0
})

// Initialize or update invoice data
watch(() => props.invoice, (newInvoice) => {
  if (newInvoice) {
    Object.assign(invoiceData, {
      invoice_number: newInvoice.invoice_number,
      client_name: newInvoice.client_name,
      status: newInvoice.status,
      issue_date: newInvoice.issue_date?.split('T')[0] || '',
      due_date: newInvoice.due_date?.split('T')[0] || '',
      amount: newInvoice.amount || 0
    })
  }
}, { immediate: true })

// Auto-update amount from summary
watch(() => props.items, () => {
  if (summaryRef.value) {
    invoiceData.amount = summaryRef.value.total
  }
}, { deep: true, immediate: true })

const handleAddItem = async (itemData) => {
  await emit('add-item', itemData)
}

const handleUpdateItem = async (itemId, itemData) => {
  await emit('update-item', itemId, itemData)
}

const handleDeleteItem = async (itemId) => {
  await emit('delete-item', itemId)
}

const validate = async () => {
  const result = await formRef.value.validate()
  return result.valid
}
const getData = () => {
  console.log('Current invoiceData:', invoiceData)
  
  // Walidacja z lepszym debugowaniem
  if (!invoiceData.invoice_number?.trim()) {
    console.error('Invoice number validation failed:', invoiceData.invoice_number)
    throw new Error('Invoice number is required')
  }
  
  if (!invoiceData.client_name?.trim()) {
    console.error('Client name validation failed:', invoiceData.client_name)
    throw new Error('Client name is required')
  }
  
  if (!invoiceData.issue_date) {
    console.error('Issue date validation failed:', invoiceData.issue_date)
    throw new Error('Issue date is required')
  }
  
  if (!invoiceData.due_date) {
    console.error('Due date validation failed:', invoiceData.due_date)
    throw new Error('Due date is required')
  }

  const issueDate = new Date(invoiceData.issue_date)
  const dueDate = new Date(invoiceData.due_date)

  if (isNaN(issueDate.getTime())) {
    throw new Error('Invalid issue date')
  }
  if (isNaN(dueDate.getTime())) {
    throw new Error('Invalid due date')
  }
  if (dueDate < issueDate) {
    throw new Error('Due date cannot be before issue date')
  }

  // Ustaw czas na 12:00:00 dla spójności
  issueDate.setHours(12, 0, 0, 0)
  dueDate.setHours(12, 0, 0, 0)

  // Upewnij się, że amount jest > 0
  const finalAmount = Number(invoiceData.amount) > 0 ? Number(invoiceData.amount) : 0.01

  const data = {
    invoice_number: invoiceData.invoice_number.trim(),
    client_name: invoiceData.client_name.trim(),
    status: invoiceData.status,
    issue_date: issueDate.toISOString(),
    due_date: dueDate.toISOString(),
    amount: finalAmount
  }

  console.log('Final data to submit:', data)
  return data
}


defineExpose({
  validate,
  getData,
  formValid
})
</script>