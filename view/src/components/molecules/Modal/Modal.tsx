import Button from "@/components/atoms/Button/Button";
import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  message: string;
  show: boolean;
  onClick: () => void;
};

const Modal = ({ css, message, show, onClick }: Props) => {
  const handleClick = () => onClick();
  return (
    show && (
      <div className={css}>
        <div className={backgroundCss}>
          <div className={modalCss}>
            <div className={messageCss}>
              <p>{message}</p>
            </div>
            <Button
              css={buttonCss}
              visual="text"
              size="sm"
              onClick={handleClick}
            >
              Close
            </Button>
          </div>
        </div>
      </div>
    )
  );
};

const backgroundCss = css({
  width: "100vw",
  height: "100vh",
  bgColor: "rgba(0, 0, 0, 0.7)",
  position: "fixed",
  top: 0,
  right: 0,
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  zIndex: 9999,
});
const modalCss = css({
  width: "300px",
  minHeight: "120px",
  borderRadius: "lg",
  bgColor: "WHITE",
  boxShadow: "rgba(0, 0, 0, 0.35) 0px 5px 15px",
  display: "flex",
  flexDirection: "column",
  paddingX: "25px",
  paddingTop: "15px",
});
const messageCss = css({
  width: "250px",
  minHeight: "68px",
  whiteSpace: "pre-wrap",
});
const buttonCss = css({ marginX: "auto" });

export default Modal;
