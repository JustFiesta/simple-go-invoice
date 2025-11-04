<template>
  <v-app>
    <v-app-bar color="surface" elevation="0" class="px-4 app-bar">
      <v-toolbar-title class="ml-2">
        <router-link to="/" class="text-decoration-none d-flex align-center">
          <v-icon icon="mdi-file-document-multiple" size="32" color="primary" class="mr-3"></v-icon>
          <span class="text-h5 font-weight-bold">Invoice Manager</span>
        </router-link>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <v-btn
        v-if="$route.name !== 'InvoiceList'"
        :to="{ name: 'InvoiceList' }"
        variant="text"
        prepend-icon="mdi-view-list"
        color="secondary"
        class="mr-2"
      >
        All Invoices
      </v-btn>

      <v-btn
        v-if="$route.name !== 'InvoiceCreate'"
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

    <v-main class="main-content">
      <v-container fluid class="pa-6">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </v-container>
    </v-main>

    <v-snackbar
      v-model="notification.show"
      :color="notification.color"
      :timeout="3000"
      location="top right"
      elevation="8"
      rounded="lg"
    >
      {{ notification.message }}
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
import { reactive, provide } from 'vue'

const notification = reactive({
  show: false,
  message: '',
  color: 'success'
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
</style>