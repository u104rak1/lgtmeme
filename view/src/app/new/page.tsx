"use client";

import { useState, useRef, useEffect } from "react";
import Button from "@/components/atoms/Button/Button";
import CheckBox from "@/components/atoms/CheckBox/CheckBox";
import InputColor from "@/components/atoms/InputColor/InputColor";
import InputFile from "@/components/atoms/InputFile/InputFile";
import InputText from "@/components/atoms/InputText/InputText";
import Loading from "@/components/atoms/Loading/Loading";
import SelectBox from "@/components/atoms/SelectBox/SelectBox";
import Svg from "@/components/atoms/Svg/Svg";
import Form from "@/components/molecules/Form/Form";
import Modal from "@/components/molecules/Modal/Modal";
import { PAGE_ENDPOINTS, IMAGE_SIZE } from "@/utils/constants";
import { ImageService } from "@/services/image.service";
import { css } from "@@/styled-system/css";

const FONT_FAMILY_OPTIONS = [
  { value: "Arial", label: "Arial" },
  { value: "Verdana", label: "Verdana" },
  { value: "Times New Roman", label: "Times New Roman" },
];

const TEXT_SIZE_SMALL = 36;
const TEXT_SIZE_MEDIUM = 60;
const TEXT_SIZE_LARGE = 84;

const TEXT_SIZE_OPTIONS = [
  { value: TEXT_SIZE_SMALL, label: "Small" },
  { value: TEXT_SIZE_MEDIUM, label: "Medium" },
  { value: TEXT_SIZE_LARGE, label: "Large" },
];

const SIZE_MAP: SizeMap = new Map([
  [
    TEXT_SIZE_SMALL,
    { size: TEXT_SIZE_SMALL, width: 102, height: TEXT_SIZE_SMALL, diff: 30 },
  ],
  [
    TEXT_SIZE_MEDIUM,
    { size: TEXT_SIZE_MEDIUM, width: 169, height: TEXT_SIZE_MEDIUM, diff: 50 },
  ],
  [
    TEXT_SIZE_LARGE,
    { size: TEXT_SIZE_LARGE, width: 235, height: TEXT_SIZE_LARGE, diff: 71 },
  ],
]);

