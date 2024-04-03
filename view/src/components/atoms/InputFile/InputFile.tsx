import { useRef } from "react";
import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  onChange: (file: File) => void;
};

const FileInput = ({ css, onChange }: Props) => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) onChange(file);
  };
  return (
    <div className={css}>
      <label className={labelCss}>
        Select file
        <input
          className={inputCss}
          type="file"
          accept="image/*"
          ref={inputRef}
          onChange={handleChange}
        />
      </label>
    </div>
  );
};

const labelCss = css({
  paddingY: "1.5",
  paddingX: "10",
  color: "WHITE",
  bgColor: "BLACK",
  cursor: "pointer",
  borderRadius: "sm",
  _hover: { opacity: "0.8" },
});
const inputCss = css({ display: "none" });

export default FileInput;
