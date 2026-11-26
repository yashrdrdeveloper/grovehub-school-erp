import path from 'node:path';

import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig, type ProxyOptions } from 'vite';

const root = process.cwd();
const tunnelOrigin = process.env.GROVEHUB_FRONTEND_TUNNEL_URL ?? '';
const getTunnelHost = (origin: string): string => {
  if (!/^https:\/\/[a-z0-9-]+\.trycloudflare\.com$/.test(origin)) return '';
  return new URL(origin).hostname;
};
const tunnelHost = getTunnelHost(tunnelOrigin);
const getAllowedHosts = (): string[] => {
  if (!tunnelHost) return [];
  return [tunnelHost];
};
const localOrigin = process.env.GROVEHUB_FRONTEND_LOCAL_ORIGIN ?? 'http://127.0.0.1:4173';

// Only the managed development tunnel's same-origin browser requests are
// translated to the existing loopback origin. Never forward arbitrary origins
// as trusted, and never alter cookies, CSRF proofs or gateway permissions.
const createProxy = (): ProxyOptions => ({
  target: 'http://127.0.0.1:8080',
  changeOrigin: false,
  configure: (proxy) => {
    proxy.on('proxyReq', (proxyReq, request) => {
      if (!tunnelHost || request.headers.host !== tunnelHost || request.headers.origin !== tunnelOrigin) return;
      proxyReq.setHeader('Origin', localOrigin);
    });
  },
});

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@Base': path.resolve(root, 'packages/base'),
      '@Business': path.resolve(root, 'packages/business'),
      '@Components': path.resolve(root, 'packages/components'),
      '@Domain': path.resolve(root, 'packages/domain'),
      '@Shell': path.resolve(root, 'packages/shell'),
      '@Styles': path.resolve(root, 'packages/styles'),
      '@Assets': path.resolve(root, 'packages/assets'),
    },
  },
  server: {
    host: '127.0.0.1',
    port: 4173,
    strictPort: true,
    allowedHosts: getAllowedHosts(),
    proxy: {
      '/identity': createProxy(),
      '/foundation': createProxy(),
      '/integrations': createProxy(),
      '/health': createProxy(),
    },
  },
});
