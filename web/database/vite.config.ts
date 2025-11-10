import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    hmr: {
      host: '127.0.0.1',
      port: 5173 // 可以尝试调整端口
    },
    proxy: {
      // 示例：将以 /api 开头的请求代理到目标服务器
      '/api': {
        target: 'http://127.0.0.1:8080', // 目标服务器地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})