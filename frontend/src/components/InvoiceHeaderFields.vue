<template>
  <v-card elevation="0" class="header-card">
    <v-card-title class="pa-6 pb-4 text-h6 font-weight-bold">
      <v-icon icon="mdi-file-document-edit" color="primary" class="mr-3"></v-icon>
      Invoice Information
    </v-card-title>

    <v-divider :thickness="1" color="surface-lighter"></v-divider>

    <v-card-text class="pa-6">
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field
            :model-value="invoiceNumber"
            @update:model-value="$emit('update:invoice-number', $event)"
            label="Invoice Number *"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-file-document"
            variant="outlined"
            bg-color="surface-light"
            hint="Unique identifier for this invoice"
            persistent-hint
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="6">
          <v-text-field
            :model-value="clientName"
            @update:model-value="$emit('update:client-name', $event)"
            label="Client Name *"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-account"
            variant="outlined"
            bg-color="surface-light"
            hint="Full name or company name"
            persistent-hint
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-select
            :model-value="status"
            @update:model-value="$emit('update:status', $event)"
            label="Status *"
            :items="statusOptions"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-tag"
            variant="outlined"
            bg-color="surface-light"
          >
            <template #selection="{ item }">
              <InvoiceStatusBadge :status="item.value" :size="'small'" :show-icon="false" />
            </template>
          </v-select>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            :model-value="issueDate"
            @update:model-value="$emit('update:issue-date', $event)"
            label="Issue Date *"
            type="date"
            :rules="[rules.required]"
            prepend-inner-icon="mdi-calendar"
            variant="outlined"
            bg-color="surface-light"
          ></v-text-field>
        </v-col>

        <v-col cols="12" md="4">
          <v-text-field
            :model-value="dueDate"
            @update:model-value="$emit('update:due-date', $event)"
            label="Due Date *"
            type="date"
            :rules="[rules.required, validateDueDate]"
            prepend-inner-icon="mdi-calendar-clock"
            variant="outlined"
            bg-color="surface-light"
            :error-messages="dueDateError"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { computed } from 'vue'
import InvoiceStatusBadge from './InvoiceStatusBadge.vue'

const props = defineProps({
  invoiceNumber: String,
  clientName: String,
  status: String,
  issueDate: String,
  dueDate: String
})

const emit = defineEmits([
  'update:invoice-number',
  'update:client-name', 
  'update:status',
  'update:issue-date',
  'update:due-date'
])

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Sent', value: 'sent' },
  { title: 'Paid', value: 'paid' }
]

const rules = {
  required: v => (v !== null && v !== undefined && v !== '') || 'This field is required'
}

const dueDateError = computed(() => {
  if (!props.dueDate || !props.issueDate) return ''
  const dueDate = new Date(props.dueDate)
  const issueDate = new Date(props.issueDate)
  return dueDate < issueDate ? 'Due date must be after issue date' : ''
})

const validateDueDate = (v) => {
  if (!v || !props.issueDate) return true
  const dueDate = new Date(v)
  const issueDate = new Date(props.issueDate)
  return dueDate >= issueDate || 'Due date must be after issue date'
}
</script>

<style scoped>
.header-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
}
</style>