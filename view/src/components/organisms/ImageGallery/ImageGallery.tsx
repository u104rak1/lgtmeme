import { useEffect, useState } from "react";
import ReportModal from "@/components/organisms/ReportModal/ReportModal";
import Button from "@/components/atoms/Button/Button";
import InputText from "@/components/atoms/InputText/InputText";
import Loading from "@/components/atoms/Loading/Loading";
import Svg from "@/components/atoms/Svg/Svg";
import Tabs from "@/components/atoms/Tabs/Tabs";
import ImageCard from "@/components/molecules/ImageCard/ImageCard";
import Modal from "@/components/molecules/Modal/Modal";
import {
  ACTIVE_TAB_ID,
  MAX_IMAGES_FETCH_COUNT,
  PATCH_IMAGE_REQUEST_TYPE,
} from "@/utils/constants";
import { ImageService } from "@/services/image.service";
import { css } from "@@/styled-system/css";

export const LOCAL_STORAGE_KEY_FAVORITE_IMAGE_IDS = "favoriteImageIds";

type Props = {
  css?: string;
};

const ImageGallery = ({ css }: Props) => {
  const [images, setImages] = useState<Image[]>([]);
  const [page, setPage] = useState(0);
  const [keyword, setKeyword] = useState("");
  const [activeTabId, setActiveTabId] = useState<ActiveTabId>(
    ACTIVE_TAB_ID.latest
  );
  const [favoriteImageIds, setFavoriteImageIds] = useState<string[]>([]);
  const [isFull, setIsFull] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [modal, setModal] = useState({ message: "", show: false });

  useEffect(() => {
    const favoriteImageIds = JSON.parse(
      localStorage.getItem(LOCAL_STORAGE_KEY_FAVORITE_IMAGE_IDS) || "[]"
    ) as string[];
    setFavoriteImageIds(favoriteImageIds);
    handleGetImages([], 0, "", ACTIVE_TAB_ID.latest, favoriteImageIds);
  }, []);

  const tabs: { id: ActiveTabId; label: string }[] = [
    { id: ACTIVE_TAB_ID.latest, label: "Latest" },
    { id: ACTIVE_TAB_ID.popular, label: "Popular" },
    { id: ACTIVE_TAB_ID.favorite, label: "Favorite" },
  ];

  const handleGetImages = async (
    images: Image[],
    page: number,
    keyword: string,
    activeTabId: ActiveTabId,
    favoriteImageIds: string[]
  ) => {
    setIsLoading(true);
    setPage(page);
    setActiveTabId(activeTabId);
    const service = new ImageService();
    const result = await service.getImages({
      page,
      keyword,
      sort:
        activeTabId === ACTIVE_TAB_ID.popular
          ? ACTIVE_TAB_ID.popular
          : ACTIVE_TAB_ID.latest,
      favoriteImageIds:
        activeTabId === ACTIVE_TAB_ID.favorite ? favoriteImageIds : [],
      authCheck: false,
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

  const handleClickTab = (id: string) => {
    setIsFull(false);
    setImages([]);
    handleGetImages(images, 0, keyword, id as ActiveTabId, favoriteImageIds);
  };

  const handleSetKeyword = (value: string) => setKeyword(value);

  const handlePressEnter = () => {
    setIsFull(false);
    setImages([]);
    handleGetImages(
      images,
      0,
      keyword,
      activeTabId as ActiveTabId,
      favoriteImageIds
    );
  };

  const handleCopyToClipboard = async (image: Image) => {
    try {
      await navigator.clipboard.writeText(`![LGTM](${image.url})`);
      setModal({ message: "Copied to clipboard!", show: true });
    } catch {
      setModal({
        message: "Failed to copy clipboard. Please try again later.",
        show: true,
      });
    }
    const service = new ImageService();
    service.patchImage(image.id, {
      type: PATCH_IMAGE_REQUEST_TYPE.used,
    });
  };

  const handleToggleFavorite = (isFavorite: boolean, image: Image) => {
    const newIsFavorite = !isFavorite;
    const newFavoriteImageIds = newIsFavorite
      ? [...favoriteImageIds, image.id]
      : favoriteImageIds.filter((id) => id !== image.id);
    localStorage.setItem(
      LOCAL_STORAGE_KEY_FAVORITE_IMAGE_IDS,
      JSON.stringify(newFavoriteImageIds)
    );
    setFavoriteImageIds(newFavoriteImageIds);
  };

  const handleCloseModal = () => setModal({ message: "", show: false });

  const [reportImage, setReportImage] = useState<Image | null>(null);
  const handleOpenReportModal = async (image: Image) => setReportImage(image);
  const handleCloseReportModal = () => {
    if (!reportImage) return;
    const imageId = reportImage.id;
    setImages(images.filter((i) => i.id !== imageId));
    setReportImage(null);
  };

  return (
    <div className={css}>
      <Tabs
        css={tabCss}
        tabs={tabs}
        value={activeTabId}
        onClick={handleClickTab}
      />
      <InputText
        css={textBoxCss}
        value={keyword}
        placeholder="Keyword"
        icon={<Svg icon="search" />}
        onChange={handleSetKeyword}
        onEnterPress={handlePressEnter}
      />
      <div className={imageCardsCss}>
        {images.map((i) => (
          <ImageCard
            css={imageCardCss}
            key={i.id}
            image={i}
            isFavorite={favoriteImageIds.some((id) => id === i.id)}
            onClickCopy={() => handleCopyToClipboard(i)}
            onClickFavorite={(isFavorite: boolean) =>
              handleToggleFavorite(isFavorite, i)
            }
            onClickReport={() => handleOpenReportModal(i)}
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
          onClick={() =>
            handleGetImages(
              images,
              page + 1,
              keyword,
              activeTabId,
              favoriteImageIds
            )
          }
        >
          See more
        </Button>
      )}
      <Modal {...modal} onClick={handleCloseModal} />
      <ReportModal image={reportImage} onClickClose={handleCloseReportModal} />
    </div>
  );
};

const tabCss = css({ paddingTop: "8", paddingBottom: "4" });
const textBoxCss = css({ paddingBottom: "8" });
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

export default ImageGallery;
