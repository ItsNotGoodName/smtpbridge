import { defineConfig } from "vite"
import FullReload from 'vite-plugin-full-reload'
import inject from '@rollup/plugin-inject'
import UnoCSS from 'unocss/vite'
import { presetUno } from 'unocss'
import presetIcons from '@unocss/preset-icons'

export default defineConfig({
  plugins: [
    FullReload(['views/**/*', "controllers/**"]),
    inject({
      htmx: 'htmx.org',
      htmx_loading_states: "htmx.org/dist/ext/loading-states.js"
    }),
    UnoCSS({
      content: {
        filesystem: ["views/**/*"]
      },
      rules: [
        [/^bg-(.*)$/, ([, c], { theme }) => {
          //@ts-ignore
          if (theme.colors[c]) {
            return {
              //@ts-ignore
              color: theme.colors[c],
            }
          }
        }],
      ],
      theme: {
        colors: {
          'pico': {
            'card': 'var(--card-background-color)',
            'background': 'var(--background-color)',
          }
        },
      },
      presets: [
        presetUno(),
        presetIcons(),
      ],
    }),
  ],
  build: {
    // generate manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      // DEPS: ./views/layouts/index.html
      input: 'src/main.ts',
    },
  },
})
