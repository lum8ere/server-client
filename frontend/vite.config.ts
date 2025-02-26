import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import viteTsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
    base: '/',
    build: {
        outDir: 'build'
    },
    plugins: [react(), viteTsconfigPaths()],
    server: {
        host: true,
        open: true,
        port: 3000,
        proxy: {
            '/api': 'http://localhost:4000',
            '/uploads': 'http://localhost:4000',
            '/ws': { target: 'ws://localhost:4000', ws: true }
        }
    }
});
