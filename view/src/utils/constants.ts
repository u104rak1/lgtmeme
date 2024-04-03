export const AUTH_ENDPOINTS = {
  login: "/auth-api/login",
} as const;

export const CLIENT_ENDPOINTS = {
  auth: "/client-api/auth",
  images: "/client-api/images",
} as const;

export const PAGE_ENDPOINTS = {
  home: "/",
  createImage: "/create-image",
  privacyPolicy: "/privacy-policy",
  termsOfService: "/terms-of-service",
} as const;
