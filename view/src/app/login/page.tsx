"use client";

import { useState } from "react";
import Button from "@/components/atoms/Button/Button";
import TextBox from "@/components/atoms/InputText/InputText";
import { css } from "@@/styled-system/css";
import { LoginService } from "@/services/login.service";
import Modal from "@/components/molecules/Modal/Modal";

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [modal, setModal] = useState({ message: "", show: false });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const service = new LoginService();
    const result = await service.postLogin(username, password);
    if (result.ok) {
      window.location.href = result.redirectURL;
    } else {
      setModal({ message: result.errorMessage, show: true });
    }
  };

  return (
    <form className={formCss} onSubmit={handleSubmit}>
      <TextBox
        type="text"
        placeholder="username"
        value={username}
        onChange={(v) => setUsername(v)}
      />
      <TextBox
        css={marginCss}
        type="password"
        placeholder="password"
        value={password}
        onChange={(v) => setPassword(v)}
      />
      <Button css={marginCss} type="submit">
        Login
      </Button>
      <Modal
        {...modal}
        onClick={() => setModal({ message: "", show: false })}
      />
    </form>
  );
};

const formCss = css({
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  margin: "100px auto",
  padding: "5",
  border: "1px solid #ccc",
  boxShadow: "0 2px 4px rgba(0,0,0,0.1)",
  width: "300px",
  boxSizing: "border-box",
});

const marginCss = css({
  marginTop: "5",
});

export default LoginPage;
