import { defineRecipe } from '@pandacss/dev';

export const buttonRecipe = defineRecipe({
	className: 'button',
	description: 'Styles for button component',
	base: {
		padding: '8px 16px',
    height: '40px',
		boxShadow: 'token(shadows.realistic)',
    cursor: 'pointer',
    backgroundColor: 'colorPalette.500',
    borderRadius: '4px',
    borderColor: 'colorPalette.500',
    color: 'white',
    '&:hover': {
      backgroundColor: 'colorPalette.700'
    }
	},
	variants: {
		color: {
			blue: {
        colorPalette: 'blue'
			},
      green: {
        colorPalette: 'green'
      },
      red: {
        colorPalette: 'red'
      }
		},
    visual: {
      outlined: {
        border: '1px solid',
        backgroundColor: 'white',
        color: 'colorPalette.500',
      },
      solid: {
        
      }
    },
    shape: {
      boxy: {
        borderRadius: '4px'
      },
      rounded: {
        borderRadius: 'full',
        padding: '10px'
      }
    }

	},
	defaultVariants: {
    color: 'blue',
    visual: 'solid',
    shape: 'boxy',
	}
});


