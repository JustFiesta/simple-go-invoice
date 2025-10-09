<template>
  <v-app>
    <v-main class="d-flex flex-column align-center justify-center text-center">
      <h1>Invoice SPA</h1>
      <v-btn color="primary" @click="fetchMessage">Ping Backend</v-btn>
      <p>{{ message }}</p>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'

const message = ref('')

async function fetchMessage() {
  const res = await fetch('http://localhost:8080/api/hello')
  const data = await res.json()
  message.value = data.message
}
</script>
