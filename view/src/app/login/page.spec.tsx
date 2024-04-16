import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import LoginPage from "@/app/login/page";
import { LoginService } from "@/services/login.service";
import userEvent from "@testing-library/user-event";

jest.mock("@/services/login.service", () => ({
  LoginService: jest.fn().mockImplementation(() => ({
    postLogin: jest.fn(),
  })),
}));

describe("LoginPage", () => {
  beforeEach(() => {
    Object.defineProperty(window, "location", {
      value: {
        ...window.location,
        href: "http://localhost",
        search: "?scopes=read,write&descriptions=Read data,Write data",
      },
      writable: true,
    });
  });
  afterEach(() => {
    jest.clearAllMocks();
  });

  test("LoginPage is rendered", () => {
    render(<LoginPage />);

    const usernameTextBox = screen.getByPlaceholderText("username");
    const passwordTextBox = screen.getByPlaceholderText("password");
    const consentCheckbox = screen.getByLabelText("Agree to scopes");
    const loginButton = screen.getByRole("button", { name: "Login" });
    expect(usernameTextBox).toBeInTheDocument();
    expect(passwordTextBox).toBeInTheDocument();
    expect(consentCheckbox).toBeInTheDocument();
    expect(loginButton).toBeInTheDocument();
  });
  test("The login is successful and you will be redirected to the destination", async () => {
    const redirectUrl = "http://example.com";
    const mockPostLogin = jest.fn().mockResolvedValue({
      ok: true,
      redirectURL: redirectUrl,
    });
    (LoginService as jest.Mock).mockImplementation(() => ({
      postLogin: mockPostLogin,
    }));

    render(<LoginPage />);

    expect(window.location.href).toBe("http://localhost");

    const usernameTextBox = screen.getByPlaceholderText("username");
    const passwordTextBox = screen.getByPlaceholderText("password");
    const consentCheckbox = screen.getByLabelText("Agree to scopes");
    const loginButton = screen.getByRole("button", { name: "Login" });

    userEvent.type(usernameTextBox, "username");
    userEvent.type(passwordTextBox, "password");
    userEvent.click(consentCheckbox);
    userEvent.click(loginButton);

    await waitFor(() => {
      expect(window.location.href).toBe(redirectUrl);
    });
  });
  test("Login fails and an error message is displayed", async () => {
    const mockPostLogin = jest.fn().mockResolvedValue({
      ok: false,
      errorCode: "invalid_credentials",
      errorMessage: "Invalid credentials",
    });
    (LoginService as jest.Mock).mockImplementation(() => ({
      postLogin: mockPostLogin,
    }));

    render(<LoginPage />);

    const usernameTextBox = screen.getByPlaceholderText("username");
    const passwordTextBox = screen.getByPlaceholderText("password");
    const consentCheckbox = screen.getByLabelText("Agree to scopes");
    const loginButton = screen.getByRole("button", { name: "Login" });

    userEvent.type(usernameTextBox, "username");
    userEvent.type(passwordTextBox, "password");
    userEvent.click(consentCheckbox);
    userEvent.click(loginButton);

    await waitFor(() => {
      const modalMessage = screen.getByText("Invalid credentials");
      expect(modalMessage).toBeInTheDocument();
    });
  });
});
