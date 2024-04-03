"use client";

import Link from "next/link";
import Button from "@/components/atoms/Button/Button";
import { CLIENT_ENDPOINTS, PAGE_ENDPOINTS } from "@/utils/constants";
import packageJson from "@@/package.json";
import { css } from "@@/styled-system/css";

const Footer = () => {
  return (
    <footer className={footerCss}>
      <div className={buttonsCss}>
        {/* <Link href={PAGE_ENDPOINTS.home}>
          <Button visual="text">Home</Button>
        </Link>
        <Link href={PAGE_ENDPOINTS.termsOfService}>
          <Button visual="text">Terms of service</Button>
        </Link>
        <Link href={PAGE_ENDPOINTS.privacyPolicy}>
          <Button visual="text">Privacy policy</Button>
        </Link> */}
        <Button
          visual="text"
          onClick={() => (window.location.href = CLIENT_ENDPOINTS.auth)}
        >
          Management
        </Button>
      </div>
      <div className={copyrightCss}>
        ©2024 ~ LGTMeme version {packageJson.version}
      </div>
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
