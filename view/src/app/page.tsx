"use client";

import Link from "next/link";
import Button from "@/components/atoms/Button/Button";
import Svg from "@/components/atoms/Svg/Svg";
import ImageGallery from "@/components/organisms/ImageGallery/ImageGallery";
import { ImageService } from "@/services/image.service";
import { css } from "@@/styled-system/css";
import { useEffect, useState } from "react";
import { PAGE_ENDPOINTS } from "@/utils/constants";

const HomePage = () => {
  const [images, setImages] = useState<Image[]>([]);

  const getImages = async () => {
    const service = new ImageService();
    const result = await service.getImages({
      page: 0,
      keyword: "",
      sort: "latest",
      favoriteImageIds: [],
      authCheck: false,
    });
    if (result.ok) {
      setImages(result.images);
    }
  };

  useEffect(() => {
    getImages();
  }, []);

  return (
    <div>
      <ImageGallery css={imageGalleryCss} initImages={images} />
      <Link href={PAGE_ENDPOINTS.new}>
        <Button css={buttonCss} icon={<Svg icon="plus" color="white" />}>
          New Image
        </Button>
      </Link>
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
