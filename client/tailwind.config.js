/** @type {import('tailwindcss').Config} */
export default {
    content: ['index.html', './src/**/*.{js,jsx,ts,tsx}'],
    theme: {
        colors: {
            transparent: '#00000000',
            black: '#000000',
            white: '#FFFFFF',
            primary: '#00ACE6',
            primaryLight: '#CEEAF7',
            light: '#657985',
            lighter: '#8C9CA6',
            lightest: '#C7CFD3',
            background: '#F5F5F5',
            backgroundDark: '#D4DBDF',
            error: '#EA609C',
            gold: '#F9BC1F',
            errorDark: '#C63A77',
            goldDark: '#E2AA17',
            primaryDark: '#009ACE',
            bronze: '#B28511',
        },
        fontFamily: {
            text: ['Rubik', 'Franklin Gothic Medium', 'Arial Narrow', 'Arial', 'sans-serif'],
            mono: ['Roboto Mono', 'Courier New', 'Courier', 'monospace'],
        },
        extend: {
            keyframes: {
                wiggle: {
                    '0%, 100%': {
                        transform: 'skew(2deg, 5deg)',
                    },
                    '50%': {
                        transform: 'skew(-5deg, -2deg)',
                    },
                },
            },
            animation: {
                wiggle: 'wiggle 0.3s cubic-bezier(.23,.7,.3,.85) infinite',
            },
        },
    },
    plugins: [require('@tailwindcss/forms'), require('@tailwindcss/line-clamp')],
};
