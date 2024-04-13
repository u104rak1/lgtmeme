import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import AdminImageCard from "@/components/molecules/AdminImageCard/AdminImageCard";

describe("AdminImageCard", () => {
  const imageMock = {
    id: "1",
    url: "https://placehold.jp/300x300.png",
    keyword: "",
    usedCount: 0,
    reported: false,
    confirmed: false,
    createdAt: "2021-01-01T00:00:00.000Z",
  };
  let onClickConfirm: jest.Mock<any, any, any>;
  let onClickDelete: jest.Mock<any, any, any>;
  beforeEach(() => {
    onClickConfirm = jest.fn();
    onClickDelete = jest.fn();
  });
  afterEach(() => {
    onClickConfirm.mockReset();
    onClickDelete.mockReset();
  });
  test("AdminImageCard is rendered", () => {
    render(
      <AdminImageCard
        image={imageMock}
        onClickConfirm={onClickConfirm}
        onClickDelete={onClickDelete}
      />
    );
    const image = screen.getByAltText("LGTM");
    const confirmButton = screen.getAllByRole("button")[0];
    const deleteButton = screen.getAllByRole("button")[1];
    expect(image).toBeInTheDocument();
    expect(confirmButton).toBeInTheDocument();
    expect(deleteButton).toBeInTheDocument();
  });
  test("onClickConfirm is called", async () => {
    render(
      <AdminImageCard
        image={imageMock}
        onClickConfirm={onClickConfirm}
        onClickDelete={onClickDelete}
      />
    );
    const confirmButton = screen.getAllByRole("button")[0];
    await userEvent.click(confirmButton);
    expect(onClickConfirm).toHaveBeenCalledTimes(1);
  });
  test("onClickDelete is called", async () => {
    render(
      <AdminImageCard
        image={imageMock}
        onClickConfirm={onClickConfirm}
        onClickDelete={onClickDelete}
      />
    );
    const deleteButton = screen.getAllByRole("button")[1];
    await userEvent.click(deleteButton);
    expect(onClickDelete).toHaveBeenCalledTimes(1);
  });
});
