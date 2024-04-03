"use client";

import ImageGallery from "@/components/organisms/ImageGallery/ImageGallery";
import { ImageService } from "@/services/image.service";
import { css } from "@@/styled-system/css";
import { useEffect, useState } from "react";

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
