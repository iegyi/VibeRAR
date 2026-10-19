/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        winrar: {
          bg: '#f0f0f0',
          panel: '#ffffff',
          toolbar: '#ece9d8',
          header: '#f5f6f7',
          border: '#d0d0d0',
          darkborder: '#919b9c',
          selected: '#cce8ff',
          selectedBorder: '#99d1ff',
          highlight: '#0078d7',
          text: '#222222',
          muted: '#666666',
        }
      },
      fontFamily: {
        sans: ['"Segoe UI"', '-apple-system', 'BlinkMacSystemFont', 'Tahoma', 'sans-serif'],
        mono: ['"Consolas"', '"Courier New"', 'monospace']
      }
    },
  },
  plugins: [],
}
