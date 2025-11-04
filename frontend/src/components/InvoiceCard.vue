<template>
  <v-card class="invoice-card" elevation="2">
    <v-card-title class="d-flex justify-space-between align-center pb-2">
      <span class="text-h6 font-weight-bold">{{ invoice.invoice_number }}</span>
      <v-chip
        :color="getStatusColor(invoice.status)"
        size="small"
        label
      >
        {{ invoice.status.toUpperCase() }}
      </v-chip>
    </v-card-title>

    <v-card-subtitle class="pb-3">
      <v-icon icon="mdi-account" size="16" class="mr-1"></v-icon>
      {{ invoice.client_name }}
    </v-card-subtitle>

    <v-card-text class="pt-0">
      <v-divider class="mb-3"></v-divider>

      <div class="d-flex justify-space-between mb-2">
        <span class="text-body-2 text-grey">Amount:</span>
        <span class="text-h6 font-weight-bold text-primary">
          {{ formatCurrency(invoice.amount) }}
        </span>
      </div>

      <div class="d-flex justify-space-between mb-2">
        <span class="text-body-2 text-grey">Issue Date:</span>
        <span class="text-body-2">{{ formatDate(invoice.issue_date) }}</span>
      </div>

      <div class="d-flex justify-space-between">
        <span class="text-body-2 text-grey">Due Date:</span>
        <span
          class="text-body-2"
          :class="{ 'text-error': isOverdue(invoice.due_date) }"
        >
          {{ formatDate(invoice.due_date) }}
          <v-icon
            v-if="isOverdue(invoice.due_date) && invoice.status !== 'paid'"
            icon="mdi-alert-circle"
            size="16"
            color="error"
            class="ml-1"
          ></v-icon>
        </span>
      </div>
    </v-card-text>

    <v-card-actions class="px-4 pb-4">
      <v-btn
        size="small"
        variant="text"
        prepend-icon="mdi-eye"
        @click="$emit('view', invoice.id)"
      >
        View
      </v-btn>

      <v-spacer></v-spacer>

      <v-btn
        icon="mdi-download"
        size="small"
        variant="text"
        @click="$emit('download', invoice)"
      ></v-btn>

      <v-btn
        icon="mdi-pencil"
        size="small"
        variant="text"
        @click="$emit('edit', invoice.id)"
      ></v-btn>

      <v-btn
        icon="mdi-delete"
        size="small"
        variant="text"
        color="error"
        @click="$emit('delete', invoice)"
      ></v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup>
defineProps({
  invoice: {
    type: Object,
    required: true
  }
})

defineEmits(['view', 'edit', 'delete', 'download'])

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
  return new Date(date).toLocaleDateString('pl-PL')
}

const isOverdue = (dueDate) => {
  return new Date(dueDate) < new Date() && new Date(dueDate).toDateString() !== new Date().toDateString()
}
</script>

<style scoped>
.invoice-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.invoice-card:hover {
  cursor: pointer;
}
</style>