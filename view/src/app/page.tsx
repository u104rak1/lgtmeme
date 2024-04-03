"use client";

import { ImageService } from "@/services/image.service";
import { useEffect, useState } from "react";

const HomePage = () => {
  const [images, setImages] = useState<Image[]>([]);

  const getImages = async () => {
    const service = new ImageService();
    const result = await service.getImages({
      page: 0,
      keyword: "",
      sort: "latest",
      favaoriteImageIds: [],
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
      <h1>HomePage</h1>
      {images.map((image) => (
        <img key={image.id} src={image.url} alt="image" />
      ))}
    </div>
  );
};

export default HomePage;
