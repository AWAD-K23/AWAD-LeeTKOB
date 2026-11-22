// Starts React and forwards API requests to the local Go server.
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// Use IPv4 explicitly so localhost resolution cannot select another service.
const apiTarget = "http://127.0.0.1:8080";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": new URL("./src", import.meta.url).pathname.replace(/^\/(.:)/, "$1"),
    },
  },
  server: {
    // Send these paths to Go while developing locally.
    proxy: {
      "/api": {
        target: apiTarget,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
