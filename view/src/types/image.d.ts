type Image = {
  id: string;
  url: string;
};

/** GET */
type GetImagesQuery = {
  page: number;
  keyword: string;
  sort: "popular" | "latest";
  favaoriteImageIds: string[];
  authCheck: boolean;
};
type GetImagesRespBody = {
  images: Image[];
};
type GetImagesSuccessResult = {
  images: Image[];
  ok: true;
};
type GetImagesResult = GetImagesSuccessResult | ErrResult;

/** POST */
type PostImageReqBody = {
  image: string;
  keyword: string;
};
type PostImageRespBody = {
  imageUrl: string;
};
type PostImageSuccessResult = {
  imageUrl: string;
  ok: true;
};
type PostImageResult = PostImageSuccessResult | ErrResult;

/** PATCH */
type PatchRequestType = "used" | "report" | "confirmed";
type PatchImageReqBody = {
  type: PatchRequestType;
};
type PatchImageSuccessResult = {
  ok: true;
};
type PatchImageResult = PatchImageSuccessResult | ErrResult;

/** DELETE */
type DeleteImageSuccessResult = {
  ok: true;
};
type DeleteImageResult = DeleteImageSuccessResult | ErrResult;

/** Image editor */
type TextStyle = {
  left: number;
  top: number;
  color: string;
  fontSize: SizeMapKey;
  width: number;
  height: number;
  lineHeight: string;
  fontFamily: string;
};
type TextSizeSmall = 36;
type TextSizeMedium = 60;
type TextSizeLarge = 84;
type SizeMapKey = TextSizeSmall | TextSizeMedium | TextSizeLarge;
type SizeMap = Map<
  SizeMapKey,
  {
    size: SizeMapKey;
    width: number;
    height: number;
    diff: number;
  }
>;
