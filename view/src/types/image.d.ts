type Image = {
  id: string;
  url: string;
};

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
