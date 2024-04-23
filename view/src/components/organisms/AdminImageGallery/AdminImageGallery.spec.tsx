import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import AdminImageGallery from "@/components/organisms/AdminImageGallery/AdminImageGallery";
import { ImageService } from "@/services/image.service";

jest.mock("@/services/image.service", () => ({
  ImageService: jest.fn().mockImplementation(() => ({
    getImages: jest.fn(),
    patchImage: jest.fn(),
  })),
}));

const image = {
  id: "a2128761-21a8-53c6-b6cd-1578eaf12c14",
  url: "https://placehold.jp/300x300.png",
  keyword: "",
  usedCount: 0,
  reported: false,
  confirmed: false,
  createdAt: "2021-01-01T00:00:00.000Z",
};

describe("AdminImageGallery", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  test("AdminImageGallery is rendered", async () => {
    const mockGetImages = jest.fn().mockResolvedValue({
      ok: true,
      images: [image],
    });
    (ImageService as jest.Mock).mockImplementation(() => ({
      getImages: mockGetImages,
      patchImage: jest.fn(),
    }));

    render(<AdminImageGallery />);

    await waitFor(() => {
      const imageCard = screen.getByAltText("LGTM");
      const seeMoreButton = screen.getByRole("button", { name: "See more" });
      expect(imageCard).toBeInTheDocument();
      expect(seeMoreButton).toBeInTheDocument();
    });
  });
});
