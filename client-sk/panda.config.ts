import { defineConfig } from "@pandacss/dev"
import { textareaRecipe } from "./recipies/textarea"
import { buttonRecipe } from "./recipies/button"
import { textButtonRecipe } from "./recipies/textButton"
export default defineConfig({
    // Whether to use css reset
    preflight: true,
    
    // Where to look for your css declarations
    include: ["./src/**/*.{js,ts,svelte}", "./pages/**/*.{js,svelte,ts}"],

    // Files to exclude
    exclude: [],

    // Useful for theme customization
    theme: {
      extend: {
        recipes: {
          textarea: textareaRecipe,
          button: buttonRecipe,
          textButton: textButtonRecipe
        },
        tokens: {
          shadows: {
            // string value
            subtle: { value: '0 1px 2px 0 rgba(0, 0, 0, 0.05)' },
            // composite value
            accent: {
              value: {
                offsetX: 0,
                offsetY: 4,
                blur: 4,
                spread: 0,
                color: 'rgba(0, 0, 0, 0.1)'
              }
            },
            // multiple string values
            realistic: {
              value: [
                '0 1px 2px 0 rgba(0, 0, 0, 0.15)',
                '0 1px 4px 0 rgba(0, 0, 0, 0.2)'
              ]
            }
          }
        }
      }
    },

    // The output directory for your css system
    outdir: "styled-system",
    
    
})