<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12">
        <h1 class="text-h4 font-weight-bold text-primary">
          {{ isEditMode ? 'Edit Invoice' : 'Create Invoice' }}
        </h1>
      </v-col>
    </v-row>

    <v-form ref="formRef" v-model="formValid" @submit.prevent="handleSubmit">
      <v-card elevation="0" class="form-card mb-6">
        <v-card-title class="pa-6 pb-4 text-h6 font-weight-bold">
          <v-icon icon="mdi-file-document-edit" color="primary" class="mr-2"></v-icon>
          Invoice Details
        </v-card-title>

        <v-divider :thickness="1" color="surface-lighter"></v-divider>

        <v-card-text class="pa-6">
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.invoice_number"
                label="Invoice Number *"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-file-document"
                bg-color="surface-light"
              ></v-text-field>
            </v-col>

            <v-col cols="12" md="6">
              <v-text-field
                v-model="form.client_name"
                label="Client Name *"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-account"
                bg-color="surface-light"
              ></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
              <v-text-field
                v-model.number="form.amount"
                label="Amount *"
                type="number"
                step="0.01"
                :rules="[rules.required, rules.positiveNumber]"
                prepend-inner-icon="mdi-currency-usd"
                bg-color="surface-light"
              ></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
              <v-select
                v-model="form.status"
                label="Status *"
                :items="statusOptions"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-tag"
                bg-color="surface-light"
              ></v-select>
            </v-col>

            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.issue_date"
                label="Issue Date *"
                type="date"
                :rules="[rules.required]"
                prepend-inner-icon="mdi-calendar"
                bg-color="surface-light"
              ></v-text-field>
            </v-col>

            <v-col cols="12" md="4">
              <v-text-field
                v-model="form.due_date"
                label="Due Date *"
                type="date"
                :rules="[rules.required, validateDueDate]"
                prepend-inner-icon="mdi-calendar-clock"
                bg-color="surface-light"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <v-card v-if="isEditMode && items.length > 0" elevation="0" class="form-card mb-6">
        <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
          <div class="d-flex align-center">
            <v-icon icon="mdi-format-list-bulleted" color="primary" class="mr-2"></v-icon>
            <span class="text-h6 font-weight-bold">Invoice Items</span>
          </div>
          <v-btn
            color="secondary"
            prepend-icon="mdi-plus"
            @click="showItemDialog()"
          >
            Add Item
          </v-btn>
        </v-card-title>

        <v-divider :thickness="1" color="surface-lighter"></v-divider>

        <v-card-text class="pa-0">
          <div class="table-wrapper">
            <v-table class="custom-table">
              <thead>
                <tr>
                  <th>Description</th>
                  <th class="text-center">Quantity</th>
                  <th class="text-right">Unit Price</th>
                  <th class="text-center">VAT %</th>
                  <th class="text-right">Total</th>
                  <th class="text-center">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in items" :key="item.id">
                  <td>{{ item.description }}</td>
                  <td class="text-center">{{ item.quantity }}</td>
                  <td class="text-right">{{ formatCurrency(item.unit_price) }}</td>
                  <td class="text-center">{{ item.vat_rate }}%</td>
                  <td class="text-right font-weight-bold">
                    {{ formatCurrency(item.quantity * item.unit_price) }}
                  </td>
                  <td class="text-center">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      variant="text"
                      color="secondary"
                      @click="showItemDialog(item)"
                    ></v-btn>
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      @click="deleteItem(item)"
                    ></v-btn>
                  </td>
                </tr>
              </tbody>
            </v-table>
          </div>
        </v-card-text>
      </v-card>

      <div class="mt-6 d-flex justify-space-between">
        <v-btn
          variant="outlined"
          prepend-icon="mdi-arrow-left"
          @click="$router.back()"
        >
          Cancel
        </v-btn>

        <v-btn
          type="submit"
          color="primary"
          :loading="loading"
          prepend-icon="mdi-content-save"
        >
          {{ isEditMode ? 'Update' : 'Create' }} Invoice
        </v-btn>
      </div>
    </v-form>

    <!-- Item Dialog -->
    <v-dialog v-model="itemDialog.show" max-width="600">
      <v-card>
        <v-card-title>
          {{ itemDialog.item ? 'Edit Item' : 'Add Item' }}
        </v-card-title>

        <v-card-text>
          <v-form ref="itemFormRef" v-model="itemDialog.valid">
            <v-text-field
              v-model="itemDialog.form.description"
              label="Description *"
              :rules="[rules.required]"
              class="mb-2"
            ></v-text-field>

            <v-text-field
              v-model.number="itemDialog.form.quantity"
              label="Quantity *"
              type="number"
              :rules="[rules.required, rules.positiveNumber]"
              class="mb-2"
            ></v-text-field>

            <v-text-field
              v-model.number="itemDialog.form.unit_price"
              label="Unit Price *"
              type="number"
              step="0.01"
              :rules="[rules.required, rules.positiveNumber]"
              class="mb-2"
            ></v-text-field>

            <v-text-field
              v-model.number="itemDialog.form.vat_rate"
              label="VAT Rate % *"
              type="number"
              step="0.01"
              :rules="[rules.required]"
            ></v-text-field>
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="itemDialog.show = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="itemDialog.loading"
            :disabled="!itemDialog.valid"
            @click="saveItem"
          >
            Save
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

