import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

const customTheme = {
  dark: true,
  colors: {
    primary: '#FF7B3E',
    secondary: '#4BF0F0',
    accent: '#9935C1',
    background: '#1E1E1E',
    surface: '#2A2A2A',
    'surface-light': '#353535',
    'surface-lighter': '#404040',
    error: '#FF5252',
    info: '#4BF0F0',
    success: '#4CAF50',
    warning: '#FF7B3E',
    'on-surface': '#FFFFFF',
    'on-background': '#FFFFFF'
  }
}

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'customTheme',
    themes: {
      customTheme
    }
  },
  defaults: {
    VBtn: {
      color: 'primary',
      variant: 'flat'
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable'
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable'
    }
  }
})