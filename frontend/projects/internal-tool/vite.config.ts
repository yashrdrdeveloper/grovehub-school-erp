import { fileURLToPath, URL } from 'node:url';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: { alias: {
    '@Base': fileURLToPath(new URL('../../libs/base', import.meta.url)),
    '@Components': fileURLToPath(new URL('../../libs/components', import.meta.url)),
    '@Styles': fileURLToPath(new URL('../../libs/styles', import.meta.url)),
  }},
  server: {
    host: '127.0.0.1', port: 4174, strictPort: true,
    proxy: { '/internal': { target: 'http://127.0.0.1:8082', changeOrigin: false } },
  },
});
