import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  value: string;
  onChange: (value: string) => void;
};

const ColorInput = ({ css, value, onChange }: Props) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value);
  };
  return (
    <div className={css}>
      <input
        className={inputCss}
        type="color"
        list="datalist"
        value={value}
        onChange={handleChange}
      />
      <datalist id="datalist">
        <option value="#ffffff" />
        {/** White */}
        <option value="#f43f5e" />
        {/** Rose */}
        <option value="#ec4899" />
        {/** Pink */}
        <option value="#a855f7" />
        {/** Purple */}
        <option value="#3730a3" />
        {/** Indigo */}
        <option value="#3b82f6" />
        {/** Blue */}
        <option value="#0ea5e9" />
        {/** Light blue */}
        <option value="#14b8a6" />
        {/** Teal */}
        <option value="#22c55e" />
        {/** Green */}
        <option value="#a3e635" />
        {/** Light green */}
        <option value="#fde047" />
        {/** Yellow */}
        <option value="#f97316" />
        {/** Orange */}
        <option value="#ef4444" />
        {/** Red */}
        <option value="#64748b" />
        {/** Gray */}
        <option value="#000000" />
        {/** Black */}
      </datalist>
    </div>
  );
};

const inputCss = css({
  cursor: "pointer",
  width: "150px",
  height: "8",
  marginTop: "9px",
  border: "1px solid #cccccc",
  borderRadius: "sm",
  bgColor: "WHITE",
});

export default ColorInput;
