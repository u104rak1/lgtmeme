import { render, screen } from "@testing-library/react";
import Header from "@/components/organisms/Header/Header";
import { PAGE_ENDPOINTS } from "@/utils/constants";

describe("Header", () => {
  test("Header is rendered", () => {
    render(<Header />);
    const h1 = screen.getByRole("heading", { name: "LGTMeme" });
    const link = screen.getByText("LGTMeme");
    expect(h1).toBeInTheDocument();
    expect(link).toHaveAttribute("href", PAGE_ENDPOINTS.home);
  });
});
