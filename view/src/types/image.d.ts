type Image = {
  id: string;
  url: string;
  keyword: string;
  usedCount: number;
  reported: boolean;
  confirmed: boolean;
  createdAt: string;
};

/** GET */
type GetImagesQuery = {
  page: number;
  keyword: string;
  sort: "popular" | "latest";
  favoriteImageIds: string[];
  authCheck: boolean;
  activeTabId: ActiveTabId;
};
type GetImagesRespBody = {
  images: Image[];
};
type GetImagesSuccessResult = GetImagesRespBody & { ok: true };
type GetImagesResult = GetImagesSuccessResult | ErrResult;

/** POST */
type PostImageReqBody = {
  base64image: string;
  keyword: string;
};
type PostImageRespBody = {
  imageUrl: string;
};
type PostImageSuccessResult = PostImageRespBody & { ok: true };
type PostImageResult = PostImageSuccessResult | ErrResult;

/** PATCH */
type PatchRequestType = "used" | "report" | "confirm";
type PatchImageReqBody = {
  type: PatchRequestType;
};
type PatchImageSuccessResult = { ok: true };
type PatchImageResult = PatchImageSuccessResult | ErrResult;

/** DELETE */
type DeleteImageSuccessResult = {
  ok: true;
};
type DeleteImageResult = DeleteImageSuccessResult | ErrResult;

/** Other image related types */
type ActiveTabId = "latest" | "popular" | "favorite";
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
