<template>
  <tr :class="{ 'deleting-row': deleting }">
    <td class="text-secondary">{{ index + 1 }}</td>
    <td>{{ item.description }}</td>
    <td class="text-center">{{ item.quantity }}</td>
    <td class="text-right">{{ formatCurrency(item.unit_price) }}</td>
    <td class="text-center">{{ item.vat_rate }}%</td>
    <td class="text-right font-weight-medium">
      {{ formatCurrency(item.quantity * item.unit_price) }}
    </td>
    <td class="text-right font-weight-bold text-primary">
      {{ formatCurrency(calculateGrossTotal) }}
    </td>
    <td class="text-center">
      <div class="d-flex justify-center gap-1">
        <v-btn
          icon
          size="x-small"
          variant="text"
          color="primary"
          :disabled="deleting"
          @click="$emit('edit', item)"
        >
          <v-icon>mdi-pencil</v-icon>
        </v-btn>
        <v-btn
          icon
          size="x-small"
          variant="text"
          color="error"
          :loading="deleting"
          @click="$emit('delete', item)"
        >
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </div>
    </td>
  </tr>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    required: true
  },
  deleting: {
    type: Boolean,
    default: false
  }
})

defineEmits(['edit', 'delete'])

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

const calculateGrossTotal = computed(() => {
  const net = (props.item.quantity || 0) * (props.item.unit_price || 0)
  const vat = net * ((props.item.vat_rate || 0) / 100)
  return net + vat
})
</script>

<style scoped>
.gap-1 {
  gap: 4px;
}

.deleting-row {
  opacity: 0.6;
  background-color: rgba(244, 67, 54, 0.1);
}
</style>