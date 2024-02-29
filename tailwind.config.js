/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./template/**/*.templ"
    ],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/forms')
    ],
}

