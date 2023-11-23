import { defineRecipe } from "@pandacss/dev";

export const textareaRecipe = defineRecipe({
  className: "textarea",
  description: "Styles for textarea component",
  base: {
    resize: "none",
    borderRadius: 4,
    border: "1px solid",
    borderColor: "gray.400",
    "&::placeholder": {
      color: "gray.800",
    },
    backgroundColor: "gray.50",
    padding: "0.25em 0.5em"
  }
}
)