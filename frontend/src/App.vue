<template>
  <v-app>
    <!-- App Bar -->
    <v-app-bar color="surface" elevation="2" class="px-4">
      <v-toolbar-title class="ml-2">
        <router-link to="/" class="text-decoration-none d-flex align-center">
          <v-icon icon="mdi-file-document-multiple" size="32" color="primary" class="mr-3"></v-icon>
          <span class="text-h5 font-weight-bold">Invoice Manager</span>
        </router-link>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <!-- Navigation Buttons -->
      <v-btn
        :to="{ name: 'Dashboard' }"
        :variant="$route.name === 'Dashboard' ? 'flat' : 'text'"
        :color="$route.name === 'Dashboard' ? 'primary' : 'secondary'"
        prepend-icon="mdi-view-dashboard"
        class="mr-2"
      >
        Dashboard
      </v-btn>

      <v-btn
        :to="{ name: 'InvoiceList' }"
        :variant="$route.name === 'InvoiceList' ? 'flat' : 'text'"
        :color="$route.name === 'InvoiceList' ? 'primary' : 'secondary'"
        prepend-icon="mdi-format-list-bulleted"
        class="mr-2"
      >
        Invoices
      </v-btn>

      <v-btn
        :to="{ name: 'InvoiceCreate' }"
        color="primary"
        variant="flat"
        prepend-icon="mdi-plus"
        size="large"
        elevation="2"
      >
        New Invoice
      </v-btn>
    </v-app-bar>

    <!-- Main Content -->
    <v-main class="main-content">
      <v-container fluid class="pa-6">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </v-container>
    </v-main>

    <!-- Global Notification -->
    <v-snackbar
      v-model="notification.show"
      :color="notification.color"
      :timeout="3000"
      location="top right"
      elevation="8"
      rounded="lg"
    >
      <div class="d-flex align-center">
        <v-icon
          :icon="notificationIcon"
          class="mr-3"
        ></v-icon>
        {{ notification.message }}
      </div>
      <template #actions>
        <v-btn
          variant="text"
          icon="mdi-close"
          @click="notification.show = false"
        ></v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { reactive, computed, provide } from 'vue'

const notification = reactive({
  show: false,
  message: '',
  color: 'success'
})

const notificationIcon = computed(() => {
  const icons = {
    success: 'mdi-check-circle',
    error: 'mdi-alert-circle',
    warning: 'mdi-alert',
    info: 'mdi-information'
  }
  return icons[notification.color] || 'mdi-information'
})

const showNotification = (message, color = 'success') => {
  notification.message = message
  notification.color = color
  notification.show = true
}

provide('showNotification', showNotification)
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

a {
  text-decoration: none;
  color: inherit;
}

.main-content {
  background: linear-gradient(180deg, #1E1E1E 0%, #1A1A1A 100%);
  min-height: 100vh;
}
</style>