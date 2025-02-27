import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import viteTsconfigPaths from 'vite-tsconfig-paths';;

export default defineConfig({
  base: "/",
  build: {
    outDir: 'build',
  },
  plugins: [react(), viteTsconfigPaths()],
  server: {
    host: true,
    open: true,
    port: 3000,
  },
});
