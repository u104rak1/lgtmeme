import { css, cva } from "@@/styled-system/css";

type Props = {
  css?: string;
  visual?: "solid" | "text";
  color?: "black" | "yellow" | "lightPink" | "red";
  size?: "sm" | "md" | "lg";
  disabled?: boolean;
  icon?: JSX.Element;
  children: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

const Button = ({
  css,
  visual,
  color,
  size,
  disabled,
  icon,
  children,
  onClick,
}: Props) => {
  return (
    <div className={css}>
      <button
        className={buttonRecipe({ visual, color, size, disabled })}
        disabled={disabled}
        onClick={onClick}
      >
        {icon}
        {icon && <div className={iconCss} />}
        {children}
      </button>
    </div>
  );
};

const buttonRecipe = cva({
  base: {
    cursor: "pointer",
    display: "inline-flex",
    alignItems: "center",
    padding: "4",
    border: "none",
    outline: "none",
    textDecoration: "none",
    fontWeight: "bold",
    borderRadius: "lg",
    _hover: { opacity: 0.8 },
    _active: { opacity: 1 },
  },
  variants: {
    visual: {
      solid: {},
      text: {},
    },
    color: {
      black: { color: "BLACK" },
      yellow: { bgColor: "YELLOW", color: "BLACK" },
      lightPink: { bgColor: "LIGHT_PINK", color: "BLACK" },
      red: { bgColor: "RED", color: "WHITE" },
    },
    size: {
      sm: { padding: "2", fontSize: "sm", minWidth: "15" },
      md: { padding: "3", fontSize: "md", minWidth: "18" },
      lg: { padding: "4", fontSize: "lg", minWidth: "21" },
    },
    disabled: {
      true: { opacity: 0.5, cursor: "not-allowed" },
    },
  },
  compoundVariants: [
    {
      visual: "solid",
      color: "black",
      css: { bgColor: "BLACK", color: "WHITE" },
    },
  ],
  defaultVariants: {
    visual: "solid",
    color: "black",
    size: "md",
    disabled: false,
  },
});
const iconCss = css({ width: "1" });

export default Button;
