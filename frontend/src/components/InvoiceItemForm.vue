<template>
  <tr class="item-form-row">
    <td colspan="8" class="pa-0">
      <v-form ref="formRef" v-model="valid">
        <table class="w-100">
            <tbody>
                <tr>
                    <td class="text-secondary"></td>
                    <td>
                    <v-text-field
                        v-model="localItem.description"
                        label="Description *"
                        :rules="[rules.required]"
                        density="compact"
                        variant="outlined"
                        hide-details="auto"
                    />
                    </td>
                    <td>
                    <v-text-field
                        v-model.number="localItem.quantity"
                        label="Quantity *"
                        type="number"
                        :rules="[rules.required, rules.positiveNumber]"
                        density="compact"
                        variant="outlined"
                        hide-details="auto"
                    />
                    </td>
                    <td>
                    <v-text-field
                        v-model.number="localItem.unit_price"
                        label="Unit Price *"
                        type="number"
                        step="0.01"
                        min="0"
                        :rules="[rules.required, rules.nonNegative]"
                        density="compact"
                        variant="outlined"
                        hide-details="auto"
                    />
                    </td>
                    <td>
                    <v-text-field
                        v-model.number="localItem.vat_rate"
                        label="VAT % *"
                        type="number"
                        :rules="[rules.required, rules.nonNegative]"
                        density="compact"
                        variant="outlined"
                        hide-details="auto"
                    />
                    </td>
                    <td class="text-right font-weight-medium">
                    {{ formatCurrency(localItem.quantity * localItem.unit_price) }}
                    </td>
                    <td class="text-right font-weight-bold text-primary">
                    {{ formatCurrency(calculateGrossTotal) }}
                    </td>
                    <td class="text-center">
                    <div class="d-flex justify-center gap-1">
                        <v-btn
                        color="success"
                        variant="flat"
                        size="small"
                        :disabled="!valid"
                        :loading="saving"
                        @click="handleSave"
                        >
                        <v-icon icon="mdi-check"></v-icon>
                        </v-btn>
                        <v-btn
                        color="error"
                        variant="outlined"
                        size="small"
                        @click="handleCancel"
                        >
                        <v-icon icon="mdi-close"></v-icon>
                        </v-btn>
                    </div>
                    </td>
                </tr>
            </tbody>
        </table>
      </v-form>
    </td>
  </tr>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'

const props = defineProps({
  item: {
    type: Object,
    default: () => ({
      description: '',
      quantity: 1,
      unit_price: 0,
      vat_rate: 23
    })
  },
  saving: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['save', 'cancel'])

const formRef = ref(null)
const valid = ref(false)

const localItem = reactive({ ...props.item })

watch(() => props.item, (newItem) => {
  Object.assign(localItem, newItem)
}, { deep: true })

const rules = {
  required: v => (v !== null && v !== undefined && v !== '') || 'Required',
  positiveNumber: v => v > 0 || 'Must be greater than 0',
  nonNegative: v => v >= 0 || 'Must be 0 or greater'
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

const calculateGrossTotal = computed(() => {
  const net = (localItem.quantity || 0) * (localItem.unit_price || 0)
  const vat = net * ((localItem.vat_rate || 0) / 100)
  return net + vat
})

const handleSave = async () => {
  const { valid: isValid } = await formRef.value.validate()
  if (isValid) {
    emit('save', { ...localItem })
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.item-form-row {
  background-color: rgba(255, 123, 62, 0.05);
  border-left: 3px solid var(--color-primary);
}

.gap-1 {
  gap: 4px;
}
</style>