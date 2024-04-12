"use client";

import Button from "@/components/atoms/Button/Button";
import Svg from "@/components/atoms/Svg/Svg";
import ImageGallery from "@/components/organisms/ImageGallery/ImageGallery";
import { css } from "@@/styled-system/css";
import { PAGE_ENDPOINTS } from "@/utils/constants";

const HomePage = () => {
  return (
    <div>
      <ImageGallery css={imageGalleryCss} />
      <a href={PAGE_ENDPOINTS.new}>
        <Button css={buttonCss} icon={<Svg icon="plus" color="white" />}>
          New Image
        </Button>
      </a>
    </div>
  );
};

const imageGalleryCss = css({
  margin: "auto",
  paddingX: "3",
  lg: { maxWidth: "1090px" },
  md: { maxWidth: "730px" },
  sm: { maxWidth: "370px" },
});
const buttonCss = css({ position: "fixed", bottom: "10", right: "5" });

export default HomePage;
