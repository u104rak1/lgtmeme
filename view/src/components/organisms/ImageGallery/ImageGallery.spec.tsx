import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import ImageGallery from "@/components/organisms/ImageGallery/ImageGallery";

global.fetch = jest.fn(() =>
  Promise.resolve(
    new Response(
      JSON.stringify({
        images: [],
        ok: true,
      })
    )
  )
);

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

Object.defineProperty(window, "localStorage", {
  value: mockLocalStorage,
});

describe("ImageGallery", () => {
  beforeEach(() => {
    window.localStorage.clear();
    (global.fetch as jest.Mock).mockClear();
  });

  test("コンポーネントが正しくレンダリングされる", async () => {
    (global.fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            images: [
              {
                id: "1",
                url: "https://placehold.jp/300x300.png",
                keyword: "",
                usedCount: 0,
                reported: false,
                confirmed: false,
                createdAt: "2021-01-01T00:00:00.000Z",
              },
            ],
            ok: true,
          })
        )
      )
    );

    render(<ImageGallery />);

    await waitFor(() => {
      expect(screen.getByText("Latest")).toBeInTheDocument();
      expect(screen.getByText("Popular")).toBeInTheDocument();
      expect(screen.getByText("Favorite")).toBeInTheDocument();

      // 結果のアサート
      const imageElement = screen.getByAltText("http://example.com/image1.png");
      expect(imageElement).toBeInTheDocument();
    });
  });
});
