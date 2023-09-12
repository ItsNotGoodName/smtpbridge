/** @type {import('tailwindcss').Config} */
export default {
  content: ["./pages/**/*.templ", "./components/**/*.templ", "./icons/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    darkTheme: "dracula", 
    themes: [
      "garden",
      "dracula",
    ],
  },
}

