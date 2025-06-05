/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#EBF3FF',
          100: '#D6E7FF',
          200: '#ADC9FF',
          300: '#85ABFF',
          400: '#5C8DFF',
          500: '#1A365D',
          600: '#0F2D5A',
          700: '#0A234E',
          800: '#051A42',
          900: '#021036',
        },
        accent: {
          50: '#FFFAEB',
          100: '#FFF1C8',
          200: '#FFE390',
          300: '#FFD559',
          400: '#FFC721',
          500: '#F59E0B',
          600: '#D97706',
          700: '#B45309',
          800: '#92400E',
          900: '#783010',
        },
        success: {
          500: '#10B981',
        },
        warning: {
          500: '#FBBF24',
        },
        error: {
          500: '#EF4444',
        },
      },
      fontFamily: {
        serif: ['Georgia', 'Cambria', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'card': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
      },
    },
  },
  plugins: [],
};