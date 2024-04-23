import { ACTIVE_TAB_ID, PATCH_IMAGE_REQUEST_TYPE } from "@/utils/constants";
import { ImageService } from "@/services/image.service";

const image = {
  id: "a2128761-21a8-53c6-b6cd-1578eaf12c14",
  url: "https://placehold.jp/300x300.png",
  keyword: "",
  usedCount: 0,
  reported: false,
  confirmed: false,
  createdAt: "2021-01-01T00:00:00.000Z",
};

describe("ImageService", () => {
  let imageService: ImageService;
  const sampleImageUrl = "https://placehold.jp/300x300.png";
  beforeAll(() => {
    imageService = new ImageService();
  });
  describe("getImages", () => {
    let resImages: Image[];
    let query: GetImagesQuery;
    beforeEach(() => {
      resImages = [image];
      query = {
        page: 0,
        keyword: "keyword",
        sort: ACTIVE_TAB_ID.latest,
        favoriteImageIds: [],
        authCheck: false,
      };
    });
    test("positive: Return images", async () => {
      const mockResponse = { images: resImages };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: true,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      const result = await imageService.getImages(query);
      expect(result).toEqual({ ...mockResponse, ok: true });
    });
    test("negative: Return error, when bad request error occurs", async () => {
      const mockResponse = {
        errorCode: "bad_request",
        errorMessage: "bad request error occured",
      };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: false,
        status: 400,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      query.keyword = "a".repeat(51);
      const result = await imageService.getImages(query);
      expect(result).toEqual({ ...mockResponse, ok: false });
    });
    test("negative: Return error, when unknown error occurs", async () => {
      const mockFetch = jest.fn().mockRejectedValue(new Error());
      global.fetch = mockFetch;
      const result = await imageService.getImages(query);
      expect(result).toEqual({
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      });
    });
  });
  describe("postImage", () => {
    let body: PostImageReqBody;
    beforeEach(() => {
      body = {
        base64image: "data:image/png;base64,iVBORw0KGgoA",
        keyword: "keyword",
      };
    });
    test("positive: Return imageUrl", async () => {
      const mockResponse = { imageUrl: sampleImageUrl };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: true,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      const result = await imageService.postImage(body);
      expect(result).toEqual({ ...mockResponse, ok: true });
    });
    test("negative: Return error, when bad request error occurs", async () => {
      const mockResponse = {
        errorCode: "bad_request",
        errorMessage: "bad request error occured",
      };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: false,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      const result = await imageService.postImage(body);
      expect(result).toEqual({ ...mockResponse, ok: false });
    });
    test("negative: Return error, when unknown error occurs", async () => {
      const mockFetch = jest.fn().mockRejectedValue(new Error());
      global.fetch = mockFetch;
      const result = await imageService.postImage(body);
      expect(result).toEqual({
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      });
    });
  });
  describe("patchImage", () => {
    let id: string;
    let body: PatchImageReqBody;
    beforeEach(() => {
      id = "a2128761-21a8-53c6-b6cd-1578eaf12c14";
      body = { type: PATCH_IMAGE_REQUEST_TYPE.used };
    });
    test("positive: Return ok", async () => {
      const mockFetch = jest.fn().mockResolvedValue({
        ok: true,
        json: async () => {},
      });
      global.fetch = mockFetch;
      const result = await imageService.patchImage(id, body);
      expect(result).toEqual({ ok: true });
    });
    test("negative: Return error, when bad request error occurs", async () => {
      const mockResponse = {
        errorCode: "bad_request",
        errorMessage: "bad request error occured",
      };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: false,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      const result = await imageService.patchImage(id, body);
      expect(result).toEqual({ ...mockResponse, ok: false });
    });
    test("negative: Return error, when unknown error occurs", async () => {
      const mockFetch = jest.fn().mockRejectedValue(new Error());
      global.fetch = mockFetch;
      const result = await imageService.patchImage(id, body);
      expect(result).toEqual({
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      });
    });
  });
  describe("deleteImage", () => {
    let id: string;
    beforeEach(() => {
      id = "a2128761-21a8-53c6-b6cd-1578eaf12c14";
    });
    test("positive: Return ok", async () => {
      const mockFetch = jest.fn().mockResolvedValue({
        ok: true,
        json: async () => {},
      });
      global.fetch = mockFetch;
      const result = await imageService.deleteImage(id);
      expect(result).toEqual({ ok: true });
    });
    test("negative: Return error, when not found error occurs", async () => {
      const mockResponse = {
        errorCode: "not_found",
        errorMessage: "not found error occured",
      };
      const mockFetch = jest.fn().mockResolvedValue({
        ok: false,
        json: async () => mockResponse,
      });
      global.fetch = mockFetch;
      const result = await imageService.deleteImage(id);
      expect(result).toEqual({ ...mockResponse, ok: false });
    });
    test("positive: Return error, when unknown error occurs", async () => {
      const mockFetch = jest.fn().mockRejectedValue(new Error());
      global.fetch = mockFetch;
      const result = await imageService.deleteImage(id);
      expect(result).toEqual({
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      });
    });
  });
});
