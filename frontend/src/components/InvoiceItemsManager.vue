<template>
  <v-card elevation="0" class="items-manager-card">
    <v-card-title class="pa-6 pb-4 d-flex justify-space-between align-center">
      <div class="d-flex align-center">
        <v-icon icon="mdi-format-list-bulleted" color="primary" class="mr-3"></v-icon>
        <span class="text-h6 font-weight-bold">Invoice Items</span>
        <v-chip v-if="displayedItems.length > 0" size="small" class="ml-3" color="secondary">
          {{ displayedItems.length }}
        </v-chip>
      </div>
      <v-btn
        color="secondary"
        prepend-icon="mdi-plus"
        variant="flat"
        @click="startAddingItem"
        :disabled="isAdding || isEditing"
      >
        Add Item
      </v-btn>
    </v-card-title>

    <v-divider :thickness="1" color="surface-lighter"></v-divider>

    <v-card-text class="pa-0">
      <div v-if="displayedItems.length > 0 || isAdding || isEditing" class="table-wrapper">
        <v-table class="custom-table">
          <thead>
            <tr>
              <th style="width: 50px">#</th>
              <th>Description</th>
              <th class="text-center" style="width: 100px">Quantity</th>
              <th class="text-right" style="width: 130px">Unit Price</th>
              <th class="text-center" style="width: 100px">VAT %</th>
              <th class="text-right" style="width: 130px">Net Total</th>
              <th class="text-right" style="width: 150px">Gross Total</th>
              <th class="text-center" style="width: 120px">Actions</th>
            </tr>
          </thead>
          <tbody>
            <!-- New Item Form - zawsze na górze -->
            <InvoiceItemForm
              v-if="isAdding"
              :item="newItem"
              :saving="operationLoading"
              @save="handleAddItem"
              @cancel="cancelAdding"
            />

            <!-- Dynamiczna lista - przełącza między display a edit mode -->
            <template v-for="(item, index) in props.items" :key="item.id">
              <!-- Tryb edycji -->
              <InvoiceItemForm
                v-if="editingItemId === item.id"
                :item="editingItem"
                :saving="operationLoading"
                @save="handleUpdateItem"
                @cancel="cancelEditing"
              />
              <!-- Tryb wyświetlania -->
              <InvoiceItemRow
                v-else
                :item="item"
                :index="getItemIndex(index)"
                :deleting="deletingItemId === item.id"
                @edit="startEditingItem"
                @delete="handleDeleteItem"
              />
            </template>
          </tbody>
        </v-table>
      </div>

      <div v-else class="empty-state pa-12 text-center">
        <v-icon icon="mdi-package-variant" size="64" color="surface-lighter"></v-icon>
        <h4 class="text-h6 mt-4 mb-2">No items added yet</h4>
        <p class="text-body-2 text-secondary mb-4">Add items to this invoice to calculate totals</p>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import InvoiceItemForm from './InvoiceItemForm.vue'
import InvoiceItemRow from './InvoiceItemRow.vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  invoiceId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['add', 'update', 'delete'])

const isAdding = ref(false)
const editingItemId = ref(null)
const editingItem = reactive({})
const operationLoading = ref(false)
const deletingItemId = ref(null)

const newItem = reactive({
  description: '',
  quantity: 1,
  unit_price: 0,
  vat_rate: 23
})

// Computed properties
const isEditing = computed(() => editingItemId.value !== null)
const displayedItems = computed(() => {
  return props.items.filter(item => item.id !== editingItemId.value)
})

const getItemIndex = (index) => {
  // Jeśli jest formularz dodawania na górze, przesuwamy indeksy
  if (isAdding.value) {
    return index + 1
  }
  return index + 1
}

const startAddingItem = () => {
  console.log('Starting to add new item')
  isAdding.value = true
  editingItemId.value = null
  Object.assign(newItem, {
    description: '',
    quantity: 1,
    unit_price: 0,
    vat_rate: 23
  })
}

const cancelAdding = () => {
  console.log('Canceling add')
  isAdding.value = false
}

const handleAddItem = async (itemData) => {
  console.log('Adding item:', itemData)
  operationLoading.value = true
  try {
    await emit('add', itemData)
    isAdding.value = false
    Object.assign(newItem, {
      description: '',
      quantity: 1,
      unit_price: 0,
      vat_rate: 23
    })
  } finally {
    operationLoading.value = false
  }
}

const startEditingItem = (item) => {
  console.log('Starting edit for item:', item)
  editingItemId.value = item.id
  isAdding.value = false
  Object.assign(editingItem, { ...item })
}

const cancelEditing = () => {
  console.log('Canceling edit')
  editingItemId.value = null
  // Czyścimy editingItem
  Object.keys(editingItem).forEach(key => {
    delete editingItem[key]
  })
}

const handleUpdateItem = async (itemData) => {
  console.log('Updating item:', itemData)
  operationLoading.value = true
  try {
    await emit('update', editingItemId.value, itemData)
    editingItemId.value = null
    // Czyścimy editingItem
    Object.keys(editingItem).forEach(key => {
      delete editingItem[key]
    })
  } finally {
    operationLoading.value = false
  }
}

const handleDeleteItem = async (item) => {
  console.log('Deleting item:', item)
  if (!confirm(`Delete item "${item.description}"?`)) return
  
  deletingItemId.value = item.id
  try {
    await emit('delete', item.id)
  } finally {
    deletingItemId.value = null
  }
}
</script>

<style scoped>
.items-manager-card {
  background: linear-gradient(135deg, #2A2A2A 0%, #252525 100%);
  border: 1px solid #404040;
}

.table-wrapper {
  overflow-x: auto;
}

.custom-table {
  background: transparent;
}

.custom-table thead th {
  background-color: rgba(75, 240, 240, 0.1);
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.5px;
  padding: 16px 12px;
  border-bottom: 2px solid #404040;
}

.custom-table tbody td {
  padding: 12px;
  border-bottom: 1px solid #353535;
}

.empty-state {
  min-height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>