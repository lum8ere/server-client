import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import viteTsconfigPaths from 'vite-tsconfig-paths';
import svgr from '@svgr/rollup';

export default defineConfig({
    base: '/devices',
    build: {
        outDir: 'build'
    },
    preview: {
        port: 3000,
        strictPort: true
    },
    plugins: [react(), viteTsconfigPaths(), svgr()],
    server: {
        host: true,
        open: true,
        port: 3000
    }
});
