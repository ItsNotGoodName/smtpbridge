import { defineConfig } from "vite"
import FullReload from 'vite-plugin-full-reload'
import inject from '@rollup/plugin-inject'

export default defineConfig({
  plugins: [
    FullReload(['views/**/*', "controllers/**"]),
    inject({
      htmx: 'htmx.org',
      htmx_loading_states: "htmx.org/dist/ext/loading-states.js"
    }),
  ],
  build: {
    // generate manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      // also change in ./views/layouts/index.html
      input: 'src/main.ts',
    },
  },
})
