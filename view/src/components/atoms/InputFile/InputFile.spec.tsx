import { render, screen, fireEvent } from "@testing-library/react";
import InputFile from "@/components/atoms/InputFile/InputFile";

describe("InputFile", () => {
  let onChangeMock: jest.Mock<any, any, any>;
  beforeEach(() => {
    onChangeMock = jest.fn();
  });
  afterEach(() => {
    onChangeMock.mockReset();
  });
  test("InputFile is rendered", () => {
    render(<InputFile onChange={onChangeMock} />);
    const fileInput = screen.getByLabelText("Select file");
    expect(fileInput).toBeInTheDocument();
  });
  test("onChange is called", () => {
    render(<InputFile onChange={onChangeMock} />);
    const fileInput = screen.getByLabelText("Select file");
    const file = new File(["file content"], "file.png", {
      type: "image/png",
    });
    fireEvent.change(fileInput, { target: { files: [file] } });
    expect(onChangeMock).toHaveBeenCalledWith(file);
  });
});
