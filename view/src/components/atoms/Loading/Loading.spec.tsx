import { render, screen } from "@testing-library/react";
import Loading from "@/components/atoms/Loading/Loading";

describe("Loading", () => {
  test("Loading is rendered", async () => {
    render(<Loading />);
    const loading = await screen.findByTestId("loading");
    expect(loading).toBeInTheDocument();
  });
});
