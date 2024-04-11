"use client";

import AdminImageGallery from "@/components/organisms/AdminImageGallery/AdminImageGallery";
import { css } from "@@/styled-system/css";

const AdminPage = () => {
  return (
    <div>
      <AdminImageGallery css={imageGalleryCss} />
    </div>
  );
};

const imageGalleryCss = css({
  margin: "auto",
  paddingX: "3",
  lg: { maxWidth: "1090px" },
  md: { maxWidth: "730px" },
  sm: { maxWidth: "370px" },
});

export default AdminPage;
