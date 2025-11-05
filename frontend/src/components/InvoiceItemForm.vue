<template>
  <tr class="item-form-row">
    <td colspan="7" class="pa-0">
      <v-form ref="formRef" v-model="valid" class="item-form">
        <v-row no-gutters class="pa-3">
          <v-col cols="12" md="4" class="px-2">
            <v-text-field
              v-model="localItem.description"
              label="Description *"
              :rules="[rules.required]"
              density="compact"
              variant="outlined"
              hide-details="auto"
              bg-color="surface-light"
            ></v-text-field>
          </v-col>

          <v-col cols="6" md="2" class="px-2">
            <v-text-field
              v-model.number="localItem.quantity"
              label="Quantity *"
              type="number"
              min="1"
              :rules="[rules.required, rules.positiveNumber]"
              density="compact"
              variant="outlined"
              hide-details="auto"
              bg-color="surface-light"
            ></v-text-field>
          </v-col>

          <v-col cols="6" md="2" class="px-2">
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
              bg-color="surface-light"
            ></v-text-field>
          </v-col>

          <v-col cols="6" md="2" class="px-2">
            <v-text-field
              v-model.number="localItem.vat_rate"
              label="VAT % *"
              type="number"
              step="0.01"
              min="0"
              max="100"
              :rules="[rules.required, rules.nonNegative]"
              density="compact"
              variant="outlined"
              hide-details="auto"
              bg-color="surface-light"
            ></v-text-field>
          </v-col>

          <v-col cols="6" md="2" class="px-2 d-flex align-center">
            <v-btn
              color="success"
              variant="flat"
              size="small"
              :disabled="!valid"
              :loading="saving"
              @click="handleSave"
              class="mr-2"
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
          </v-col>
        </v-row>
      </v-form>
    </td>
  </tr>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

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

.item-form {
  width: 100%;
}
</style>