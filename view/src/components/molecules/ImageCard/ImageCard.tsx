import Button from "@/components/atoms/Button/Button";
import Svg from "@/components/atoms/Svg/Svg";
import { IMAGE_SIZE } from "@/utils/constants";
import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
  image: Image;
  isFavorite: boolean;
  onClickCopy: () => void;
  onClickFavorite: (isFavorite: boolean) => void;
  onClickReport: () => void;
};

const ImageCard = ({
  css,
  image,
  isFavorite,
  onClickCopy,
  onClickFavorite,
  onClickReport,
}: Props) => {
  const handleClickCopy = () => onClickCopy();
  const handleClickFavorite = () => onClickFavorite(isFavorite);
  const handleClickReport = () => onClickReport();
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
            icon={<Svg icon="copy" color="white" />}
            onClick={handleClickCopy}
          >
            Copy
          </Button>
          <Button
            css={buttonCss}
            size="sm"
            color="lightPink"
            icon={
              <Svg
                icon="heart"
                color="pink"
                fillStyle={isFavorite ? "solid" : "outline"}
              />
            }
            onClick={handleClickFavorite}
          >
            Favorite
          </Button>
          {!image.reported && (
            <Button
              css={buttonCss}
              size="sm"
              color="yellow"
              icon={<Svg icon="flag" />}
              onClick={handleClickReport}
            >
              Report
            </Button>
          )}
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

export default ImageCard;
