import { cva } from "@@/styled-system/css";

type Icon = "search" | "copy" | "flag" | "heart" | "plus" | "upload";

type Props = {
  css?: string;
  icon: Icon;
  color?: "black" | "white" | "pink";
  size?: "md" | "lg";
  fillStyle?: "solid" | "outline";
};

const Svg = ({ css, icon, color, size, fillStyle }: Props) => {
  const iconMap: Map<Icon, string> = new Map([
    [
      "search",
      "M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z",
    ],
    [
      "copy",
      "M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 01-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 011.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 00-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 01-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 00-3.375-3.375h-1.5a1.125 1.125 0 01-1.125-1.125v-1.5a3.375 3.375 0 00-3.375-3.375H9.75",
    ],
    [
      "flag",
      "M3 3v1.5M3 21v-6m0 0l2.77-.693a9 9 0 016.208.682l.108.054a9 9 0 006.086.71l3.114-.732a48.524 48.524 0 01-.005-10.499l-3.11.732a9 9 0 01-6.085-.711l-.108-.054a9 9 0 00-6.208-.682L3 4.5M3 15V4.5",
    ],
    [
      "heart",
      "M11.645 20.91l-.007-.003-.022-.012a15.247 15.247 0 01-.383-.218 25.18 25.18 0 01-4.244-3.17C4.688 15.36 2.25 12.174 2.25 8.25 2.25 5.322 4.714 3 7.688 3A5.5 5.5 0 0112 5.052 5.5 5.5 0 0116.313 3c2.973 0 5.437 2.322 5.437 5.25 0 3.925-2.438 7.111-4.739 9.256a25.175 25.175 0 01-4.244 3.17 15.247 15.247 0 01-.383.219l-.022.012-.007.004-.003.001a.752.752 0 01-.704 0l-.003-.001z",
    ],
    ["plus", "M12 4.5v15m7.5-7.5h-15"],
    [
      "upload",
      "M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5",
    ],
  ]);
  const d = iconMap.get(icon);
  const fill = fillStyle === "solid" ? "currentColor" : "none";
  return (
    <div className={css}>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill={fill}
        strokeWidth={1.5}
        stroke="currentColor"
        viewBox="0 0 25 25"
        className={svgRecipe({ color, size })}
        data-testid="icon"
      >
        <path strokeLinecap="round" strokeLinejoin="round" d={d} />
      </svg>
    </div>
  );
};

const svgRecipe = cva({
  variants: {
    color: {
      black: { color: "BLACK" },
      white: { color: "WHITE" },
      pink: { color: "PINK" },
    },
    size: {
      xs: { width: "15px", height: "15px" },
      sm: { width: "15px", height: "15px" },
      md: { width: "20px", height: "20px" },
      lg: { width: "25px", height: "25px" },
    },
  },
  defaultVariants: { color: "black", size: "md" },
});

export default Svg;
