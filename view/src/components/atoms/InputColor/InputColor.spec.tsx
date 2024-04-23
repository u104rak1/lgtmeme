import { fireEvent, render, screen } from "@testing-library/react";
import InputColor from "@/components/atoms/InputColor/InputColor";

describe("InputColor", () => {
  let onChangeMock: jest.Mock<any, any, any>;
  const value = "#ffffff";
  beforeEach(() => {
    onChangeMock = jest.fn();
  });
  afterEach(() => {
    onChangeMock.mockReset();
  });
  test("InputColor is rendered", () => {
    render(<InputColor value={value} onChange={onChangeMock} />);
    const colorInput = screen.getByDisplayValue(value);
    expect(colorInput).toBeInTheDocument();
  });
  test("onChange is called", () => {
    const newColor = "#f43f5e";
    render(<InputColor value={value} onChange={onChangeMock} />);
    const colorInput = screen.getByDisplayValue(value);
    fireEvent.change(colorInput, { target: { value: newColor } });
    expect(onChangeMock).toHaveBeenCalledWith(newColor);
  });
});
