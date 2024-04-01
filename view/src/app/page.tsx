"use client";

import Button from "@/components/atoms/Button/Button";
import { CLIENT_ENDPOINTS } from "@/utils/constants";

const HomePage = () => {
  return (
    <div>
      <h1>HomePage</h1>
      <Button onClick={() => (window.location.href = CLIENT_ENDPOINTS.auth)}>
        Auth
      </Button>
    </div>
  );
};

export default HomePage;
