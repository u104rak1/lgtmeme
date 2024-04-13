import React from "react";
import { render, screen } from "@testing-library/react";
import AdminPage from "@/app/admin/page";

jest.mock("@/components/organisms/AdminImageGallery/AdminImageGallery", () => ({
  __esModule: true,
  default: jest.fn(() => <div>MockAdminImageGallery</div>),
}));

describe("AdminPage", () => {
  it("Page is rendered", () => {
    render(<AdminPage />);
    expect(screen.getByText("MockAdminImageGallery")).toBeInTheDocument();
  });
});
