"use client";

import Button from "@/components/atoms/Button/Button";
import AdminImageGallery from "@/components/organisms/AdminImageGallery/AdminImageGallery";
import { AUTH_ENDPOINTS } from "@/utils/constants";
import { css } from "@@/styled-system/css";

const AdminPage = () => {
  return (
    <div>
      <div className={buttonCss}>
        <a href={AUTH_ENDPOINTS.logout}>
          <Button visual="text">Logout</Button>
        </a>
      </div>
      <AdminImageGallery css={imageGalleryCss} />
    </div>
  );
};

const buttonCss = css({
  display: "flex",
  flexDirection: "row",
  justifyContent: "flex-end",
});

const imageGalleryCss = css({
  margin: "auto",
  paddingX: "3",
  lg: { maxWidth: "1090px" },
  md: { maxWidth: "730px" },
  sm: { maxWidth: "370px" },
});

export default AdminPage;
