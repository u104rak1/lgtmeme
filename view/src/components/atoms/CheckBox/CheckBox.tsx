import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  label?: string;
  checked: boolean;
  onChange: () => void;
};

const CheckBox = ({ css, label, checked, onChange }: Props) => {
  const handleChange = () => onChange();
  return (
    <div className={css}>
      <input
        className={inputCss}
        type="checkbox"
        id="check"
        checked={checked}
        onChange={handleChange}
      />
      <label htmlFor="check" className={a}>
        <span className={checkCss}>{checked ? "✓" : ""}</span>
        <div className={labelCss}>{label}</div>
      </label>
    </div>
  );
};

const inputCss = css({ display: "none" });
const a = css({ height: "22px", display: "flex" });
const checkCss = css({
  display: "inline-block",
  width: "22px",
  height: "22px",
  border: "1px solid #cccccc",
  borderRadius: "sm",
  cursor: "pointer",
  paddingLeft: "3px",
});
const labelCss = css({ height: "22px", marginLeft: "2", cursor: "pointer" });

export default CheckBox;