const NewPage = () => {
  const [imageInfo, setImageInfo] = useState({ url: "", width: 0, height: 0 });
  const [textStyle, setTextStyle] = useState<TextStyle>({
    left: 66,
    top: 120,
    color: "#ffffff",
    fontSize: 60,
    width: 169,
    height: 60,
    lineHeight: "60px",
    fontFamily: "Arial",
  });
  const [keyword, setKeyword] = useState("");
  const [isDragging, setIsDragging] = useState(false);
  const [dragStartX, setDragStartX] = useState(0);
  const [dragStartY, setDragStartY] = useState(0);
  const [isUpload, setIsUpload] = useState(false);
  const [checked, setChecked] = useState(false);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  const handleImageUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const dataURL = e.target?.result as string;
      const img = new Image();
      img.src = dataURL;
      img.onload = () => {
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext("2d");
        if (ctx) {
          const maxWidth = IMAGE_SIZE;
          const maxHeight = IMAGE_SIZE;
          let canvasWidth = img.width;
          let canvasHeight = img.height;
          if (img.width > maxWidth || img.height > maxHeight) {
            const widthRatio = maxWidth / img.width;
            const heightRatio = maxHeight / img.height;
            const minRatio = Math.min(widthRatio, heightRatio);
            canvasWidth = img.width * minRatio;
            canvasHeight = img.height * minRatio;
          }
          canvas.width = maxWidth;
          canvas.height = maxHeight;
          const canvasX = (maxWidth - canvasWidth) / 2;
          const canvasY = (maxHeight - canvasHeight) / 2;
          ctx.clearRect(0, 0, maxWidth, maxHeight);
          ctx.drawImage(img, canvasX, canvasY, canvasWidth, canvasHeight);
          const newImageInfo = {
            url: canvas.toDataURL("image/jpeg"),
            width: canvasWidth,
            height: canvasHeight,
          };
          setImageInfo(newImageInfo);
        }
      };
    };
    reader.readAsDataURL(file);
  };

  const handleTextSizeChange = (value: string) => {
    const map = SIZE_MAP.get(Number(value) as SizeMapKey)!;
    setTextStyle((prev) => ({
      ...prev,
      fontSize: map.size,
      width: map.width,
      height: map.height,
      lineHeight: `${map.size}px`,
    }));
  };

  const handleTextColorChange = (value: string) => {
    setTextStyle((prev) => ({ ...prev, color: value }));
  };

  const handleFontFamilyChange = (value: string) => {
    setTextStyle((prev) => ({ ...prev, fontFamily: value }));
  };

  /** Start text drag */
  const startTextDrag = (eventX: number, eventY: number) => {
    setIsDragging(true);
    setDragStartX(eventX - textStyle.left);
    setDragStartY(eventY - textStyle.top);
  };
  const handleMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
    if (isDragging) return;
    startTextDrag(e.clientX, e.clientY);
  };
  const handleTouchStart = (e: React.TouchEvent<HTMLDivElement>) => {
    if (isDragging) return;
    const touch = e.touches[0];
    startTextDrag(touch.clientX, touch.clientY);
  };

  /** Moving text */
  const movingText = (newLeft: number, newTop: number) => {
    const parentElement = document.getElementById("canvas-border");
    if (!parentElement) return;
    const parentRect = parentElement.getBoundingClientRect();
    const maxX = parentRect.width - textStyle.width;
    const maxY = parentRect.height - textStyle.height;
    const clampedLeft = Math.min(Math.max(newLeft, 0), maxX);
    const clampedTop = Math.min(Math.max(newTop, 0), maxY);
    setTextStyle((prev) => ({ ...prev, left: clampedLeft, top: clampedTop }));
  };
  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!isDragging) return;
    const newLeft = e.clientX - dragStartX;
    const newTop = e.clientY - dragStartY;
    movingText(newLeft, newTop);
  };
  const handleTouchMove = (e: React.TouchEvent<HTMLDivElement>) => {
    if (!isDragging) return;
    const touch = e.touches[0];
    const newLeft = touch.clientX - dragStartX;
    const newTop = touch.clientY - dragStartY;
    movingText(newLeft, newTop);
  };

  /** Prevent scroll */
  useEffect(() => {
    const preventScroll = (e: TouchEvent) => {
      if (isDragging) e.preventDefault();
    };
    document.addEventListener("touchstart", preventScroll, { passive: false });
    document.addEventListener("touchmove", preventScroll, { passive: false });
    return () => {
      document.removeEventListener("touchstart", preventScroll);
      document.removeEventListener("touchmove", preventScroll);
    };
  }, [isDragging]);

  /** Stop text move */
  const handleStopTextMove = () => {
    setIsDragging(false);
  };

  const handletoggleChecked = () => setChecked((prev) => !prev);

  const [modal, setModal] = useState({ message: "", show: false });
  const handleCreateImage = async () => {
    setIsUpload(true);
    const diff = SIZE_MAP.get(textStyle.fontSize)?.diff;
    if (!canvasRef.current || !diff) return;
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.fillStyle = textStyle.color;
    ctx.font = `${textStyle.fontSize}px ${textStyle.fontFamily}`;
    const text = "LGTM";
    const textX = textStyle.left;
    const textY = textStyle.top + diff;
    ctx.fillText(text, textX, textY);

    const image = canvas.toDataURL("image/webp");
    const service = new ImageService();
    const res = await service.postImage({ base64image: image, keyword });
    if (!res.ok) {
      setIsUpload(false);
      setModal({ message: res.errorMessage, show: true });
      return;
    }

    try {
      await navigator.clipboard.writeText(`![LGTM](${res.imageUrl})`);
      setModal({
        message: "Created an image and copied it to the clipboard!",
        show: true,
      });
    } catch {
      setModal({
        message: `Created an image!\nPlease copy markdown below:\n\n![LGTM](${res.imageUrl})`,
        show: true,
      });
    } finally {
      setIsUpload(false);
    }
  };

  const handleCloseModal = async () => {
    window.location.href = "http://localhost:8080" + PAGE_ENDPOINTS.home;
  };
  return (
    <div>
      <div className={imageEditorCss}>
        <div className={gridCss}>
          <div>
            <div className={borderCss} id="canvas-border">
              <canvas ref={canvasRef}></canvas>
              {imageInfo.url !== "" && !isUpload && (
                <div
                  className={lgtmCss}
                  style={textStyle}
                  onMouseDown={handleMouseDown}
                  onMouseMove={handleMouseMove}
                  onMouseUp={handleStopTextMove}
                  onTouchStart={handleTouchStart}
                  onTouchMove={handleTouchMove}
                  onTouchEnd={handleStopTextMove}
                >
                  LGTM
                </div>
              )}
            </div>
            <InputFile css={inputFileCss} onChange={handleImageUpload} />
          </div>
          <div>
            <Form css={topFormCss} label="Keyword">
              <InputText
                css={inputTextCss}
                value={keyword}
                placeholder="ex: anime, human"
                onChange={setKeyword}
              />
            </Form>
            <Form label="Size">
              <SelectBox
                value={textStyle.fontSize}
                options={TEXT_SIZE_OPTIONS}
                onChange={handleTextSizeChange}
              />
            </Form>
            <Form label="Font family">
              <SelectBox
                value={textStyle.fontFamily}
                options={FONT_FAMILY_OPTIONS}
                onChange={handleFontFamilyChange}
              />
            </Form>
            <Form css={bottomFormCss} label="Color" isUnderLine>
              <InputColor
                css={inputColorCss}
                value={textStyle.color}
                onChange={handleTextColorChange}
              />
            </Form>
            <div className={termsOfServiceLinkWrapCss}>
              <a
                className={termsOfServiceLinkCss}
                href={
                  process.env.NEXT_PUBLIC_APP_URL +
                  PAGE_ENDPOINTS.termsOfService
                }
                target="_blank"
              >
                Please agree to the terms of service
              </a>
            </div>
            <CheckBox
              css={checkBoxCss}
              label="Agree"
              checked={checked}
              onChange={handletoggleChecked}
            />
            {isUpload ? (
              <Loading css={loadingCss} />
            ) : (
              <Button
                css={uploadButtonCss}
                icon={<Svg icon="upload" color="white" size="lg" />}
                disabled={imageInfo.url === "" || !checked}
                onClick={handleCreateImage}
              >
                Create LGTM image
              </Button>
            )}
          </div>
        </div>
        <Modal {...modal} onClick={handleCloseModal} />
      </div>
    </div>
  );
};

