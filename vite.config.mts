import { defineConfig } from "vite";
import laravel from "laravel-vite-plugin";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [
    laravel({
      input: ["resources/js/app.jsx", "resources/css/app.css"],
      publicDirectory: "dist",
      buildDirectory: ".",
      refresh: true,
    }),
    react({ include: /\.(mdx|js|jsx|ts|tsx)$/ }),
    tailwindcss(),
  ],
  esbuild: {
    jsx: "automatic",
  },
  build: {
    manifest: true, // Generate manifest.json file
    outDir: "dist",
    rollupOptions: {
      input: "resources/js/app.jsx",
      output: {
        entryFileNames: "assets/[name].[hash].js",
        chunkFileNames: "assets/[name].[hash].js",
        assetFileNames: "assets/[name].[hash].[ext]",
        manualChunks: undefined, // Disable automatic chunk splitting
      },
    },
  },
  server: {
    hmr: {
      host: "localhost",
    },
  },
  // test: {
  //   browser: {
  //     enabled: true,
  //     name: "chromium",
  //     provider: "playwright",
  //     headless: true,
  //   },
  //   setupFiles: ["./vitest.setup.tsx"],
  // },
});
