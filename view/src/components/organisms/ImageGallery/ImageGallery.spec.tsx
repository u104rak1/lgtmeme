import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import ImageGallery from "@/components/organisms/ImageGallery/ImageGallery";
import { ImageService } from "@/services/image.service";

jest.mock("@/services/image.service", () => ({
  ImageService: jest.fn().mockImplementation(() => ({
    getImages: jest.fn(),
    patchImage: jest.fn(),
  })),
}));

const mockLocalStorage = (() => {
  let store: { [key: string]: string } = {};
  return {
    getItem: function (key: string) {
      return store[key] || null;
    },
    setItem: function (key: string, value: string) {
      store[key] = value.toString();
    },
    removeItem: function (key: string) {
      delete store[key];
    },
    clear: function () {
      store = {};
    },
  };
})();

const image = {
  id: "a2128761-21a8-53c6-b6cd-1578eaf12c14",
  url: "https://placehold.jp/300x300.png",
  keyword: "",
  usedCount: 0,
  reported: false,
  confirmed: false,
  createdAt: "2021-01-01T00:00:00.000Z",
};

Object.defineProperty(window, "localStorage", {
  value: mockLocalStorage,
});

describe("ImageGallery", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    window.localStorage.clear();
  });

  test("ImageGallery is rendered", async () => {
    const mockGetImages = jest.fn().mockResolvedValue({
      ok: true,
      images: [image],
    });
    (ImageService as jest.Mock).mockImplementation(() => ({
      getImages: mockGetImages,
      patchImage: jest.fn(),
    }));

    render(<ImageGallery />);

    await waitFor(() => {
      const tabComp = screen.getByText("Latest");
      const keywordInput = screen.getByPlaceholderText("Keyword");
      const imageCard = screen.getByAltText("LGTM");
      const seeMoreButton = screen.getByRole("button", { name: "See more" });
      expect(tabComp).toBeInTheDocument();
      expect(keywordInput).toBeInTheDocument();
      expect(imageCard).toBeInTheDocument();
      expect(seeMoreButton).toBeInTheDocument();
    });
  });
});
