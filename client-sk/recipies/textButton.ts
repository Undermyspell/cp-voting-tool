import { defineRecipe } from "@pandacss/dev";

export const textButtonRecipe = defineRecipe({
  className: 'textButton',
  description: 'Style for a textButton component',
  base: {
    cursor: 'pointer',
    padding: '4px 8px',
    color: 'colorPalette.600',
    '&:hover': {
      backgroundColor: 'colorPalette.700',
      color: 'white',
      borderRadius: '4px'
    }
  },
  variants: {
    color: {
      neutral: {
        colorPalette: 'gray'
      },
      green: {
        colorPalette: 'green'
      },
      red: {
        colorPalette: 'red'
      }
    }
  },
  defaultVariants: {
    color: 'neutral'
  }
})