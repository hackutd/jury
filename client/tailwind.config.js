/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./src/**/*.{js,jsx,ts,tsx}'],
    theme: {
        colors: {
            transparent: '#00000000',
            primary: '#00ACE6',
            light: '#657985',
            lighter: '#8C9CA6',
            lightest: '#C7CFD3',
            background: '#F5F5F5',
            backgroundDark: '#D4DBDF',
            error: '#EA609C'
        },
        fontFamily: {
            text: ['Rubik', 'Franklin Gothic Medium', 'Arial Narrow', 'Arial', 'sans-serif'],
            mono: ['Roboto Mono', 'Courier New', 'Courier', 'monospace'],
        },
        extend: {},
    },
    plugins: [],
};
