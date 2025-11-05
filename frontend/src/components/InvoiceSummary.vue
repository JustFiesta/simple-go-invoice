<template>
  <v-card elevation="0" class="summary-card">
    <v-card-title class="text-h6 font-weight-bold pa-4">
      <v-icon icon="mdi-calculator" color="primary" class="mr-2"></v-icon>
      Invoice Summary
    </v-card-title>

    <v-divider :thickness="1" color="surface-lighter"></v-divider>

    <v-card-text class="pa-4">
      <div class="summary-row mb-3">
        <span class="text-body-1 text-secondary">Subtotal (Net):</span>
        <span class="text-h6 font-weight-medium">{{ formatCurrency(subtotal) }}</span>
      </div>

      <div class="summary-row mb-3">
        <span class="text-body-1 text-secondary">Total VAT:</span>
        <span class="text-h6 font-weight-medium">{{ formatCurrency(totalVat) }}</span>
      </div>

      <v-divider class="my-3" :thickness="1"></v-divider>

      <div class="summary-row">
        <span class="text-h6 font-weight-bold">Total (Gross):</span>
        <span class="text-h4 font-weight-bold text-primary">{{ formatCurrency(total) }}</span>
      </div>

      <div v-if="itemCount > 0" class="text-caption text-secondary mt-3 text-center">
        {{ itemCount }} {{ itemCount === 1 ? 'item' : 'items' }}
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  }
})

const itemCount = computed(() => props.items.length)

const subtotal = computed(() => {
  return props.items.reduce((sum, item) => {
    const itemTotal = (item.quantity || 0) * (item.unit_price || 0)
    return sum + itemTotal
  }, 0)
})

const totalVat = computed(() => {
  return props.items.reduce((sum, item) => {
    const itemNet = (item.quantity || 0) * (item.unit_price || 0)
    const itemVat = itemNet * ((item.vat_rate || 0) / 100)
    return sum + itemVat
  }, 0)
})

const total = computed(() => subtotal.value + totalVat.value)

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

defineExpose({ total })
</script>

<style scoped>
.summary-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
  position: sticky;
  top: 20px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>