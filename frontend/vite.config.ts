// Starts React and forwards API requests to the local Go server.
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

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
      "/api": "http://localhost:8080",
      "/auth": "http://localhost:8080",
      "/labs": "http://localhost:8080",
      "/classes": "http://localhost:8080",
      "/me": "http://localhost:8080",
      "/submissions": "http://localhost:8080",
    },
  },
});
