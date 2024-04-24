import Button from "@/components/atoms/Button/Button";
import Svg from "@/components/atoms/Svg/Svg";
import { IMAGE_SIZE } from "@/utils/constants";
import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  image: Image;
  onClickConfirm: () => void;
  onClickDelete: () => void;
};

const AdminImageCard = ({
  css,
  image,
  onClickConfirm,
  onClickDelete,
}: Props) => {
  const handleClickConfirm = () => onClickConfirm();
  const handleClickDelete = () => onClickDelete();
  if (!image) return <></>;
  return (
    <div className={css}>
      <div className={cardCss}>
        <div className={imageContainerCss}>
          <img
            className={imageCss}
            src={image.url}
            height={IMAGE_SIZE}
            width={IMAGE_SIZE}
            alt="LGTM"
          />
        </div>
        <div className={buttonsCss}>
          <Button
            css={buttonCss}
            size="sm"
            icon={<Svg icon="thumbUp" color="white" />}
            onClick={handleClickConfirm}
          >
            No problem
          </Button>
          <Button
            css={buttonCss}
            size="sm"
            color="red"
            icon={<Svg icon="trash" color="white" />}
            onClick={handleClickDelete}
          >
            Delete
          </Button>
        </div>
      </div>
    </div>
  );
};

const cardCss = css({
  width: "350px",
  height: "400px",
  padding: "25px",
  boxShadow: "lg",
  bgColor: "WHITE",
});
const imageContainerCss = css({
  position: "relative",
  width: "300px",
  height: "300px",
  border: "1px solid #cccccc",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  overflow: "hidden",
});
const imageCss = css({
  position: "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  objectFit: "cover",
});
const buttonsCss = css({
  display: "flex",
  marginTop: "3",
  justifyContent: "center",
});
const buttonCss = css({ marginX: "1" });

export default AdminImageCard;
