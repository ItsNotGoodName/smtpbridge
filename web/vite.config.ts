import { defineConfig } from "vite"
import FullReload from 'vite-plugin-full-reload'

export default defineConfig({
  server: {
    // Disable HMR to prevent page reloads whem *.templ files change. 
    // This is because Tailwind is watching those files in PostCSS.
    // We already do hot reload via the reload-vite.local file.
    hmr: false,
  },
  plugins: [
    FullReload("reload-vite.local")
  ],
  build: {
    // generate manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: 'src/main.ts',
    },
  },
})
