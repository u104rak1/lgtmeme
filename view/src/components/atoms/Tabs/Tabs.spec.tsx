import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Tabs from "@/components/atoms/Tabs/Tabs";

describe("Tabs", () => {
  const tabs = [
    { id: "latest", label: "Latest" },
    { id: "popular", label: "Popular" },
    { id: "favorite", label: "Favorite" },
  ];
  const id = "latest";
  let onClickMock: jest.Mock<any, any, any>;
  beforeEach(() => {
    onClickMock = jest.fn();
  });
  afterEach(() => {
    onClickMock.mockReset();
  });
  test("Tabs is rendered", () => {
    render(<Tabs tabs={tabs} value={id} onClick={onClickMock} />);
    const latest = screen.getByText("Latest");
    const popular = screen.getByText("Popular");
    const favorite = screen.getByText("Favorite");
    expect(latest).toBeInTheDocument();
    expect(popular).toBeInTheDocument();
    expect(favorite).toBeInTheDocument();
  });
  test("onClick is called", async () => {
    render(<Tabs tabs={tabs} value={id} onClick={onClickMock} />);
    const popular = screen.getByText("Popular");
    await userEvent.click(popular);
    expect(onClickMock).toHaveBeenCalledWith("popular");
  });
});
