/* 
 * STEP 1: Export PostCSS Configuration
 * Using ES Module syntax (export default) required because "type": "module" 
 * is declared in package.json for Tailwind v4 compatibility.
 */
export default {
  plugins: {
    /* 
     * STEP 2: Configure Tailwind CSS v4 Plugin
     * Delegates styling transformations to the @tailwindcss/postcss engine.
     */
    '@tailwindcss/postcss': {},

    /* 
     * STEP 3: Configure Autoprefixer Plugin
     * Automatically adds vendor prefixes to CSS rules for browser compatibility.
     */
    autoprefixer: {},
  },
}