import { defineConfig } from 'unocss'
import { presetUno } from 'unocss'
import presetIcons from '@unocss/preset-icons'

export default defineConfig({
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
  theme: {
    colors: {
      'pico': {
        'background': 'var(--background-color)',
      }
    },
  }
})
