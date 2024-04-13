import React from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ReportModal from "@/components/organisms/ReportModal/ReportModal";
import { ImageService } from "@/services/image.service";

const mockImage = {
  id: "a2128761-21a8-53c6-b6cd-1578eaf12c14",
  url: "https://placehold.jp/300x300.png",
  keyword: "",
  usedCount: 0,
  reported: false,
  confirmed: false,
  createdAt: "2021-01-01T00:00:00.000Z",
};

describe("ReportModal", () => {
  let onClickCloseMock: jest.Mock<any, any, any>;
  beforeEach(() => {
    onClickCloseMock = jest.fn();
  });
  afterEach(() => {
    onClickCloseMock.mockReset();
  });
  test("ReportModal is rendered", () => {
    render(<ReportModal image={mockImage} onClickClose={onClickCloseMock} />);
    const image = screen.getByAltText("LGTM");
    const message = screen.getByText(
      "Would you like to report an image that may be inappropriate or violate copyright/privacy?"
    );
    const closeButton = screen.getByRole("button", { name: "Close" });
    const sendButton = screen.getByRole("button", { name: "Send" });
    expect(image).toBeInTheDocument();
    expect(message).toBeInTheDocument();
    expect(closeButton).toBeInTheDocument();
    expect(sendButton).toBeInTheDocument();
  });
  test("onClickCloseMock is called", async () => {
    render(<ReportModal image={mockImage} onClickClose={onClickCloseMock} />);
    const closeButton = screen.getByRole("button", { name: "Close" });
    await userEvent.click(closeButton);
    expect(onClickCloseMock).toHaveBeenCalledTimes(1);
  });
  test("Show success modal, when sending a report successfully", async () => {
    ImageService.prototype.patchImage = jest.fn(async () => {
      const response: PatchImageSuccessResult = { ok: true };
      return response;
    });
    render(<ReportModal image={mockImage} onClickClose={onClickCloseMock} />);
    const sendButton = screen.getByRole("button", { name: "Send" });
    await userEvent.click(sendButton);
    const modalMessage = screen.getByText(
      "The report was successful! Please wait a moment for the operator to confirm."
    );
    expect(modalMessage).toBeInTheDocument();
  });
  test("Show failure modal, when sending a report fails", async () => {
    ImageService.prototype.patchImage = jest.fn(async () => {
      const response: ErrResult = {
        errorCode: "unknown",
        errorMessage: "an unknown error occured",
        ok: false,
      };
      return response;
    });
    render(<ReportModal image={mockImage} onClickClose={onClickCloseMock} />);
    const sendButton = screen.getByRole("button", { name: "Send" });
    await userEvent.click(sendButton);
    const modalMessage = screen.getByText("an unknown error occured");
    expect(modalMessage).toBeInTheDocument();
  });
});
