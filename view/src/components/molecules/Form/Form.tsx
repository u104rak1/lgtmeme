import { css, cva } from "@@/styled-system/css";

type Props = {
  css?: string;
  label: string;
  isUnderLine?: boolean;
  children: JSX.Element;
};

const Form = ({ css, label, children, isUnderLine }: Props) => {
  return (
    <div className={css}>
      <div className={formCss({ isUnderLine })}>
        <div className={labelCss}>{label}</div>
        <div className={inputCss}>{children}</div>
      </div>
    </div>
  );
};

const formCss = cva({
  base: {
    borderX: "1px solid #e5e5e5",
    borderTop: "1px solid #e5e5e5",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    margin: "auto",
    textAlign: "center",
    width: "350px",
    height: "50px",
    lineHeight: "50px",
    bgColor: "WHITE",
  },
  variants: {
    isUnderLine: {
      true: { borderBottom: "1px solid #e5e5e5" },
    },
  },
  defaultVariants: { isUnderLine: false },
});
const labelCss = css({
  width: "40%",
  textAlign: "center",
  borderRight: "1px solid #e5e5e5",
});
const inputCss = css({ width: "60%", textAlign: "center" });

export default Form;
