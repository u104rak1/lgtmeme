import { LoginService } from "@/services/login.service";

describe("LoginService", () => {
  let loginService = new LoginService();
  beforeEach(() => {
    loginService = new LoginService();
  });
  it("positive: Return redirectURL", async () => {
    const mockResponse = { redirectURL: "http://example.com/home" };
    const mockFetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => mockResponse,
    });
    global.fetch = mockFetch;

    const result = await loginService.postLogin("user", "password");

    expect(result).toEqual({ ...mockResponse, ok: true });
  });
  it("negative: Return error, when bad request error occurs", async () => {
    const mockResponse = {
      errorCode: "bad_request",
      errorMessage: "bad request error occured",
    };
    const mockFetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => mockResponse,
    });
    global.fetch = mockFetch;

    const result = await loginService.postLogin("user", "password");

    expect(result).toEqual({ ...mockResponse, ok: true });
  });
  it("negative: Return error, when unknown error occurs", async () => {
    const mockFetch = jest.fn().mockRejectedValue(new Error());
    global.fetch = mockFetch;
    global.fetch = mockFetch;
    const result = await loginService.postLogin("user", "password");
    expect(result).toEqual({
      errorCode: "unknown",
      errorMessage: "an unknown error occured",
      ok: false,
    });
  });
});
