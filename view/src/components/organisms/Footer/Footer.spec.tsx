import { render, screen } from "@testing-library/react";
import Footer from "@/components/organisms/Footer/Footer";
import { PAGE_ENDPOINTS } from "@/utils/constants";

describe("Footer", () => {
  test("Footer is rendered", () => {
    render(<Footer />);
    const homeButton = screen.getByRole("button", { name: "Home" });
    const homeLink = homeButton.closest("a");
    const termsOfServiceButton = screen.getByRole("button", {
      name: "Terms of service",
    });
    const termsOfServiceLink = termsOfServiceButton.closest("a");
    const privacyPolicyButton = screen.getByRole("button", {
      name: "Privacy policy",
    });
    const privacyPolicyLink = privacyPolicyButton.closest("a");
    const adminButton = screen.getByRole("button", {
      name: "Admin",
    });
    const adminLink = adminButton.closest("a");
    const copyrightText = screen.getByText(`©2024 ~ LGTMeme`);
    expect(homeButton).toBeInTheDocument();
    expect(homeLink).toHaveAttribute("href", PAGE_ENDPOINTS.home);
    expect(termsOfServiceButton).toBeInTheDocument();
    expect(termsOfServiceLink).toHaveAttribute(
      "href",
      PAGE_ENDPOINTS.termsOfService
    );
    expect(privacyPolicyButton).toBeInTheDocument();
    expect(privacyPolicyLink).toHaveAttribute(
      "href",
      PAGE_ENDPOINTS.privacyPolicy
    );
    expect(adminLink).toHaveAttribute("href", PAGE_ENDPOINTS.admin);
    expect(copyrightText).toBeInTheDocument();
  });
});
