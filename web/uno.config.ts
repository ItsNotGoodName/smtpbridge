import { defineConfig } from 'unocss'
import { presetUno } from 'unocss'
import presetIcons from '@unocss/preset-icons'

export default defineConfig({
  rules: [],
  presets: [
    presetUno(),
    presetIcons(),
  ],
  cli: {
    entry: {
      patterns: ["views/**/*.html"],
      outFile: "src/uno.css"
    },
  },
})
