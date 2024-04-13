import { render, screen } from "@testing-library/react";
import Svg from "@/components/atoms/Svg/Svg";

describe("Svg", () => {
  test("Svg is rendered", async () => {
    render(<Svg icon="search" />);
    const icon = await screen.findByTestId("icon");
    expect(icon).toBeInTheDocument();
  });
});
