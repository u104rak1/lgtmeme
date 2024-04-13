import { CommonService } from "@/services";

describe("CommonService", () => {
  let service: CommonService;
  beforeEach(() => {
    service = new CommonService();
  });
  describe("createPathWithQuery", () => {
    it("should generate a path with single query parameter", () => {
      const result = service.createPathWithQuery("/api/data", { userId: 123 });
      expect(result).toBe("/api/data?user_id=123");
    });
    it("should handle multiple query parameters", () => {
      const result = service.createPathWithQuery("/api/data", {
        userId: 123,
        userName: "John",
      });
      expect(result).toBe("/api/data?user_id=123&user_name=John");
    });
    it("should handle array query parameters", () => {
      const result = service.createPathWithQuery("/api/data", {
        tags: ["news", "article"],
      });
      expect(result).toBe("/api/data?tags=news,article");
    });
  });
  describe("createConfig", () => {
    it("should create a config for a GET request", () => {
      const result = service.createConfig("GET");
      expect(result).toEqual({
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        cache: "no-store",
        body: undefined,
      });
    });
    it("should create a config with a body for POST requests", () => {
      const body = { data: "test" };
      const result = service.createConfig("POST", body);
      expect(result).toEqual({
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        cache: "no-store",
        body: JSON.stringify(body),
      });
    });
  });
  describe("returnUnknownError", () => {
    it("should return a standard error object", () => {
      const result = service.returnUnknownError();
      expect(result).toEqual({
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      });
    });
  });
});
