import { render, screen } from "@testing-library/react";
import TermsOfService from "@/app/terms-of-service/page";

describe("Terms of service page", () => {
  test("Page is rendered", () => {
    render(<TermsOfService />);
    const headings = screen.getAllByRole("heading");
    expect(headings.length).toBeGreaterThanOrEqual(1);
  });
});