const formRef = ref(null)
const itemFormRef = ref(null)
const formValid = ref(false)
const loading = ref(false)
const items = ref([])

const isEditMode = computed(() => route.name === 'InvoiceEdit')

const form = reactive({
  invoice_number: '',
  client_name: '',
  amount: 0,
  status: 'draft',
  issue_date: new Date().toISOString().split('T')[0],
  due_date: ''
})

const itemDialog = reactive({
  show: false,
  valid: false,
  loading: false,
  item: null,
  form: {
    description: '',
    quantity: 1,
    unit_price: 0,
    vat_rate: 23
  }
})

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Sent', value: 'sent' },
  { title: 'Paid', value: 'paid' }
]

const rules = {
  required: v => !!v || 'This field is required',
  positiveNumber: v => v > 0 || 'Must be greater than 0'
}

const validateDueDate = (v) => {
  if (!v || !form.issue_date) return true
  return new Date(v) >= new Date(form.issue_date) || 'Due date must be after issue date'
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount)
}

const loadInvoice = async () => {
  if (!isEditMode.value) return

  loading.value = true
  try {
    const [invoiceRes, itemsRes] = await Promise.all([
      invoiceService.getInvoice(route.params.id),
      invoiceService.getInvoiceItems(route.params.id)
    ])

    const invoice = invoiceRes.data.data
    form.invoice_number = invoice.invoice_number
    form.client_name = invoice.client_name
    form.amount = invoice.amount
    form.status = invoice.status
    form.issue_date = invoice.issue_date.split('T')[0]
    form.due_date = invoice.due_date.split('T')[0]

    items.value = itemsRes.data.data || []
  } catch (error) {
    showNotification(error.message, 'error')
    router.push({ name: 'InvoiceList' })
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  try {
    const data = {
      ...form,
      issue_date: new Date(form.issue_date).toISOString(),
      due_date: new Date(form.due_date).toISOString()
    }

    if (isEditMode.value) {
      await invoiceService.updateInvoice(route.params.id, data)
      showNotification('Invoice updated successfully', 'success')
    } else {
      await invoiceService.createInvoice(data)
      showNotification('Invoice created successfully', 'success')
    }

    router.push({ name: 'InvoiceList' })
  } catch (error) {
    showNotification(error.message, 'error')
  } finally {
    loading.value = false
  }
}

const showItemDialog = (item = null) => {
  if (item) {
    itemDialog.item = item
    itemDialog.form = { ...item }
  } else {
    itemDialog.item = null
    itemDialog.form = {
      description: '',
      quantity: 1,
      unit_price: 0,
      vat_rate: 23
    }
  }
  itemDialog.show = true
}

const saveItem = async () => {
  const { valid } = await itemFormRef.value.validate()
  if (!valid) return

  itemDialog.loading = true
  try {
    if (itemDialog.item) {
      await invoiceService.updateInvoiceItem(
        route.params.id,
        itemDialog.item.id,
        itemDialog.form
      )
      showNotification('Item updated successfully', 'success')
    } else {
      await invoiceService.addInvoiceItem(route.params.id, itemDialog.form)
      showNotification('Item added successfully', 'success')
    }

    itemDialog.show = false
    loadInvoice()
  } catch (error) {
    showNotification(error.message, 'error')
  } finally {
    itemDialog.loading = false
  }
}

const deleteItem = async (item) => {
  if (!confirm(`Delete item "${item.description}"?`)) return

  try {
    await invoiceService.deleteInvoiceItem(route.params.id, item.id)
    showNotification('Item deleted successfully', 'success')
    loadInvoice()
  } catch (error) {
    showNotification(error.message, 'error')
  }
}

onMounted(() => {
  loadInvoice()
})
</script>