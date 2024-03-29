"use client";

import Button from "@/components/atoms/Button/Button";
import { ENDPOINTS } from "@/utils/constants";

const HomePage = () => {
  return (
    <div>
      <h1>HomePage</h1>
      <Button onClick={() => (window.location.href = ENDPOINTS.auth)}>
        Auth
      </Button>
    </div>
  );
};

export default HomePage;
