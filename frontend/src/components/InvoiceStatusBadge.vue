<template>
  <v-chip
    :color="statusColor"
    :variant="variant"
    :size="size"
    label
    :class="classes"
  >
    <v-icon v-if="showIcon" :icon="statusIcon" size="16" class="mr-1"></v-icon>
    {{ statusText }}
  </v-chip>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['draft', 'sent', 'paid'].includes(value)
  },
  size: {
    type: String,
    default: 'default'
  },
  variant: {
    type: String,
    default: 'flat'
  },
  showIcon: {
    type: Boolean,
    default: true
  }
})

const statusConfig = {
  draft: {
    color: 'grey',
    text: 'DRAFT',
    icon: 'mdi-file-document-edit-outline'
  },
  sent: {
    color: 'warning',
    text: 'SENT',
    icon: 'mdi-send'
  },
  paid: {
    color: 'success',
    text: 'PAID',
    icon: 'mdi-check-circle'
  }
}

const statusColor = computed(() => statusConfig[props.status]?.color || 'grey')
const statusText = computed(() => statusConfig[props.status]?.text || props.status.toUpperCase())
const statusIcon = computed(() => statusConfig[props.status]?.icon || 'mdi-file-document')

const classes = computed(() => ({
  'font-weight-bold': true
}))
</script>

<style scoped>
.v-chip {
  letter-spacing: 0.5px;
}
</style>