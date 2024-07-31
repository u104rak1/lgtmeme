import { useEffect, useState } from "react";
import Button from "@/components/atoms/Button/Button";
import Loading from "@/components/atoms/Loading/Loading";
import AdminImageCard from "@/components/molecules/AdminImageCard/AdminImageCard";
import Modal from "@/components/molecules/Modal/Modal";
import {
  ACTIVE_TAB_ID,
  MAX_IMAGES_FETCH_COUNT,
  PATCH_IMAGE_REQUEST_TYPE,
} from "@/utils/constants";
import { ImageService } from "@/services/image.service";
import { css } from "@@/styled-system/css";

type Props = {
  css?: string;
};

const AdminImageGallery = ({ css }: Props) => {
  const [images, setImages] = useState<Image[]>([]);
  const [page, setPage] = useState(0);
  const [isFull, setIsFull] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [modal, setModal] = useState({ message: "", show: false });

  useEffect(() => {
    handleGetImages([], 0);
    setIsLoading(false);
  }, []);

  const handleGetImages = async (images: Image[], page: number) => {
    setIsLoading(true);
    setPage(page);
    const service = new ImageService();
    const result = await service.getImages({
      page,
      keyword: "",
      sort: ACTIVE_TAB_ID.latest,
      favoriteImageIds: [],
      authCheck: true,
      activeTabId: ACTIVE_TAB_ID.latest,
    });
    if (!result.ok) {
      setIsLoading(false);
      setModal({ message: result.errorMessage, show: true });
      return;
    }
    if (page === 0) {
      setImages(result.images);
    } else {
      /** I am using Map to avoid duplicate images. */
      const imageMap = new Map();
      images.forEach((image) => imageMap.set(image.id, image));
      result.images.forEach((image) => imageMap.set(image.id, image));
      setImages(Array.from(imageMap.values()));
    }
    if (result.images.length < MAX_IMAGES_FETCH_COUNT) setIsFull(true);
    setIsLoading(false);
  };

  const handleConfirm = async (imageId: string) => {
    const service = new ImageService();
    const result = await service.patchImage(imageId, {
      type: PATCH_IMAGE_REQUEST_TYPE.confirm,
    });
    if (result.ok) {
      setImages(images.filter((i) => i.id !== imageId));
      setModal({ message: "Confirmation succeeded.", show: true });
    } else {
      setModal({ message: result.errorMessage, show: true });
    }
  };

  const handleDelete = async (imageId: string) => {
    const service = new ImageService();
    const result = await service.deleteImage(imageId);
    if (result.ok) {
      setImages(images.filter((i) => i.id !== imageId));
      setModal({ message: "Deletion succeeded.", show: true });
    } else {
      setModal({ message: result.errorMessage, show: true });
    }
  };

  const handleCloseModal = () => setModal({ message: "", show: false });

  return (
    <div className={css}>
      <div className={imageCardsCss}>
        {images.map((i) => (
          <AdminImageCard
            css={imageCardCss}
            key={i.id}
            image={i}
            onClickConfirm={() => handleConfirm(i.id)}
            onClickDelete={() => handleDelete(i.id)}
          />
        ))}
      </div>
      {isLoading ? (
        <Loading css={loadingCss} />
      ) : (
        <Button
          css={buttonCss}
          size="lg"
          disabled={isFull}
          onClick={() => handleGetImages(images, page + 1)}
        >
          See more
        </Button>
      )}
      <Modal {...modal} onClick={handleCloseModal} />
    </div>
  );
};

const imageCardsCss = css({
  display: "grid",
  gap: "10px",
  lg: { gridTemplateColumns: "1fr 1fr 1fr" },
  md: { gridTemplateColumns: "1fr 1fr" },
  sm: { gridTemplateColumns: "1fr" },
});
const imageCardCss = css({ marginX: "auto" });
const buttonCss = css({ paddingY: "30px", textAlign: "center" });
const loadingCss = css({ marginTop: "30px", marginX: "auto" });

export default AdminImageGallery;
