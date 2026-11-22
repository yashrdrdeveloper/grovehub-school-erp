import path from 'node:path';

import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

const root = process.cwd();

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@Base': path.resolve(root, 'packages/base'),
      '@Business': path.resolve(root, 'packages/business'),
      '@Components': path.resolve(root, 'packages/components'),
      '@Domain': path.resolve(root, 'packages/domain'),
      '@Shell': path.resolve(root, 'packages/shell'),
    },
  },
  server: {
    host: '127.0.0.1',
    port: 4173,
    strictPort: true,
    proxy: {
      '/identity': { target: 'http://127.0.0.1:8080', changeOrigin: false },
      '/foundation': { target: 'http://127.0.0.1:8080', changeOrigin: false },
      '/integrations': { target: 'http://127.0.0.1:8080', changeOrigin: false },
      '/health': { target: 'http://127.0.0.1:8080', changeOrigin: false },
    },
  },
});
