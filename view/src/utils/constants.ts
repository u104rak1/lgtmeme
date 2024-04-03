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

export const IMAGE_SIZE = 300;
export const MAX_KEYWORD_LENGTH = 50;
