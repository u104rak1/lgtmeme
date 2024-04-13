import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import InputText from "@/components/atoms/InputText/InputText";

describe("InputText", () => {
  let onEnterPressMock: jest.Mock<any, any, any>;
  let onChangeMock: jest.Mock<any, any, any>;
  const placeholder = "Search";
  beforeEach(() => {
    onEnterPressMock = jest.fn();
    onChangeMock = jest.fn();
  });
  afterEach(() => {
    onEnterPressMock.mockReset();
    onChangeMock.mockReset();
  });
  test("InputText is rendered", () => {
    render(
      <InputText
        placeholder={placeholder}
        onChange={onChangeMock}
        onEnterPress={onEnterPressMock}
      />
    );
    const inputText = screen.getByPlaceholderText(placeholder);
    expect(inputText).toBeInTheDocument();
  });
  test("onChange is called", async () => {
    render(
      <InputText
        placeholder={placeholder}
        onChange={onChangeMock}
        onEnterPress={onEnterPressMock}
      />
    );
    const inputText = screen.getByPlaceholderText(placeholder);
    await userEvent.type(inputText, "test");
    expect(onChangeMock).toHaveBeenCalledTimes(4);
  });
  test("onEnterPress is called", async () => {
    render(
      <InputText
        placeholder={placeholder}
        onChange={onChangeMock}
        onEnterPress={onEnterPressMock}
      />
    );
    const inputText = screen.getByPlaceholderText(placeholder);
    await userEvent.type(inputText, "{enter}");
    expect(onEnterPressMock).toHaveBeenCalledTimes(1);
  });
});
