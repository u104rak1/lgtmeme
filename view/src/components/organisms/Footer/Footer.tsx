"use client";

import Button from "@/components/atoms/Button/Button";
import { PAGE_ENDPOINTS } from "@/utils/constants";
import { css } from "@@/styled-system/css";

const Footer = () => {
  return (
    <footer className={footerCss}>
      <div className={buttonsCss}>
        <a href={PAGE_ENDPOINTS.home}>
          <Button visual="text">Home</Button>
        </a>
        <a href={PAGE_ENDPOINTS.termsOfService}>
          <Button visual="text">Terms of service</Button>
        </a>
        <a href={PAGE_ENDPOINTS.privacyPolicy}>
          <Button visual="text">Privacy policy</Button>
        </a>
        <a href={PAGE_ENDPOINTS.admin}>
          <Button visual="text">Admin</Button>
        </a>
      </div>
      <div className={copyrightCss}>©2024 ~ LGTMeme version</div>
    </footer>
  );
};

const footerCss = css({
  bgColor: "GHOUST_WHITE",
  color: "BLACK",
  maxWidth: "100vw",
  height: "210px",
  md: { height: "140px" },
});
const copyrightCss = css({ textAlign: "center" });
const buttonsCss = css({
  display: "flex",
  justifyContent: "center",
  paddingTop: "8",
});

export default Footer;