const imageEditorCss = css({
  paddingTop: "30px",
  marginX: "auto",
  width: "350px",
  md: { width: "700px" },
});
const gridCss = css({
  display: "grid",
  gap: "10px",
  height: "760px",
  md: {
    height: "370px",
    gridTemplateColumns: "1fr 1fr",
    gridTemplateRows: "1fr",
  },
});
const borderCss = css({
  marginTop: "3",
  border: "2px dashed #737373",
  height: "304px",
  width: "304px",
  marginX: "auto",
  marginBottom: "3",
  position: "relative",
});
const lgtmCss = css({
  position: "absolute",
  userSelect: "none",
  _hover: { cursor: "grab" },
  _active: { cursor: "grabbing" },
});
const inputFileCss = css({ textAlign: "center" });
const topFormCss = css({ marginTop: "3" });
const bottomFormCss = css({ marginBottom: "3" });
const inputColorCss = css({ height: "50px" });
const inputTextCss = css({ width: "90%", marginX: "auto" });
const uploadButtonCss = css({ marginTop: "4", textAlign: "center" });
const termsOfServiceLinkWrapCss = css({ textAlign: "center" });
const termsOfServiceLinkCss = css({
  color: "LIGHT_BLUE",
  textAlign: "center",
  _hover: {
    opacity: "0.8",
    borderBottom: "1px solid",
    borderColor: "LIGHT_BLUE",
  },
});
const checkBoxCss = css({
  marginTop: "3",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
});
const loadingCss = css({ marginTop: "4", marginX: "auto" });

export default NewPage;
