import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Button from "@/components/atoms/Button/Button";

describe("Button", () => {
  let onClickMock: jest.Mock<any, any, any>;
  const text = "Click me";
  beforeEach(() => {
    onClickMock = jest.fn();
  });
  afterEach(() => {
    onClickMock.mockReset();
  });
  test("Button is rendered", () => {
    render(<Button onClick={onClickMock}>{text}</Button>);
    const button = screen.getByRole("button", { name: text });
    expect(button).toBeInTheDocument();
  });
  test("onClick is called", async () => {
    render(<Button onClick={onClickMock}>{text}</Button>);
    const button = screen.getByRole("button", { name: text });
    await userEvent.click(button);
    expect(onClickMock).toHaveBeenCalledTimes(1);
  });
});
