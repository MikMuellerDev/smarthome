import { resolve } from 'path'
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    rollupOptions: {
      input: {
        login: resolve(__dirname, 'html/login.html'),
        dash: resolve(__dirname, 'html/dash.html'),
      }
    },
  },
})
