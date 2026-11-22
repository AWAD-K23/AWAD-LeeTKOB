// Configures browser tests and starts the frontend for them.
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  reporter: "list",
  use: { baseURL: "http://127.0.0.1:5173", ...devices["Desktop Chrome"] },
  // Playwright starts the dev server before running browser tests.
  webServer: {
    command: "npm run dev -- --host 127.0.0.1",
    url: "http://127.0.0.1:5173",
    reuseExistingServer: !process.env.CI,
    timeout: 30_000,
  },
});
