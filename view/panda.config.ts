import { defineConfig } from "@pandacss/dev";

export default defineConfig({
  // Whether to use css reset
  preflight: true,

  // Where to look for your css declarations
  include: [
    "./src/components/**/*.{ts,tsx,js,jsx}",
    "./src/app/**/*.{ts,tsx,js,jsx}",
  ],

  // Files to exclude
  exclude: [],

  jsxFramework: "react",

  // Useful for theme customization
  theme: {
    extend: {
      tokens: {
        colors: {
          BLACK: { value: "#292524" },
          WHITE: { value: "#fff" },
          GHOUST_WHITE: { value: "#f5f5f5" },
          RED: { value: "#ef4444" },
          PINK: { value: "#f472b6" },
          LIGHT_PINK: { value: "#fdf2f8" },
          YELLOW: { value: "#fef08a" },
          GRAY: { value: "#e5e5e5" },
          LIGHT_BLUE: { value: "#38bdf8" },
        },
      },
      breakpoints: {
        md: "720px",
        lg: "1090px",
      },
    },
  },

  // The output directory for your css system
  outdir: "styled-system",
});
