"use client";

import { useEffect, useState } from "react";
import Button from "@/components/atoms/Button/Button";
import TextBox from "@/components/atoms/InputText/InputText";
import Checkbox from "@/components/atoms/CheckBox/CheckBox";
import { css } from "@@/styled-system/css";
import { LoginService } from "@/services/login.service";
import Modal from "@/components/molecules/Modal/Modal";

type ScopesWithDesc = {
  scope: string;
  description: string;
};

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [scopeConsent, setScopeConsent] = useState(true);
  const [scopesWithDescriptions, setScopesWithDescriptions] = useState<
    ScopesWithDesc[]
  >([]);
  const [modal, setModal] = useState({ message: "", show: false });

  useEffect(() => {
    if (window) {
      const queryParams = new URLSearchParams(window.location.search);
      const scopes = queryParams.get("scopes");
      const descriptions = queryParams.get("descriptions");

      if (!scopes || !descriptions) return;
      const scopesArray = scopes.split(",");
      const descriptionsArray = descriptions.split(",");
      const combinedScopes = scopesArray.map((scope, i) => ({
        scope,
        description: descriptionsArray[i] || "",
      }));
      if (combinedScopes.length > 0) {
        setScopeConsent(false);
      }
      setScopesWithDescriptions(combinedScopes);
    }
  }, []);

  const handletoggleChecked = () => setScopeConsent((prev) => !prev);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const service = new LoginService();
    const result = await service.postLogin(username, password, scopeConsent);
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
      {scopesWithDescriptions.length > 0 && (
        <>
          <Checkbox
            css={checkBoxCss}
            label="Agree to scopes"
            checked={scopeConsent}
            onChange={handletoggleChecked}
          />
          <div className={tableCss}>
            {scopesWithDescriptions.map((item) => (
              <div className={scopeRowCss} key={item.scope}>
                <span>{item.scope}</span>
                <span>{item.description}</span>
              </div>
            ))}
          </div>
        </>
      )}
      <Button css={marginCss} disabled={!scopeConsent} type="submit">
        Login
      </Button>
      <Modal
        {...modal}
        onClick={() => setModal({ message: "", show: false })}
      />
    </form>
  );
};

const tableCss = css({
  display: "grid",
  gridTemplateColumns: "auto auto",
  alignItems: "center",
  gap: "1",
  marginTop: "5",
});

const scopeRowCss = css({
  fontSize: "12px",
  display: "contents",
  padding: "10px 0",
});

const checkBoxCss = css({
  marginTop: "5",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
});

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
