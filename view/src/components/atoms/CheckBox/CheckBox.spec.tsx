import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import CheckBox from "@/components/atoms/CheckBox/CheckBox";

describe("CheckBox", () => {
  let onChangeMock: jest.Mock<any, any, any>;
  const checked = false;
  const label = "Label";
  beforeEach(() => {
    onChangeMock = jest.fn();
  });
  afterEach(() => {
    onChangeMock.mockReset();
  });
  test("CheckBox is rendered", () => {
    render(
      <CheckBox label={label} checked={checked} onChange={onChangeMock} />
    );
    const checkbox = screen.getByRole("checkbox", { name: label });
    expect(checkbox).toBeInTheDocument();
  });
  test("onClick is called", async () => {
    render(
      <CheckBox label={label} checked={checked} onChange={onChangeMock} />
    );
    const checkbox = screen.getByRole("checkbox", { name: label });
    await userEvent.click(checkbox);
    expect(onChangeMock).toHaveBeenCalledTimes(1);
  });
});
